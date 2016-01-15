package log

import (
	"flag"
	"fmt"
	"io"
	"os"
	"sync"
	"text/template"
	"time"
)

type LevelType int

const (
	NOTSET LevelType = iota
	DEBUG
	INFO
	WARN
	FATA
)

var (
	globalLevel    = NOTSET
	logLevel       string
	defaultFlagSet = flag.CommandLine
)

var LevelName = map[LevelType]string{
	DEBUG: "DEBUG",
	INFO:  "INFO",
	WARN:  "WARN",
	FATA:  "FATA",
}

var LevelColor = map[LevelType]Color{
	DEBUG: Blue,
	INFO:  Green,
	WARN:  Yellow,
	FATA:  Red,
}

var levelFlag = map[string]LevelType{
	"debug": DEBUG,
	"info":  INFO,
	"warn":  WARN,
	"fata":  FATA,
}

type logger struct {
	sync.Mutex
	wg       sync.WaitGroup
	name     string
	lv       LevelType
	tpl      *template.Template
	handlers map[Handler]bool
	appID    string
}

func New(name string) Logger {
	return NewWithWriter(name, os.Stdout)
}

func NewWithWriter(name string, w io.Writer) Logger {
	l := new(logger)
	l.name = name
	l.lv = defaultLevel
	l.handlers = make(map[Handler]bool)
	if w != nil {
		hdr, err := NewStreamHandler(w, defaultTpl)
		if err != nil {
			panic(err)
		}
		l.AddHandler(hdr)
	}
	return l
}

func SetGlobalLevel(lv LevelType) {
	globalLevel = lv
}

func GlobalLevel() LevelType {
	return globalLevel
}

// AttachFlagSet set some flag, if flagSet is nil, will use flag.CommandLine
func AttachFlagSet(flagSet *flag.FlagSet) {
	if flagSet != nil {
		defaultFlagSet = flagSet
	}
	defaultFlagSet.StringVar(&logLevel, "logLevel", "info", "logs at or above this level to the logging output: debug, info, warn, fata")
}

func ParseFlag() bool {
	lvl, ok := levelFlag[logLevel]
	if ok {
		SetGlobalLevel(lvl)
		return true
	}
	return false
}

func (l *logger) AddHandler(h Handler) {
	l.Lock()
	defer l.Unlock()
	if !l.handlers[h] {
		l.handlers[h] = true
	}
}

func (l *logger) RemoveHandler(h Handler) {
	l.Lock()
	defer l.Unlock()
	delete(l.handlers, h)
}

func (l *logger) Level() LevelType {
	if globalLevel == NOTSET {
		return l.lv
	}
	return globalLevel
}

func (l *logger) SetLevel(lv LevelType) {
	l.lv = lv
}

func (l *logger) Output(lv LevelType, s string) {
	if lv < l.Level() {
		return
	}
	r := &Record{
		name:  l.name,
		now:   time.Now(),
		lv:    lv,
		msg:   s,
		appID: l.appID,
	}
	var wg sync.WaitGroup
	l.Lock()
	for h := range l.handlers {
		wg.Add(1)
		go func(h Handler, r *Record) {
			defer wg.Done()
			h.Log(r)
		}(h, r)
	}
	l.Unlock()
	wg.Wait()
}

// Debug APIs
func (l *logger) Debug(a ...interface{}) {
	l.Output(DEBUG, fmt.Sprint(a...))
}

func (l *logger) Debugf(format string, a ...interface{}) {
	l.Output(DEBUG, fmt.Sprintf(format, a...))
}

// Print APIs output logs with default level
func (l *logger) Print(a ...interface{}) {
	l.Output(l.Level(), fmt.Sprint(a...))
}

func (l *logger) Println(a ...interface{}) {
	l.Print(a...)
}

func (l *logger) Printf(f string, a ...interface{}) {
	l.Output(l.Level(), fmt.Sprintf(f, a...))
}

// Info APIs
func (l *logger) Info(a ...interface{}) {
	l.Output(INFO, fmt.Sprint(a...))
}

func (l *logger) Infof(f string, a ...interface{}) {
	l.Output(INFO, fmt.Sprintf(f, a...))
}

// Warn APIs
func (l *logger) Warn(a ...interface{}) {
	l.Output(WARN, fmt.Sprint(a...))
}

func (l *logger) Warnf(f string, a ...interface{}) {
	l.Output(WARN, fmt.Sprintf(f, a...))
}

// Fatal APIs
func (l *logger) Fatal(a ...interface{}) {
	l.Output(FATA, fmt.Sprint(a...))
	os.Exit(1)
}

func (l *logger) Fatalf(f string, a ...interface{}) {
	l.Output(FATA, fmt.Sprintf(f, a...))
	os.Exit(1)
}

func (l *logger) SetAppID(appID string) {
	l.appID = appID
}
