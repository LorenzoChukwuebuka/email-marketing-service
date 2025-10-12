package logger

import (
	"os"
	"sync"
)

type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
	FATAL
)

type Logger struct {
	out      *os.File
	minLevel LogLevel
}

func NewLogger(filename string, minLevel LogLevel) (*Logger, error) {
	out, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	return &Logger{out: out, minLevel: minLevel}, nil
}

// Singleton logger implementation
var (
	defaultLogger *Logger
	once          sync.Once
	initErr       error
)

// InitDefaultLogger initializes the default logger exactly once
func InitDefaultLogger(filename string, level LogLevel) error {
	once.Do(func() {
		defaultLogger, initErr = NewLogger(filename, level)
	})
	return initErr
}

// GetDefaultLogger returns the initialized default logger
// Panics if called before InitDefaultLogger
func GetDefaultLogger() *Logger {
	if defaultLogger == nil {
		panic("default logger not initialized - call InitDefaultLogger first")
	}
	return defaultLogger
}
