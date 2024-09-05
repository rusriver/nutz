### PLEASE LOOK AT 20240826 AREC, ABOUT METRICS WRAPPER OVER THE PROMETHEUS LIB

Expected usage:

```go

    // create a metrics id (a metric, from the POV of user)
    myMetric := promethabs.NewMetricId("my_metric")

    // init a metrics service
    metricsService := ...

    // use the metric, e.g. increment it
    metricsService.Metric(myMetric.Id).Inc()

```
