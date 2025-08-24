package utils

import (
	"fmt"
	"os"
	"time"
)

type LogLevel int

const (
	LogLevelError LogLevel = iota
	LogLevelInfo
	LogLevelDebug
)

type Logger struct {
	level     LogLevel
	verbose   bool
	timestamp bool
}

func NewLogger(verbose bool) *Logger {
	level := LogLevelInfo
	if verbose {
		level = LogLevelDebug
	}

	return &Logger{
		level:     level,
		verbose:   verbose,
		timestamp: verbose,
	}
}

func NewLoggerWithLevel(level LogLevel) *Logger {
	return &Logger{
		level:     level,
		verbose:   level >= LogLevelDebug,
		timestamp: level >= LogLevelDebug,
	}
}

func (l *Logger) log(level LogLevel, prefix, msg string, args ...interface{}) {
	if level > l.level {
		return
	}

	var output *os.File
	if level == LogLevelError {
		output = os.Stderr
	} else {
		output = os.Stdout
	}

	message := fmt.Sprintf(msg, args...)

	if l.timestamp {
		timestamp := time.Now().Format("15:04:05")
		fmt.Fprintf(output, "[%s] %s: %s\n", timestamp, prefix, message)
	} else {
		fmt.Fprintf(output, "%s: %s\n", prefix, message)
	}
}

func (l *Logger) Info(msg string, args ...interface{}) {
	l.log(LogLevelInfo, "INFO", msg, args...)
}

func (l *Logger) Error(msg string, args ...interface{}) {
	l.log(LogLevelError, "ERROR", msg, args...)
}

func (l *Logger) Debug(msg string, args ...interface{}) {
	l.log(LogLevelDebug, "DEBUG", msg, args...)
}

func (l *Logger) Fatal(msg string, args ...interface{}) {
	l.Error(msg, args...)
	os.Exit(1)
}

func (l *Logger) Success(msg string, args ...interface{}) {
	l.log(LogLevelInfo, "SUCCESS", msg, args...)
}

func (l *Logger) Warning(msg string, args ...interface{}) {
	l.log(LogLevelInfo, "WARNING", msg, args...)
}

// LogUserFriendlyError logs an error with user-friendly formatting
func (l *Logger) LogUserFriendlyError(err error) {
	friendlyMsg := GetUserFriendlyError(err)
	l.Error(friendlyMsg)

	if l.verbose {
		l.Debug("Technical error details: %v", err)
	}
}
