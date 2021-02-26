package trace

import (
	"time"
	"trace/random"
)

type TracerOption struct {
	// 0-bit no use for extends
	// 41-bit millisecond timestamp
	// 10-bit projectId
	// 12-bit sequenceId
	ProjectId uint16
}

type Tracer struct {
	TraceId   uint64
	TraceName string // other key message

	Stop      bool
	StartTime time.Time
	EndTime   time.Time

	Summary string
}

func NewTracer(option TracerOption) *Tracer {
	var traceId uint64
	traceId = uint64(time.Now().UnixNano()/1e6) << 22
	traceId += uint64(option.ProjectId)<<12 + uint64(random.RandomUint12())

	return &Tracer{
		TraceId:   traceId,
		StartTime: time.Now(),
	}
}

func (t *Tracer) End() {
	t.EndTime = time.Now()
	t.Stop = true
}
