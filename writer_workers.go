package log

import (
	"io"
	"sync"

	"github.com/tevino/abool"
)

const (
	maxRecordChanSize = 100000
)

type writerWorker struct {
	ch      chan func()
	closing chan bool
	closed  *abool.AtomicBool
}

func (w *writerWorker) Push(f func()) {
	if w.closed.IsSet() {
		return
	}
	w.l.RUnlock()
	select {
	case w.ch <- f:
	default:
		//throw message if full
	}
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
	w.closed.Set()
	close(w.ch)
	<-w.closing
}

type writerSupervisor struct {
	m  map[io.Writer]*writerWorker
	mu sync.RWMutex
	closed
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
			closed:  abool.New(),
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
