package memory

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"strings"
	"syscall"
	"unsafe"

	"github.com/lxn/win"
	"golang.org/x/sys/windows"
)

var (
	kernel32                  = syscall.NewLazyDLL("kernel32.dll")
	ntdll                     = syscall.NewLazyDLL("ntdll.dll")
	modpsapi                  = windows.NewLazySystemDLL("psapi.dll")
	user32                    = syscall.NewLazyDLL("user32.dll")
	procOpenProcess           = kernel32.NewProc("OpenProcess")
	procVirtualQueryEx        = kernel32.NewProc("VirtualQueryEx")
	procVirtualAllocEx        = kernel32.NewProc("VirtualAllocEx")
	procNtUnlockVirtualMemory = ntdll.NewProc("NtUnlockVirtualMemory")
	procSendMessage           = user32.NewProc("SendMessageA")
	ntReadVirtualMemory       = ntdll.NewProc("NtReadVirtualMemory")
	procWriteProcessMemory    = kernel32.NewProc("WriteProcessMemory")
	procQueryWorkingSetEx     = modpsapi.NewProc("QueryWorkingSetEx")
)

type Processes struct {
	Name string
	Pid  uint32
}

const (
	PROCESS_ALL_ACCESS = 0x1F0FFF

	PAGE_READONLY          = 0x02
	PAGE_READWRITE         = 0x04
	PAGE_EXECUTE_READ      = 0x20
	PAGE_EXECUTE_READWRITE = 0x40

	MEM_COMMIT  = 0x1000
	MEM_RESERVE = 0x2000
	MEM_RELEASE = 0x8000

	MEM_DECOMMIT = 0x4000
)

var ALLOWED_PROTECTIONS = []uint32{
	PAGE_READONLY,
	PAGE_READWRITE,
	PAGE_EXECUTE_READ,
	PAGE_EXECUTE_READWRITE,
}

type MEMORY_BASIC_INFORMATION struct {
	BaseAddress       uintptr
	AllocationBase    uintptr
	AllocationProtect uint32
	RegionSize        uintptr
	State             uint32
	Protect           uint32
	Type              uint32
}

var pg = syscall.Getpagesize()

type Luna struct {
	ProcessHandle syscall.Handle
	Is64Bit       bool
	RobloxBase    uintptr
	AllocAddr     uintptr
	Pid           uint32
}

func SendMessage(hwnd win.HWND, msg uint32, wParam, lParam uintptr) error {
	ret, _, err := procSendMessage.Call(
		uintptr(hwnd),
		uintptr(msg),
		wParam,
		lParam,
	)
	if ret == 0 {
		return err
	}
	return nil
}

func IsHandleValid(h syscall.Handle) bool {
	var exitCode uint32
	err := syscall.GetExitCodeProcess(h, &exitCode)
	if err != nil || exitCode != 259 {
		return false
	}
	return true
}

func NewLuna(pid uint32) (*Luna, error) {
	handle, _, err := procOpenProcess.Call(
		0x1F0FFF,
		0,
		uintptr(pid),
	)
	if handle == 0 {
		return nil, err
	}

	new := &Luna{
		ProcessHandle: syscall.Handle(handle),
		Is64Bit:       true,
		Pid:           pid,
	}

	base, err := new.GetBaseAddr(pid, "RobloxPlayerBeta", "Windows10Universal", "eurotrucks2")
	if err != nil {
		return nil, nil
	}
	new.RobloxBase = base

	return new, nil
}

func (m *Luna) IsWorkingSet(address uintptr) bool {
	var wsInfo = struct {
		VirtualAddress    uintptr
		VirtualAttributes uintptr
	}{VirtualAddress: address & ^(uintptr(pg) - 1), VirtualAttributes: 0}
	procQueryWorkingSetEx.Call(
		uintptr(m.ProcessHandle),
		uintptr(unsafe.Pointer(&wsInfo)),
		uintptr(unsafe.Sizeof(wsInfo)),
	)
	return (wsInfo.VirtualAttributes & 0x1) != 0
}

func (m *Luna) WaitUntilPossiblyReadable(address uintptr) {
	var i = 0
	for !m.IsWorkingSet(address) {
		i++
		if i > 300 {
			break
		}
	}
}

func (m *Luna) MemRead(address uintptr, buffer unsafe.Pointer, size uintptr) error {

	if m == nil || !IsHandleValid(m.ProcessHandle) {
		return errors.New("Invalid memory address")
	}

	m.WaitUntilPossiblyReadable(address)

	var mbi MEMORY_BASIC_INFORMATION
	mbiSize := uintptr(unsafe.Sizeof(mbi))

	procVirtualQueryEx.Call(
		uintptr(m.ProcessHandle),
		address,
		uintptr(unsafe.Pointer(&mbi)),
		mbiSize,
	)

	ntReadVirtualMemory.Call(
		uintptr(m.ProcessHandle),
		address,
		uintptr(buffer),
		size,
		0,
	)

	baddr := mbi.AllocationBase
	s := mbi.RegionSize

	procNtUnlockVirtualMemory.Call(
		uintptr(m.ProcessHandle),
		uintptr(unsafe.Pointer(&baddr)),
		uintptr(unsafe.Pointer(&s)),
		1,
	)

	return nil
}

