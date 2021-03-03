package storage

import (
	"trace"
	"trace/storage/mysql"
)

type API interface {
	LoadTracer(*trace.Tracer) error
	SaveTracer(*trace.Tracer) error
	LoadSpanner(*trace.Spanner) error
	SaveSpanner(*trace.Spanner) error
}

func DefaultAPI() API {
	return &mysql.Mysql{}
}
