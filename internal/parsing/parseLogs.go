package parsing

import (
	"LogParsing/internal/loading"
	"LogParsing/internal/storage"
	"fmt"
	"log/slog"
	"regexp"
	"sync"
	"time"
)

var waitingGroup sync.WaitGroup

var waitingLock sync.Mutex

const (
	numberOfWorkers = 20

	dateFormat = "2006-01-02 15:04:05"
)

var logChannel = loading.LogChannel

type LogEntry struct {
	Timestamp time.Time
	Level     string
	Message   string
}
type ParsedLogs struct {
	logEntries []LogEntry
}

func Parse(storage storage.Storage) {

	parsedLogs := &ParsedLogs{}

	for worker := 0; worker < numberOfWorkers; worker++ {
		waitingGroup.Add(1)
		go parsedLogs.writeLogs(storage)
	}

	waitingGroup.Wait()

	slog.Info("Parsing of logs...done")
}

func (parsedLogs *ParsedLogs) writeLogs(storage storage.Storage) {

	defer waitingGroup.Done()

	for item := range logChannel {

		regex := regexp.MustCompile("(\\d+-\\d+-\\d+\\s+\\d+:\\d+:\\d+)\\s+(\\S+)\\s+(.*)")

		matches := regex.FindStringSubmatch(item)

		if len(matches) == 4 {

			timestamp, err := time.Parse(dateFormat, matches[1])

			if err != nil {
				slog.Error("Error parsing timestamp")
			}

			logEntry := LogEntry{Timestamp: timestamp, Level: matches[2], Message: matches[3]}

			waitingLock.Lock()

			parsedLogs.logEntries = append(parsedLogs.logEntries, logEntry)

			lastId, err := storage.AddLog(
				logEntry.Timestamp,
				logEntry.Level,
				logEntry.Message)

			slog.Info("Log added successfully", slog.String("Log Id : ", fmt.Sprint(lastId)))

			waitingLock.Unlock()
		}
	}
}
