package mysql

import (
	"encoding/json"
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
	if tracer.TraceId == 0 {
		return errors.New("trace.Tracer.TraceId must be setting")
	}

	params := map[string]interface{}{
		"table": trace.GlobalConfig.Mysql.TraceTableName,
		"data": &trace.Trace{
			TraceId:   tracer.TraceId,
			TraceName: tracer.TraceName,
			StartTime: tracer.StartTime.Format(trace.GlobalConfig.Server.DefaultTimeLayout),
			EndTime:   tracer.EndTime.Format(trace.GlobalConfig.Server.DefaultTimeLayout),
			Summary:   tracer.Summary,
		},
	}

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

func (mysql *Mysql) SaveSpanner(spanner *trace.Spanner) error {
	if spanner.TraceId == 0 {
		return errors.New("trace.Spanner.TraceId must be setting")
	}

	// InsertSpan
	if spanner.Summary == "" {
		spanner.Summary = spanner.Strategy.Summary(spanner)
	}

	_, err := InsertTable(map[string]interface{}{
		"table": trace.GlobalConfig.Mysql.SpanTableName,
		"data": &trace.Span{
			SpanId:       spanner.SpanId,
			ParentSpanId: spanner.ParentSpanId,
			SpanName:     spanner.SpanName,
			TraceId:      spanner.TraceId,
			StartTime:    spanner.StartTime.Format(trace.GlobalConfig.Server.DefaultTimeLayout),
			EndTime:      spanner.StartTime.Format(trace.GlobalConfig.Server.DefaultTimeLayout),
			Summary:      spanner.Summary,
			Flags:        spanner.Flags,
		},
	})
	if err != nil {
		return err
	}

	// Insert Tag
	tags := make([]trace.Tag, len(spanner.Tags))
	index := 0
	for key, value := range spanner.Tags {
		tags[index].SpanId = spanner.SpanId

		tags[index].Field = key
		bytes, err := json.Marshal(value)
		if err != nil {
			tags[index].Value = err.Error()
		} else {
			tags[index].Value = string(bytes)
		}

		index++
	}

	_, err = InsertTable(map[string]interface{}{
		"table": trace.GlobalConfig.Mysql.TagTableName,
		"data":  &tags,
	})
	if err != nil {
		return err
	}

	// Insert Log
	logs := make([]trace.Log, len(spanner.Logs))
	index = 0
	for key, value := range spanner.Logs {
		logs[index].SpanId = spanner.SpanId

		logs[index].Field = key
		logs[index].Time = value.Time.Format(trace.GlobalConfig.Server.LogTimeLayout)
		bytes, err := json.Marshal(value.Value)
		if err != nil {
			logs[index].Value = err.Error()
		} else {
			logs[index].Value = string(bytes)
		}

		index++
	}

	_, err = InsertTable(map[string]interface{}{
		"table": trace.GlobalConfig.Mysql.LogTableName,
		"data":  &logs,
	})
	if err != nil {
		return err
	}

	// Insert Baggage
	baggages := make([]trace.Baggage, len(spanner.Baggages))
	index = 0
	for key, value := range spanner.Baggages {
		baggages[index].SpanId = spanner.SpanId

		baggages[index].Field = key
		baggages[index].Time = value.Time.Format(trace.GlobalConfig.Server.BaggageTimeLayout)
		bytes, err := json.Marshal(value.Value)
		if err != nil {
			baggages[index].Value = err.Error()
		} else {
			baggages[index].Value = string(bytes)
		}

		index++
	}

	_, err = InsertTable(map[string]interface{}{
		"table": trace.GlobalConfig.Mysql.LogTableName,
		"data":  &baggages,
	})

	return err
}
