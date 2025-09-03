package sqlite

import (
	"LogParsing/internal/config"
	"LogParsing/internal/types"
	"database/sql"
	"time"

	_ "modernc.org/sqlite"
)

type Sqlite struct {
	Db *sql.DB
}

func NewSqlite(config *config.Config) (*Sqlite, error) {

	db, err := sql.Open("sqlite", config.StoragePath)

	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS logs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    dateTime DATETIME,
    logType TEXT,
    logMessage TEXT
	)`)

	if err != nil {
		return nil, err
	}

	return &Sqlite{Db: db}, nil
}

func (sqlite *Sqlite) AddLog(dateTime time.Time, logType string, logMessage string) (int64, error) {

	stmt, err := sqlite.Db.Prepare("INSERT INTO logs (dateTime, logType, logMessage) VALUES (?, ?, ?)")

	if err != nil {
		return 0, err
	}

	defer stmt.Close()

	result, err := stmt.Exec(dateTime, logType, logMessage)

	if err != nil {
		return 0, err
	}

	lastId, err := result.LastInsertId()

	if err != nil {
		return 0, err
	}

	return lastId, nil
}

func (sqlite *Sqlite) GetLog(logType string) ([]types.LogParsing, error) {

	stmt, err := sqlite.Db.Prepare("SELECT * FROM logs WHERE logType = ?")

	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	rows, err := stmt.Query(logType)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var logs []types.LogParsing

	for rows.Next() {
		var log types.LogParsing

		err := rows.Scan(&log.Id, &log.DateTime, &log.LogMessage, &log.LogType)

		if err != nil {
			return nil, err
		}

		logs = append(logs, log)
	}

	return logs, nil
}

func (sqlite *Sqlite) GetAllLog() ([]types.LogParsing, error) {

	stmt, err := sqlite.Db.Prepare("SELECT * FROM logs")

	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	rows, err := stmt.Query()

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var logs []types.LogParsing

	for rows.Next() {
		var log types.LogParsing

		err := rows.Scan(&log.Id, &log.DateTime, &log.LogMessage, &log.LogType)

		if err != nil {
			return nil, err
		}

		logs = append(logs, log)
	}

	return logs, nil
}
