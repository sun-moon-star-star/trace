package mysql

import (
	"errors"
	"time"
	"trace"
)

// type API interface {
// 	SaveTracer(*trace.Tracer) error
// 	LoadTracer(traceId string) (*trace.Tracer, error)
// 	SaveSpanner(*trace.Spanner) error
// 	LoadSpanner(spanId uint64) (*trace.Spanner, error)
// }

const TIME_LAYOUT = "2006-01-02 15:04:05.000000"

type Mysql struct{}

const (
	traceTableName = "trace"
)

func (mysql *Mysql) SaveTracer(tracer *trace.Tracer) error {
	params := make(map[string]interface{})
	trace := &trace.Trace{}

	if tracer.TraceId == 0 {
		return errors.New("tracer.TraceId must be setting")
	}

	trace.TraceId = tracer.TraceId
	trace.TraceName = tracer.TraceName
	trace.StartTime = tracer.StartTime.Format(TIME_LAYOUT)
	trace.EndTime = tracer.EndTime.Format(TIME_LAYOUT)
	trace.Summary = tracer.Summary

	params["table"] = traceTableName
	params["data"] = trace

	_, err := InsertTable(params)
	return err
}

func (mysql *Mysql) LoadTracer(tracer *trace.Tracer) (*trace.Tracer, error) {
	params := make(map[string]interface{})
	where, data := &trace.Trace{}, &trace.Trace{}

	where.TraceId = tracer.TraceId
	where.TraceName = tracer.TraceName

	params["table"] = traceTableName
	params["num"] = 1
	params["where"] = where
	params["data"] = data

	err := SelectTableLimit(params)
	if err != nil {
		return nil, err
	}

	startTime, err := time.Parse(TIME_LAYOUT, data.StartTime)
	if err != nil {
		return nil, err
	}

	endTime, err := time.Parse(TIME_LAYOUT, data.EndTime)
	if err != nil {
		return nil, err
	}

	res := &trace.Tracer{
		TraceId:   data.TraceId,
		TraceName: data.TraceName,

		Flags:     data.Flags,
		StartTime: startTime,
		EndTime:   endTime,
		Summary:   data.Summary,
	}

	return res, nil
}

func (mysql *Mysql) SaveSpanner(*trace.Spanner) error {

}
