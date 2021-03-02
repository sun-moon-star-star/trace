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
//  SaveSpannerWithChildren(*trace.Spanner) error
// 	LoadSpanner(spanId uint64) (*trace.Spanner, error)
//	LoadSpannerWithChildren(spanId uint64) (*trace.Spanner, error)
// }

type Mysql struct{}

func (mysql *Mysql) SaveTracer(tracer *trace.Tracer) error {
	params := make(map[string]interface{})
	data := &trace.Trace{}

	if tracer.TraceId == 0 {
		return errors.New("tracer.TraceId must be setting")
	}

	data.TraceId = tracer.TraceId
	data.TraceName = tracer.TraceName
	data.StartTime = tracer.StartTime.Format(trace.GlobalConfig.Server.DefaultTimeLayout)
	data.EndTime = tracer.EndTime.Format(trace.GlobalConfig.Server.DefaultTimeLayout)
	data.Summary = tracer.Summary

	params["table"] = trace.GlobalConfig.Mysql.TraceTableName
	params["data"] = data

	_, err := InsertTable(params)
	return err
}

func (mysql *Mysql) LoadTracer(tracer *trace.Tracer) (*trace.Tracer, error) {
	params := make(map[string]interface{})
	where, data := &trace.Trace{}, &trace.Trace{}

	where.TraceId = tracer.TraceId
	where.TraceName = tracer.TraceName

	params["table"] = trace.GlobalConfig.Mysql.TraceTableName
	params["num"] = 1
	params["where"] = where
	params["data"] = data

	err := SelectTableLimit(params)
	if err != nil {
		return nil, err
	}

	startTime, err := time.Parse(trace.GlobalConfig.Server.DefaultTimeLayout, data.StartTime)
	if err != nil {
		return nil, err
	}

	endTime, err := time.Parse(trace.GlobalConfig.Server.DefaultTimeLayout, data.EndTime)
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
