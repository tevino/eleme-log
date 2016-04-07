package log

import "log/syslog"

type SyslogHandler struct {
	*Formatter
	w *syslog.Writer
}

func NewSyslogHandler(w *syslog.Writer) (*SyslogHandler, error) {
	return NewSyslogHandlerWithFormat(w, syslogTpl)
}

func NewSyslogHandlerWithFormat(w *syslog.Writer, f string) (*SyslogHandler, error) {
	h := new(SyslogHandler)
	h.w = w
	formatter, err := NewFormatter(f, false)
	h.Formatter = formatter
	return h, err
}

func (sh *SyslogHandler) Log(r *Record) {
	b := sh.Formatter.Format(r)
	switch r.lv {
	case DEBUG:
		sh.w.Debug(b)
	case INFO:
		sh.w.Info(b)
	case WARN:
		sh.w.Warning(b)
	case FATA:
		sh.w.Err(b)
	}
}
