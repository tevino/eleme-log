package handler

import (
	"io"
	"log/syslog"

	"github.com/eleme/log"
)

// SyslogHandler can send log to syslog
type SyslogHandler struct {
	log.Formatter
	w *syslog.Writer
}

// NewSyslogHandler creates a SyslogHandler with given syslog.Writer which
// could be created by syslog.New, the log format as follows:
//
//	"[{{app_id}} {{rpc_id}} {{request_id}}] ## {{}}"
func NewSyslogHandler(w *syslog.Writer) (*SyslogHandler, error) {
	f := log.NewBaseFormatter(false)
	if err := f.ParseFormat(log.TplSyslog); err != nil {
		return nil, err
	}
	return NewSyslogHandlerWithFormat(w, f), nil
}

// NewSyslogHandlerWithFormat is just like NewSyslogHandler but with customized
// format string
func NewSyslogHandlerWithFormat(w *syslog.Writer, f log.Formatter) *SyslogHandler {
	h := new(SyslogHandler)
	h.w = w
	h.Formatter = f
	return h
}

// Log prints the Record info syslog writer
func (sh *SyslogHandler) Log(r log.Record) {
	b := string(sh.Formatter.Format(r))
	switch r.Level() {
	case log.DEBUG:
		sh.w.Debug(b)
	case log.INFO:
		sh.w.Info(b)
	case log.WARN:
		sh.w.Warning(b)
	case log.ERRO:
		sh.w.Err(b)
	case log.FATA:
		sh.w.Crit(b)
	}
}

// Writer return the writer
func (sh *SyslogHandler) Writer() io.Writer {
	return sh.w
}
