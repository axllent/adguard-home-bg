package app

import (
	"fmt"
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

	if Level == "v" {
		fmt.Println("v")
		logLevel = logger.WarningLevel
	} else if Level == "vv" {
		logLevel = logger.NoticeLevel
	} else if Level == "vvv" {
		logLevel = logger.InfoLevel
	} else if Level == "vvvv" {
		logLevel = logger.DebugLevel
	}

	// if LogToFile != "" {
	//
	// }

	l, _ = logger.New("log", 1, os.Stdout, logLevel)
	l.SetFormat("%{time} [%{level}] %{message}")

	return l
}
