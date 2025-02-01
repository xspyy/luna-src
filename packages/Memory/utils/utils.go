package utils

import (
	"bufio"
	"bytes"
	"fmt"
	"io/fs"
	"io/ioutil"
	"main/packages/Memory/memory"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"unsafe"
)

type Offsets struct {
	RenderViewFromRenderJob uint64
	DataModelHolder         uint64
	DataModel               uint64
	VisualDataModel         uint64

	Name            uint64
	Children        uint64
	Parent          uint64
	ClassDescriptor uint64
	LocalPlayer     uint64

	ValueBase   uint64
	ModuleFlags uint64
	IsCore      uint64

	PlaceID uint64

	BytecodeSize uint64
	Bytecode     map[string]uint64

	OffsetTaskScheduler uint64
	OffsetJobsContainer uint64
}

var OffsetsDataPlayer = Offsets{
	RenderViewFromRenderJob: 0x1E8,
	DataModelHolder:         0x118,
	DataModel:               0x1A8,
	VisualDataModel:         0x720,

	Name:            0x68,
	Children:        0x70,
	Parent:          0x50,
	ClassDescriptor: 0x18,
	ValueBase:       0xC8,
	LocalPlayer:     0x118,

	ModuleFlags: 0x1B0 - 0x4,
	IsCore:      0x1B0,

	PlaceID: 0x170,

	BytecodeSize: 0xA8,
	Bytecode: map[string]uint64{
		"LocalScript":  0x1C0,
		"ModuleScript": 0x168,
	},

	OffsetTaskScheduler: 0x5C71FC8, //0x5C19D28,
	OffsetJobsContainer: 0x1C8,
}

var OffsetsDataUwp = Offsets{
	DataModelHolder: 0x118,
	DataModel:       0x1A8,
	VisualDataModel: 0x720,

	Name:            0x68,
	Children:        0x70,
	Parent:          0x50,
	ClassDescriptor: 0x18,
	ValueBase:       0xC8,
	LocalPlayer:     0x118,

	ModuleFlags: 0x1B0 - 0x4,
	IsCore:      0x1B0,

	BytecodeSize: 0xA8,
	Bytecode: map[string]uint64{
		"LocalScript":  0x1C0,
		"ModuleScript": 0x168,
	},
}

var Capabilities map[int]uint64 = map[int]uint64{
	9: 0xEFFFFFFFFFFFFFFF,
	8: 0xEFFFFFFFFFFFFFFF,
	7: 0xEFFFFFFFFFFFFFFF,
	6: 0xEFFFFFFFFFFFFFFF,
	5: 0xEFFFFFFFFFFFFFFF,
	4: 0xEFFFFFFFFFFFFFFF,
	3: 0xEFFFFFFFFFFFFFFF,
	2: 0xEFFFFFFFFFFFFFFF,
	1: 0xEFFFFFFFFFFFFFFF,
	0: 0xEFFFFFFFFFFFFFFF,
}

var appdata = os.Getenv("LOCALAPPDATA")
var uwpLogs string

func init() {
	go func() {
		uwpLogs = func() string {
			var data map[string]any = make(map[string]any)
			filepath.Walk(filepath.Join(appdata, `\Packages\`), func(path string, info fs.FileInfo, err error) error {
				if strings.Contains(path, "ROBLOXCORPORATION.ROBLOX") {
					data[strings.Split(strings.Split(path, `\Packages\`)[1], `\`)[0]] = info.ModTime()
					return nil
				}
				return nil
			})
			var rbxLog string
			for name, _ := range data {
				rbx := filepath.Join(appdata, "Packages", name, "LocalState", "logs")
				return rbx
			}
			return rbxLog
		}()
	}()
}

func retrieveRV(mem *memory.Luna, offsets Offsets, logs string) (uint64, uint64) {
	logsDir := logs
	if stat, err := os.Stat(logsDir); os.IsNotExist(err) || !stat.IsDir() {
		fmt.Println(stat, err)
		return 0, 0
	}

	var logFiles []string

	files, err := ioutil.ReadDir(logsDir)
	if err != nil {
		fmt.Println(err)
		return 0, 0
	}

	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".log" {
			logFiles = append(logFiles, filepath.Join(logsDir, file.Name()))
		}
	}

	if len(logFiles) == 0 {
		fmt.Println("Logfile = 0 Error")
		return 0, 0
	}

	sort.Slice(logFiles, func(i, j int) bool {
		info1, err := os.Stat(logFiles[i])
		if err != nil {
			return false
		}
		info2, err := os.Stat(logFiles[j])
		if err != nil {
			return false
		}
		return info1.ModTime().After(info2.ModTime())
	})

	var lockedFiles []string
	for _, logPath := range logFiles {
		if err := os.Remove(logPath); err != nil {
			lockedFiles = append(lockedFiles, logPath)
		}
	}

	if len(lockedFiles) == 0 {
		fmt.Println("Locked files Error")
		return 0, 0
	}

	for _, logPath := range lockedFiles {

		data, err := os.ReadFile(logPath)
		if err != nil {
			continue
		}

		scanner := bufio.NewScanner(bytes.NewBuffer(data))
		for scanner.Scan() {
			line := scanner.Text()
			match := strings.Index(line, "view(")
			if match != -1 {
				address := line[match+5 : match+5+16]
				renderView, err := strconv.ParseUint(address, 16, 64)
				if err != nil {
					continue
				}

				var fakeDatamodel uint64
				var realDatamodel uint64

				mem.MemRead(uintptr(renderView+offsets.DataModelHolder), unsafe.Pointer(&fakeDatamodel), unsafe.Sizeof(fakeDatamodel))
				mem.MemRead(uintptr(fakeDatamodel+offsets.DataModel), unsafe.Pointer(&realDatamodel), unsafe.Sizeof(realDatamodel))

				if realDatamodel != 0 {
					namePointer, _ := mem.ReadPointer(uintptr(realDatamodel + offsets.Name))
					name, _ := mem.ReadRbxStr(namePointer)
					if name == "Ugc" ||
						name == "LuaApp" ||
						name == "App" ||
						name == "Game" {
						return renderView, realDatamodel
					}
				}
			}
		}
	}
	return 0, 0
}

func getRV(mem *memory.Luna, offsets Offsets, UWP bool) (uint64, uint64) {
	if UWP {
		return retrieveRV(mem, offsets, uwpLogs)
	}
	return retrieveRV(mem, offsets, filepath.Join(appdata, "Roblox", "logs"))
}

func GetRenderVDM(pid uint32, mem *memory.Luna, offsets Offsets, UWP bool) uint64 {
	if mem == nil {
		return 0
	}

	task, _ := mem.ReadPointer(mem.RobloxBase + uintptr(offsets.OffsetTaskScheduler))
	task_list, _ := mem.ReadPointer(task + uintptr(offsets.OffsetJobsContainer))
	for i := 0; i < 0x500; i = i + 0x10 {
		task_container, _ := mem.ReadPointer(task_list + uintptr(i))
		name, _ := mem.ReadString(task_container+0x90, 10)
		if strings.EqualFold(name, "RenderJob") {
			return uint64(task_container)
		}
	}
	return 0
	/*
		renderView, dm := getRV(mem, offsets, UWP)

		var visualEngine uint64
		err := mem.MemRead(uintptr(renderView+0x10), unsafe.Pointer(&visualEngine), unsafe.Sizeof(visualEngine))
		if err != nil {
			return renderView, dm, visualEngine
		}

		return renderView, dm, visualEngine
	*/
}
