package log

import "time"

// Record stands for a single record of log, usually a single line
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

// String returns the raw message of the Record
func (r *Record) String() string {
	return r.msg
}
