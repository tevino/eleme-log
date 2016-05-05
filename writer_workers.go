package log

import (
	"io"
	"sync"
)

const (
	maxRecordChanSize = 100000
)

type writerWorker struct {
	ch      chan func()
	closing chan bool
	closed  bool
	l       sync.RWMutex
}

func (w *writerWorker) Push(f func()) {
	w.l.RLock()
	if w.closed {
		w.l.RUnlock()
		return
	}
	select {
	case w.ch <- f:
	default:
		//throw message if full
	}
	w.l.RUnlock()

}

func (w *writerWorker) Start() {
	go func() {
		for ff := range w.ch {
			ff()
		}
		w.closing <- true
	}()
}

func (w *writerWorker) WaitClose() {
	w.l.Lock()
	w.closed = true
	w.l.Unlock()

	close(w.ch)
	<-w.closing
}

type writerSupervisor struct {
	m      map[io.Writer]*writerWorker
	mu     sync.RWMutex
	closed bool
}

func (ws *writerSupervisor) WaitClose() {
	ws.mu.RLock()
	if ws.closed {
		ws.mu.RUnlock()
		return
	}
	ws.mu.RUnlock()

	ws.mu.Lock()
	ws.closed = true
	ws.mu.Unlock()

	ws.mu.RLock()
	defer ws.mu.RUnlock()
	for _, worker := range ws.m {
		worker.WaitClose()
	}
}

func (ws *writerSupervisor) Do(w io.Writer, f func()) {
	ws.mu.RLock()
	if ws.closed {
		ws.mu.RUnlock()
		return
	}
	worker, ok := ws.m[w]
	ws.mu.RUnlock()

	if !ok {
		worker = &writerWorker{
			ch:      make(chan func(), maxRecordChanSize),
			closing: make(chan bool),
			closed:  false,
			l:       sync.RWMutex{},
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

	worker.Push(f)
}

func newWriterSupervisor() *writerSupervisor {
	return &writerSupervisor{
		m:      make(map[io.Writer]*writerWorker),
		mu:     sync.RWMutex{},
		closed: false,
	}
}
