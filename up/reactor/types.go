package reactor

import (
	upcommon "github.com/rusriver/nutz/up/up-common"
)

const (
	Signal_End uint8 = iota
	Signal_IsUp
	Signal_IsFailed
	Signal_Cede
	Signal_CallTheInjectFunc
)

type Msg struct {
	Signal     uint8
	Name       string                    // process node Id / OTDG name
	OTDG       upcommon.IRecursiveWindUp // the OTDG itself - can't use real type because of stupid Go language cycle import limitation
	InjectFunc func(r *Reactor)
}

const (
	Status_Undef uint8 = iota
	Status_Fresh
	Status_Up
	Status_Failed
	Status_Ceded
	RequiredStatus_Up
	RequiredStatus_UpOrCeded
)

type Subsystems map[string]*Subsystem

type Subsystem struct {
	Name   string
	Itself upcommon.IRecursiveWindUp
	Status uint8
}

type SubsystemsPattern map[string]uint8 // the required dependencies pattern, map[<dep_name>]<required_status>
