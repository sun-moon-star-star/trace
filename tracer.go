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

	First *Spanner
}

func NewTracer() *Tracer {
	return &Tracer{
		TraceId:   NewUUID(),
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

	if t.First == nil {
		t.First = s
	}

	return s
}