func (m *Luna) MemWrite(address uintptr, buffer unsafe.Pointer, size uintptr) error {

	if m == nil || !IsHandleValid(m.ProcessHandle) {
		return errors.New("Invalid memory address")
	}

	var bytesWritten uintptr
	status, _, err := procWriteProcessMemory.Call(
		uintptr(m.ProcessHandle),
		address,
		uintptr(buffer),
		size,
		uintptr(unsafe.Pointer(&bytesWritten)),
	)
	if status == 0 {
		return fmt.Errorf("WriteProcessMemory failed at address %#x: %v", address, err)
	}
	/*
		m.WaitUntilPossiblyReadable(address)

		var mbi MemoryBasicInformation
		mbiSize := uintptr(unsafe.Sizeof(mbi))

		procVirtualQueryEx.Call(
			uintptr(m.ProcessHandle),
			address,
			uintptr(unsafe.Pointer(&mbi)),
			mbiSize,
		)

		var bytesWritten uintptr
		status, _, err := procWriteProcessMemory.Call(
			uintptr(m.ProcessHandle),
			address,
			uintptr(buffer),
			size,
			uintptr(unsafe.Pointer(&bytesWritten)),
		)

		baddr := mbi.AllocationBase
		s := mbi.RegionSize

		m.UnlockMemory(baddr, s)

		if status == 0 {
			return fmt.Errorf("WriteProcessMemory failed at address %#x: %v", address, err)
		}

		if bytesWritten < size {
			return fmt.Errorf("only wrote %d bytes out of %d requested to address %#x", bytesWritten, size, address)
		}
	*/
	return nil
}

func (m *Luna) PointerSize() uintptr {

	if m == nil {
		return 8
	}

	if m.Is64Bit {
		return 8
	}
	return 4
}

func (m *Luna) AllocateMemory(size uintptr, address uintptr) (uintptr, error) {
	if m == nil || !IsHandleValid(m.ProcessHandle) {
		return 0, errors.New("Invalid memory address")
	}
	addr, _, err := procVirtualAllocEx.Call(
		uintptr(m.ProcessHandle),
		0,
		size,
		windows.MEM_COMMIT,
		windows.PAGE_EXECUTE_READWRITE,
	)
	if addr == 0 {
		return 0, err
	}
	return addr, nil
}

func (m *Luna) ReadMemory(address uintptr, size uintptr) ([]byte, error) {

	if m == nil || !IsHandleValid(m.ProcessHandle) {
		return []byte{}, errors.New("Invalid memory address")
	}

	buffer := make([]byte, size)
	err := m.MemRead(address, unsafe.Pointer(&buffer[0]), size)
	if err != nil {
		return nil, err
	}
	return buffer, nil
}

func (m *Luna) WriteMemory(address uintptr, data []byte) error {
	if len(data) == 0 {
		return errors.New("Byte size is empty!")
	}
	return m.MemWrite(address, unsafe.Pointer(&data[0]), uintptr(len(data)))
}

func (m *Luna) ReadByte(address uintptr) (byte, error) {

	var result byte
	err := m.MemRead(address, unsafe.Pointer(&result), 1)
	if err != nil {
		return 0, err
	}
	return result, nil
}

func (m *Luna) ReadBytes(address uintptr, length uintptr) ([]byte, error) {

	return m.ReadMemory(address, length)
}

func (m *Luna) ReadString(address uintptr, maxLength uintptr) (string, error) {

	if maxLength > 1000 || maxLength == 0 {
		maxLength = 100
	}

	buffer := make([]byte, maxLength)
	err := m.MemRead(address, unsafe.Pointer(&buffer[0]), maxLength)
	if err != nil {
		return "", err
	}
	idx := bytes.IndexByte(buffer, 0)
	if idx != -1 {
		buffer = buffer[:idx]
	}
	return string(buffer), nil
}

func (m *Luna) ReadDouble(address uintptr) (float64, error) {
	var result float64
	err := m.MemRead(address, unsafe.Pointer(&result), 8)
	if err != nil {
		return 0, err
	}
	return result, nil
}

func (m *Luna) ReadFloat(address uintptr) (float32, error) {
	var result float32
	err := m.MemRead(address, unsafe.Pointer(&result), 4)
	if err != nil {
		return 0, err
	}
	return result, nil
}

