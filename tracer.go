package trace

type TraceIdOption {
	// 0-bit no use for extends
	// 41-bit millisecond timestamp
	// 10-bit projectId
	// 12-bit sequenceId
	ProjectId uint16
}

type Tracer struct {
	TraceId   uint64
	TraceName string // other key message

	Stop      bool
	StartTime string
	EndTime   string

	Summary string
}





