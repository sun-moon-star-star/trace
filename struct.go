package trace

type Trace struct {
	TraceId      uint64 `json:"trace_id" gorm:"type:bigint(20) unsigned not null primaryKey autoIncrement;"`
	TraceName    string `json:"trace_name" gorm:"varchar(255)"`
	EndTimestamp uint64 `json:"end_timestamp" gorm:"type:bigint(20) unsigned"`
	Summary      string `json:"summary" gorm:"varchar(4096)"`
	Flags        uint64 `json:"flags" gorm:"type:bigint(20) unsigned"`
}

type Span struct {
	SpanId       uint64 `json:"span_id" gorm:"type:bigint(20) unsigned not null primaryKey autoIncrement;"`
	ParentSpanId uint64 `json:"parent_span_id" gorm:"type:bigint(20) unsigned;"`
	SpanName     string `json:"span_name" gorm:"varchar(255)"`
	TraceId      uint64 `json:"trace_id" gorm:"type:bigint(20) unsigned not null"`
	EndTimestamp uint64 `json:"end_timestamp" gorm:"type:bigint(20) unsigned"`
	Summary      string `json:"summary" gorm:"varchar(4096)"`
	Flags        uint64 `json:"flags" gorm:"type:bigint(20) unsigned"`
}

type Tag struct {
	TagId  uint64 `json:"tag_id" gorm:"type:bigint(20) unsigned not null primaryKey autoIncrement;"`
	Field  string `json:"field" gorm:"varchar(255)"`
	Value  string `json:"value" gorm:"varchar(16384)"`
	SpanId uint64 `json:"span_id" gorm:"type:bigint(20) unsigned not null;"`
}

type Log struct {
	LogId  uint64 `json:"log_id" gorm:"type:bigint(20) unsigned not null primaryKey autoIncrement;"`
	Field  string `json:"field" gorm:"varchar(255)"`
	Value  string `json:"value" gorm:"varchar(16384)"`
	SpanId uint64 `json:"span_id" gorm:"type:bigint(20) unsigned not null;"`
}

type Baggage struct {
	BaggageId uint64 `json:"baggage_id" gorm:"type:bigint(20) unsigned not null primaryKey autoIncrement;"`
	Field     string `json:"field" gorm:"varchar(255)"`
	Value     string `json:"value" gorm:"varchar(16384)"`
	SpanId    uint64 `json:"span_id" gorm:"type:bigint(20) unsigned not null;"`
}
