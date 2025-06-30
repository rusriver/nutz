package upcommon

import loggeterrapi "github.com/rusriver/nutz/loggeterr-api"

type IRecursiveWindUp interface {
	// Use specific call-side tracer here
	RecursiveWindUp(tracer loggeterrapi.ITracer)
}

type ISystem interface {
	Start()
	Stop()
}

type IRemovable interface {
	WindUp(tracer loggeterrapi.ITracer)
}
