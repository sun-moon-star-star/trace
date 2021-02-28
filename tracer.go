package trace

import (
	"time"
)

type Tracer struct {
	TraceId   uint64
	TraceName string // other key message

	Stop      bool
	StartTime time.Time
	EndTime   time.Time

	Summary string
}

func NewTracer(option TracerIdOption) *Tracer {
	return &Tracer{
		TraceId:   GlobalTraceIdGenerator.GenerateTraceId(option),
		StartTime: time.Now(),
	}
}

func (t *Tracer) End() {
	t.EndTime = time.Now()
	t.Stop = true
}
