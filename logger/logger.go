/*
 * Copyright (C) 2022 Marian Micek
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
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

// Initialize.
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

// Log clean mesage.
func (logger *Logger) Echo(source string, reason string) {
	logger.write("", "exeQute", reason, nil)
}

// Log info message.
func (logger *Logger) Info(source string, reason string) {
	logger.write(INFO, source, reason, nil)
}

// Log error.
func (logger *Logger) Error(source string, reason string, err error) {
	logger.write(ERROR, source, reason, err)
}

// Log error and panic.
func (logger *Logger) Fatal(source string, reason string, err error) {
	logger.write(ERROR, source, reason, err)
	os.Exit(1)
}

// Write to log file.
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
