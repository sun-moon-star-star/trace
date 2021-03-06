package mysql_test

import (
	"errors"
	"testing"
	"trace"
	"trace/storage"
	"trace/storage/mysql"
)

func TestCRUD(t *testing.T) {
	var saver storage.CRUDAPI
	saver = &mysql.Mysql{}

	var err error

	tracer := trace.NewTracer("tracer_test")

	a := tracer.NewSpanner("span_test_a", 0)
	a.SpanName = "A"
	a.Tag("A-Tag", 1)
	a.Log("A-Log", "alog")
	a.Baggage("A-Baggage", errors.New("unknown error"))
	a.End()

	err = saver.SaveSpanner(a)
	if err != nil {
		t.Fatal(err)
	}

	b := a.NewSpanner("span_test_b")
	b.SpanName = "B"
	b.ParentSpanId = a.SpanId
	b.Tag("B-Tag", a)
	b.Log("B-Log", "bbbb")
	b.Baggage("B-Baggage", 15.90)
	b.End()

	err = saver.SaveSpanner(b)
	if err != nil {
		t.Fatal(err)
	}

	tracer.End()
	tracer.Summary = "love babe"
	err = saver.SaveTracer(tracer)
	if err != nil {
		t.Fatal(err)
	}
}
