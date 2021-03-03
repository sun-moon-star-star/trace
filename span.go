package trace

import (
	"fmt"
	"time"
)

type SpanMap map[string]ValueInfo

type ValueInfo struct {
	Id    uint64
	Value interface{}
}

type Strategy struct {
	// Output Format
	Tag     func(SpanMap) string
	Log     func(SpanMap) string
	Baggage func(SpanMap) string
	// Output Format
	Spanner func(*Spanner) string
	// Generate Summary To Save
	Summary func(*Spanner) string
}

var DefaultStrategy *Strategy

func init() {
	DefaultStrategy = &Strategy{
		Tag: func(maps SpanMap) string {
			var info string
			for key, value := range maps {
				info += fmt.Sprintf(" [%s(%s): %+v]", key, TimeFormatFromUUID(value.Id), value.Value)
			}
			return info
		},
		Log: func(maps SpanMap) string {
			var info string
			for key, value := range maps {
				info += fmt.Sprintf(" [%s(%s): %+v]", key, TimeFormatFromUUID(value.Id), value.Value)
			}
			return info
		},
		Baggage: func(maps SpanMap) string {
			var info string
			for key, value := range maps {
				info += fmt.Sprintf(" [%s(%s): %+v]", key, TimeFormatFromUUID(value.Id), value.Value)
			}
			return info
		},
		Spanner: func(s *Spanner) string {
			var info string

			if s.Flags&SpannerStop > 0 {
				info = fmt.Sprintf("[%s, %s]",
					TimeFormatFromUUID(s.SpanId),
					time.Unix(int64(s.EndTimestamp)/1e3, int64(s.EndTimestamp)%1e3*1e6).Format("2006-01-02 15:04:05.000"))
			} else {
				info = fmt.Sprintf("[%s, %s]",
					TimeFormatFromUUID(s.SpanId),
					time.Now().Format("2006-01-02 15:04:05.000"))
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

	Flags        uint64
	EndTimestamp uint64

	Summary string

	Tags     SpanMap
	Logs     SpanMap
	Baggages SpanMap

	Strategy *Strategy `json:"-"`

	Father   *Spanner `json:"-"`
	Children *Spanner `json:"-"`
	Next     *Spanner `json:"-"`
}

func NewSpanner() *Spanner {
	return &Spanner{
		SpanId: NewUUID(),

		Tags:     make(SpanMap),
		Logs:     make(SpanMap),
		Baggages: make(SpanMap),

		Strategy: DefaultStrategy,
	}
}

func (s *Spanner) End() {
	s.EndTimestamp = uint64(time.Now().UnixNano() / 1e6)
	s.Flags |= SpannerStop
}

func (s *Spanner) Tag(key string, value interface{}) {
	s.Tags[key] = ValueInfo{
		Id:    NewUUID(),
		Value: value,
	}
}

func (s *Spanner) Log(key string, value interface{}) {
	s.Logs[key] = ValueInfo{
		Id:    NewUUID(),
		Value: value,
	}
}

func (s *Spanner) Baggage(key string, value interface{}) {
	s.Baggages[key] = ValueInfo{
		Id:    NewUUID(),
		Value: value,
	}
}
