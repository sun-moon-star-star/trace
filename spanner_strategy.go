package trace

import (
	"fmt"
	"time"
	"trace/uuid"
)

type SpannerStrategy interface {
	// Output Format
	SpannerTagStrategy(SpanMap) string
	SpannerLogStrategy(SpanMap) string
	SpannerBaggageStrategy(SpanMap) string

	// Output Format
	SpannerStrategy(*Spanner) string

	// Generate Summary To Save
	SpannerSummaryStrategy(*Spanner) string
}

type UserDefineSpannerStrategy struct {
	Tag     func(SpanMap) string
	Log     func(SpanMap) string
	Baggage func(SpanMap) string
	// Output Format
	Spanner func(*Spanner) string
	// Generate Summary To Save
	Summary func(*Spanner) string
}

func NewUserDefineSpannerStrategy() *UserDefineSpannerStrategy {
	return &UserDefineSpannerStrategy{
		Tag:     GlobalDefaultSpannerStrategy.SpannerTagStrategy,
		Log:     GlobalDefaultSpannerStrategy.SpannerLogStrategy,
		Baggage: GlobalDefaultSpannerStrategy.SpannerBaggageStrategy,
		Spanner: GlobalDefaultSpannerStrategy.SpannerStrategy,
		Summary: GlobalDefaultSpannerStrategy.SpannerSummaryStrategy,
	}
}

func (s *UserDefineSpannerStrategy) SpannerTagStrategy(maps SpanMap) string {
	return s.Tag(maps)
}

func (s *UserDefineSpannerStrategy) SpannerLogStrategy(maps SpanMap) string {
	return s.Log(maps)
}

func (s *UserDefineSpannerStrategy) SpannerBaggageStrategy(maps SpanMap) string {
	return s.Baggage(maps)
}

func (s *UserDefineSpannerStrategy) SpannerStrategy(spanner *Spanner) string {
	return s.Spanner(spanner)
}

func (s *UserDefineSpannerStrategy) SpannerSummaryStrategy(spanner *Spanner) string {
	return s.Summary(spanner)
}

type DefaultSpannerStrategy struct{}

var GlobalDefaultSpannerStrategy SpannerStrategy = &DefaultSpannerStrategy{}

func (*DefaultSpannerStrategy) SpannerTagStrategy(maps SpanMap) string {
	var info string
	for key, value := range maps {
		info += fmt.Sprintf(" [%s(%s): %+v]", key, uuid.TimeFormatFromUUID(value.Id), value.Value)
	}
	return info
}

func (*DefaultSpannerStrategy) SpannerLogStrategy(maps SpanMap) string {
	var info string
	for key, value := range maps {
		info += fmt.Sprintf(" [%s(%s): %+v]", key, uuid.TimeFormatFromUUID(value.Id), value.Value)
	}
	return info
}

func (*DefaultSpannerStrategy) SpannerBaggageStrategy(maps SpanMap) string {
	var info string
	for key, value := range maps {
		info += fmt.Sprintf(" [%s(%s)->: %+v]", key, uuid.TimeFormatFromUUID(value.Id), value.Value)
	}
	return info
}

func (strategy *DefaultSpannerStrategy) SpannerStrategy(s *Spanner) string {
	var info string

	if s.Flags&SpannerStop > 0 {
		info = fmt.Sprintf("[%s, %s]",
			uuid.TimeFormatFromUUID(s.SpanId),
			time.Unix(int64(s.EndTimestamp)/1e3, int64(s.EndTimestamp)%1e3*1e6).Format("2006-01-02 15:04:05.000"))
	} else {
		info = fmt.Sprintf("[%s, %s]",
			uuid.TimeFormatFromUUID(s.SpanId),
			time.Now().Format("2006-01-02 15:04:05.000"))
	}

	info += fmt.Sprintf(" [TraceId: %d] [SpanId: %d] [SpanName: %s]",
		s.TraceId, s.SpanId, s.SpanName)

	info += s.TagString()
	info += s.LogString()
	info += s.BaggageString()

	return info
}

func (strategy *DefaultSpannerStrategy) SpannerSummaryStrategy(s *Spanner) string {
	return s.String()
}
