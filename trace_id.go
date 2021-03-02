package trace

import (
	"sync"
	"time"
)

type TracerIdOption struct {
	// 0-bit no use for extends
	// 41-bit millisecond timestamp
	// 10-bit projectId
	// 12-bit sequenceId
	ProjectId uint16
}

type ClockCallbackStrategy func(millis uint64, seq uint32) (uint64, uint32)

var DefaultClockCallbackStrategy ClockCallbackStrategy = func(millis uint64, seq uint32) (uint64, uint32) {
	if seq == 4095 {
		return millis + 1, 0
	}
	return millis, seq + 1
}

type TraceIdGenerator struct {
	millisecondTimestamp uint64 // 41 bits
	sequenceId           uint32 // 12 bits

	ClockCallbackStrategy ClockCallbackStrategy

	lock sync.Mutex
}

var GlobalTraceIdGenerator TraceIdGenerator = TraceIdGenerator{
	ClockCallbackStrategy: DefaultClockCallbackStrategy,
}

func (g *TraceIdGenerator) GenerateTraceId(option TracerIdOption) uint64 {
	millis := uint64(time.Now().UnixNano() / 1e6)
	seq := uint32(0)

	g.lock.Lock()

	if millis > g.millisecondTimestamp {
		g.millisecondTimestamp, g.sequenceId = millis, 0
	} else if millis == g.millisecondTimestamp {
		g.sequenceId++
		seq = g.sequenceId
	} else {
		// 时钟回拨
		millis, seq = g.ClockCallbackStrategy(
			g.millisecondTimestamp, g.sequenceId)
		g.millisecondTimestamp, g.sequenceId = millis, seq
	}

	g.lock.Unlock()

	traceId := millis << 22                   // 41-bit millisecond timestamp
	traceId += uint64(option.ProjectId) << 12 // 10-bit projectId
	traceId += uint64(seq)                    // 12-bit sequenceId

	return traceId
}