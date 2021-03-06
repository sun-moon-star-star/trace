package uuid

import (
	"sync"
	"time"
)

// 0-bit no use for extends
// 41-bit millisecond timestamp
// 10-bit projectId
// 12-bit sequenceId

type ClockCallbackStrategy func(millis uint64, seq uint32) (uint64, uint32)

const (
	ProjectIdBits  = 10
	SequenceIdBits = 12
	MaxProjectId   = (1 << ProjectIdBits) - 1
	MaxSequenceId  = (1 << SequenceIdBits) - 1
)

var DefaultClockCallbackStrategy ClockCallbackStrategy = func(millis uint64, seq uint32) (uint64, uint32) {
	if seq == MaxSequenceId {
		return millis + 1, 0
	}
	return millis, seq + 1
}

type UUIdGenerator struct {
	ProjectId uint64

	millisecondTimestamp uint64 // 41 bits
	sequenceId           uint32 // 12 bits

	ClockCallbackStrategy ClockCallbackStrategy

	lock sync.Mutex
}

var GlobalUUIDGenerator UUIdGenerator = UUIdGenerator{
	ClockCallbackStrategy: DefaultClockCallbackStrategy,
}

func (g *UUIdGenerator) NewUUID() uint64 {
	millis := uint64(time.Now().UnixNano() / 1e6)
	seq := uint32(0)

	g.lock.Lock()

	if millis > g.millisecondTimestamp {
		g.millisecondTimestamp, g.sequenceId = millis, seq
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

	traceId := millis << (SequenceIdBits + ProjectIdBits) // 41-bit millisecond timestamp
	traceId += uint64(g.ProjectId) << SequenceIdBits      // 10-bit projectId
	traceId += uint64(seq)                                // 12-bit sequenceId

	return traceId
}

func NewUUID() uint64 {
	return GlobalUUIDGenerator.NewUUID()
}

func TimestampFromUUID(uuid uint64) uint64 {
	return uuid >> (SequenceIdBits + ProjectIdBits)
}

func TimeFormatFromUUID(uuid uint64) string {
	return time.Unix(int64(uuid>>(SequenceIdBits+ProjectIdBits))/1e3,
		int64(uuid>>(SequenceIdBits+ProjectIdBits))%1e3*1e6).Format("2006-01-02 15:04:05.000")
}
