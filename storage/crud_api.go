package storage

import (
	"trace"
	"trace/storage/mysql"
)

type CRUDAPI interface {
	LoadTracer(*trace.Tracer) error
	SaveTracer(*trace.Tracer) error
	LoadSpanner(*trace.Spanner) error
	SaveSpanner(*trace.Spanner) error
}

func DefaultCRUDAPI() CRUDAPI {
	return &mysql.Mysql{}
}
