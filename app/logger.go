package app

import (
	"os"

	"github.com/apsdehal/go-logger"
)

var (
	// Log global
	log *logger.Logger
	// Level - Log level
	Level = "vvvv" //string
	// LogToFile - log to file
	LogToFile string
)

// Log handler
func Log() *logger.Logger {
	if log == nil {
		log = initLogger()
	}

	return log
}

func initLogger() *logger.Logger {
	var l *logger.Logger

	logLevel := logger.ErrorLevel

	switch Level {
	case "v":
		logLevel = logger.WarningLevel
	case "vv":
		logLevel = logger.NoticeLevel
	case "vvv":
		logLevel = logger.InfoLevel
	case "vvvv":
		logLevel = logger.DebugLevel
	}

	l, _ = logger.New("log", 1, os.Stdout, logLevel)
	l.SetFormat("%{time} [%{level}] %{message}")

	return l
}
