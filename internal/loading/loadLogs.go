package loading

import (
	"bufio"
	"log/slog"
	"os"
)

var LogChannel = make(chan string)

type LoadLogs struct {
	FileName string
	logs     []string
}

func (loadLogs *LoadLogs) Load() error {

	file, err := os.Open(loadLogs.FileName)

	if err != nil {
		close(LogChannel)
		return err
	}

	scanner := bufio.NewScanner(file)

	go loadLogs.readLogs(scanner, file)

	slog.Info("Logs loaded")

	return nil
}

func (loadLogs *LoadLogs) readLogs(scanner *bufio.Scanner, file *os.File) {

	defer file.Close()

	defer close(LogChannel)

	for scanner.Scan() {
		line := scanner.Text()
		loadLogs.logs = append(loadLogs.logs, line)
		LogChannel <- line
	}
}
