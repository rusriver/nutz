package readinesssensor

type Probe struct {
	Sensor *Sensor
	Name   string
}

func (p *Probe) Ready() {
	p.Sensor.Report(p.Name, SubsystemStatus_Ready)
}

func (p *Probe) Cede() {
	p.Sensor.Report(p.Name, SubsystemStatus_Ceded)
}
