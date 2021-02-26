package trace

import (
	"fmt"
	"time"
)

type SpanType uint64

const NoSet SpanType = SpanType(0)
const ChildOf SpanType = SpanType(1)
const FollowsFrom SpanType = SpanType(2)

type TagMap map[string]interface{}

func NewTagMap() TagMap { return make(TagMap) }

type FormatTagMapStrategy func(TagMap) string

var DefaultFormatMapStrategy FormatMapStrategy = func(maps Map) string {
	info := fmt.Sprintf("[%s] :", time.Now().Format("2006-01-02 15:04:05.000000"))
	for key, value := range maps {
		info += fmt.Sprintf(" [%s: %+v]", key, value)
	}
	return info
}

type FormatSpannerStrategy func(*Spanner) string

var DefaultFormatSpannerStrategy FormatSpannerStrategy = func(s *Spanner) string {
	var info string
	if s.Stop {
		info = fmt.Sprintf("[%s, %s] :",
			s.StartTime.Format("2006-01-02 15:04:05.000000"),
			s.EndTime.Format("2006-01-02 15:04:05.000000"))
	} else {
		info = fmt.Sprintf("[%s, %s] :",
			s.StartTime.Format("2006-01-02 15:04:05.000000"),
			time.Now().Format("2006-01-02 15:04:05.000000"))
	}

	info += " {Baggages:"
	for key, value := range s.Baggages {
		info += fmt.Sprintf(" [%s: %+v]", key, value)
	}
	info += "}"

	info += " {Tags:"
	for key, value := range s.Tags {
		info += fmt.Sprintf(" [%s: %+v]", key, value)
	}
	info += "}"

	info += " {Logs:"
	for key, value := range s.Logs {
		info += fmt.Sprintf(" [%s: %+v]", key, value)
	}
	info += "}"

	return info
}

type Spanner struct {
	TraceId      uint64
	SpanId       uint64
	ParentSpanId uint64
	SpanType     SpanType

	Stop      bool
	StartTime time.Time
	EndTime   time.Time

	Tags     Map
	Logs     Map
	Baggages Map

	FormatMapStrategy     FormatMapStrategy
	FormatSpannerStrategy FormatSpannerStrategy

	Father   *Spanner
	Children *Spanner
	Next     *Spanner
}

func NewSpanner() *Spanner {
	return &Spanner{
		StartTime:             time.Now(),
		Tags:                  NewMap(),
		Logs:                  NewMap(),
		Baggages:              NewMap(),
		FormatMapStrategy:     DefaultFormatMapStrategy,
		FormatSpannerStrategy: DefaultFormatSpannerStrategy,
	}
}

func (s *Spanner) End() {
	s.EndTime = time.Now()
	s.Stop = true
}
