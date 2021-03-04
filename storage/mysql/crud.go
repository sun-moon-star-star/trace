package mysql

import (
	"errors"
	"fmt"
	"strings"
	"trace"

	_ "github.com/go-sql-driver/mysql"
)

func (mysql *Mysql) LoadTags(SpanId uint64) (tags trace.SpanMap, err error) {
	if err = setDB(); err != nil {
		return nil, err
	}

	var tags_data []trace.Tag
	res := db.Table(trace.Config.Mysql.TagTableName).Where(&trace.Tag{
		SpanId: SpanId,
	}).Find(&tags_data)

	if res.Error != nil {
		return nil, res.Error
	}

	tags = make(trace.SpanMap)
	for _, tag := range tags_data {
		tags[tag.Field] = trace.ValueInfo{
			Id:    tag.TagId,
			Value: tag.Value,
		}
	}

	return
}

func (mysql *Mysql) LoadLogs(SpanId uint64) (logs trace.SpanMap, err error) {
	if err = setDB(); err != nil {
		return nil, err
	}

	var logs_data []trace.Log
	res := db.Table(trace.Config.Mysql.LogTableName).Where(&trace.Log{
		SpanId: SpanId,
	}).Find(&logs_data)

	if res.Error != nil {
		return nil, res.Error
	}

	logs = make(trace.SpanMap)
	for _, log := range logs_data {
		logs[log.Field] = trace.ValueInfo{
			Id:    log.LogId,
			Value: log.Value,
		}
	}

	return
}

func (mysql *Mysql) LoadBaggages(SpanId uint64) (baggages trace.SpanMap, err error) {
	if err = setDB(); err != nil {
		return nil, err
	}

	var baggages_data []trace.Baggage
	res := db.Table(trace.Config.Mysql.BaggageTableName).Where(&trace.Baggage{
		SpanId: SpanId,
	}).Find(&baggages_data)

	if res.Error != nil {
		return nil, res.Error
	}

	baggages = make(trace.SpanMap)
	for _, baggage := range baggages_data {
		baggages[baggage.Field] = trace.ValueInfo{
			Id:    baggage.BaggageId,
			Value: baggage.Value,
		}
	}

	return
}

func (mysql *Mysql) SaveTracer(tracer *trace.Tracer) (err error) {
	if tracer.TraceId == 0 {
		return errors.New("trace.Tracer.TraceId must be setting")
	}

	if err = setDB(); err != nil {
		return err
	}

	res := db.Table(trace.Config.Mysql.TraceTableName).Create(&trace.Trace{
		TraceId:      tracer.TraceId,
		TraceName:    tracer.TraceName,
		EndTimestamp: tracer.EndTimestamp,
		Summary:      tracer.Summary,
	})

	if res.Error != nil {
		return res.Error
	}

	return nil
}

func (mysql *Mysql) LoadTracer(tracer *trace.Tracer) (err error) {
	if err = setDB(); err != nil {
		return err
	}

	trace_data := &trace.Trace{}

	res := db.Table(trace.Config.Mysql.TraceTableName).Where(&trace.Trace{
		TraceId:   tracer.TraceId,
		TraceName: tracer.TraceName,
	}).Find(trace_data).Limit(1)

	if res.Error != nil {
		return res.Error
	}

	tracer.TraceId = trace_data.TraceId
	tracer.TraceName = trace_data.TraceName
	tracer.Flags = trace_data.Flags
	tracer.Summary = trace_data.Summary
	tracer.EndTimestamp = trace_data.EndTimestamp

	tracer.First = &trace.Spanner{
		TraceId: trace_data.TraceId,
	}

	return mysql.LoadSpanner(tracer.First)
}

func (mysql *Mysql) SaveSpanner(spanner *trace.Spanner) (err error) {
	if spanner.TraceId == 0 {
		return errors.New("trace.Spanner.TraceId must be setting")
	}

	if err = setDB(); err != nil {
		return err
	}

	// InsertSpan
	if spanner.Summary == "" {
		spanner.Summary = spanner.SummaryString()
	}

	if spanner.SpanId == 0 {
		spanner.SpanId = trace.NewUUID()
	}

	if spanner.ParentSpanId == 0 {
		spanner.ParentSpanId = spanner.SpanId
	}

	tx := db.Begin()

	trace_data := &trace.Span{
		SpanId:       spanner.SpanId,
		ParentSpanId: spanner.ParentSpanId,
		SpanName:     spanner.SpanName,
		TraceId:      spanner.TraceId,
		EndTimestamp: spanner.EndTimestamp,
		Summary:      spanner.Summary,
		Flags:        spanner.Flags,
	}

	res := tx.Table(trace.Config.Mysql.SpanTableName).Create(trace_data)

	if res.Error != nil {
		tx.Rollback()
		return res.Error
	}

	// Insert Tag
	for key, value := range spanner.Tags {
		res := tx.Table(trace.Config.Mysql.TagTableName).Create(&trace.Tag{
			TagId:  value.Id,
			SpanId: spanner.SpanId,
			Field:  key,
			Value:  fmt.Sprintf("%+v", value.Value),
		})

		if res.Error != nil {
			tx.Rollback()
			return res.Error
		}
	}

	// Insert Log
	for key, value := range spanner.Logs {
		res := tx.Table(trace.Config.Mysql.LogTableName).Create(&trace.Log{
			LogId:  value.Id,
			SpanId: spanner.SpanId,
			Field:  key,
			Value:  fmt.Sprintf("%+v", value.Value),
		})

		if res.Error != nil {
			tx.Rollback()
			return res.Error
		}
	}

	// Insert Baggage
	for key, value := range spanner.Baggages {
		res := tx.Table(trace.Config.Mysql.BaggageTableName).Create(&trace.Baggage{
			BaggageId: value.Id,
			SpanId:    spanner.SpanId,
			Field:     key,
			Value:     fmt.Sprintf("%+v", value.Value),
		})

		if res.Error != nil {
			tx.Rollback()
			return res.Error
		}
	}

	tx.Commit()
	return nil
}

func (mysql *Mysql) LoadSpanner(spanner *trace.Spanner) error {
	if err := setDB(); err != nil {
		return err
	}

	span_data := &trace.Span{}

	res := db.Table(trace.Config.Mysql.SpanTableName).Where(&trace.Span{
		SpanId:       spanner.SpanId,
		ParentSpanId: spanner.ParentSpanId,
		SpanName:     spanner.SpanName,
		TraceId:      spanner.TraceId,
	}).Find(span_data).Order("span_id").Limit(1)

	if res.Error != nil {
		return res.Error
	}

	spanner.SpanId = span_data.SpanId
	spanner.ParentSpanId = span_data.ParentSpanId
	spanner.SpanName = span_data.SpanName
	spanner.TraceId = span_data.TraceId
	spanner.Summary = span_data.Summary
	spanner.Flags = span_data.Flags
	spanner.EndTimestamp = span_data.EndTimestamp

	var errors []string
	var err error

	spanner.Tags, err = mysql.LoadTags(span_data.SpanId)
	if err != nil {
		errors = append(errors, err.Error())
	}

	spanner.Logs, err = mysql.LoadLogs(span_data.SpanId)
	if err != nil {
		errors = append(errors, err.Error())
	}

	spanner.Baggages, err = mysql.LoadBaggages(span_data.SpanId)
	if err != nil {
		errors = append(errors, err.Error())
	}

	if len(errors) == 0 {
		return nil
	}
	return fmt.Errorf(strings.Join(errors, "\n"))
}
