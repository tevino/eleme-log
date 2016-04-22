package log

import (
	"bytes"
	"io"
	"testing"
)

func newRPCLogger(t *testing.T, w io.Writer, f string) RPCLogger {
	l := &RPC{
		Logger: NewWithWriter("test", nil),
	}
	h, err := NewStreamHandler(w, f)
	if err != nil {
		t.Error("error creating stream handler: ", err)
		t.FailNow()
	}
	h.Colored(false)
	l.AddHandler(h)
	return l
}

func TestSetRPCID(t *testing.T) {
	var buf bytes.Buffer
	l := newRPCLogger(t, &buf, "[{{rpc_id}}] ## {{}}")
	rpcLog := l.(RPCLogger)

	expectedNil := "[-] ## InfoLog\n"
	rpcLog.Info("InfoLog")
	if buf.String() != expectedNil {
		t.Errorf("Expected:\n%s\nGot:\n%s", expectedNil, buf.String())
	}

	buf.Reset()
	expected := "[test.rpcid] ## InfoLog\n"
	rpcLog.WithRPCID("test.rpcid").Info("InfoLog")
	if buf.String() != expected {
		t.Errorf("Expected:\n%s\nGot:\n%s", expected, buf.String())
	}
}

func TestSetRequestID(t *testing.T) {
	var buf bytes.Buffer
	l := newRPCLogger(t, &buf, "[{{request_id}}] ## {{}}")
	rpcLog := l.(RPCLogger)

	expectedNil := "[-] ## InfoLog\n"
	rpcLog.Info("InfoLog")
	if buf.String() != expectedNil {
		t.Errorf("Expected:\n%s\nGot:\n%s", expectedNil, buf.String())
	}

	buf.Reset()
	expected := "[test.request_id] ## InfoLog\n"
	rpcLog.WithRequestID("test.request_id").Info("InfoLog")
	if buf.String() != expected {
		t.Errorf("Expected:\n%s\nGot:\n%s", expected, buf.String())
	}
}
