package trace

type Trace struct {
	// projectId - timeId - randomId
	Id        uint64 `json:"id" gorm:"type:bigint(20) unsigned not null primaryKey autoIncrement;"`
	TraceId   string `json:"trace_id" gorm:"varchar(32) not null"`
	TraceName string `json:"trace_name" gorm:"varchar(255)"`
	StartTime string `json:"start_time" gorm:"datetime(6)"`
	EndTime   string `json:"end_time" gorm:"datetime(6)"`
	Summary   string `json:"summary" gorm:"datetime(4096)"`
}

type Span struct {
	SpanId    uint64 `json:"span_id" gorm:"type:bigint(20) unsigned not null primaryKey autoIncrement;"`
	SpanName  string `json:"span_name" gorm:"varchar(255)"`
	StartTime string `json:"start_time" gorm:"datetime(6)"`
	EndTime   string `json:"end_time" gorm:"datetime(6)"`
	TraceId   string `json:"trace_id" gorm:"varchar(32) not null"`
}

type SpanReference struct {
	SpanReferenceId uint64 `json:"span_reference_id" gorm:"type:bigint(20) unsigned not null primaryKey autoIncrement;"`
	FatherSpanId    uint64 `json:"father_span_id" gorm:"type:bigint(20) unsigned not null;"`
	ChildSpanId     uint64 `json:"child_span_id" gorm:"type:bigint(20) unsigned not null;"`
	Type            uint64 `json:"type" gorm:"type:int(11);"`
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
