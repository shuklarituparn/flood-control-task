package logger

import (
	"log"
	"os"
	"path/filepath"
)

var FileLogger *log.Logger

func SetupLogger() *log.Logger {
	relativeLogsDir := "./logs"
	absLogsDir, err := filepath.Abs(relativeLogsDir)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := os.Stat(absLogsDir); os.IsNotExist(err) {
		if err := os.MkdirAll(absLogsDir, 0777); err != nil {
			log.Fatal(err)
		}
	}
	logFilePath := filepath.Join(absLogsDir, "app.log")
	logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	FileLogger = log.New(logFile, "SERVER_LOGS: ", log.LstdFlags)
	return FileLogger
}
