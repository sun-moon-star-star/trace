package trace_test

import (
	"errors"
	"testing"
	"trace"
)

func TestSpan(t *testing.T) {
	Spanner := trace.NewSpanner()
	Spanner.TraceId = "trace_test"
	Spanner.SpanName = "span_test"

	Spanner.Tag("id", 123456789)
	Spanner.Tag("name", "zhaolu")

	Spanner.Log("error", errors.New("unknown error"))
	Spanner.Log("action", "success")

	Spanner.Baggage("data-access", "(0, ok)")

	t.Log(Spanner.FormatSpannerStrategy(Spanner))

	Spanner.End()

	t.Log(Spanner.FormatTagMapStrategy(Spanner.Tags))
	t.Log(Spanner.FormatLogMapStrategy(Spanner.Logs))
	t.Log(Spanner.FormatBaggageMapStrategy(Spanner.Baggages))

	t.Log(Spanner.FormatSpannerStrategy(Spanner))
	t.Log(Spanner.FormatSpannerStrategy(Spanner))
}
