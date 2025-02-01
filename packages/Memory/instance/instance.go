package instance

import (
	"fmt"
	"main/packages/Memory/classdescriptor"
	"main/packages/Memory/memory"
	"main/packages/Memory/utils"
	"strings"
	"time"
	"unsafe"

	"github.com/shirou/gopsutil/v3/process"
	"golang.org/x/sys/windows"
)

type Instance struct {
	Address uintptr
	Mem     *RobloxInstances
}

type RobloxInstances struct {
	Error    bool
	Injected bool
	Username string
	Pid      int64
	ExeName  string
	//Ingame    bool
	//Inmenu    bool
	Avatar    string
	Mem       *memory.Luna
	Instances Instances
	Offsets   utils.Offsets
}
type Instances struct {
	RenderView uint64
	RobloxBase uint64
}

type Modules struct {
	Name     string
	Address  uint64
	LuaState uint64
}

type Level8 struct {
	CoreGuiContainer, CoreGuiToModules, ModulesToInstances, InstancesToChildren, ModuleScript uint64
	//
	ToIdentity, ToCapabilities, LuaState uint64
}

var Offsets Level8 = Level8{
	CoreGuiContainer:    0x390,
	CoreGuiToModules:    0x8,
	ModulesToInstances:  0x78,
	InstancesToChildren: 0x10,
	ModuleScript:        0x50,

	ToIdentity:     0x30,
	ToCapabilities: 0x48,
	LuaState:       0x6b0,
}

func NewInstance(address uintptr, Mem *RobloxInstances) Instance {
	return Instance{
		Address: address,
		Mem:     Mem,
	}
}

func (inst *Instance) ClassDescriptor() *classdescriptor.ClassDescriptor {

	if inst == nil || inst.Mem == nil || inst.Mem.Mem == nil || inst.Address < 1000 {
		return &classdescriptor.ClassDescriptor{}
	}

	mem := inst.Mem.Mem

	addr, _ := mem.ReadPointer(inst.Address + uintptr(inst.Mem.Offsets.ClassDescriptor))
	return classdescriptor.NewClassDescriptor(addr)
}

func (inst *Instance) ClassName() string {
	if inst == nil || inst.Mem == nil || inst.Mem.Mem == nil || inst.Address < 1000 {
		return "None"
	}

	mem := inst.Mem.Mem

	return inst.ClassDescriptor().Name(mem)
}

func (inst *Instance) Name() string {
	if inst == nil || inst.Mem == nil || inst.Mem.Mem == nil || inst.Address < 1000 {
		return "None"
	}

	mem := inst.Mem.Mem

	namePointer, _ := mem.ReadPointer(inst.Address + uintptr(inst.Mem.Offsets.Name))
	name, _ := mem.ReadRbxStr(namePointer)
	return name
}

func (inst *Instance) Parent() Instance {

	if inst == nil || inst.Mem == nil || inst.Mem.Mem == nil || inst.Address < 1000 {
		return NewInstance(0, nil)
	}

	mem := inst.Mem.Mem

	parentAddr, _ := mem.ReadPointer(inst.Address + uintptr(inst.Mem.Offsets.Parent))
	return NewInstance(parentAddr, inst.Mem)
}

func (inst *Instance) LocalPlayer() Instance {

	if inst == nil || inst.Mem == nil || inst.Mem.Mem == nil || inst.Address < 1000 {
		return NewInstance(0, inst.Mem)
	}

	mem := inst.Mem.Mem

	if inst.Address < 1000 {
		return NewInstance(0, inst.Mem)
	}
	localplayer, _ := mem.ReadPointer(inst.Address + uintptr(inst.Mem.Offsets.LocalPlayer))
	return NewInstance(localplayer, inst.Mem)
}

