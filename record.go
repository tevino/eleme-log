package log

import (
	"fmt"
	"runtime"
	"strconv"
	"sync"
	"time"
)

// Record stands for a single record of log, usually a single line
type Record struct {
	logger *Logger

	fileLine  string
	name      string
	now       time.Time
	lv        LevelType
	msg       string
	rpcID     string
	requestID string
	appID     string
}

// String returns the raw message of the Record
func (r *Record) String() string {
	return r.msg
}

// NewRecord creates a Record with given Logger and calldepth.
func NewRecord(logger *Logger, calldepth int) *Record {
	fileLine := ""
	_, file, line, ok := runtime.Caller(calldepth)
	if !ok {
		file = "???"
		line = 0
	}
	fileLine = file + ":" + strconv.Itoa(line)

	r := &Record{
		logger:   logger,
		fileLine: fileLine,
		name:     logger.name,
		now:      time.Now(),
		appID:    globalAppID,
	}
	return r
}

// output writes a log to all writers
func (r Record) output() {
	if r.lv < r.logger.Level() {
		return
	}

	var wg sync.WaitGroup
	r.logger.RLock()
	defer r.logger.RUnlock()
	for h := range r.logger.handlers {
		wg.Add(1)
		go func(h Handler, r *Record) {
			defer wg.Done()
			h.Log(r)
		}(h, &r)
	}
	wg.Wait()
}

// Debug APIs

// Debug calls Output to log with DEBUG level
func (r *Record) Debug(a ...interface{}) {
	r.lv = DEBUG
	r.msg = fmt.Sprint(a...)
	r.output()
}

// Debugf calls Output to log with DEBUG level and given format
func (r *Record) Debugf(format string, a ...interface{}) {
	r.lv = DEBUG
	r.msg = fmt.Sprintf(format, a...)
	r.output()
}

// Print APIs

// Print calls Output to log with default level
func (r *Record) Print(a ...interface{}) {
	r.lv = r.logger.Level()
	r.msg = fmt.Sprint(a...)
	r.output()
}

// Println calls Output to log with default level
func (r *Record) Println(a ...interface{}) {
	r.lv = r.logger.Level()
	r.msg = fmt.Sprint(a...)
	r.output()
}

// Printf calls Output to log with default level and given format
func (r *Record) Printf(f string, a ...interface{}) {
	r.lv = r.logger.Level()
	r.msg = fmt.Sprintf(f, a...)
	r.output()
}

// Info APIs

// Info calls Output to log with INFO level
func (r *Record) Info(a ...interface{}) {
	r.lv = INFO
	r.msg = fmt.Sprint(a...)
	r.output()
}

// Infof calls Output to log with INFO level and given format
func (r *Record) Infof(f string, a ...interface{}) {
	r.lv = INFO
	r.msg = fmt.Sprintf(f, a...)
	r.output()
}

// Warn APIs

// Warn calls Output to log with WARN level
func (r *Record) Warn(a ...interface{}) {
	r.lv = WARN
	r.msg = fmt.Sprint(a...)
	r.output()
}

// Warnf calls Output to log with WARN level and given format
func (r *Record) Warnf(f string, a ...interface{}) {
	r.lv = WARN
	r.msg = fmt.Sprintf(f, a...)
	r.output()
}

// Error APIs

// Error calls Output to log with ERRO level
func (r *Record) Error(a ...interface{}) {
	r.lv = ERRO
	r.msg = fmt.Sprint(a...)
	r.output()
}

// Errorf calls Output to log with ERRO level and given format
func (r *Record) Errorf(f string, a ...interface{}) {
	r.lv = ERRO
	r.msg = fmt.Sprintf(f, a...)
	r.output()
}

// Fatal APIs

// Fatal calls Output to log with FATA level followed by a call to os.Exit(1)
func (r *Record) Fatal(a ...interface{}) {
	r.lv = FATA
	r.msg = fmt.Sprint(a...)
	r.output()
}

// Fatalf calls Output to log with FATA level with given format, followed by a call to os.Exit(1)
func (r *Record) Fatalf(f string, a ...interface{}) {
	r.lv = FATA
	r.msg = fmt.Sprintf(f, a...)
	r.output()
}
