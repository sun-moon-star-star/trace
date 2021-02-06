package trace

type Trace struct {
	Id uint64 `json:"id" gorm:"type:bigint(20) unsigned not null primaryKey autoIncrement;"`
	TraceId string `json:"trace_id" gorm:"varchar(32) not null"`
	TraceName string `json:"trace_name" gorm:"varchar(255)"`
	StartTime string `json:"start_time" gorm:"datetime(3)"`
	FinishTime string `json:"finish_time" gorm:"datetime(3)"`
	Summary string `json:"summary" gorm:"datetime(4096)"`
}

type Span struct {
	SpanId uint64 `json:"span_id" gorm:"type:bigint(20) unsigned not null primaryKey autoIncrement;"`
	SpanName string `json:"span_name" gorm:"varchar(255)"`
	StartTime string `json:"start_time" gorm:"datetime(3)"`
	FinishTime string `json:"finish_time" gorm:"datetime(3)"`
	TraceId string `json:"trace_id" gorm:"varchar(32) not null"`
}

type SpanReference struct {
	SpanReferenceId uint64 `json:"span_reference_id" gorm:"type:bigint(20) unsigned not null primaryKey autoIncrement;"`
	FatherSpanId uint64 `json:"father_span_id" gorm:"type:bigint(20) unsigned not null;"`
	ChildSpanId uint64 `json:"child_span_id" gorm:"type:bigint(20) unsigned not null;"`
	Relationship uint64 `json:"relationship" gorm:"type:int(11);"`
}

type Tag struct {
	TagId uint64 `json:"tag_id" gorm:"type:bigint(20) unsigned not null primaryKey autoIncrement;"`
	Field string `json:"field" gorm:"varchar(255)"`
	Value string `json:"value" gorm:"varchar(16384)"`
	SpanId uint64 `json:"span_id" gorm:"type:bigint(20) unsigned not null;"`
}

type Log struct {
	LogId uint64 `json:"log_id" gorm:"type:bigint(20) unsigned not null primaryKey autoIncrement;"`
	Field string `json:"field" gorm:"varchar(255)"`
	Value string `json:"value" gorm:"varchar(16384)"`
	LogTime string `json:"log_time" gorm:"datetime(3)"`
	SpanId uint64 `json:"span_id" gorm:"type:bigint(20) unsigned not null;"`
}

type Baggage struct {
	BaggageId uint64 `json:"baggage_id" gorm:"type:bigint(20) unsigned not null primaryKey autoIncrement;"`
	Field string `json:"field" gorm:"varchar(255)"`
	Value string `json:"value" gorm:"varchar(16384)"`
	SpanId uint64 `json:"span_id" gorm:"type:bigint(20) unsigned not null;"`
}