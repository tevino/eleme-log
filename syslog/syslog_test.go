package syslog

import (
	"bytes"
	"strings"
	"testing"

	"github.com/eleme/log"
	"github.com/eleme/log/rpc"
)

func TestSyslogtpl(t *testing.T) {
	buf := bytes.NewBuffer(make([]byte, 0, 100))

	elog := rpc.NewELogger("elog")

	ef := rpc.NewELogFormatter(false)
	if err := ef.ParseFormat(log.TplSyslog); err != nil {
		t.Error("error creating stream handler: ", err)
		t.FailNow()
	}
	h := log.NewStreamHandler(buf, ef)
	h.Colored(false)
	elog.AddHandler(h)

	log.SetGlobalAppID("samaritan.test")
	defer log.SetGlobalAppID("")
	elog.Info("TEST_TEST")

	strs := strings.Split(buf.String(), " ")

	if strs[0] != "[samaritan.test" || strs[4] != "TEST_TEST\n" ||
		strs[1] != "-" || strs[2] != "-]" {
		t.Errorf("SyslogTpl Error: %s", buf.String())
	}
}
