package log

import (
	"io"
	"sync"
)

type WriterLocks struct {
	m  map[io.Writer]*sync.Mutex
	mu sync.RWMutex
}

func (wl *WriterLocks) Lock(w io.Writer) {
	wl.mu.RLock()
	if l, ok := wl.m[w]; ok {
		l.Lock()
		wl.mu.RUnlock()
	} else {
		wl.mu.RUnlock()
		// add new lock to map
		var newLock sync.Mutex
		wl.mu.Lock()
		wl.m[w] = &newLock
		wl.mu.Unlock()
		// lock it
		newLock.Lock()
	}
}

func (wl *WriterLocks) Unlock(w io.Writer) {
	wl.mu.RLock()
	if l, ok := wl.m[w]; ok {
		l.Unlock()
	}
	defer wl.mu.RUnlock()
}
