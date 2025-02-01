package memory

import (
	"errors"
	"fmt"
	"strings"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	procCreateToolhelp32 = kernel32.NewProc("CreateToolhelp32Snapshot")
	procModule32FirstW   = kernel32.NewProc("Module32FirstW")
	procModule32NextW    = kernel32.NewProc("Module32NextW")
)

const (
	TH32CS_SNAPMODULE   = 0x00000008
	TH32CS_SNAPMODULE32 = 0x00000010
)

type MODULEENTRY32 struct {
	DwSize        uint32
	Th32ModuleID  uint32
	Th32ProcessID uint32
	GlblcntUsage  uint32
	ProccntUsage  uint32
	ModBaseAddr   *byte
	ModBaseSize   uint32
	HModule       windows.Handle
	SzModule      [256]uint16
	SzExePath     [260]uint16
}

func createToolhelp32Snapshot(flags, pid uint32) (windows.Handle, error) {
	r1, _, e1 := syscall.Syscall(procCreateToolhelp32.Addr(), 2,
		uintptr(flags),
		uintptr(pid),
		0)
	handle := windows.Handle(r1)
	if handle == windows.InvalidHandle {
		return 0, error(e1)
	}
	return handle, nil
}

func module32First(snapshot windows.Handle, me *MODULEENTRY32) error {
	r1, _, e1 := syscall.Syscall(procModule32FirstW.Addr(), 2,
		uintptr(snapshot),
		uintptr(unsafe.Pointer(me)),
		0)
	if r1 == 0 {
		return error(e1)
	}
	return nil
}

func module32Next(snapshot windows.Handle, me *MODULEENTRY32) error {
	r1, _, e1 := syscall.Syscall(procModule32NextW.Addr(), 2,
		uintptr(snapshot),
		uintptr(unsafe.Pointer(me)),
		0)
	if r1 == 0 {
		return error(e1)
	}
	return nil
}

func utf16PtrToString(u16 *uint16) string {
	if u16 == nil {
		return ""
	}
	var length int
	for p := u16; *p != 0; p = (*uint16)(unsafe.Pointer(uintptr(unsafe.Pointer(p)) + unsafe.Sizeof(*p))) {
		length++
	}
	s := make([]uint16, length)
	for i := 0; i < length; i++ {
		s[i] = *(*uint16)(unsafe.Pointer(uintptr(unsafe.Pointer(u16)) + uintptr(i)*2))
	}
	return syscall.UTF16ToString(s)
}

func (m *Luna) GetBaseAddr(pid uint32, target ...string) (uintptr, error) {
	snapshot, err := createToolhelp32Snapshot(TH32CS_SNAPMODULE|TH32CS_SNAPMODULE32, pid)
	if err != nil {
		return 0, fmt.Errorf("CreateToolhelp32Snapshot failed: %v", err)
	}
	defer windows.CloseHandle(snapshot)

	var me MODULEENTRY32
	me.DwSize = uint32(unsafe.Sizeof(me))

	if err := module32First(snapshot, &me); err != nil {
		return 0, fmt.Errorf("module32First failed: %v", err)
	}

	for {
		modName := utf16PtrToString(&me.SzModule[0])
		for _, name := range target {
			if strings.Contains(modName, name) {
				base := uintptr(unsafe.Pointer(me.ModBaseAddr))
				return base, nil
			}
		}
		if err := module32Next(snapshot, &me); err != nil {
			if err == syscall.ERROR_NO_MORE_FILES {
				break
			}
			return 0, fmt.Errorf("module32Next failed: %v", err)
		}
	}
	return 0, errors.New("RobloxPlayerBeta.exe not found in module list")
}
