package log

import "io"

var writerLocks *writerLocker
var wSupervisor *writerSupervisor

func init() {
	writerLocks = newWriterLocker()
	wSupervisor = newWriterSupervisor()
}

// Handler represents a handler of Record
type Handler interface {
	Log(r *Record)
	AsyncLog(r *Record)
}

// StreamHandler is a Handler of Stream writer e.g. console
type StreamHandler struct {
	writer io.Writer
	*Formatter
}

// NewStreamHandler creates a StreamHandler with given writer(usually os.Stdout)
// and format string, whether to color the output is determined by the type of
// writer
func NewStreamHandler(w io.Writer, f string) (*StreamHandler, error) {
	h := new(StreamHandler)
	h.writer = w

	formatter, err := NewFormatter(f, IsTerminal(w))
	h.Formatter = formatter

	return h, err
}

// Colored enable or disable the color function of internal format, usually
// this is determined automatically
//
// When called with no argument, it returns the current state of color function
func (sw *StreamHandler) Colored(ok ...bool) bool {
	if len(ok) > 0 {
		sw.Formatter.colored = ok[0]
	}
	return sw.Formatter.colored
}

// Log print the Record to the internal writer
func (sw *StreamHandler) Log(r *Record) {
	b := sw.Formatter.Format(r)
	writerLocks.Lock(sw.writer)
	defer writerLocks.Unlock(sw.writer)
	sw.writer.Write(b)
}

// AsyncLog print the Record by worker
func (sw *StreamHandler) AsyncLog(r *Record) {
	wSupervisor.Write(sw.writer, func() {
		b := sw.Formatter.Format(r)
		writerLocks.Lock(sw.writer)
		defer writerLocks.Unlock(sw.writer)
		sw.writer.Write([]byte(b))
	})
}
