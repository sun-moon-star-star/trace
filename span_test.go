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

	spanner.FormatBaggageMapStrategy = func(maps trace.BaggageMap) string {
		var info string
		for key, value := range maps {
			info += fmt.Sprintf(" [%s(%s)->: %+v]", key, value.Time.Format("15:04:05.000000"), value.Value)
		}
		return info
	}

	spanner.TraceName = "trace_test"
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

	t.Log(spanner.FormatSpannerStrategy(spanner))

	spanner.End()

	t.Log(spanner.FormatTagMapStrategy(spanner.Tags))
	t.Log(spanner.FormatLogMapStrategy(spanner.Logs))
	t.Log(spanner.FormatBaggageMapStrategy(spanner.Baggages))

	t.Log(spanner.FormatSpannerStrategy(spanner))
	t.Log(spanner.FormatSpannerStrategy(spanner))
}
