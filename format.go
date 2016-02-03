package log

import (
	"bytes"
	"os"
	"regexp"
	"strconv"
	"strings"
	"text/template"
)

type Formatter struct {
	colored bool
	tpl     *template.Template
}

func NewFormatter(format string, colored bool) (*Formatter, error) {
	fm := new(Formatter)
	fm.colored = colored
	if err := fm.SetFormat(format); err != nil {
		return nil, err
	}
	return fm, nil
}

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

func (f *Formatter) Format(r *Record) string {
	var buf bytes.Buffer
	f.tpl.Execute(&buf, r)
	return buf.String()
}

// TODO: the 'if color then paint' is ugly!!

func (f *Formatter) LevelType(r *Record) string {
	s := LevelName[r.lv]
	if f.colored {
		s = f.paint(r.lv, s)
	}
	return s
}

func (f *Formatter) l(r *Record) string {
	s := LevelName[r.lv][0:1]
	if f.colored {
		s = f.paint(r.lv, s)
	}
	return s
}

func (f *Formatter) datetime(r *Record) string {
	s := r.now.Format("2006-01-02 15:04:05.999")
	if f.colored {
		s = f.paint(r.lv, s)
	}
	return s
}

func (f *Formatter) date(r *Record) string {
	s := r.now.Format("2006-01-02")
	if f.colored {
		s = f.paint(r.lv, s)
	}
	return s
}

func (f *Formatter) time(r *Record) string {
	s := r.now.Format("15:04:05")
	if f.colored {
		s = f.paint(r.lv, s)
	}
	return s
}

func (f *Formatter) name(r *Record) string {
	s := r.name
	if f.colored {
		s = f.paint(r.lv, s)
	}
	return s
}

func (f *Formatter) pid(r *Record) string {
	s := strconv.Itoa(os.Getpid())
	if f.colored {
		s = f.paint(r.lv, s)
	}
	return s
}

func (f *Formatter) rpcID(r *Record) string {
	s := r.rpcID
	if s == "" {
		s = "-"
	}
	if f.colored {
		s = f.paint(r.lv, s)
	}
	return s
}

func (f *Formatter) requestID(r *Record) string {
	s := r.requestID
	if s == "" {
		s = "-"
	}
	if f.colored {
		s = f.paint(r.lv, s)
	}
	return s
}

func (f *Formatter) appID(r *Record) string {
	s := r.appID
	if s == "" {
		s = "-"
	}
	if f.colored {
		s = f.paint(r.lv, s)
	}
	return s
}

func (f *Formatter) fileLine(r *Record) string {
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
		"date":       f.date,
		"time":       f.time,
		"l":          f.l,
		"level":      f.LevelType,
		"name":       f.name,
		"pid":        f.pid,
		"rpc_id":     f.rpcID,
		"request_id": f.requestID,
		"app_id":     f.appID,
		"datetime":   f.datetime,
		"file_line":  f.fileLine,
	}
}

var rTagLong = regexp.MustCompile("{{ *([a-zA-Z]+) *}}")
var tagShort = []byte("{{$1}}")
var tagReplacer = strings.NewReplacer(
	"{{}}", "{{.String}}",
	"{{level}}", "{{level .}}",
	"{{l}}", "{{l .}}",
	"{{date}}", "{{date .}}",
	"{{time}}", "{{time .}}",
	"{{name}}", "{{name .}}",
	"{{pid}}", "{{pid .}}",
	"{{rpc_id}}", "{{rpc_id .}}",
	"{{request_id}}", "{{request_id .}}",
	"{{app_id}}", "{{app_id .}}",
	"{{datetime}}", "{{datetime .}}",
	"{{file_line}}", "{{file_line .}}",
)

func (f *Formatter) paint(lv LevelType, s string) string {
	return painter(LevelColor[lv], s)
}
