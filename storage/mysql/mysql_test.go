package mysql_test

import (
	"errors"
	"testing"
	"trace"
	"trace/random"
	"trace/storage"
	"trace/storage/mysql"
)

func TestMysql(t *testing.T) {
	var saver storage.API
	saver = &mysql.Mysql{}

	var err error

	tracer := trace.NewTracer()
	tracer.TraceName, err = random.RandomUUID()

	if err != nil {
		t.Fatal(err)
	}

	a := tracer.NewSpanner()
	a.SpanName = "A"
	a.Tag("A-Tag", 1)
	a.Log("A-Log", "alog")
	a.Baggage("A-Baggage", errors.New("unknown error"))
	a.End()

	err = saver.SaveSpanner(a)
	if err != nil {
		t.Fatal(err)
	}

	b := tracer.NewSpanner()
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
