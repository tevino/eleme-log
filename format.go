package log

import (
	"bytes"
	"os"
	"regexp"
	"strconv"
	"strings"
	"text/template"
)

// Formatter describes the format of outputting log
type Formatter struct {
	colored bool
	tpl     *template.Template
}

// NewFormatter creates a Formatter with given format string and whether to
// color the output
func NewFormatter(format string, colored bool) (*Formatter, error) {
	fm := new(Formatter)
	fm.colored = colored
	if err := fm.SetFormat(format); err != nil {
		return nil, err
	}
	return fm, nil
}

var rTagLong = regexp.MustCompile("{{ *([a-zA-Z]+) *}}")
var tagShort = []byte("{{$1}}")
var tagReplacer = strings.NewReplacer(
	"{{}}", "{{.String}}",
	"{{level}}", "{{level .}}",
	"{{l}}", "{{l .}}",
	"{{date}}", "{{date .}}",
	"{{time}}", "{{time .}}",
	"{{datetime}}", "{{datetime .}}",
	"{{name}}", "{{name .}}",
	"{{pid}}", "{{pid .}}",
	"{{file_line}}", "{{file_line .}}",

	"{{rpc_id}}", "{{rpc_id .}}",
	"{{request_id}}", "{{request_id .}}",
	"{{app_id}}", "{{app_id .}}",
)

// SetFormat set the format of outputting log
//
// The default format is "{{ level }} {{ date }} {{ time }} {{ name }} {{}}"
//
// {{this is a placeholder}} which will be replaced by the actual content
//
// Available placeholders:
//	{{}}            The message provided by you e.g. l.Info(message)
//	{{ level }}     Log level in four UPPER-CASED letters e.g. INFO, WARN
//	{{ l }}         Log level in one UPPER-CASED letter e.g. I, W
//	{{ data }}      Date in format "2006-01-02"
//	{{ time }}      Time in format "15:04:05")
//	{{ datetime }}  Date and time in format "2006-01-02 15:04:05.999"
//	{{ name }}      Logger name
//	{{ pid }}       Current process ID
//	{{ file_line }} Filename and line number in format "file.go:12"
func (f *Formatter) SetFormat(tpl string) error {
	// {{ tag }} -> {{tag}}
	tpl = string(rTagLong.ReplaceAll([]byte(tpl), tagShort))

	tpl = tagReplacer.Replace(tpl)

	if !strings.HasSuffix(tpl, "\n") {
		tpl += "\n"
	}

	t, err := template.New("tpl").Funcs(f.funcMap()).Parse(tpl)
	if err != nil {
		return err
	}

	// TODO: validation

	f.tpl = t
	return nil
}

// Format formats a Record with set format
func (f *Formatter) Format(r *Record) string {
	var buf bytes.Buffer
	f.tpl.Execute(&buf, r)
	return buf.String()
}

// TODO: the 'if color then paint' is ugly!!

func (f *Formatter) _level(r *Record) string {
	s := LevelName[r.lv]
	if f.colored {
		s = f.paint(r.lv, s)
	}
	return s
}

func (f *Formatter) _l(r *Record) string {
	s := LevelName[r.lv][0:1]
	if f.colored {
		s = f.paint(r.lv, s)
	}
	return s
}

func (f *Formatter) _datetime(r *Record) string {
	s := r.now.Format("2006-01-02 15:04:05.999")
	if f.colored {
		s = f.paint(r.lv, s)
	}
	return s
}

func (f *Formatter) _date(r *Record) string {
	s := r.now.Format("2006-01-02")
	if f.colored {
		s = f.paint(r.lv, s)
	}
	return s
}

func (f *Formatter) _time(r *Record) string {
	s := r.now.Format("15:04:05")
	if f.colored {
		s = f.paint(r.lv, s)
	}
	return s
}

func (f *Formatter) _name(r *Record) string {
	s := r.name
	if f.colored {
		s = f.paint(r.lv, s)
	}
	return s
}

func (f *Formatter) _pid(r *Record) string {
	s := strconv.Itoa(os.Getpid())
	if f.colored {
		s = f.paint(r.lv, s)
	}
	return s
}

func (f *Formatter) _rpcID(r *Record) string {
	s := r.rpcID
	if s == "" {
		s = "-"
	}
	if f.colored {
		s = f.paint(r.lv, s)
	}
	return s
}

func (f *Formatter) _requestID(r *Record) string {
	s := r.requestID
	if s == "" {
		s = "-"
	}
	if f.colored {
		s = f.paint(r.lv, s)
	}
	return s
}

func (f *Formatter) _appID(r *Record) string {
	s := r.appID
	if s == "" {
		s = "-"
	}
	if f.colored {
		s = f.paint(r.lv, s)
	}
	return s
}

func (f *Formatter) _fileLine(r *Record) string {
	s := r.fileLine
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == '/' {
			s = s[i+1:]
			break
		}
	}
	if f.colored {
		s = f.paint(r.lv, s)
	}
	return s
}

func (f *Formatter) funcMap() template.FuncMap {
	return template.FuncMap{
		"date":      f._date,
		"time":      f._time,
		"datetime":  f._datetime,
		"l":         f._l,
		"level":     f._level,
		"name":      f._name,
		"pid":       f._pid,
		"file_line": f._fileLine,

		"rpc_id":     f._rpcID,
		"request_id": f._requestID,
		"app_id":     f._appID,
	}
}

func (f *Formatter) paint(lv LevelType, s string) string {
	return painter(levelColor[lv], s)
}
