package classdescriptor

import (
	"main/packages/Memory/memory"
	"main/packages/Memory/propertydescriptor"
)

type ClassDescriptor struct {
	Address uintptr
}

func NewClassDescriptor(address uintptr) *ClassDescriptor {
	return &ClassDescriptor{Address: address}
}

func (cd *ClassDescriptor) Name(mem *memory.Luna) string {
	namePointer, _ := mem.ReadPointer(cd.Address + 0x8)
	name, _ := mem.ReadRbxStr(namePointer)
	return name
}

func (cd *ClassDescriptor) PropertyDescriptors(mem *memory.Luna) *propertydescriptor.PropertyDescriptorContainer {
	return propertydescriptor.NewPropertyDescriptorContainer(cd.Address)
}
