package trace_test

import (
	"errors"
	"fmt"
	"testing"
	"time"
	"trace"
	"trace/random"
)

type TestInfo struct {
	Id   uint64
	Desc string
	Time time.Time
}

func TestSpan(t *testing.T) {
	spanner := trace.NewSpanner()

	spanner.Strategy.Baggage = func(maps trace.SpanMap) string {
		var info string
		for key, value := range maps {
			info += fmt.Sprintf(" [%s(%s)->: %+v]", key, trace.TimeFormatFromUUID(value.Id), value.Value)
		}
		return info
	}

	spanner.SpanName = "span_test"

	spanner.Tag("id", 123456789)
	spanner.Tag("name", "zhaolu")

	spanner.Log("error", errors.New("unknown error"))
	spanner.Log("action", "success")

	info := TestInfo{
		Id:   random.RandomUint64(),
		Desc: "test_info",
		Time: time.Now(),
	}
	spanner.Log("info", info)

	spanner.Baggage("data-access", "(0, ok)")

	t.Log(spanner.Strategy.Spanner(spanner))

	spanner.End()

	t.Log(spanner.Strategy.Tag(spanner.Tags))
	t.Log(spanner.Strategy.Log(spanner.Logs))
	t.Log(spanner.Strategy.Baggage(spanner.Baggages))

	t.Log(spanner.Strategy.Spanner(spanner))
	t.Log(spanner.Strategy.Spanner(spanner))
}
