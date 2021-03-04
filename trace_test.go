package trace_test

import (
	"testing"
	"trace"
	"trace/random"
)

func TestTracer(t *testing.T) {
	TraceName, err := random.RandomUUID()
	if err != nil {
		t.Fatal(err)
	}

	tracer := trace.NewTracer(TraceName)

	tracer.End()
	tracer.Summary = "love babe"

	t.Log(tracer)
}
