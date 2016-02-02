package log

import "time"

type Record struct {
	fileLine  string
	name      string
	now       time.Time
	lv        LevelType
	msg       string
	rpcID     string
	requestID string
	appID     string
}

func (r *Record) String() string {
	return r.msg
}
