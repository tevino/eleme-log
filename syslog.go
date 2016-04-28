package log

import "log/syslog"

// SyslogHandler can send log to syslog
type SyslogHandler struct {
	Formatter
	w *syslog.Writer
}

// NewSyslogHandler creates a SyslogHandler with given syslog.Writer which
// could be created by syslog.New, the log format as follows:
//
//	"[{{app_id}} {{rpc_id}} {{request_id}}] ## {{}}"
func NewSyslogHandler(w *syslog.Writer) (*SyslogHandler, error) {
	f := NewBaseFormatter(false)
	if err := f.ParseFormat(syslogTpl); err != nil {
		return nil, err
	}
	return NewSyslogHandlerWithFormat(w, f), nil
}

// NewSyslogHandlerWithFormat is just like NewSyslogHandler but with customized
// format string
func NewSyslogHandlerWithFormat(w *syslog.Writer, f Formatter) *SyslogHandler {
	h := new(SyslogHandler)
	h.w = w
	h.Formatter = f
	return h
}

// Log prints the Record info syslog writer
func (sh *SyslogHandler) Log(r Record) {
	b := string(sh.Formatter.Format(r))
	switch r.Level() {
	case DEBUG:
		sh.w.Debug(b)
	case INFO:
		sh.w.Info(b)
	case WARN:
		sh.w.Warning(b)
	case ERRO:
		sh.w.Err(b)
	case FATA:
		sh.w.Crit(b)
	}
}
