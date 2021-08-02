package fio

import "log"

// IsLoggingEnabled allows to disable all logging done by package fio by setting it to false.
var IsLoggingEnabled = true

// Logger is the logger used to write logs by package fio when IsLoggingEnabled is true.
// If Logger is nil and IsLoggingEnabled is true, package fio will write logs with log.Default().
var Logger *log.Logger = nil

func logger() *log.Logger {
	if !IsLoggingEnabled {
		return nil
	}
	if Logger == nil {
		return log.Default()
	}
	return Logger
}
