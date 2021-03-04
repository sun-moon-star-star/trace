package trace_test

import (
	"errors"
	"fmt"
	"testing"
	"time"
	"trace"
	"trace/random"
	"trace/uuid"
)

type TestInfo struct {
	Id   uint64
	Desc string
	Time time.Time
}

func TestSpan(t *testing.T) {
	tracer := trace.NewTracer("tracer_test")
	spanner := trace.NewSpanner("spanner_test", tracer.TraceId, 0)

	if spanner.SpanId != spanner.ParentSpanId {
		t.Fatal("spanner.SpanId != spanner.ParentSpanId")
	}

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

	t.Log(spanner.String())

	spanner.End()

	t.Log(spanner.TagString())
	t.Log(spanner.LogString())
	t.Log(spanner.BaggageString())

	spanner.BaggageStrategy(func(maps trace.SpanMap) string {
		var info string
		for key, value := range maps {
			info += fmt.Sprintf(" [%s(%s)baby: %+v]", key, uuid.TimeFormatFromUUID(value.Id), value.Value)
		}
		return info
	})

	t.Log(spanner.String())
	t.Log(spanner.String())
}
