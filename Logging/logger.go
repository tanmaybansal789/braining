package Logging

import (
	"fmt"
	"os"
	"time"
)

type LogSeverity int

const (
	DEBUG LogSeverity = iota
	TRACE
	INFO
	WARNING
	ERROR
)

type Logger struct {
	minSeverity LogSeverity
	colors      map[LogSeverity]string
	name        string
}

func NewLogger(name string, minSeverity LogSeverity, colors map[LogSeverity]string) *Logger {
	return &Logger{name: name, minSeverity: minSeverity, colors: colors}
}

func NewLoggerWithDefaultColors(name string, minSeverity LogSeverity) *Logger {
	defaultColors := map[LogSeverity]string{
		DEBUG:   "\033[36m", // Cyan
		TRACE:   "\033[34m", // Blue
		INFO:    "\033[32m", // Green
		WARNING: "\033[33m", // Yellow
		ERROR:   "\033[31m", // Red
	}
	return NewLogger(name, minSeverity, defaultColors)
}

func (l *Logger) log(severity LogSeverity, format string, args ...interface{}) {
	if severity >= l.minSeverity {
		msg := fmt.Sprintf(format, args...)
		timestamp := time.Now().Format(time.RFC3339)
		color := l.colors[severity]
		resetColor := "\033[0m"
		fmt.Printf("%s[%s] [%s] %s%s\n", color, timestamp, severity.String(), msg, resetColor)
		if severity == ERROR {
			os.Exit(1)
		}
	}
}

func (l *Logger) Debug(format string, args ...interface{}) {
	l.log(DEBUG, format, args...)
}

func (l *Logger) Trace(format string, args ...interface{}) {
	l.log(TRACE, format, args...)
}

func (l *Logger) Info(format string, args ...interface{}) {
	l.log(INFO, format, args...)
}

func (l *Logger) Warning(format string, args ...interface{}) {
	l.log(WARNING, format, args...)
}

func (l *Logger) Error(format string, args ...interface{}) {
	l.log(ERROR, format, args...)
}

func (severity LogSeverity) String() string {
	switch severity {
	case DEBUG:
		return "DEBUG"
	case TRACE:
		return "TRACE"
	case INFO:
		return "INFO"
	case WARNING:
		return "WARNING"
	case ERROR:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}
