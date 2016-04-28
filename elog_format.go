package log

import (
	"text/template"
)

var elogTags = []string{
	"{{rpc_id}}", "{{rpc_id .}}",
	"{{request_id}}", "{{request_id .}}",
}

// ELogForamtter is the formatter for elog.
type ELogForamtter struct {
	*BaseFormatter
}

// NewELogFormatter create a ELogFormatter with colored.
func NewELogFormatter(colored bool) *ELogForamtter {
	ef := new(ELogForamtter)
	ef.BaseFormatter = NewBaseFormatter(colored)
	ef.colored = colored
	ef.AddTags(elogTags...)
	ef.AddFuncMap(template.FuncMap{
		"rpc_id":     ef._rpcID,
		"request_id": ef._requestID,
	})
	return ef
}

// Format formats a Record with set format
func (f *ELogForamtter) Format(record Record) []byte {
	return f.BaseFormatter.Format(record)
}

func (f *ELogForamtter) _rpcID(r *ELogRecord) string {
	s := r.rpcID
	if s == "" {
		s = "-"
	}
	if f.Colored() {
		s = f.paint(r.lv, s)
	}
	return s
}

func (f *ELogForamtter) _requestID(r *ELogRecord) string {
	s := r.requestID
	if s == "" {
		s = "-"
	}
	if f.Colored() {
		s = f.paint(r.lv, s)
	}
	return s
}
