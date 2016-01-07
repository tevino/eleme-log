package log

import "time"

type Record struct {
	name  string
	now   time.Time
	lv    LevelType
	msg   string
	appID string
}

func (r *Record) String() string {
	return r.msg
}
