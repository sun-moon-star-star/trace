package trace

import (
	"time"
)

const (
	FlagStop = iota
)

type Tracer struct {
	TraceId   uint64
	TraceName string

	Flags     uint64
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
	t.Flags |= FlagStop
}

func (t *Tracer) NewSpanner() *Spanner {
	s := NewSpanner()
	s.TraceId = t.TraceId
	return s
}
