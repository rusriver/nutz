package reactor

const (
	SubsystemStatus_Fresh uint8 = iota
	SubsystemStatus_Ready
	SubsystemStatus_Ceded
)

type subsystems map[string]uint8
