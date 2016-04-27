package log

import (
	"io"
	"sync"
)

const (
	maxRecordChanSize = 100000
)

type writerWorker struct {
	ch chan func()
}

func (w *writerWorker) Start() {
	go func() {
		for ff := range w.ch {
			ff()
		}
	}()
}

type writerSupervisor struct {
	m  map[io.Writer]*writerWorker
	mu sync.RWMutex
}

func (ws *writerSupervisor) Do(w io.Writer, f func()) {
	ws.mu.RLock()
	worker, ok := ws.m[w]
	ws.mu.RUnlock()

	if !ok {
		worker = &writerWorker{
			ch: make(chan func(), maxRecordChanSize),
		}

		ws.mu.Lock()
		if currentWorker, ok := ws.m[w]; ok {
			worker = currentWorker
		} else {
			ws.m[w] = worker
			worker.Start()
		}
		ws.mu.Unlock()
	}

	select {
	case worker.ch <- f:
	default:
		//throw message if full
	}
}

func newWriterSupervisor() *writerSupervisor {
	return &writerSupervisor{
		m:  make(map[io.Writer]*writerWorker),
		mu: sync.RWMutex{},
	}
}