func (inst *Instance) GetChildren() (Data []Instance) {
	if inst == nil || inst.Mem == nil || inst.Mem.Mem == nil || inst.Address < 1000 {
		return []Instance{NewInstance(0, inst.Mem)}
	}

	mem := inst.Mem.Mem

	childrenPointer, _ := mem.ReadPointer(inst.Address + uintptr(inst.Mem.Offsets.Children))
	top, _ := mem.ReadPointer(childrenPointer)
	end, _ := mem.ReadPointer(childrenPointer + 0x8)
	var ContinueIfFound bool
	for childAddr := top; childAddr < end; childAddr += 0x10 {
		if ContinueIfFound {
			ContinueIfFound = false
			continue
		}
		childPtr, _ := mem.ReadPointer(childAddr)
		if childPtr < 1000 {
			continue
		}
		child := NewInstance(childPtr, inst.Mem)
		n := child.Name()
		if n == "MarketplaceService" {
			ContinueIfFound = true
			continue
		}
		Data = append(Data, child)
	}
	return Data
}

func mostFrequentUint64(nums []uint64) (uint64, int) {
	frequencyMap := make(map[uint64]int)
	for _, num := range nums {
		frequencyMap[num]++
	}
	var mostFrequentNum uint64
	maxFrequency := 0
	for num, count := range frequencyMap {
		if count > maxFrequency {
			maxFrequency = count
			mostFrequentNum = num
		}
	}

	return mostFrequentNum, maxFrequency
}

func (inst *Instance) GetIdentity(timeout int, className ...string) uint64 {

	if inst == nil || inst.Mem == nil || inst.Mem.Mem == nil || inst.Address < 1000 {
		return 8
	}

	mem := inst.Mem.Mem

	Modules := inst.GetRunningScripts(timeout, className...)
	var (
		Identities []uint64
	)
	if len(Modules) > 0 {
		for _, Module := range Modules {
			var identity int
			mem.MemRead(uintptr(Module.LuaState+Offsets.ToIdentity), unsafe.Pointer(&identity), unsafe.Sizeof(identity))
			Identities = append(Identities, uint64(identity))
		}
	}
	most, _ := mostFrequentUint64(Identities)
	return most
}

func (inst *Instance) GetRunningScripts(timeout int, className ...string) (M []Modules) {

	if inst == nil || inst.Mem == nil || inst.Mem.Mem == nil || inst.Address < 1000 {
		return
	}

	mem := inst.Mem.Mem

	for i := 0; i < timeout*10; i++ {
		var (
			PointerToList      uint64
			PointerToModules   uint64
			PointerToInstances uint64
			Looping            uint64
			BlankCounter       = 0
			ModulePointer      uint64
			ModuleScript       uint64
		)

		// find jestglobals
		mem.MemRead(inst.Address+uintptr(Offsets.CoreGuiContainer), unsafe.Pointer(&PointerToList), unsafe.Sizeof(PointerToList))
		mem.MemRead(uintptr(PointerToList+Offsets.CoreGuiToModules), unsafe.Pointer(&PointerToModules), unsafe.Sizeof(PointerToModules))
		mem.MemRead(uintptr(PointerToModules+Offsets.ModulesToInstances), unsafe.Pointer(&PointerToInstances), unsafe.Sizeof(PointerToInstances))
		mem.MemRead(uintptr(PointerToInstances+Offsets.InstancesToChildren), unsafe.Pointer(&Looping), unsafe.Sizeof(Looping))

		for i := 0x10; i < 0x10*1500; i = i + 0x10 {
			mem.MemRead(uintptr(Looping+0x10), unsafe.Pointer(&ModulePointer), unsafe.Sizeof(ModulePointer))

			Looping = ModulePointer
			mem.MemRead(uintptr(Looping+Offsets.ModuleScript), unsafe.Pointer(&ModuleScript), unsafe.Sizeof(ModuleScript))
			Module := NewInstance(uintptr(ModuleScript), inst.Mem)
			name := Module.Name()
			if name == "" {
				BlankCounter++
			} else {
				for _, classN := range className {
					if strings.Contains(name, classN) {
						M = append(M, Modules{
							Name:     classN,
							Address:  uint64(Module.Address),
							LuaState: ModulePointer,
						})
					}
				}
				if BlankCounter > 0 {
					BlankCounter = 0
				}
			}
			if BlankCounter >= 20 {
				break
			}
		}
		if len(M) > 0 {
			break
		}
		time.Sleep(time.Millisecond * 100)
	}
	return M
}

