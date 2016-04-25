package log

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"text/template"
	"time"
)

// LevelType identifies the level of a logger
type LevelType int

const (
	// NOTSET indicates the logger level not set
	NOTSET LevelType = iota
	// DEBUG indicates the logger level DEBUG
	DEBUG
	// INFO indicates the logger level INFO
	INFO
	// WARN indicates the logger level WARNING
	WARN
	// ERRO indicates the logger level ERROR
	ERRO
	// FATA indicates the logger level FATAL
	FATA
)

var (
	globalLevel = NOTSET
	logLevel    string
	globalAppID = ""
)

// LevelName maps LevelType to human-readable string
var LevelName = map[LevelType]string{
	DEBUG: "DEBUG",
	INFO:  "INFO",
	WARN:  "WARN",
	ERRO:  "ERRO",
	FATA:  "FATA",
}

var levelColor = map[LevelType]color{
	DEBUG: colorBlue,
	INFO:  colorGreen,
	WARN:  colorYellow,
	ERRO:  colorRed,
	FATA:  colorRed,
}

var levelFlag = map[string]LevelType{
	"debug": DEBUG,
	"info":  INFO,
	"warn":  WARN,
	"erro":  ERRO,
	"fata":  FATA,
}

// Logger is an object for logging with a set of configurations, including
// name, level, logging format, and multiple handlers
type Logger struct {
	sync.RWMutex
	wg        sync.WaitGroup
	name      string
	lv        LevelType
	tpl       *template.Template
	handlers  map[Handler]bool
	rpcID     string
	requestID string
	async     bool
}

// New creates a Logger with Stdout as default output
func New(name string) *Logger {
	return NewWithWriter(name, os.Stdout)
}

