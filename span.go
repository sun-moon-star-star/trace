package trace

import (
	"fmt"
	"time"
	"trace/random"
)

type SpanType uint64

const NoSet SpanType = SpanType(0)
const ChildOf SpanType = SpanType(1)
const FollowsFrom SpanType = SpanType(2)

type TagMap map[string]interface{}

type ValueWithTime struct {
	Time  time.Time
	Value interface{}
}

type LogMap map[string]ValueWithTime
type BaggageMap map[string]ValueWithTime

type FormatTagMapStrategy func(maps TagMap) string
type FormatLogMapStrategy func(maps LogMap) string
type FormatBaggageMapStrategy func(maps BaggageMap) string

var DefaultFormatTagMapStrategy FormatTagMapStrategy = func(maps TagMap) string {
	var info string
	for key, value := range maps {
		info += fmt.Sprintf(" [%s: %+v]", key, value)
	}
	return info
}

var DefaultFormatLogMapStrategy FormatLogMapStrategy = func(maps LogMap) string {
	var info string
	for key, value := range maps {
		info += fmt.Sprintf(" [%s(%s): %+v]", key, value.Time.Format("15:04:05.000000"), value.Value)
	}
	return info
}

var DefaultFormatBaggageMapStrategy FormatBaggageMapStrategy = func(maps BaggageMap) string {
	var info string
	for key, value := range maps {
		info += fmt.Sprintf(" [%s(%s): %+v]", key, value.Time.Format("15:04:05.000000"), value.Value)
	}
	return info
}

type FormatSpannerStrategy func(*Spanner) string

var DefaultFormatSpannerStrategy FormatSpannerStrategy = func(s *Spanner) string {
	var info string

	if s.Stop {
		info = fmt.Sprintf("[%s, %s] ",
			s.StartTime.Format("2006-01-02 15:04:05.000000"),
			s.EndTime.Format("2006-01-02 15:04:05.000000"))
	} else {
		info = fmt.Sprintf("[%s, %s] ",
			s.StartTime.Format("2006-01-02 15:04:05.000000"),
			time.Now().Format("2006-01-02 15:04:05.000000"))
	}

	info += fmt.Sprintf("TraceId: %s, SpanId: %d, SpanName: %s", s.TraceId, s.SpanId, s.SpanName)

	info += ", Tags: {" + s.FormatTagMapStrategy(s.Tags) + " }"
	info += ", Logs: {" + s.FormatLogMapStrategy(s.Logs) + " }"
	info += ", Baggages: {" + s.FormatBaggageMapStrategy(s.Baggages) + " }"

	return info
}

type Spanner struct {
	TraceId string

	SpanId   uint64
	SpanName string

	ParentSpanId uint64
	SpanType     SpanType

	Stop      bool
	StartTime time.Time
	EndTime   time.Time

	Tags     TagMap
	Logs     LogMap
	Baggages BaggageMap

	FormatTagMapStrategy     FormatTagMapStrategy
	FormatLogMapStrategy     FormatLogMapStrategy
	FormatBaggageMapStrategy FormatBaggageMapStrategy

	FormatSpannerStrategy FormatSpannerStrategy

	Father   *Spanner
	Children *Spanner
	Next     *Spanner
}

func NewSpanner() *Spanner {
	return &Spanner{
		SpanId: random.RandomUint64(),

		StartTime: time.Now(),

		Tags:     make(TagMap),
		Logs:     make(LogMap),
		Baggages: make(BaggageMap),

		FormatTagMapStrategy:     DefaultFormatTagMapStrategy,
		FormatLogMapStrategy:     DefaultFormatLogMapStrategy,
		FormatBaggageMapStrategy: DefaultFormatBaggageMapStrategy,

		FormatSpannerStrategy: DefaultFormatSpannerStrategy,
	}
}

func (s *Spanner) End() {
	s.EndTime = time.Now()
	s.Stop = true
}

func (s *Spanner) Tag(key string, value interface{}) {
	s.Tags[key] = value
}

func (s *Spanner) Log(key string, value interface{}) {
	s.Logs[key] = ValueWithTime{
		Time:  time.Now(),
		Value: value,
	}
}

func (s *Spanner) Baggage(key string, value interface{}) {
	s.Baggages[key] = ValueWithTime{
		Time:  time.Now(),
		Value: value,
	}
}
