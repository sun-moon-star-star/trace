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
	Spanner := trace.NewSpanner()

	Spanner.FormatBaggageMapStrategy = func(maps trace.BaggageMap) string {
		var info string
		for key, value := range maps {
			info += fmt.Sprintf(" [%s(%s)->: %+v]", key, value.Time.Format("15:04:05.000000"), value.Value)
		}
		return info
	}

	Spanner.TraceName = "trace_test"
	Spanner.SpanName = "span_test"

	Spanner.Tag("id", 123456789)
	Spanner.Tag("name", "zhaolu")

	Spanner.Log("error", errors.New("unknown error"))
	Spanner.Log("action", "success")

	info := TestInfo{
		Id:   random.RandomUint64(),
		Desc: "test_info",
		Time: time.Now(),
	}
	Spanner.Log("info", info)

	Spanner.Baggage("data-access", "(0, ok)")

	t.Log(Spanner.FormatSpannerStrategy(Spanner))

	Spanner.End()

	t.Log(Spanner.FormatTagMapStrategy(Spanner.Tags))
	t.Log(Spanner.FormatLogMapStrategy(Spanner.Logs))
	t.Log(Spanner.FormatBaggageMapStrategy(Spanner.Baggages))

	t.Log(Spanner.FormatSpannerStrategy(Spanner))
	t.Log(Spanner.FormatSpannerStrategy(Spanner))
}
