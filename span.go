package trace

import (
	"fmt"
	"time"
)

type TagMap map[string]interface{}

type ValueWithTime struct {
	Time  time.Time
	Value interface{}
}

type LogMap map[string]ValueWithTime
type BaggageMap map[string]ValueWithTime

type Strategy struct {
	// Output Format
	Tag     func(maps TagMap) string
	Log     func(maps LogMap) string
	Baggage func(maps BaggageMap) string
	// Output Format
	Spanner func(*Spanner) string
	// Generate Summary To Save
	Summary func(*Spanner) string
}

var DefaultStrategy *Strategy

func init() {
	DefaultStrategy = &Strategy{
		Tag: func(maps TagMap) string {
			var info string
			for key, value := range maps {
				info += fmt.Sprintf(" [%s: %+v]", key, value)
			}
			return info
		},
		Log: func(maps LogMap) string {
			var info string
			for key, value := range maps {
				info += fmt.Sprintf(" [%s(%s): %+v]", key, value.Time.Format(Config.Server.LogTimeLayout), value.Value)
			}
			return info
		},
		Baggage: func(maps BaggageMap) string {
			var info string
			for key, value := range maps {
				info += fmt.Sprintf(" [%s(%s): %+v]", key, value.Time.Format(Config.Server.BaggageTimeLayout), value.Value)
			}
			return info
		},
		Spanner: func(s *Spanner) string {
			var info string

			if s.Flags&SpannerStop > 0 {
				info = fmt.Sprintf("[%s, %s]",
					s.StartTime.Format(Config.Server.DefaultTimeLayout),
					s.EndTime.Format(Config.Server.DefaultTimeLayout))
			} else {
				info = fmt.Sprintf("[%s, %s]",
					s.StartTime.Format(Config.Server.DefaultTimeLayout),
					time.Now().Format(Config.Server.DefaultTimeLayout))
			}

			info += fmt.Sprintf(" [TraceId: %d] [SpanId: %d] [SpanName: %s]",
				s.TraceId, s.SpanId, s.SpanName)

			info += s.Strategy.Tag(s.Tags)
			info += s.Strategy.Log(s.Logs)
			info += s.Strategy.Baggage(s.Baggages)

			return info
		},
		Summary: func(s *Spanner) string {
			return DefaultStrategy.Spanner(s)
		},
	}
}

const (
	SpannerStop = 1 << (iota)
	SpanType    // not set => FollowsFrom, set => ChildOf
)

type Spanner struct {
	SpanId       uint64
	SpanName     string
	TraceId      uint64
	ParentSpanId uint64

	Flags     uint64
	StartTime time.Time
	EndTime   time.Time

	Summary string

	Tags     TagMap
	Logs     LogMap
	Baggages BaggageMap

	Strategy *Strategy `json:"-"`

	Father   *Spanner `json:"-"`
	Children *Spanner `json:"-"`
	Next     *Spanner `json:"-"`
}

func NewSpanner() *Spanner {
	return &Spanner{
		SpanId:    NewUUID(),
		StartTime: time.Now(),

		Tags:     make(TagMap),
		Logs:     make(LogMap),
		Baggages: make(BaggageMap),

		Strategy: DefaultStrategy,
	}
}

func (s *Spanner) End() {
	s.EndTime = time.Now()
	s.Flags |= SpannerStop
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
