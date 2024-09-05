package promethabs

type IMetricsService interface {
	/* Inventory of metrics and their types must be specified in the constructor
	   of the service, id->type, e.g. histogram, etc; If you use unknown metric ID,
	   it's recorded as a Gauge */

	Metric(id uint16, labels map[string]string) IMetric
}

type IMetric interface {
	Set(float64)
	Inc()
	Dec()
	Add(float64)
	Sub(float64)
}
