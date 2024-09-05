package promethabs

var IdCounter uint16

type MetricId struct {
	Id          uint16
	Name        string
	Description string
}

// Non-thread-safe
func NewMetricId(name string) (id *MetricId) {
	IdCounter++
	id = &MetricId{
		Id:   IdCounter,
		Name: name,
	}
	return
}
