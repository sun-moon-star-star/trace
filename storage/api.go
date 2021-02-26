package storage

import "trace"

type API interface {
	SaveSpanner(*trace.Spanner) error
	LoadSpanner(spanId uint64) (*trace.Spanner, error)
	LoadTracer(traceId string) (*trace.Tracer, error)
	SaveTracer(*trace.Tracer) error
}
