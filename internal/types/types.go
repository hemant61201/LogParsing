package types

import "time"

type LogParsing struct {
	Id         int64
	DateTime   time.Time
	LogType    string
	LogMessage string
}
