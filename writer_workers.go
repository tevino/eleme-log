package log

import (
	"io"
	"sync"
)

const (
	maxRecordChanSize = 5000
)

type writerWorker struct {
	w  io.Writer
	ch chan func()
}

type writerSupervisor struct {
	m  map[io.Writer]*writerWorker
	mu sync.RWMutex
}

func (ws *writerSupervisor) Write(w io.Writer, f func()) {
	ws.mu.RLock()
	worker, ok := ws.m[w]
	ws.mu.RUnlock()

	if !ok {
		worker = &writerWorker{
			w:  w,
			ch: make(chan func(), maxRecordChanSize),
		}
		go func(w io.Writer, ch chan func()) {
			for ff := range ch {
				ff()
			}
		}(worker.w, worker.ch)

		ws.mu.Lock()
		ws.m[w] = worker
		ws.mu.Unlock()
	}

	select {
	case worker.ch <- f:
	default:
		go f()
	}
}

func newWriterSupervisor() *writerSupervisor {
	return &writerSupervisor{
		m:  make(map[io.Writer]*writerWorker),
		mu: sync.RWMutex{},
	}
}
