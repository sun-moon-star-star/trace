package trace

type Trace struct {
	TraceId   uint64 `json:"trace_id" gorm:"type:bigint(20) unsigned not null primaryKey;"`
	TraceName string `json:"trace_name" gorm:"varchar(255)"`
	StartTime string `json:"start_time" gorm:"datetime(6)"`
	EndTime   string `json:"end_time" gorm:"datetime(6)"`
	Summary   string `json:"summary" gorm:"varchar(4096)"`
	Flags     uint64 `json:"flags" gorm:"type:int(11) unsigned"`
}

type Span struct {
	SpanId       uint64 `json:"span_id" gorm:"type:bigint(20) unsigned not null primaryKey autoIncrement;"`
	FatherSpanId uint64 `json:"father_span_id" gorm:"type:bigint(20) unsigned not null;"`
	SpanName     string `json:"span_name" gorm:"varchar(255)"`
	TraceId      uint64 `json:"trace_id" gorm:"type:bigint(20)"`
	StartTime    string `json:"start_time" gorm:"datetime(6)"`
	EndTime      string `json:"end_time" gorm:"datetime(6)"`
	Summary      string `json:"summary" gorm:"varchar(4096)"`
	Flags        uint64 `json:"flags" gorm:"type:int(11) unsigned"`
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
	Time   string `json:"time" gorm:"datetime(6)"`
	SpanId uint64 `json:"span_id" gorm:"type:bigint(20) unsigned not null;"`
}

type Baggage struct {
	BaggageId uint64 `json:"baggage_id" gorm:"type:bigint(20) unsigned not null primaryKey autoIncrement;"`
	Field     string `json:"field" gorm:"varchar(255)"`
	Value     string `json:"value" gorm:"varchar(16384)"`
	Time      string `json:"time" gorm:"datetime(6)"`
	SpanId    uint64 `json:"span_id" gorm:"type:bigint(20) unsigned not null;"`
}
