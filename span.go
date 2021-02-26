package trace

import (
	"fmt"
	"time"
)

type SpanType uint64

const NoSet SpanType = SpanType(0)
const ChildOf SpanType = SpanType(1)
const FollowsFrom SpanType = SpanType(2)

type Map map[string]interface{}

func NewMap() Map { return make(Map) }

type FormatMapStrategy func(Map) string

var DefaultFormatMapStrategy FormatMapStrategy = func(maps Map) string {
	info := fmt.Sprintf("[%s] :", time.Now().Format("2006-01-02 15:04:05.000000"))
	for key, value := range maps {
		info += fmt.Sprintf(" [%s: %+v]", key, value)
	}
	return info
}

type FormatSpanerStrategy func(*Spaner) string

var DefaultFormatSpanerStrategy FormatSpanerStrategy = func(s *Spaner) string {
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

type Spaner struct {
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

	FormatMapStrategy    FormatMapStrategy
	FormatSpanerStrategy FormatSpanerStrategy

	Father   *Spaner
	Children *Spaner
	Next     *Spaner
}

func NewSpaner() *Spaner {
	return &Spaner{
		StartTime:            time.Now(),
		Tags:                 NewMap(),
		Logs:                 NewMap(),
		Baggages:             NewMap(),
		FormatMapStrategy:    DefaultFormatMapStrategy,
		FormatSpanerStrategy: DefaultFormatSpanerStrategy,
	}
}

func (s *Spaner) End() {
	s.EndTime = time.Now()
	s.Stop = true
}
