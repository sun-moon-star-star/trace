package mysql

import (
	"strings"
	"trace"
)

func GetTableType(table string) interface{} {
	table = strings.ToLower(table)

	if table == "trace" {
		return &trace.Trace{}
	} else if table == "span" {
		return &trace.Span{}
	} else if table == "span_reference" {
		return &trace.SpanReference{}
	} else if table == "log" {
		return &trace.Log{}
	} else if table == "tag" {
		return &trace.Tag{}
	} else if table == "baggage" {
		return &trace.Baggage{}
	}

	return nil
}

func GetTableArrayType(table string) interface{} {
	table = strings.ToLower(table)

	if table == "trace" {
		return &[]trace.Trace{}
	} else if table == "span" {
		return &[]trace.Span{}
	} else if table == "span_reference" {
		return &[]trace.SpanReference{}
	} else if table == "log" {
		return &[]trace.Log{}
	} else if table == "tag" {
		return &[]trace.Tag{}
	} else if table == "baggage" {
		return &[]trace.Baggage{}
	}

	return nil
}
