package mysql

import (
	"errors"
	"fmt"
	"strings"
	"time"
	"trace"

	_ "github.com/go-sql-driver/mysql"

	"github.com/jinzhu/gorm"
)

var db *gorm.DB

type Mysql struct{}

func setDB() (err error) {
	if db != nil {
		if err = db.DB().Ping(); err == nil {
			return nil
		}
		db.DB().Close()
		db = nil
	}

	hostname := trace.GlobalConfig.Mysql.Hostname
	port := trace.GlobalConfig.Mysql.Port
	username := trace.GlobalConfig.Mysql.Username
	password := trace.GlobalConfig.Mysql.Password
	network := trace.GlobalConfig.Mysql.Network
	database := trace.GlobalConfig.Mysql.Database

	db_desc := fmt.Sprintf("%v:%v@%v(%v:%v)/%v",
		username, password, network, hostname, port, database)

	db, err = gorm.Open("mysql", db_desc)
	if err != nil {
		return err
	}

	err = db.DB().Ping()
	if err != nil {
		return err
	}

	db.DB().SetConnMaxLifetime(
		time.Duration(trace.GlobalConfig.Mysql.ConnMaxLifeTime) * time.Second)
	db.DB().SetMaxIdleConns(trace.GlobalConfig.Mysql.MaxIdleConns)
	db.DB().SetMaxOpenConns(trace.GlobalConfig.Mysql.MaxOpenConns)

	return nil
}

func (mysql *Mysql) LoadTags(SpanId uint64) (tags trace.TagMap, err error) {
	if err = setDB(); err != nil {
		return nil, err
	}

	var tags_data []trace.Tag
	res := db.Table(trace.GlobalConfig.Mysql.TagTableName).Where(&trace.Tag{
		SpanId: SpanId,
	}).Find(&tags_data)

	if res.Error != nil {
		return nil, res.Error
	}

	tags = make(trace.TagMap)
	for _, tag := range tags_data {
		tags[tag.Field] = tag.Value
	}

	return
}

func (mysql *Mysql) LoadLogs(SpanId uint64) (logs trace.LogMap, err error) {
	if err = setDB(); err != nil {
		return nil, err
	}

	var logs_data []trace.Log
	res := db.Table(trace.GlobalConfig.Mysql.LogTableName).Where(&trace.Log{
		SpanId: SpanId,
	}).Find(&logs_data)

	if res.Error != nil {
		return nil, res.Error
	}

	logs = make(trace.LogMap)
	for _, log := range logs_data {
		log_time, _ := time.Parse(
			trace.GlobalConfig.Server.DefaultTimeLayout, log.Time)

		logs[log.Field] = trace.ValueWithTime{
			Time:  log_time,
			Value: log.Value,
		}
	}

	return
}

func (mysql *Mysql) LoadBaggages(SpanId uint64) (baggages trace.BaggageMap, err error) {
	if err = setDB(); err != nil {
		return nil, err
	}

	var baggages_data []trace.Baggage
	res := db.Table(trace.GlobalConfig.Mysql.BaggageTableName).Where(&trace.Baggage{
		SpanId: SpanId,
	}).Find(&baggages_data)

	if res.Error != nil {
		return nil, res.Error
	}

	baggages = make(trace.BaggageMap)
	for _, baggage := range baggages_data {
		baggage_time, _ := time.Parse(
			trace.GlobalConfig.Server.DefaultTimeLayout, baggage.Time)

		baggages[baggage.Field] = trace.ValueWithTime{
			Time:  baggage_time,
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

	res := db.Table(trace.GlobalConfig.Mysql.TraceTableName).Create(&trace.Trace{
		TraceId:   tracer.TraceId,
		TraceName: tracer.TraceName,
		StartTime: tracer.StartTime.Format(trace.GlobalConfig.Server.DefaultTimeLayout),
		EndTime:   tracer.EndTime.Format(trace.GlobalConfig.Server.DefaultTimeLayout),
		Summary:   tracer.Summary,
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

	res := db.Table(trace.GlobalConfig.Mysql.TraceTableName).Where(&trace.Trace{
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

	startTime, err := time.Parse(
		trace.GlobalConfig.Server.DefaultTimeLayout, trace_data.StartTime)
	if err == nil {
		tracer.StartTime = startTime
	}

	endTime, err := time.Parse(
		trace.GlobalConfig.Server.DefaultTimeLayout, trace_data.EndTime)
	if err != nil {
		tracer.EndTime = endTime
	}

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
		spanner.Summary = spanner.Strategy.Summary(spanner)
	}

	if spanner.SpanId == 0 {
		spanner.SpanId = trace.GlobalTraceIdGenerator.GenerateTraceId(trace.TracerIdOption{ProjectId: 993})
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
		StartTime:    spanner.StartTime.Format(trace.GlobalConfig.Server.DefaultTimeLayout),
		EndTime:      spanner.StartTime.Format(trace.GlobalConfig.Server.DefaultTimeLayout),
		Summary:      spanner.Summary,
		Flags:        spanner.Flags,
	}

	res := tx.Table(trace.GlobalConfig.Mysql.SpanTableName).Create(trace_data)

	if res.Error != nil {
		tx.Rollback()
		return res.Error
	}

	// Insert Tag
	for key, value := range spanner.Tags {
		res := tx.Table(trace.GlobalConfig.Mysql.TagTableName).Create(&trace.Tag{
			SpanId: spanner.SpanId,
			Field:  key,
			Value:  fmt.Sprintf("%+v", value),
		})

		if res.Error != nil {
			tx.Rollback()
			return res.Error
		}
	}

	// Insert Log
	for key, value := range spanner.Logs {
		res := tx.Table(trace.GlobalConfig.Mysql.LogTableName).Create(&trace.Log{
			SpanId: spanner.SpanId,
			Field:  key,
			Value:  fmt.Sprintf("%+v", value.Value),
			Time:   value.Time.Format(trace.GlobalConfig.Server.LogTimeLayout),
		})

		if res.Error != nil {
			tx.Rollback()
			return res.Error
		}
	}

	// Insert Baggage
	for key, value := range spanner.Baggages {
		res := tx.Table(trace.GlobalConfig.Mysql.BaggageTableName).Create(&trace.Baggage{
			SpanId: spanner.SpanId,
			Field:  key,
			Value:  fmt.Sprintf("%+v", value.Value),
			Time:   value.Time.Format(trace.GlobalConfig.Server.LogTimeLayout),
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

	res := db.Table(trace.GlobalConfig.Mysql.SpanTableName).Where(&trace.Span{
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

	var errors []string

	startTime, err := time.Parse(
		trace.GlobalConfig.Server.DefaultTimeLayout, span_data.StartTime)
	if err == nil {
		spanner.StartTime = startTime
	} else {
		errors = append(errors, err.Error())
	}

	endTime, err := time.Parse(
		trace.GlobalConfig.Server.DefaultTimeLayout, span_data.EndTime)
	if err == nil {
		spanner.EndTime = endTime
	} else {
		errors = append(errors, err.Error())
	}

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