func (inst *Instance) ApplyCapacity(identity int, caps uint64, timeout int, className ...string) {

	if inst == nil || inst.Mem == nil || inst.Mem.Mem == nil || inst.Address < 1000 {
		return
	}

	mem := inst.Mem.Mem

	time.Sleep(time.Second * 3)

	var (
		Holder  uint64
		State   uint64
		NewCaps uint64 = caps
	)
	Modules := inst.GetRunningScripts(timeout, className...)
	if len(Modules) > 0 {

		for _, Modules := range Modules {
			fmt.Println(mem.MemWrite(uintptr(Modules.LuaState+Offsets.ToIdentity), unsafe.Pointer(&identity), unsafe.Sizeof(identity)))
			fmt.Println(mem.MemWrite(uintptr(Modules.LuaState+Offsets.ToCapabilities), unsafe.Pointer(&NewCaps), unsafe.Sizeof(NewCaps)))
		}

		mem.MemRead(inst.Address+uintptr(Offsets.LuaState), unsafe.Pointer(&Holder), unsafe.Sizeof(Holder))
		mem.MemRead(uintptr(Holder+0x8), unsafe.Pointer(&State), unsafe.Sizeof(State))
		mem.MemWrite(uintptr(State+0x10), unsafe.Pointer(&NewCaps), unsafe.Sizeof(NewCaps))
		mem.MemWrite(uintptr(State+0x18), unsafe.Pointer(&NewCaps), unsafe.Sizeof(NewCaps))

		mem.MemRead(inst.Address+uintptr(Offsets.LuaState), unsafe.Pointer(&Holder), unsafe.Sizeof(Holder))
		mem.MemRead(uintptr(Holder), unsafe.Pointer(&State), unsafe.Sizeof(State))
		mem.MemWrite(uintptr(State+0x10), unsafe.Pointer(&NewCaps), unsafe.Sizeof(NewCaps))
		mem.MemWrite(uintptr(State+0x18), unsafe.Pointer(&NewCaps), unsafe.Sizeof(NewCaps))
	}
}

func (inst *Instance) WaitForChild(name string, timeout int) Instance {

	if inst == nil || inst.Mem == nil || inst.Mem.Mem == nil || inst.Address < 1000 {
		return NewInstance(0, inst.Mem)
	}

	for times := 0; times < timeout*10; times++ {
		child := inst.FindFirstChild(name, false)
		if child.Address != 0 {
			n := child.Name()
			if n == name {
				return child
			}
		}
		time.Sleep(time.Millisecond * time.Duration(100))
	}

	return Instance{Address: 0}
}

func (inst *Instance) WaitForClass(name string, timeout int) Instance {

	if inst == nil || inst.Mem == nil || inst.Mem.Mem == nil || inst.Address < 1000 {
		return NewInstance(0, inst.Mem)
	}

	for times := 0; times < timeout*10; times++ {
		child := inst.FindFirstChildOfClass(name, false)
		if child.Address != 0 {
			n := child.ClassName()
			if n == name {
				return child
			}
		}
		time.Sleep(time.Millisecond * time.Duration(100))
	}

	return Instance{Address: 0}
}

func (inst *Instance) FindFirstChild(name string, recursive bool) Instance {

	if inst == nil || inst.Mem == nil || inst.Mem.Mem == nil || inst.Address < 1000 {
		return NewInstance(0, inst.Mem)
	}

	mem := inst.Mem.Mem

	childrenPointer, _ := mem.ReadPointer(inst.Address + uintptr(inst.Mem.Offsets.Children))
	top, _ := mem.ReadPointer(childrenPointer)
	end, _ := mem.ReadPointer(childrenPointer + mem.PointerSize()*2)

	var ContinueIfFound bool

	for childAddr := top; childAddr < end; childAddr += mem.PointerSize() * 2 {

		if ContinueIfFound {
			ContinueIfFound = false
			continue
		}

		childPtr, _ := mem.ReadPointer(childAddr)
		if childPtr < 1000 {
			continue
		}

		child := NewInstance(childPtr, inst.Mem)

		n := child.Name()

		if n == name {
			return child
		}
		if n == "MarketplaceService" && name != "Players" {
			ContinueIfFound = true
			continue
		}
		if recursive {
			descendantChild := child.FindFirstChild(name, true)
			if descendantChild.Address > 1000 {
				return descendantChild
			}
		}
	}
	return NewInstance(0, inst.Mem)
}

