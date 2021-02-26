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

/**
=== RUN   TestSpan
    span_test.go:22: [2021-02-26 14:32:44.526820, 2021-02-26 14:32:44.526953] TraceId: trace_test, SpanId: 2211634155081678851, SpanName: span_test, Tags: { [id: 123456789] [name: zhaolu] }, Logs: { [error(14:32:44.526828): unknown error] [action(14:32:44.526828): success] }, Baggages: { [data-access(14:32:44.526828): (0, ok)] }
    span_test.go:26:  [id: 123456789] [name: zhaolu]
    span_test.go:27:  [action(14:32:44.526828): success] [error(14:32:44.526828): unknown error]
    span_test.go:28:  [data-access(14:32:44.526828): (0, ok)]
    span_test.go:30: [2021-02-26 14:32:44.526820, 2021-02-26 14:32:44.527090] TraceId: trace_test, SpanId: 2211634155081678851, SpanName: span_test, Tags: { [id: 123456789] [name: zhaolu] }, Logs: { [error(14:32:44.526828): unknown error] [action(14:32:44.526828): success] }, Baggages: { [data-access(14:32:44.526828): (0, ok)] }
    span_test.go:31: [2021-02-26 14:32:44.526820, 2021-02-26 14:32:44.527090] TraceId: trace_test, SpanId: 2211634155081678851, SpanName: span_test, Tags: { [id: 123456789] [name: zhaolu] }, Logs: { [error(14:32:44.526828): unknown error] [action(14:32:44.526828): success] }, Baggages: { [data-access(14:32:44.526828): (0, ok)] }
--- PASS: TestSpan (0.00s)
PASS
ok      trace   0.013s
*/
