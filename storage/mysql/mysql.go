package mysql

import "trace"

// type API interface {
// 	SaveTracer(*trace.Tracer) error
// 	LoadTracer(traceId string) (*trace.Tracer, error)
// 	SaveSpanner(*trace.Spanner) error
// 	LoadSpanner(spanId uint64) (*trace.Spanner, error)
// }

type Mysql struct{}

func (mysql *Mysql) SaveTracer(*trace.Tracer) error {

}
