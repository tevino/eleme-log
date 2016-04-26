package log

import (
	"strings"
	"text/template"
)

var rpcTagReplacer = strings.NewReplacer(
	"{{rpc_id}}", "{{rpc_id .}}",
	"{{request_id}}", "{{request_id .}}",
)

// RPCFormatter is the formatter for RPC.
type RPCFormatter struct {
	*BaseFormatter
}

// NewRPCFormatter return
func NewRPCFormatter(format string, colored bool) (Formatter, error) {
	rpcFormatter := &RPCFormatter{
		BaseFormatter: new(BaseFormatter),
	}
	rpcFormatter.colored = colored

	if err := rpcFormatter.SetFormat(format); err != nil {
		return nil, err
	}
	return rpcFormatter, nil
}

// SetFormat parse foramt string to template.
func (f *RPCFormatter) SetFormat(format string) error {
	// {{ tag }} -> {{tag}}
	format = string(rTagLong.ReplaceAll([]byte(format), tagShort))

	format = tagReplacer.Replace(format)
	format = rpcTagReplacer.Replace(format)

	if !strings.HasSuffix(format, "\n") {
		format += "\n"
	}

	t, err := template.New("tpl").Funcs(f.funcMap()).Parse(format)
	if err != nil {
		return err
	}

	// TODO: validation

	f.tpl = t
	return nil
}

// Format formats a Record with set format
func (f *RPCFormatter) Format(record Record) []byte {
	return f.BaseFormatter.Format(record)
}

func (f *RPCFormatter) _rpcID(r *RPCRecord) string {
	s := r.rpcID
	if s == "" {
		s = "-"
	}
	if f.Colored() {
		s = f.paint(r.lv, s)
	}
	return s
}

func (f *RPCFormatter) _requestID(r *RPCRecord) string {
	s := r.requestID
	if s == "" {
		s = "-"
	}
	if f.Colored() {
		s = f.paint(r.lv, s)
	}
	return s
}

func (f *RPCFormatter) funcMap() template.FuncMap {
	return template.FuncMap{
		"date":       f._date,
		"time":       f._time,
		"datetime":   f._datetime,
		"l":          f._l,
		"level":      f._level,
		"name":       f._name,
		"pid":        f._pid,
		"file_line":  f._fileLine,
		"app_id":     f._appID,
		"rpc_id":     f._rpcID,
		"request_id": f._requestID,
	}
}
