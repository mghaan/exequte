package logger

import (
	"fmt"
	"os"
	"time"
)

const (
	INFO         = "INFO"
	ERROR        = "ERROR"
	SYSTEM       = "SYSTEM"
	MQTT         = "MQTT"
	DEBUG_HOST   = "127.0.0.1"
	DEBUG_PORT   = 1883
	DEBUG_USER   = "test"
	DEBUG_PASS   = "pokus"
	DEBUG_CLIENT = "exequte"
)

type Logger struct {
	logfile *os.File
}

// Initialize
func New() *Logger {
	logger := &Logger{}

	return logger
}

// Create new empty log file.
func (logger *Logger) Create(path string) error {
	var err error

	if len(path) > 2 {
		logger.logfile, err = os.Create(path)
	}

	return err
}

// Log events
func (logger *Logger) Info(source string, reason string) {
	logger.write(INFO, source, reason, nil)
}

func (logger *Logger) Error(source string, reason string, err error) {
	logger.write(ERROR, source, reason, err)
}

func (logger *Logger) Fatal(source string, reason string, err error) {
	logger.write(ERROR, source, reason, err)
	os.Exit(1)
}

func (logger *Logger) write(category string, source string, reason string, err error) {
	t := time.Now()
	stamp := t.Format("02 Jan 2006 15:04:05")
	message := fmt.Sprintf("%s [%s] %s: %s", stamp, category, source, reason)
	if err != nil {
		message += "\n" + err.Error()
	}

	fmt.Fprintln(os.Stdout, message)
	if logger.logfile != nil {
		fmt.Fprintln(logger.logfile, message)
	}
}