func (m *Luna) ReadInt32(address uintptr) (int32, error) {
	var result int32
	err := m.MemRead(address, unsafe.Pointer(&result), 4)
	if err != nil {
		return 0, err
	}
	return result, nil
}

func (m *Luna) ReadInt64(address uintptr) (int64, error) {
	var result int64
	err := m.MemRead(address, unsafe.Pointer(&result), 8)
	if err != nil {
		return 0, err
	}
	return result, nil
}

func (mem *Luna) ReadRbxStr(address uintptr) (string, error) {
	var strCheck uint64
	err := mem.MemRead(address+0x10, unsafe.Pointer(&strCheck), unsafe.Sizeof(strCheck))
	if err != nil {
		return "", err
	}

	if strCheck > 15 {
		var strPointer uint64
		err = mem.MemRead(address, unsafe.Pointer(&strPointer), unsafe.Sizeof(strPointer))
		if err != nil {
			return "", err
		}
		return mem.ReadString(uintptr(strPointer), uintptr(strCheck))
	}

	return mem.ReadString(address, uintptr(strCheck))
}

func (m *Luna) ReadPointer(address uintptr) (uintptr, error) {
	if m.Is64Bit {
		val, err := m.ReadUint64(address)
		return uintptr(val), err
	} else {
		val, err := m.ReadUint32(address)
		return uintptr(val), err
	}
}

func (m *Luna) ReadUint32(address uintptr) (uint32, error) {
	var result uint32
	err := m.MemRead(address, unsafe.Pointer(&result), unsafe.Sizeof(result))
	if err != nil {
		return 0, err
	}
	return result, nil
}

func (m *Luna) ReadUint64(address uintptr) (uint64, error) {
	var result uint64
	err := m.MemRead(address, unsafe.Pointer(&result), unsafe.Sizeof(result))
	if err != nil {
		return 0, err
	}
	return result, nil
}

func (m *Luna) WriteByte(address uintptr, value byte) error {
	return m.MemWrite(address, unsafe.Pointer(&value), 1)
}

func (m *Luna) WriteBytes(address uintptr, data []byte) error {
	return m.WriteMemory(address, data)
}

func (m *Luna) WriteString(address uintptr, value string) error {
	data := append([]byte(value), 0)
	return m.WriteBytes(address, data)
}

func (m *Luna) WriteDouble(address uintptr, value float64) error {
	return m.MemWrite(address, unsafe.Pointer(&value), 8)
}

func (m *Luna) WriteFloat(address uintptr, value float32) error {
	return m.MemWrite(address, unsafe.Pointer(&value), 4)
}

func (m *Luna) WriteInt32(address uintptr, value int32) error {
	return m.MemWrite(address, unsafe.Pointer(&value), 4)
}

func (m *Luna) WriteInt64(address uintptr, value int64) error {
	return m.MemWrite(address, unsafe.Pointer(&value), 8)
}

func (m *Luna) WritePointer(address uintptr, value uintptr) error {
	if m.Is64Bit {
		val := uint64(value)
		return m.MemWrite(address, unsafe.Pointer(&val), 8)
	} else {
		val := uint32(value)
		return m.MemWrite(address, unsafe.Pointer(&val), 4)
	}
}

func GetProcesses() ([]Processes, error) {
	var processNames []Processes
	snapshot, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPPROCESS, 0)
	if err != nil {
		return nil, err
	}
	defer windows.CloseHandle(snapshot)

	var entry windows.ProcessEntry32
	entry.Size = uint32(unsafe.Sizeof(entry))

	if err := windows.Process32First(snapshot, &entry); err != nil {
		return nil, err
	}

	for {
		name := windows.UTF16ToString(entry.ExeFile[:])
		if strings.Contains(name, "RobloxPlayerBeta") ||
			strings.Contains(name, "Windows10Universal") ||
			strings.Contains(name, "eurotrucks2") {
			processNames = append(processNames, Processes{
				Name: name,
				Pid:  entry.ProcessID,
			})
		}
		if err := windows.Process32Next(snapshot, &entry); err != nil {
			if errors.Is(err, windows.ERROR_NO_MORE_FILES) {
				break
			}
			return nil, err
		}
	}

	return processNames, nil
}

func removeEuro(data []Processes) []Processes {
	var RemovedRobloxs []Processes
	for _, inst := range data {
		if inst.Name != "RobloxPlayerBeta.exe" && inst.Name != "Windows10Universal.exe" {
			inst.Name = "RobloxPlayerBeta.exe"
			RemovedRobloxs = append(RemovedRobloxs, inst)
		} else {
			RemovedRobloxs = append(RemovedRobloxs, inst)
		}
	}
	return RemovedRobloxs
}

func IsProcessRunning() (bool, []Processes) {
	processes, err := GetProcesses()
	if err != nil {
		log.Fatal(err)
	}
	processes = removeEuro(processes)
	return len(processes) != 0, processes
}
