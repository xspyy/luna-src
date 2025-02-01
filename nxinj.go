package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"main/packages/Memory/instance"
	"main/packages/Memory/memory"
	"main/packages/Memory/utils"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"git.sr.ht/~jackmordaunt/go-toast"
	"github.com/StackExchange/wmi"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/shirou/gopsutil/v3/process"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"golang.design/x/clipboard"
	"golang.org/x/sys/windows"
)

type Win32_Processor struct {
	ProcessorId string
}

type Win32_BIOS struct {
	SerialNumber string
}

type Win32_BaseBoard struct {
	SerialNumber string
}

type Win32_DiskDrive struct {
	SerialNumber string
}

func getHWID() (bool, string) {
	var hwid strings.Builder

	var processors []Win32_Processor
	err := wmi.Query("SELECT ProcessorId FROM Win32_Processor", &processors)
	if err != nil {
		return false, ""
	}
	for _, cpu := range processors {
		hwid.WriteString(strings.TrimSpace(cpu.ProcessorId))
	}

	var bioses []Win32_BIOS
	err = wmi.Query("SELECT SerialNumber FROM Win32_BIOS", &bioses)
	if err != nil {
		return false, ""
	}
	for _, bios := range bioses {
		hwid.WriteString(strings.TrimSpace(bios.SerialNumber))
	}

	var baseBoards []Win32_BaseBoard
	err = wmi.Query("SELECT SerialNumber FROM Win32_BaseBoard", &baseBoards)
	if err != nil {
		return false, ""
	}
	for _, board := range baseBoards {
		hwid.WriteString(strings.TrimSpace(board.SerialNumber))
	}

	var diskDrives []Win32_DiskDrive
	err = wmi.Query("SELECT SerialNumber FROM Win32_DiskDrive", &diskDrives)
	if err != nil {
		return false, ""
	}
	for _, disk := range diskDrives {
		hwid.WriteString(strings.TrimSpace(disk.SerialNumber))
	}

	hash := sha256.Sum256([]byte(hwid.String()))
	return true, hex.EncodeToString(hash[:])
}

func init() {
	clipboard.Init()
}

var (
	ErrNxNotInjected     = errors.New("failed to detect nx injection")
	ErrNxExecuteFailed   = errors.New("a error occured during execution of code")
	ErrNxAlreadyInjected = errors.New("nx already injected")
)

type RbxInstance struct {
	Address uint64
}

type InitChannel struct {
	States    RbxInstance
	Peer0     RbxInstance
	Peer1     RbxInstance
	Instances RbxInstance
}

func CreateFileWithDir(filePath string, content string) error {
	err := os.MkdirAll(filePath, os.ModePerm)
	if err != nil {
		return fmt.Errorf("error creating directory: %w", err)
	}
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}
	defer file.Close()
	_, err = file.WriteString(content)
	if err != nil {
		return fmt.Errorf("error writing to file: %w", err)
	}
	return nil
}

func AppendToFile(filePath string, content string) error {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()
	_, err = file.WriteString(content)
	if err != nil {
		return fmt.Errorf("error writing to file: %w", err)
	}

	return nil
}

func getpath() string {
	pids, err := process.Processes()
	if err == nil {
		for _, id := range pids {
			name, _ := id.Name()
			if name == "RobloxPlayerBeta.exe" {
				dir, err := id.Exe()
				if err == nil {
					return dir
				}

				return fmt.Sprintf("Error:" + err.Error())
			}
		}
	}
	return fmt.Sprintf("Error:" + err.Error())
}

type Websockets struct {
	Url          string
	Conn         *websocket.Conn
	MessagesRecv []string
}

var Cache map[string]*Websockets = make(map[string]*Websockets)

func HandleRequests(conn *websocket.Conn, url string) {
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			delete(Cache, url)
			break
		}
		Data := Cache[url]
		if Data != nil {
			Data.MessagesRecv = append(Cache[url].MessagesRecv, string(msg))
			Cache[url] = Data
		}
	}
}

var imageExtensions = []string{
	".jpg", ".jpeg", ".png", ".gif", ".bmp",
	".tiff", ".tif", ".webp", ".svg", ".ico",
	".heic", ".avif", ".raw", ".psd",
}

var cachedImages = make(map[string]string)

