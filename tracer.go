package trace

import "time"

const (
	TracerStop = 1 << iota
)

type TracerStrategy struct {
	Summary func(*Tracer) string
}

var DefaultTracerStrategy *TracerStrategy

func init() {
	DefaultTracerStrategy = &TracerStrategy{
		Summary: func(t *Tracer) string {
			if t.First == nil {
				return ""
			}
			s := t.First
			for s.Children != nil {
				s = s.Children
			}
			return s.BaggageString()
		},
	}
}

type Tracer struct {
	TraceId   uint64
	TraceName string

	Flags        uint64
	EndTimestamp uint64

	Summary string

	TracerStrategy *TracerStrategy

	First *Spanner
}

func NewTracer(TraceName string) *Tracer {
	return &Tracer{
		TraceId:   NewUUID(),
		TraceName: TraceName,

		TracerStrategy: DefaultTracerStrategy,
	}
}

func (t *Tracer) End() {
	if t.Flags&TracerStop != TracerStop {
		t.EndTimestamp = uint64(time.Now().UnixNano() / 1e6)
		t.Flags |= TracerStop
	}
}

func (t *Tracer) NewSpanner(SpanName string, ParentSpanId uint64) *Spanner {
	s := NewSpanner(SpanName, t.TraceId, ParentSpanId)

	if t.First == nil {
		t.First = s
	}

	return s
}
