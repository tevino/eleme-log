package rpc

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"sync"
	"testing"

	"github.com/eleme/log"
)

func newELogger(t *testing.T, w io.Writer, f string) log.RPCLogger {
	l := log.NewWithWriter("elog", nil)
	elog := &ELogger{Logger: l}
	elog.recordFactory = NewELogRecordFactory(elog.rpcID, elog.requestID)

	ef := NewELogFormatter(false)
	if err := ef.ParseFormat(f); err != nil {
		t.Error("error creating stream handler: ", err)
		t.FailNow()
	}
	h := log.NewStreamHandler(w, ef)
	h.Colored(false)
	elog.AddHandler(h)
	return elog
}

func TestSetRPCID(t *testing.T) {
	var buf bytes.Buffer
	elog := newELogger(t, &buf, "[{{rpc_id}}] ## {{}}")

	expectedNil := "[-] ## InfoLog\n"
	elog.Info("InfoLog")
	if buf.String() != expectedNil {
		t.Errorf("Expected:\n%s\nGot:\n%s", expectedNil, buf.String())
	}

	buf.Reset()
	expected := "[test.rpcid] ## InfoLog\n"
	elog.WithRPCID("test.rpcid").Info("InfoLog")
	if buf.String() != expected {
		t.Errorf("Expected:\n%s\nGot:\n%s", expected, buf.String())
	}
}

func TestSetRequestID(t *testing.T) {
	var buf bytes.Buffer
	elog := newELogger(t, &buf, "[{{request_id}}] ## {{}}")

	expectedNil := "[-] ## InfoLog\n"
	elog.Info("InfoLog")
	if buf.String() != expectedNil {
		t.Errorf("Expected:\n%s\nGot:\n%s", expectedNil, buf.String())
	}

	buf.Reset()
	expected := "[test.request_id] ## InfoLog\n"
	elog.WithRequestID("test.request_id").Info("InfoLog")
	if buf.String() != expected {
		t.Errorf("Expected:\n%s\nGot:\n%s", expected, buf.String())
	}
}

func TestMultiRPCLog(t *testing.T) {
	var buf bytes.Buffer
	elog := newELogger(t, &buf, "{{rpc_id}} {{}}")

	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func(flag int) {
			defer wg.Done()
			elog.WithRPCID(fmt.Sprintf("%d", flag)).Info(fmt.Sprintf("%d", flag))
		}(i)
	}
	wg.Wait()

	for {
		line, err := buf.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			t.Error(err)
			break
		}
		strs := strings.Split(line, " ")
		if strs[0] != strings.TrimSpace(strs[1]) {
			t.Errorf("rpcID: %#v, info: %#v", strs[0], strings.TrimSpace(strs[1]))
		}
	}
}
