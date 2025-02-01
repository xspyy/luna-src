package memory

import (
	"unsafe"

	"golang.org/x/sys/windows"
)

const (
	MEM_PRIVATE = 0x20000
)

type MemoryRegion struct {
	BaseAddress uintptr
	Size        uintptr
	Protect     uint32
}

func (m *Luna) QueryMemoryRegions() ([]MemoryRegion, error) {
	var regions []MemoryRegion
	var address uintptr

	if m == nil {
		return regions, nil
	}

	for {
		var mbi windows.MemoryBasicInformation
		err := windows.VirtualQueryEx(windows.Handle(m.ProcessHandle), address, &mbi, unsafe.Sizeof(mbi))
		if err != nil {
			break
		}

		if mbi.Type == MEM_PRIVATE && mbi.State == MEM_COMMIT {
			regions = append(regions, MemoryRegion{
				BaseAddress: mbi.BaseAddress,
				Size:        mbi.RegionSize,
				Protect:     mbi.Protect,
			})
		}

		address = mbi.BaseAddress + mbi.RegionSize
	}

	return regions, nil
}