func (inst *Instance) GetByteCode() ([]byte, int64) {

	if inst == nil || inst.Mem == nil || inst.Mem.Mem == nil || inst.Address < 1000 {
		return nil, 0
	}

	mem := inst.Mem.Mem

	var offset = inst.Mem.Offsets.Bytecode[inst.ClassName()]
	var size uintptr

	var btc_ptr uint64
	bytecode_pointer, _ := mem.ReadPointer(inst.Address + uintptr(offset))

	mem.MemRead(bytecode_pointer+0x10, unsafe.Pointer(&btc_ptr), unsafe.Sizeof(btc_ptr))
	mem.MemRead(bytecode_pointer+0x20, unsafe.Pointer(&size), unsafe.Sizeof(size))

	data, _ := mem.ReadBytes(uintptr(btc_ptr), size)
	return data, int64(size)
}

func (inst *Instance) SetBytecode(bytecode []byte, s uint64) {

	if inst == nil || inst.Mem == nil || inst.Mem.Mem == nil || inst.Address < 1000 {
		return
	}

	mem := inst.Mem.Mem

	size := uintptr(s)
	bytecode_pointer, _ := mem.ReadPointer(inst.Address + uintptr(inst.Mem.Offsets.Bytecode[inst.ClassName()]))
	new_ptr, _ := mem.AllocateMemory(size, bytecode_pointer)
	mem.WriteBytes(new_ptr, bytecode)
	mem.MemWrite(bytecode_pointer+0x10, unsafe.Pointer(&new_ptr), unsafe.Sizeof(new_ptr))
	mem.MemWrite(bytecode_pointer+0x20, unsafe.Pointer(&size), unsafe.Sizeof(size))
}

func (inst *Instance) FindFirstChildOfClass(className string, recursive bool) Instance {

	if inst == nil || inst.Mem == nil || inst.Mem.Mem == nil || inst.Address < 1000 {
		return NewInstance(0, inst.Mem)
	}

	mem := inst.Mem.Mem

	childrenPointer, err := mem.ReadPointer(inst.Address + uintptr(inst.Mem.Offsets.Children))
	if err != nil {
		return NewInstance(0, inst.Mem)
	}

	top, err := mem.ReadPointer(childrenPointer)
	if err != nil {
		return NewInstance(0, inst.Mem)
	}

	end, err := mem.ReadPointer(childrenPointer + mem.PointerSize())
	if err != nil {
		return NewInstance(0, inst.Mem)
	}

	var ContinueIfFound bool
	for childAddr := top; childAddr < end; childAddr += mem.PointerSize() * 2 {
		if ContinueIfFound {
			ContinueIfFound = false
			continue
		}

		childPtr, _ := mem.ReadPointer(childAddr)
		if childPtr < 1000 {
			continue
		}
		child := NewInstance(childPtr, inst.Mem)
		n := child.ClassName()
		if n == className {
			return child
		}
		if n == "MarketplaceService" && className != "Players" {
			ContinueIfFound = true
			continue
		}
		if recursive {
			descendantChild := child.FindFirstChildOfClass(className, true)
			if descendantChild.Address > 1000 {
				return descendantChild
			}
		}
	}
	return NewInstance(0, inst.Mem)
}

func (inst *Instance) SetModuleBypass() {

	if inst == nil || inst.Mem == nil || inst.Mem.Mem == nil || inst.Address < 1000 {
		return
	}

	mem := inst.Mem.Mem

	var set uint64 = 0x100000000
	var core uint64 = 0x1
	mem.MemWrite(inst.Address+uintptr(inst.Mem.Offsets.ModuleFlags), unsafe.Pointer(&set), unsafe.Sizeof(set))
	mem.MemWrite(inst.Address+uintptr(inst.Mem.Offsets.IsCore), unsafe.Pointer(&core), unsafe.Sizeof(core))
}

