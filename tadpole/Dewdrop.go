package main

import (
	"github.com/xiaonanln/goworld/engine/entity"
)

// Dewdrop type
type Dewdrop struct {
	entity.Entity
}

func (dewdrop *Dewdrop) DescribeEntityType(desc *entity.EntityTypeDesc) {
	desc.SetUseAOI(true, 10000)
	desc.DefineAttr("amount", "AllClients")
}

func (dewdrop *Dewdrop) OnEnterSpace() {
	dewdrop.setDefaultAttrs()
}

func (dewdrop *Dewdrop) setDefaultAttrs() {
	dewdrop.Attrs.SetDefaultInt("amount", 1)
}

func (dewdrop *Dewdrop) getAmount() int64 {
	return dewdrop.Attrs.GetInt("amount")
}

func (dewdrop *Dewdrop) Collect() {
	dewdrop.Destroy()
}
