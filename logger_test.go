package log

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNamedLeveler(t *testing.T) {
	ast := assert.New(t)
	l := New("tester")
	ast.Equal(l.Name(), "tester")
	ast.Equal(l.Level(), defaultLevel)

	l.SetLevel(WARN)
	ast.Equal(l.Level(), WARN)
}

func TestMultiHandler(t *testing.T) {
	ast := assert.New(t)
	l := New("tester")
	ast.Equal(len(l.Handlers()), 1)

	hdr := l.Handlers()[0]
	h, ok := hdr.(*StreamHandler)
	ast.True(ok)

	l.RemoveHandler(h)
	ast.Equal(len(l.Handlers()), 0)

	sHdr, _ := NewStreamHandler(os.Stdin, "{{level}} {{date}} {{time}} {{name}} {{file_line}} {{}}")
	l.AddHandler(sHdr)

	ast.Equal(len(l.Handlers()), 1)
	hdr = l.Handlers()[0]
	_, ok = hdr.(*StreamHandler)
	ast.True(ok)
}
