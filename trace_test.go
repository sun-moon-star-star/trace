package trace_test

import (
	"testing"
	"trace"
	"trace/random"
)

func TestTracer(t *testing.T) {
	var err error

	tracer := trace.NewTracer()
	tracer.TraceName, err = random.RandomUUID()

	if err != nil {
		t.Fatal(err)
	}

	tracer.End()
	tracer.Summary = "love babe"

	t.Log(tracer)
}
