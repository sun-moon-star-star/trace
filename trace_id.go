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

type LastValueT struct {
	MillisecondTimestamp uint64 // 41 bits
	SequenceId           uint32 // 12 bits

	ClockCallbackStrategy ClockCallbackStrategy
	Lock                  sync.Mutex
}

var LastValue LastValueT = LastValueT{
	ClockCallbackStrategy: DefaultClockCallbackStrategy,
}

func GenerateTraceId(option TracerIdOption) uint64 {
	millis := uint64(time.Now().UnixNano() / 1e6)
	seq := uint32(0)

	LastValue.Lock.Lock()

	if millis > LastValue.MillisecondTimestamp {
		LastValue.MillisecondTimestamp, LastValue.SequenceId = millis, 0
	} else if millis == LastValue.MillisecondTimestamp {
		LastValue.SequenceId++
		seq = LastValue.SequenceId
	} else {
		// 时钟回拨
		millis, seq = LastValue.ClockCallbackStrategy(
			LastValue.MillisecondTimestamp, LastValue.SequenceId)
		LastValue.MillisecondTimestamp, LastValue.SequenceId = millis, seq
	}

	LastValue.Lock.Unlock()

	traceId := millis << 22                   // 41-bit millisecond timestamp
	traceId += uint64(option.ProjectId) << 12 // 10-bit projectId
	traceId += uint64(seq)                    // 12-bit sequenceId

	return traceId
}
