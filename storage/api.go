package storage

import "trace"

type API interface {
	LoadTracer(*trace.Tracer) error
	SaveTracer(*trace.Tracer) error
	LoadSpanner(*trace.Spanner) error
	SaveSpanner(*trace.Spanner) error
}
