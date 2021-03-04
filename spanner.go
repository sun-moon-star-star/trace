package trace

import (
	"time"
	"trace/uuid"
)

type SpanMap map[string]ValueInfo

type ValueInfo struct {
	Id    uint64
	Value interface{}
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

	strategy SpannerStrategy `json:"-"`

	Father   *Spanner `json:"-"`
	Children *Spanner `json:"-"`
	Next     *Spanner `json:"-"`
}

func NewSpanner(SpanName string, TraceId, ParentSpanId uint64) *Spanner {
	s := &Spanner{
		SpanId:       uuid.NewUUID(),
		SpanName:     SpanName,
		TraceId:      TraceId,
		ParentSpanId: ParentSpanId,

		Tags:     make(SpanMap),
		Logs:     make(SpanMap),
		Baggages: make(SpanMap),

		strategy: GlobalDefaultSpannerStrategy,
	}
	if ParentSpanId == 0 {
		s.ParentSpanId = s.SpanId
	}
	return s
}

func (s *Spanner) NewSpanner(SpanName string) *Spanner {
	return NewSpanner(SpanName, s.TraceId, s.SpanId)
}

func (s *Spanner) NewPartSpanner(SpanName string) *Spanner {
	obj := NewSpanner(SpanName, s.TraceId, s.SpanId)
	obj.Flags |= SpanType
	return obj
}

func (s *Spanner) End() {
	if s.Flags&SpannerStop != SpannerStop {
		s.EndTimestamp = uint64(time.Now().UnixNano() / 1e6)
		s.Flags |= SpannerStop
	}
}

func (s *Spanner) Tag(key string, value interface{}) {
	s.Tags[key] = ValueInfo{
		Id:    uuid.NewUUID(),
		Value: value,
	}
}

func (s *Spanner) Log(key string, value interface{}) {
	s.Logs[key] = ValueInfo{
		Id:    uuid.NewUUID(),
		Value: value,
	}
}

func (s *Spanner) Baggage(key string, value interface{}) {
	s.Baggages[key] = ValueInfo{
		Id:    uuid.NewUUID(),
		Value: value,
	}
}

func (s *Spanner) TagString() string {
	if s.strategy == GlobalDefaultSpannerStrategy {
		return GlobalDefaultSpannerStrategy.SpannerTagStrategy(s.Tags)
	}
	return s.strategy.(*UserDefineSpannerStrategy).Tag(s.Tags)
}

func (s *Spanner) LogString() string {
	if s.strategy == GlobalDefaultSpannerStrategy {
		return GlobalDefaultSpannerStrategy.SpannerLogStrategy(s.Logs)
	}
	return s.strategy.(*UserDefineSpannerStrategy).Log(s.Logs)
}

func (s *Spanner) BaggageString() string {
	if s.strategy == GlobalDefaultSpannerStrategy {
		return GlobalDefaultSpannerStrategy.SpannerBaggageStrategy(s.Baggages)
	}
	return s.strategy.(*UserDefineSpannerStrategy).Baggage(s.Baggages)
}

func (s *Spanner) String() string {
	if s.strategy == GlobalDefaultSpannerStrategy {
		return GlobalDefaultSpannerStrategy.SpannerStrategy(s)
	}
	return s.strategy.(*UserDefineSpannerStrategy).Spanner(s)
}

func (s *Spanner) SummaryString() string {
	if s.strategy == GlobalDefaultSpannerStrategy {
		return GlobalDefaultSpannerStrategy.SpannerSummaryStrategy(s)
	}
	return s.strategy.(*UserDefineSpannerStrategy).Summary(s)
}

func (s *Spanner) TagStrategy(Tag func(SpanMap) string) {
	if s.strategy == GlobalDefaultSpannerStrategy {
		s.strategy = NewUserDefineSpannerStrategy()
	}
	s.strategy.(*UserDefineSpannerStrategy).Tag = Tag
}

func (s *Spanner) LogStrategy(Log func(SpanMap) string) {
	if s.strategy == GlobalDefaultSpannerStrategy {
		s.strategy = NewUserDefineSpannerStrategy()
	}
	s.strategy.(*UserDefineSpannerStrategy).Log = Log
}

func (s *Spanner) BaggageStrategy(Baggage func(SpanMap) string) {
	if s.strategy == GlobalDefaultSpannerStrategy {
		s.strategy = NewUserDefineSpannerStrategy()
	}
	s.strategy.(*UserDefineSpannerStrategy).Baggage = Baggage
}

func (s *Spanner) SpannerStrategy(spanner func(*Spanner) string) {
	if s.strategy == GlobalDefaultSpannerStrategy {
		s.strategy = NewUserDefineSpannerStrategy()
	}
	s.strategy.(*UserDefineSpannerStrategy).Spanner = spanner
}

func (s *Spanner) SummaryStrategy(summary func(*Spanner) string) {
	if s.strategy == GlobalDefaultSpannerStrategy {
		s.strategy = NewUserDefineSpannerStrategy()
	}
	s.strategy.(*UserDefineSpannerStrategy).Summary = summary
}