// NewWithWriter creates a Logger with given Writer as the only output
func NewWithWriter(name string, w io.Writer) *Logger {
	l := new(Logger)
	l.name = name
	l.lv = NOTSET
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

// NewRPCLogger creates a Logger with given name as a RPCLogger
func NewRPCLogger(name string) RPCLogger {
	// TODO: differentiate RPCLogger and Logger
	return New(name)
}

// SetGlobalLevel sets the global log level
func SetGlobalLevel(lv LevelType) {
	globalLevel = lv
}

// GlobalLevel returns the global log level
func GlobalLevel() LevelType {
	return globalLevel
}

// SetGlobalAppID sets the global AppID
func SetGlobalAppID(appID string) {
	globalAppID = appID
}

// AttachFlagSet attaches a flag to the given FlagSet indicating the global log level
//
// Passing nil flagSet for default FlagSet(flag.CommandLine)
func AttachFlagSet(flagSet *flag.FlagSet) {
	if flagSet == nil {
		flagSet = flag.CommandLine
	}
	flagSet.StringVar(&logLevel, "log", "info", "logs at or above this level to the logging output: debug, info, warn, fata")
}

// ParseFlag should be used after AttachFlagSet
func ParseFlag() error {
	lvl, ok := levelFlag[strings.ToLower(logLevel)]
	if ok {
		globalLevel = lvl
		return nil
	}
	return errors.New("unknown log level")
}

// Name returns the name of logger
func (l *Logger) Name() string {
	l.RLock()
	defer l.RUnlock()
	return l.name
}

// AddHandler adds the given handler to logger
func (l *Logger) AddHandler(h Handler) {
	l.Lock()
	defer l.Unlock()
	if !l.handlers[h] {
		l.handlers[h] = true
	}
}

// Handlers returns all handlers
func (l *Logger) Handlers() []Handler {
	l.RLock()
	defer l.RUnlock()
	var hs = make([]Handler, 0, len(l.handlers))
	for h := range l.handlers {
		hs = append(hs, h)
	}
	return hs
}

// RemoveHandler removes a handler
func (l *Logger) RemoveHandler(h Handler) {
	l.Lock()
	defer l.Unlock()
	delete(l.handlers, h)
}

// Level returns the current level of logger
//
// logger.SetLevel is always authoritative, GlobalLevel is used if SetLevel is
// not called, otherwise defaultLevel is used.
//
// Level() search priority:
//	1. logger's own level (if set)
//	2. GlobalLevel (if set)
//	3. defaultLevel (built-in, usually INFO)
func (l *Logger) Level() LevelType {
	l.RLock()
	defer l.RUnlock()
	if l.lv != NOTSET {
		return l.lv
	}
	if globalLevel != NOTSET {
		return globalLevel
	}
	return defaultLevel
}

// SetLevel set the level of logger
//
// SetLevel is always authoritative, See also logger.Level()
func (l *Logger) SetLevel(lv LevelType) {
	l.Lock()
	defer l.Unlock()
	l.lv = lv
}

// SetRPCID sets the RPCID for logger
func (l *Logger) SetRPCID(rpcID string) {
	l.Lock()
	defer l.Unlock()
	l.rpcID = rpcID
}

// RPCID returns the RPCID of logger
func (l *Logger) RPCID() string {
	l.RLock()
	defer l.RUnlock()
	return l.rpcID
}

// SetRequestID sets the RequestID for logger
func (l *Logger) SetRequestID(requestID string) {
	l.Lock()
	defer l.Unlock()
	l.requestID = requestID
}

// RequestID returns the request ID of logger
func (l *Logger) RequestID() string {
	l.RLock()
	defer l.RUnlock()
	return l.requestID
}

// SetAsync set output as async
func (l *Logger) SetAsync(async bool) {
	l.RLock()
	defer l.RUnlock()
	l.async = async
}

// Output writes a log to all writers with given calldepth and level
//
// Normally, you won't need this.
func (l *Logger) Output(calldepth int, lv LevelType, s string) {
	if lv < l.Level() {
		return
	}
	fileLine := ""
	_, file, line, ok := runtime.Caller(calldepth)
	if !ok {
		file = "???"
		line = 0
	}
	fileLine = file + ":" + strconv.Itoa(line)

	l.RLock()
	r := &Record{
		fileLine:  fileLine,
		name:      l.name,
		now:       time.Now(),
		lv:        lv,
		msg:       s,
		rpcID:     l.rpcID,
		requestID: l.requestID,
		appID:     globalAppID,
	}
	l.RUnlock()

	l.Lock()
	defer l.Unlock()
	if l.async {
		for h := range l.handlers {
			go h.Log(r)
		}
		return
	}

	var wg sync.WaitGroup
	for h := range l.handlers {
		wg.Add(1)
		go func(h Handler, r *Record) {
			defer wg.Done()
			h.Log(r)
		}(h, r)
	}
	wg.Wait()
}

// Debug APIs

// Debug calls Output to log with DEBUG level
func (l *Logger) Debug(a ...interface{}) {
	l.Output(2, DEBUG, fmt.Sprint(a...))
}

// Debugf calls Output to log with DEBUG level and given format
func (l *Logger) Debugf(format string, a ...interface{}) {
	l.Output(2, DEBUG, fmt.Sprintf(format, a...))
}

// Print APIs

// Print calls Output to log with default level
func (l *Logger) Print(a ...interface{}) {
	l.Output(2, l.Level(), fmt.Sprint(a...))
}

// Println calls Output to log with default level
func (l *Logger) Println(a ...interface{}) {
	l.Output(2, l.Level(), fmt.Sprint(a...))
}

// Printf calls Output to log with default level and given format
func (l *Logger) Printf(f string, a ...interface{}) {
	l.Output(2, l.Level(), fmt.Sprintf(f, a...))
}

// Info APIs

// Info calls Output to log with INFO level
func (l *Logger) Info(a ...interface{}) {
	l.Output(2, INFO, fmt.Sprint(a...))
}

// Infof calls Output to log with INFO level and given format
func (l *Logger) Infof(f string, a ...interface{}) {
	l.Output(2, INFO, fmt.Sprintf(f, a...))
}

// Warn APIs

// Warn calls Output to log with WARN level
func (l *Logger) Warn(a ...interface{}) {
	l.Output(2, WARN, fmt.Sprint(a...))
}

// Warnf calls Output to log with WARN level and given format
func (l *Logger) Warnf(f string, a ...interface{}) {
	l.Output(2, WARN, fmt.Sprintf(f, a...))
}

// Error APIs

// Error calls Output to log with ERRO level
func (l *Logger) Error(a ...interface{}) {
	l.Output(2, ERRO, fmt.Sprint(a...))
}

// Errorf calls Output to log with ERRO level and given format
func (l *Logger) Errorf(f string, a ...interface{}) {
	l.Output(2, ERRO, fmt.Sprintf(f, a...))
}

// Fatal APIs

// Fatal calls Output to log with FATA level followed by a call to os.Exit(1)
func (l *Logger) Fatal(a ...interface{}) {
	l.Output(2, FATA, fmt.Sprint(a...))
	os.Exit(1)
}

// Fatalf calls Output to log with FATA level with given format, followed by a call to os.Exit(1)
func (l *Logger) Fatalf(f string, a ...interface{}) {
	l.Output(2, FATA, fmt.Sprintf(f, a...))
	os.Exit(1)
}
