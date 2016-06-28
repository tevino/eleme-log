package log

import (
	"io/ioutil"
	"testing"
)

func BenchmarkWriteWorker1024(b *testing.B) {
	size := 1024
	for i := 0; i < b.N; i++ {
		ws := newWriterSupervisor(size)
		ws.Do(ioutil.Discard, func() {})
	}
}

func BenchmarkWriteWorker102400(b *testing.B) {
	size := 1024 * 100
	for i := 0; i < b.N; i++ {
		ws := newWriterSupervisor(size)
		ws.Do(ioutil.Discard, func() {})
	}
}
