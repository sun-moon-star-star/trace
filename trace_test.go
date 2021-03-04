package trace_test

import (
	"testing"
	"trace"
)

func TestTracer(t *testing.T) {
	tracer := trace.NewTracer("trace_test")

	tracer.End()
	tracer.Summary = "love babe"

	t.Log(tracer)
}
