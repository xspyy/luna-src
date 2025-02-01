package propertydescriptor

import (
	"main/packages/Memory/memory"
	"unsafe"
)

var OldAccessibleFlags = make(map[uintptr]int32)

type PropertyDescriptor struct {
	Address uintptr
}

func NewPropertyDescriptor(address uintptr) *PropertyDescriptor {
	return &PropertyDescriptor{Address: address}
}

func (pd *PropertyDescriptor) Name(mem *memory.Luna) string {
	namePointer, _ := mem.ReadPointer(pd.Address + 0x8)
	mem.ReadRbxStr(namePointer)
	name, _ := mem.ReadRbxStr(namePointer)
	return name
}

func (pd *PropertyDescriptor) Capabilities(mem *memory.Luna) int32 {
	val, _ := mem.ReadInt32(pd.Address + 0x38)
	return val
}

func (pd *PropertyDescriptor) AccessibleFlags(mem *memory.Luna) int32 {
	val, _ := mem.ReadInt32(pd.Address + 0x40)
	return val
}

func (pd *PropertyDescriptor) IsHiddenValue(mem *memory.Luna) bool {
	val := pd.AccessibleFlags(mem)
	return val < 32
}

func (pd *PropertyDescriptor) SetScriptable(mem *memory.Luna, scriptable bool) {
	if scriptable {
		if _, exists := OldAccessibleFlags[pd.Address]; !exists {
			OldAccessibleFlags[pd.Address] = pd.AccessibleFlags(mem)
			mem.WriteInt32(pd.Address+0x40, 63)
		}
	} else {
		if oldFlag, exists := OldAccessibleFlags[pd.Address]; exists {
			mem.WriteInt32(pd.Address+0x40, oldFlag)
		}
	}
}

type PropertyDescriptorContainer struct {
	Address uintptr
}

func NewPropertyDescriptorContainer(address uintptr) *PropertyDescriptorContainer {
	return &PropertyDescriptorContainer{Address: address}
}

func (pdc *PropertyDescriptorContainer) GetAllYield(mem *memory.Luna) []*PropertyDescriptor {
	var descriptors []*PropertyDescriptor
	var (
		start uint64
		end   uint64
	)

	mem.MemRead(pdc.Address+0x28, unsafe.Pointer(&start), unsafe.Sizeof(start))
	mem.MemRead(pdc.Address+0x30, unsafe.Pointer(&end), unsafe.Sizeof(end))

	for addr := start; addr < end; addr += 0x8 {
		descriptorAddr, _ := mem.ReadPointer(uintptr(addr))
		if descriptorAddr > 1000 {
			new := NewPropertyDescriptor(descriptorAddr)

			descriptors = append(descriptors, new)
		}
	}

	return descriptors
}

func (pdc *PropertyDescriptorContainer) Get(mem *memory.Luna, name string) *PropertyDescriptor {
	for _, descriptor := range pdc.GetAllYield(mem) {
		if descriptor.Name(mem) == name {
			return descriptor
		}
	}
	return NewPropertyDescriptor(0)
}
