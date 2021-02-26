package trace_test

import (
	"errors"
	"testing"
	"trace"
)

func TestSpan(t *testing.T) {
	Spanner := trace.NewSpanner()

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
    span_test.go:20: [2021-02-26 13:58:13.712473, 2021-02-26 13:58:13.712600] : { Tags: [id: 123456789] [name: zhaolu]} { Logs: [error(13:58:13.712475): unknown error] [action(13:58:13.712476): success]} { Baggage: [data-access(13:58:13.712476): (0, ok)]}
    span_test.go:24:  [id: 123456789] [name: zhaolu]
    span_test.go:25:  [error(13:58:13.712475): unknown error] [action(13:58:13.712476): success]
    span_test.go:26:  [data-access(13:58:13.712476): (0, ok)]
    span_test.go:28: [2021-02-26 13:58:13.712473, 2021-02-26 13:58:13.712647] : { Tags: [id: 123456789] [name: zhaolu]} { Logs: [error(13:58:13.712475): unknown error] [action(13:58:13.712476): success]} { Baggage: [data-access(13:58:13.712476): (0, ok)]}
    span_test.go:29: [2021-02-26 13:58:13.712473, 2021-02-26 13:58:13.712647] : { Tags: [id: 123456789] [name: zhaolu]} { Logs: [error(13:58:13.712475): unknown error] [action(13:58:13.712476): success]} { Baggage: [data-access(13:58:13.712476): (0, ok)]}
--- PASS: TestSpan (0.00s)
PASS
ok      trace   0.007s
*/