func (inst *Instance) Value() interface{} {

	if inst == nil || inst.Mem == nil || inst.Mem.Mem == nil || inst.Address < 1000 {
		return 0
	}

	mem := inst.Mem.Mem

	className := inst.ClassName()
	switch className {
	case "BoolValue":
		val, _ := mem.ReadByte(inst.Address + uintptr(inst.Mem.Offsets.ValueBase))
		return val == 1
	case "NumberValue":
		val, _ := mem.ReadDouble(inst.Address + uintptr(inst.Mem.Offsets.ValueBase))
		return val
	case "ObjectValue":
		addr, _ := mem.ReadPointer(inst.Address + uintptr(inst.Mem.Offsets.ValueBase))
		return NewInstance(addr, inst.Mem)
	case "StringValue", "":
		stringPointer := inst.Address + uintptr(inst.Mem.Offsets.ValueBase)
		stringLength, err := mem.ReadInt32(stringPointer + 0x10)
		if stringLength == 0 {
			return err
		}
		if stringLength > 15 {
			data, _ := mem.ReadInt32(stringPointer)
			stringPointer = uintptr(data)
		}
		str, _ := mem.ReadString(stringPointer, uintptr(stringLength))
		return str
	default:
		return nil
	}
}

func (inst *Instance) SetValue(value interface{}) {

	if inst == nil || inst.Mem == nil || inst.Mem.Mem == nil || inst.Address < 1000 {
		return
	}

	mem := inst.Mem.Mem
	className := inst.ClassName()

	switch className {
	case "BoolValue":
		val := byte(0)
		if value.(bool) {
			val = 1
		}
		mem.WriteByte(inst.Address+uintptr(inst.Mem.Offsets.ValueBase), val)
	case "NumberValue":
		mem.WriteDouble(inst.Address+uintptr(inst.Mem.Offsets.ValueBase), float64(value.(int)))
	case "ObjectValue":
		var addr uintptr
		if value != nil {
			addr = value.(*Instance).Address
		}
		mem.WritePointer(inst.Address+uintptr(inst.Mem.Offsets.ValueBase-0x8), addr)
	case "StringValue":
		stringAddr := inst.Address + uintptr(inst.Mem.Offsets.ValueBase)
		stringLength, _ := mem.ReadInt32(stringAddr + 0x10)
		var redirectedPtr uintptr
		if stringLength > 15 {
			val, _ := mem.ReadPointer(stringAddr)
			redirectedPtr = val
		} else {
			redirectedPtr = stringAddr
		}
		mem.WriteString(redirectedPtr, value.(string))
		mem.WriteInt32(stringAddr+0x10, int32(len(value.(string))+1))
	}
}

func (inst *Instance) String() string {
	return "(" + inst.Name() + " as " + inst.ClassName() + " | " + fmt.Sprintf("%#x", inst.Address) + ")"
}

func PatchRoblox(Memory *memory.Luna) {
	if Memory == nil {
		return
	}
	proc, err := process.NewProcess(int32(Memory.Pid))
	if err != nil {
		return
	}
	roblox_creation, err := proc.CreateTime()
	if err != nil {
		return
	}
	if time.Since(time.Unix(roblox_creation/1000, (roblox_creation%1000)*1000000)).Seconds() > 20 {
		return
	}
	for i := 0; i < 50; i++ {
		if regions, err := Memory.QueryMemoryRegions(); err == nil {
			for _, region := range regions {
				if region.Protect == windows.PAGE_READWRITE && region.Size == 0x200000 {
					var ptr = region.BaseAddress + 0x208
					if val, err := Memory.ReadDouble(ptr); err == nil && int64(val) == 0 {
						Memory.WriteDouble(ptr, 0x20)
						Memory.AllocAddr = region.BaseAddress
					}
					return
				}
			}
		} else {
			fmt.Printf("[error] Failed to query memory regions: %v\n", err)
			return
		}
		time.Sleep(100 * time.Millisecond)
	}
}
