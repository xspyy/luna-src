package humanoid

import (
	"main/packages/Memory/instance"
)

type Humanoid struct {
	Address uintptr
	Mem     *instance.Instance
}

type HumanoidOffsets struct {
	HealthDisplayDistance float64
	NameDisplayDistance   float64
	Health                float64
	MaxHealth             float64
	WalkSpeed             []float64
}

var OffsetsHumanoid = HumanoidOffsets{
	HealthDisplayDistance: 400,
	NameDisplayDistance:   436,
	Health:                396,
	MaxHealth:             428,
	WalkSpeed: []float64{
		456, 928,
	},
}

func Create(Rbx *instance.Instance) *Humanoid {

	return &Humanoid{Mem: Rbx, Address: Rbx.Address}
}

func (inst *Humanoid) GetHealth() float32 {
	if inst == nil || inst.Mem == nil || inst.Mem.Mem == nil || inst.Address < 1000 {
		return 0
	}
	mem := inst.Mem.Mem.Mem
	HP, _ := mem.ReadFloat(inst.Address + uintptr(OffsetsHumanoid.Health))
	return HP
}

func (inst *Humanoid) GetMaxHealth() float32 {
	if inst == nil || inst.Mem == nil || inst.Mem.Mem == nil || inst.Address < 1000 {
		return 0
	}
	mem := inst.Mem.Mem.Mem
	HP, _ := mem.ReadFloat(inst.Address + uintptr(OffsetsHumanoid.MaxHealth))
	return HP
}

func (inst *Humanoid) SetHealth(Health float32) {
	if inst == nil || inst.Mem == nil || inst.Mem.Mem == nil || inst.Address < 1000 {
		return
	}
	mem := inst.Mem.Mem.Mem
	mem.WriteFloat(inst.Address+uintptr(OffsetsHumanoid.Health), Health)
}

func (inst *Humanoid) SetMaxHealth(Health float32) {
	if inst == nil || inst.Mem == nil || inst.Mem.Mem == nil || inst.Address < 1000 {
		return
	}
	mem := inst.Mem.Mem.Mem
	mem.WriteFloat(inst.Address+uintptr(OffsetsHumanoid.MaxHealth), Health)
}

func (inst *Humanoid) WalkSpeed(Speed float32) {
	if inst == nil || inst.Mem == nil || inst.Mem.Mem == nil || inst.Address < 1000 {
		return
	}
	mem := inst.Mem.Mem.Mem
	mem.WriteFloat(inst.Address+uintptr(OffsetsHumanoid.WalkSpeed[0]), Speed)
	mem.WriteFloat(inst.Address+uintptr(OffsetsHumanoid.WalkSpeed[1]), Speed)
}
