package log

import (
	"io"
	"sync"
)

var writerLocks WriterLocks

func init() {
	writerLocks = WriterLocks{m: make(map[io.Writer]*sync.Mutex)}
}

type Handler interface {
	Log(r *Record)
}

type StreamHandler struct {
	writer    io.Writer
	formatter *Formatter
}

func NewStreamHandler(w io.Writer, f string) (*StreamHandler, error) {
	h := new(StreamHandler)
	h.writer = w

	formatter, err := NewFormatter(f, IsTerminal(w))
	h.formatter = formatter

	return h, err
}

func (sw *StreamHandler) Colored(ok ...bool) bool {
	if len(ok) > 0 {
		sw.formatter.colored = ok[0]
	}
	return sw.formatter.colored
}

func (sw *StreamHandler) Log(r *Record) {
	b := sw.formatter.Format(r)
	writerLocks.Lock(sw.writer)
	defer writerLocks.Unlock(sw.writer)
	sw.writer.Write([]byte(b))
}
