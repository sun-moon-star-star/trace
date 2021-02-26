package trace_test

import (
	"errors"
	"testing"
	"trace"
)

func TestSpan(t *testing.T) {
	Spanner := trace.NewSpanner()

	Spanner.Tags["id"] = 123456789
	Spanner.Tags["name"] = "zhaolu"

	Spanner.Logs["error"] = errors.New("unknown error")
	Spanner.Logs["action"] = "success"

	Spanner.Baggages["data-access"] = "(0, ok)"

	t.Log(Spanner.FormatSpannerStrategy(Spanner))

	Spanner.End()

	t.Log(Spanner.FormatMapStrategy(Spanner.Tags))
	t.Log(Spanner.FormatMapStrategy(Spanner.Logs))
	t.Log(Spanner.FormatMapStrategy(Spanner.Baggages))

	t.Log(Spanner.FormatSpannerStrategy(Spanner))
	t.Log(Spanner.FormatSpannerStrategy(Spanner))
}

/**
=== RUN   TestSpan
    span_test.go:20: [2021-02-26 11:57:52.498281, 2021-02-26 11:57:52.498394] : {Baggages: [data-access: (0, ok)]} {Tags: [id: 123456789] [name: zhaolu]} {Logs: [error: unknown error] [action: success]}
    span_test.go:24: [2021-02-26 11:57:52.498474] : [id: 123456789] [name: zhaolu]
    span_test.go:25: [2021-02-26 11:57:52.498491] : [error: unknown error] [action: success]
    span_test.go:26: [2021-02-26 11:57:52.498506] : [data-access: (0, ok)]
    span_test.go:28: [2021-02-26 11:57:52.498281, 2021-02-26 11:57:52.498474] : {Baggages: [data-access: (0, ok)]} {Tags: [id: 123456789] [name: zhaolu]} {Logs: [error: unknown error] [action: success]}
    span_test.go:29: [2021-02-26 11:57:52.498281, 2021-02-26 11:57:52.498474] : {Baggages: [data-access: (0, ok)]} {Tags: [id: 123456789] [name: zhaolu]} {Logs: [error: unknown error] [action: success]}
--- PASS: TestSpan (0.00s)
PASS
ok      trace   0.006s
*/
