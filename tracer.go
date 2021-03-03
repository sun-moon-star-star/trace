package trace

import "time"

const (
	FlagStop = iota
)

type Tracer struct {
	TraceId   uint64
	TraceName string

	Flags        uint64
	EndTimestamp uint64

	Summary string

	First *Spanner
}

func NewTracer() *Tracer {
	return &Tracer{
		TraceId: NewUUID(),
	}
}

func (t *Tracer) End() {
	t.EndTimestamp = uint64(time.Now().UnixNano() / 1e6)
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