func containsImageExtension(url string) bool {
	url = strings.ToLower(url)
	for _, ext := range imageExtensions {
		if strings.Contains(url, ext) {
			return true
		}
	}
	return false
}

func getLatestRobloxDir(rootDir string) string {

	latestFolder := ""
	var latestModTime time.Time
	filepath.WalkDir(rootDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			return nil
		}
		dirEntries, err := os.ReadDir(path)
		if err != nil {
			return err
		}
		if len(dirEntries) > 0 {
			info, err := d.Info()
			if err != nil {
				return err
			}
			if info.ModTime().After(latestModTime) {
				latestModTime = info.ModTime()
				latestFolder = path
			}
		}
		return nil
	})
	return strings.Replace(latestFolder, `\content`, "", 1)
}

var send sync.Mutex
var recv sync.Mutex

var WriteFile sync.Mutex

func safePath(baseDir, userPath string) (string, error) {
	absoluteBase, err := filepath.Abs(baseDir)
	if err != nil {
		return "", fmt.Errorf("failed to resolve base directory: %w", err)
	}

	absolutePath, err := filepath.Abs(filepath.Join(baseDir, userPath))
	if err != nil {
		return "", fmt.Errorf("failed to resolve user path: %w", err)
	}

	if !strings.HasPrefix(absolutePath, absoluteBase) {
		return "", fmt.Errorf("path escapes the workspace directory")
	}

	return absolutePath, nil
}

func getExecutableDir() (string, error) {
	execPath, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("failed to get executable path: %w", err)
	}
	return filepath.Dir(execPath), nil
}

