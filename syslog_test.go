package log

import (
	"bytes"
	"strings"
	"testing"
)

func TestSyslogtpl(t *testing.T) {
	buf := bytes.NewBuffer(make([]byte, 0, 100))

	l := new(Logger)
	l.name = "name"
	l.lv = INFO
	l.handlers = make(map[Handler]bool)

	hdr, err := NewStreamHandler(buf, syslogTpl)
	if err != nil {
		t.Fatalf("NewStreamHandler Error:%v", err)
	}
	l.AddHandler(hdr)
	SetGlobalAppID("samaritan.test")
	defer SetGlobalAppID("")
	l.Info("TEST_TEST")

	strs := strings.Split(buf.String(), " ")

	if strs[0] != "[samaritan.test" || strs[4] != "TEST_TEST\n" ||
		strs[1] != "-" || strs[2] != "-]" {
		t.Errorf("SyslogTpl Error: %s", buf.String())
	}
}
