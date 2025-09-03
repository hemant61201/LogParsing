package storage

import (
	"LogParsing/internal/types"
	"time"
)

type Storage interface {
	AddLog(dateTime time.Time, logType string, logMessage string) (int64, error)
	GetLog(logType string) ([]types.LogParsing, error)
	GetAllLog() ([]types.LogParsing, error)
}