func Start() {

	execDir, _ := getExecutableDir()
	workspace := filepath.Join(execDir, "workspace")

	router := gin.New()
	router.POST("/bridge", func(ctx *gin.Context) {
		type Websocket struct {
			URL     string `json:"url"`
			Binary  bool   `json:"binary"`
			Message string `json:"message"`
			Close   bool   `json:"closure"`
		}
		var Data struct {
			Type      string    `json:"type"`
			FileName  string    `json:"filename"`
			Source    string    `json:"source"`
			Body      any       `json:"body"`
			Websocket Websocket `json:"websocket"`
		}
		ctx.ShouldBind(&Data)

		Data.FileName = strings.Replace(Data.FileName, "workspace/", "", 1)

		switch Data.Type {
		case "getcustomasset":
			var filePath string

			Data.FileName = "workspace/" + Data.FileName

			if strings.Contains(Data.FileName, "/") {
				re := regexp.MustCompile(`/`)
				filePath = re.ReplaceAllString(Data.FileName, `\\`)
			} else {
				filePath = Data.FileName
			}
			filePath = filepath.Clean(filePath)
			fileName := filepath.Base(filePath)
			fileExt := strings.ToLower(filepath.Ext(fileName))

			if fileExt == ".txt" {
				fileContent, err := ioutil.ReadFile(filePath)
				if err != nil {
					fmt.Printf("Bridge error: failed to read file %s: %v\n", fileName, err)
					//return nil, err
				}
				ctx.JSON(200, gin.H{"status": true, "content": "rbxasset://" + string(fileContent)})
				return
			} else {
				dir := getpath()
				dir = strings.ReplaceAll(dir, `\RobloxPlayerBeta.exe`, "")
				if strings.Contains(dir, "Error:") {
					err := strings.Split(dir, "Error:")[1]
					ctx.JSON(200, gin.H{"status": false, "content": err})
					return
				}

				var parentDir string = strings.ReplaceAll(getLatestRobloxDir(dir), `\content\`, "") + `\content\`

				destPath := parentDir + fileName

				srcFile, err := os.Open(filePath)
				if err != nil {
					ctx.JSON(200, gin.H{"status": false, "content": err.Error()})
					return
				}
				defer srcFile.Close()
				destFile, err := os.Create(destPath)
				if err != nil {
					ctx.JSON(200, gin.H{"status": false, "content": err.Error()})
					return
				}
				defer destFile.Close()
				_, err = io.Copy(destFile, srcFile)
				if err != nil {
					ctx.JSON(200, gin.H{"status": false, "content": err.Error()})
					return
				}
			}
			ctx.JSON(200, gin.H{"status": true, "content": "rbxasset://" + fileName})
		case "loadstring":

			sendA.Lock()
			defer sendA.Unlock()

			if !Config.Active && !Config.OffsetsUpdated {
				return
			}

			mem := Active.Mem

			if mem == nil || Active.Error {
				ctx.JSON(200, gin.H{
					"content": "",
					"status":  false,
				})
				return
			}

			fakedm, _ := Active.Mem.ReadPointer(uintptr(Active.Instances.RenderView) + uintptr(Active.Offsets.DataModelHolder))
			realdm, _ := Active.Mem.ReadPointer(fakedm + uintptr(Active.Offsets.DataModel))
			datamodel := instance.NewInstance(realdm, Active)

			if datamodel.Address < 1000 {
				ctx.JSON(200, gin.H{"status": false, "content": "Unable to find datamodel."})
				return
			}

			CoreGui := datamodel.FindFirstChildOfClass("CoreGui", false)
			bridge := CoreGui.FindFirstChild("Bridge", false)

			ScriptHolder := bridge.FindFirstChild("LoadstringHolder", false)

			if ScriptHolder.Address > 1000 {
				Holder := ScriptHolder.Value()
				if Holder == nil {
					// handle err..
				}

				if strings.Contains(reflect.TypeOf(Holder).String(), "instance") {
					Script := Holder.(instance.Instance)
					Script.SetModuleBypass()
					code, size := bc.Compile(fmt.Sprintf(`return function(...) %v%v%v end`, "\n", Data.Source, "\n"))

					if size == -10 {
						return
					}

					Script.SetBytecode(code, uint64(size))
				}
			}

			ctx.JSON(200, gin.H{"status": true})
		case "require":
			sendA.Lock()
			defer sendA.Unlock()
			if !Config.Active && !Config.OffsetsUpdated {
				return
			}

			mem := Active.Mem

			if mem == nil || Active.Error {
				ctx.JSON(200, gin.H{
					"content": "",
					"status":  false,
				})
				return
			}

			fakedm, _ := Active.Mem.ReadPointer(uintptr(Active.Instances.RenderView) + uintptr(Active.Offsets.DataModelHolder))
			realdm, _ := Active.Mem.ReadPointer(fakedm + uintptr(Active.Offsets.DataModel))
			datamodel := instance.NewInstance(realdm, Active)

			if datamodel.Address < 1000 {
				ctx.JSON(200, gin.H{"status": false, "content": "Unable to find datamodel."})
				return
			}
			CoreGui := datamodel.FindFirstChildOfClass("CoreGui", false)
			Require := CoreGui.FindFirstChild("requireThis", false)

			RequireValue := Require.Value()
			if strings.Contains(reflect.TypeOf(RequireValue).String(), "instance") {
				Script := RequireValue.(instance.Instance)
				Script.SetModuleBypass()
			}

			ctx.JSON(200, gin.H{"status": true})
		case "isfile":

			safeResolvedPath, err := safePath(workspace, Data.FileName)
			if err != nil {
				ctx.JSON(200, gin.H{"status": false})
				return
			}

			file, err := os.Stat(safeResolvedPath)
			if err == nil {
				if file.IsDir() {
					ctx.JSON(200, gin.H{"status": false})
					return
				}
				ctx.JSON(200, gin.H{"status": true})
				return
			}

			ctx.JSON(200, gin.H{"status": os.IsExist(err)})
			return
		case "writefile":
			for fileName, test := range cachedImages {
				if strings.EqualFold(test, Data.Source) {
					Data.Source = test
					delete(cachedImages, fileName)
					break
				}
			}

			safeResolvedPath, err := safePath(workspace, Data.FileName)
			if err != nil {
				ctx.JSON(200, gin.H{"status": false, "content": err.Error()})
				return
			}

			err = os.WriteFile(safeResolvedPath, []byte(Data.Source), 0644)
			if err != nil {
				ctx.JSON(200, gin.H{
					"content": err.Error(),
					"status":  false,
				})
				return
			}
			resp, err := os.ReadFile(safeResolvedPath)
			if err != nil {
				ctx.JSON(200, gin.H{
					"content": err.Error(),
					"status":  false,
				})
				return
			}
			ctx.JSON(200, gin.H{"status": true, "content": string(resp)})
		case "appendfile":
			WriteFile.Lock()
			defer WriteFile.Unlock()

			safeResolvedPath, err := safePath(workspace, Data.FileName)
			if err != nil {
				ctx.JSON(200, gin.H{"status": false, "content": err.Error()})
				return
			}

			info, err := os.ReadFile(safeResolvedPath)
			if err == nil {
				new := string(info) + Data.Source
				os.WriteFile(safeResolvedPath, []byte(new), 0644)
				ctx.JSON(200, gin.H{"status": true})
			} else {
				ctx.JSON(200, gin.H{"status": false, "content": err.Error()})
			}
		case "delfile", "delfolder":

			safeResolvedPath, err := safePath(workspace, Data.FileName)
			if err != nil {
				ctx.JSON(200, gin.H{"status": false, "content": err.Error()})
				return
			}

			err = os.Remove(safeResolvedPath)
			if err != nil && strings.Contains(err.Error(), "The directory is not empty.") {
				os.RemoveAll(safeResolvedPath)
				os.Remove(safeResolvedPath)
			}
			ctx.JSON(200, gin.H{"status": true})
		case "request":

			var RobloxRequest struct {
				Body    string                 `json:"Body"`
				URL     string                 `json:"Url"`
				Method  string                 `json:"Method"`
				Headers map[string]interface{} `json:"Headers"`
			}
			json.Unmarshal([]byte(Data.Body.(string)), &RobloxRequest)

			req, err := http.NewRequest(RobloxRequest.Method, RobloxRequest.URL, bytes.NewBufferString(RobloxRequest.Body))
			if err != nil {
				//
			}
			for name, value := range RobloxRequest.Headers {
				req.Header.Add(name, fmt.Sprintf("%v", value))
			}

			_, hwid := getHWID()
			if req.Header.Get("User-Agent") == "" {
				req.Header.Add("User-Agent", "NX/1.0")
			}

			req.Header.Add("NX-Fingerprint", hwid)
			req.Header.Add("Exploit-Guid", hwid)
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				//
			}

			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				//
			}

			ToReturn := string(body)

			if containsImageExtension(RobloxRequest.URL) {
				cachedImages[RobloxRequest.URL] = ToReturn
			}

			ctx.JSON(200, gin.H{
				"Success":       true,
				"StatusCode":    resp.StatusCode,
				"StatusMessage": resp.Status,
				"Headers":       resp.Header,
				"Body":          ToReturn,
			})
		case "websocket":
			if Data.Websocket.URL != "" {
				_, hwid := getHWID()
				conn2, resp, err := websocket.DefaultDialer.Dial(Data.Websocket.URL, map[string][]string{
					"User-Agent":         {"NX/1.0"},
					"NX-Fingerprint":     {hwid},
					"NX-User-Identifier": {hwid},
				})
				if resp.StatusCode == 101 {
					if conn2 != nil && err == nil {
						Cache[Data.Websocket.URL] = &Websockets{
							Url:  Data.Websocket.URL,
							Conn: conn2,
						}
						go HandleRequests(conn2, Data.Websocket.URL)
						ctx.JSON(200, gin.H{"status": true})
					} else {
						ctx.JSON(200, gin.H{"status": false})
					}
					return
				}
				ctx.JSON(200, gin.H{"status": false})
			} else {
				ctx.JSON(200, gin.H{"status": false, "content": "Error: Url is empty!"})
			}
		case "websocketsend":
			ws := Cache[Data.Websocket.URL]
			if ws != nil {
				ws.Conn.WriteMessage(websocket.TextMessage, []byte(Data.Websocket.Message))
			}
		case "websocketclose":
			ws := Cache[Data.Websocket.URL]
			if ws != nil {
				ws.Conn.Close()
				delete(Cache, Data.Websocket.URL)
			}
		}
	})
	router.GET("/bridge", func(ctx *gin.Context) {
		act := ctx.Query("action")
		switch act {
		case "setthreadidentity":

			identity := ctx.Query("identity")

			mem := Active.Mem
			if mem == nil || !memory.IsHandleValid(mem.ProcessHandle) {
				ctx.JSON(200, gin.H{
					"status":  false,
					"content": 0,
				})
				return
			}

			ident, err := strconv.Atoi(identity)
			if err != nil {
				ctx.JSON(200, gin.H{
					"status":  false,
					"content": 0,
				})
				return
			}

			fakedm, _ := Active.Mem.ReadPointer(uintptr(Active.Instances.RenderView) + uintptr(Active.Offsets.DataModelHolder))
			realdm, _ := Active.Mem.ReadPointer(fakedm + uintptr(Active.Offsets.DataModel))
			game := instance.NewInstance(realdm, Active)

			ScriptContext := game.WaitForClass("ScriptContext", 5)

			if ident > 9 {
				ident = 8
			}

			go ScriptContext.ApplyCapacity(ident, utils.Capabilities[ident], 10, "JestGlobals", "LuaSocialLibrariesDeps", "Url")

			ctx.JSON(200, gin.H{
				"status": true,
			})
		case "getthreadidentity":

			mem := Active.Mem
			if mem == nil || !memory.IsHandleValid(mem.ProcessHandle) {
				ctx.JSON(200, gin.H{
					"status":  false,
					"content": 0,
				})
				return
			}

			fakedm, _ := Active.Mem.ReadPointer(uintptr(Active.Instances.RenderView) + uintptr(Active.Offsets.DataModelHolder))
			realdm, _ := Active.Mem.ReadPointer(fakedm + uintptr(Active.Offsets.DataModel))
			game := instance.NewInstance(realdm, Active)

			ScriptContext := game.WaitForClass("ScriptContext", 5)
			Identity := ScriptContext.GetIdentity(5, "JestGlobals", "LuaSocialLibrariesDeps", "Url")

			ctx.JSON(200, gin.H{
				"status":  true,
				"content": Identity,
			})
		case "gethwid":
			success, hwid := getHWID()
			if success {
				ctx.JSON(200, gin.H{
					"content": hwid,
					"status":  true,
				})
			} else {
				ctx.JSON(200, gin.H{
					"content": "Error: Unable to get HWID",
					"status":  false,
				})
			}
		case "gameentered":

			RV := Active

			fakedm, _ := RV.Mem.ReadPointer(uintptr(RV.Instances.RenderView) + uintptr(RV.Offsets.DataModelHolder))
			realdm, _ := RV.Mem.ReadPointer(fakedm + uintptr(RV.Offsets.DataModel))
			DM := instance.NewInstance(realdm, RV)

			if !RV.Injected {
				return
			}

			mem := RV.Mem

			if mem == nil || RV.Error {
				return
			}

			if DM.Address < 1000 {
				return
			}

			CoreGui := DM.FindFirstChildOfClass("CoreGui", false)
			CoreGui.WaitForChild("Bridge", 10)

			go func(DM instance.Instance) {
				ScriptContext := DM.WaitForClass("ScriptContext", 5)
				if ScriptContext.Address > 1000 {
					go ScriptContext.ApplyCapacity(8, 0xEFFFFFFFFFFFFFFF, 10, "JestGlobals", "LuaSocialLibrariesDeps", "Url")
				}
			}(DM)

			var sources []string

			filepath.WalkDir("./autoexec", func(path string, d fs.DirEntry, err error) error {
				source, _ := os.ReadFile(path)
				if string(source) != "" {
					sources = append(sources, string(source))
				}
				return nil
			})

			for _, s := range sources {
				app.Execute(s)
				time.Sleep(time.Millisecond * 250)
			}

			data := ctx.Query("name")
			(&toast.Notification{
				AppID: "NXInjector",
				Title: "Youve joined a game!",
				Body:  fmt.Sprintf("Welcome %v", data),
			}).Push()

		case "setclipboard":
			data := ctx.Query("value")
			clipboard.Write(clipboard.FmtText, []byte(data))
			ctx.JSON(200, gin.H{"status": true})
		case "readfile":
			fileName := ctx.Query("file")

			safeResolvedPath, err := safePath(workspace, fileName)
			if err != nil {
				ctx.JSON(200, gin.H{"status": false, "content": err.Error()})
				return
			}

			content, err := os.ReadFile(safeResolvedPath)
			if err != nil {
				ctx.JSON(200, gin.H{
					"content": err.Error(),
					"status":  false,
				})
				return
			}
			ctx.JSON(200, gin.H{
				"content": string(content),
				"status":  true,
			})
		case "listfiles":
			path := ctx.Query("file")

			safeResolvedPath, err := safePath(workspace, path)
			if err != nil {
				ctx.JSON(200, gin.H{"status": false, "content": err.Error()})
				return
			}

			var files []string
			items, err := ioutil.ReadDir(safeResolvedPath)
			if err != nil {
				ctx.JSON(200, gin.H{"content": err.Error(), "status": false})
				return
			}
			for _, item := range items {
				fullPath := filepath.Join(path, item.Name())
				files = append(files, fullPath)
			}
			ctx.JSON(200, gin.H{"content": files, "status": true})
		case "makefolder":
			fileName := ctx.Query("file")

			safeResolvedPath, err := safePath(workspace, fileName)
			if err != nil {
				ctx.JSON(200, gin.H{"status": false, "content": err.Error()})
				return
			}

			err = CreateFileWithDir(safeResolvedPath, "")
			if err != nil {
				if strings.Contains(err.Error(), "is a directory") {
					ctx.JSON(200, gin.H{
						"content": "",
						"status":  true,
					})
					return
				} else {
					ctx.JSON(200, gin.H{
						"content": err.Error(),
						"status":  false,
					})
				}
				return
			}
			ctx.JSON(200, gin.H{
				"content": "",
				"status":  true,
			})
		case "isfolder":

			safeResolvedPath, err := safePath(workspace, ctx.Query("file"))
			if err != nil {
				ctx.JSON(200, gin.H{"status": false, "content": err.Error()})
				return
			}

			info, err := os.Stat(safeResolvedPath)
			if os.IsNotExist(err) {
				ctx.JSON(200, gin.H{"status": false})
				return
			}
			ctx.JSON(200, gin.H{"status": info.IsDir()})
		case "setscriptable":

			mem := Active.Mem
			if mem == nil || !memory.IsHandleValid(mem.ProcessHandle) {
				ctx.JSON(200, gin.H{
					"status": false,
				})

				return
			}

			fakedm, _ := Active.Mem.ReadPointer(uintptr(Active.Instances.RenderView) + uintptr(Active.Offsets.DataModelHolder))
			realdm, _ := Active.Mem.ReadPointer(fakedm + uintptr(Active.Offsets.DataModel))
			datamodel := instance.NewInstance(realdm, Active)

			CoreGui := datamodel.FindFirstChildOfClass("CoreGui", false)
			Bridge := CoreGui.FindFirstChild("Bridge", false)
			ModuleHolder := Bridge.FindFirstChild("ModuleHolder", false)

			if Active.Injected && ModuleHolder.Address > 1000 {

				setAddrHolder := CoreGui.FindFirstChild("setaddressholder", false)
				setAddress := CoreGui.FindFirstChild("setaddressholder_bool", false)

				Holder := setAddrHolder.Value()
				IsBool := setAddress.Value()

				var Bool bool
				var Instance instance.Instance
				if IsBool != nil && reflect.TypeOf(IsBool).Kind() == reflect.Bool {
					Bool = IsBool.(bool)
				}
				if Holder != nil {
					switch true { // instance.Instance
					case reflect.TypeOf(Holder).String() == "instance.Instance":
						Instance = Holder.(instance.Instance)
					}
				}

				if Instance.Address > 1000 {
					Classes := Instance.ClassDescriptor()
					if Classes != nil {
						value := Classes.PropertyDescriptors(mem).Get(mem, ctx.Query("value"))
						if value != nil {
							value.SetScriptable(mem, Bool)
							ctx.JSON(200, gin.H{
								"status": true,
							})
							return
						}
					}
				}
			}
			ctx.JSON(200, gin.H{
				"status": false,
			})
			return
		case "getproperties":

			mem := Active.Mem
			if mem == nil || !memory.IsHandleValid(mem.ProcessHandle) {
				ctx.JSON(200, gin.H{
					"status": false,
				})
				return
			}

			fakedm, _ := Active.Mem.ReadPointer(uintptr(Active.Instances.RenderView) + uintptr(Active.Offsets.DataModelHolder))
			realdm, _ := Active.Mem.ReadPointer(fakedm + uintptr(Active.Offsets.DataModel))
			datamodel := instance.NewInstance(realdm, Active)

			CoreGui := datamodel.FindFirstChildOfClass("CoreGui", false)
			Bridge := CoreGui.FindFirstChild("Bridge", false)
			ModuleHolder := Bridge.FindFirstChild("ModuleHolder", false)

			if Active.Injected && ModuleHolder.Address > 1000 {
				setAddrHolder := CoreGui.FindFirstChild("setaddressholder", false)
				if setAddrHolder.Address > 1000 {
					classes := setAddrHolder.ClassDescriptor()
					if classes != nil && classes.Address > 1000 {
						properties := classes.PropertyDescriptors(mem)
						if properties != nil && properties.Address > 1000 {
							var data []string
							for _, inst := range properties.GetAllYield(mem) {
								data = append(data, fmt.Sprintf("%v:%v", inst.Name(mem), inst.IsHiddenValue(mem)))
							}
							ctx.JSON(200, gin.H{
								"status":  len(data) > 0,
								"content": data,
							})
						}
					}
				}
			}
			return
		case "getscriptbytecode":

			mem := Active.Mem
			if mem == nil || !memory.IsHandleValid(mem.ProcessHandle) {
				ctx.JSON(200, gin.H{
					"status": false,
				})
				return
			}

			fakedm, _ := Active.Mem.ReadPointer(uintptr(Active.Instances.RenderView) + uintptr(Active.Offsets.DataModelHolder))
			realdm, _ := Active.Mem.ReadPointer(fakedm + uintptr(Active.Offsets.DataModel))
			datamodel := instance.NewInstance(realdm, Active)

			CoreGui := datamodel.FindFirstChildOfClass("CoreGui", false)
			Bridge := CoreGui.FindFirstChild("Bridge", false)
			ModuleHolder := Bridge.FindFirstChild("ModuleHolder", false)

			if Active.Injected && ModuleHolder.Address > 1000 {

				Bytecode := CoreGui.FindFirstChild("getbytecodeobject", false)
				Holder := Bytecode.Value()
				var Instance instance.Instance
				switch true { // instance.Instance
				case reflect.TypeOf(Holder).String() == "instance.Instance":
					Instance = Holder.(instance.Instance)
				}

				if Instance.Address > 1000 {
					data, _ := Instance.GetByteCode()
					ctx.JSON(200, gin.H{
						"status":  true,
						"content": string(data),
					})
					return
				}

			}
			ctx.JSON(200, gin.H{
				"status":  false,
				"content": "Error grabbing values",
			})
			return
		default:
			switch true {
			case strings.Contains(act, "_message"):

				send.Lock()
				defer send.Unlock()

				url := strings.Split(act, "_message")[0]
				if Cache[url] == nil {
					ctx.JSON(200, gin.H{"status": false, "content": []string{}})
					return
				}
				data := Cache[url]
				if len(data.MessagesRecv) > 1 {
					ctx.JSON(200, gin.H{"status": true, "content": data.MessagesRecv})
					data.MessagesRecv = []string{}
					Cache[url] = data
					return
				}
				ctx.JSON(200, gin.H{"status": true, "content": []string{}})
			case strings.Contains(act, "_close"):

				recv.Lock()
				defer recv.Unlock()

				url := strings.Split(act, "_close")[0]
				if Cache[url] == nil {
					// return true!
					ctx.JSON(200, gin.H{"status": true})
					return
				}
				ctx.JSON(200, gin.H{"status": false})
			}
		}
	})
	router.Run("localhost:3000")
}

var Roblox []*instance.RobloxInstances
var Active *instance.RobloxInstances
var First bool

func HandleClosures() {
	for {
		var New []*instance.RobloxInstances
		for _, inst := range Roblox {
			if !inst.Error {
				New = append(New, inst)
			} else {
				runtime.EventsEmit(app.ctx, "robloxClosure", inst)
			}
		}

		Roblox = New

		time.Sleep(time.Millisecond * 250)
	}
}

var Patches map[uint32]bool = make(map[uint32]bool)

func GetRobloxProccesses() {
	go HandleClosures()
	for {
		var Pids map[int64]*instance.RobloxInstances = make(map[int64]*instance.RobloxInstances)
		for _, inst := range Roblox {
			Pids[inst.Pid] = inst
		}
		ok, pid := memory.IsProcessRunning()

		if ok {
			var WG sync.WaitGroup
			for _, instances := range pid {
				WG.Add(1)
				go func(instances memory.Processes) {
					defer WG.Done()
					if _, ok := Pids[int64(instances.Pid)]; !ok {
						if instances.Name == "Windows10Universal.exe" {
							return
						}

						getoffsets := func() utils.Offsets {
							if instances.Name == "Windows10Universal.exe" {
								return utils.OffsetsDataUwp
							}
							return utils.OffsetsDataPlayer
						}

						mem, err := memory.NewLuna(instances.Pid)
						if err != nil || mem == nil {
							return
						}

						if !Patches[instances.Pid] {
							instance.PatchRoblox(mem)
							Patches[instances.Pid] = true
							Logger.Debug(fmt.Sprintf("[%v] Roblox succesfully patched.", instances.Pid))
						}

						renderview, _ := mem.AOBSCANALL("52 65 6e 64 65 72 4a 6f 62 28 45 61 72 6c 79 52 65 6e 64 65 72 69 6e 67 3b", false, 1)

						Logger.Debug(fmt.Sprintf("[%v] Fetched renderjob %v", instances.Pid, renderview))

						if len(renderview) > 0 {

							var New *instance.RobloxInstances = &instance.RobloxInstances{
								Pid:     int64(instances.Pid),
								ExeName: instances.Name,
								Mem:     mem,
								Instances: instance.Instances{
									RobloxBase: uint64(mem.RobloxBase),
								},
								Offsets: getoffsets(),
							}

							rv, _ := mem.ReadPointer(renderview[0] + uintptr(New.Offsets.RenderViewFromRenderJob))
							New.Instances.RenderView = uint64(rv)

							fakedm, _ := mem.ReadPointer(rv + uintptr(New.Offsets.DataModelHolder))
							realdm, _ := mem.ReadPointer(fakedm + uintptr(New.Offsets.DataModel))
							DM := instance.NewInstance(realdm, New)

							name := DM.Name()

							if name == "None" {
								return
							}

							Logger.Debug(fmt.Sprintf("[%v] Got datamodel name %v", instances.Pid, name))

							check := func() bool {
								if name == "App" || name == "LuaApp" {
									return false
								}
								CoreGui := DM.FindFirstChildOfClass("CoreGui", false)
								Bridge := CoreGui.WaitForChild("Bridge", 1)
								return Bridge.Address > 1000
							}

							New.Injected = check()

							user := func() string {
								Players := DM.FindFirstChildOfClass("Players", false)
								Player := Players.LocalPlayer()
								return Player.Name()
							}

							avatar := func(ign string) string {
								return app.Avatar(ign)
							}

							New.Username = user()
							New.Avatar = avatar(New.Username)

							Roblox = append(Roblox, New)

							go DataModelHandler(New)

						}

					}
				}(instances)
			}
			WG.Wait()
			if !First || (First && Active == nil && len(Pids) > 0) {
				for name, value := range Pids {
					if value.ExeName != "Windows10Universal.exe" {
						First = true
						for _, v := range Roblox {
							if v.Pid == name {
								Active = v
								break
							}
						}
						break
					}
				}
			}
		}
		time.Sleep(time.Millisecond * 150)
	}
}

var (
	kernel32                = windows.NewLazySystemDLL("kernel32.dll")
	procWaitForSingleObject = kernel32.NewProc("WaitForSingleObject")
)

func DataModelHandler(RV *instance.RobloxInstances) {
	go func(RV *instance.RobloxInstances) {
		for !RV.Error {
			runtime.EventsEmit(app.ctx, "robloxUpdate", RV)
			time.Sleep(time.Second)
		}
	}(RV)
	for {
		r, _, _ := procWaitForSingleObject.Call(uintptr(RV.Mem.ProcessHandle), 0xFFFFFFFF)
		if r == 0 {
			RV.Error = true
			RV.Injected = false
			break
		}
	}
	/*
		for !RV.Error {

			fakedm, _ := RV.Mem.ReadPointer(uintptr(RV.Instances.RenderView) + uintptr(RV.Offsets.DataModelHolder))
			realdm, _ := RV.Mem.ReadPointer(fakedm + uintptr(RV.Offsets.DataModel))
			DM := instance.NewInstance(realdm, RV)

			name := DM.Name()

			if name == "App" || name == "LuaApp" && wait_for_detection {
			}

			if RV.Ingame && name == "App" || name == "LuaApp" {
				if RV.Ingame {
					RV.Ingame = false
				}
				RV.Inmenu = true
			} else if !RV.Ingame && name == "Ugc" || name == "Game" {
				RV.Ingame = true
				if RV.Inmenu {
					RV.Inmenu = false
				}
			}

		}
	*/
}
