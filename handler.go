package log

import (
	"io"
	"sync"
)

var writerLocks WriterLocks

func init() {
	writerLocks = WriterLocks{m: make(map[io.Writer]*sync.Mutex)}
}

// Handler represents a handler of Record
type Handler interface {
	Log(r *Record)
}

// StreamHandler is a Handler of Stream writer e.g. console
type StreamHandler struct {
	writer    io.Writer
	formatter *Formatter
}

// NewStreamHandler creates a StreamHandler with given writer(usually os.Stdout)
// and format string, whether to color the output is determined by the type of
// writer
func NewStreamHandler(w io.Writer, f string) (*StreamHandler, error) {
	h := new(StreamHandler)
	h.writer = w

	formatter, err := NewFormatter(f, IsTerminal(w))
	h.formatter = formatter

	return h, err
}

// Colored enable or disable the color function of internal format, usually
// this is determined automatically
//
// When called with no argument, it returns the current state of color function
func (sw *StreamHandler) Colored(ok ...bool) bool {
	if len(ok) > 0 {
		sw.formatter.colored = ok[0]
	}
	return sw.formatter.colored
}

// Log print the Record to the internal writer
func (sw *StreamHandler) Log(r *Record) {
	b := sw.formatter.Format(r)
	writerLocks.Lock(sw.writer)
	defer writerLocks.Unlock(sw.writer)
	sw.writer.Write([]byte(b))
}
