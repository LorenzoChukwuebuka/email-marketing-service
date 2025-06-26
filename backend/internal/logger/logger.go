package logger

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"sync"
	"time"
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
	mu       sync.Mutex
	minLevel LogLevel
}

func NewLogger(filename string, minLevel LogLevel) (*Logger, error) {
	out, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	return &Logger{out: out, minLevel: minLevel}, nil
}

func (l *Logger) log(level LogLevel, msg string, args ...interface{}) {
	if level < l.minLevel {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now().Format("2006-01-02 15:04:05")
	levelStr := getLevelString(level)
	formattedMsg := fmt.Sprintf(msg, args...)

	_, file, line, _ := runtime.Caller(2)
	logLine := fmt.Sprintf("[%s] [%s] %s:%d - %s\n", now, levelStr, file, line, formattedMsg)
	l.out.WriteString(logLine)

	if level == FATAL {
		l.out.Close()
		os.Exit(1)
	}
}

func (l *Logger) logStruct(level LogLevel, msg string, data interface{}) {
	if level < l.minLevel {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now().Format("2006-01-02 15:04:05")
	levelStr := getLevelString(level)

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		jsonData = []byte(fmt.Sprintf("Error marshaling struct: %v", err))
	}

	_, file, line, _ := runtime.Caller(2)
	logLine := fmt.Sprintf("[%s] [%s] %s:%d - %s\nData: %s\n",
		now, levelStr, file, line, msg, string(jsonData))
	l.out.WriteString(logLine)

	if level == FATAL {
		l.out.Close()
		os.Exit(1)
	}
}

func (l *Logger) Debug(msg string, args ...interface{}) {
	l.log(DEBUG, msg, args...)
}

func (l *Logger) Info(msg string, args ...interface{}) {
	l.log(INFO, msg, args...)
}

func (l *Logger) Warn(msg string, args ...interface{}) {
	l.log(WARN, msg, args...)
}

func (l *Logger) Error(msg string, args ...interface{}) {
	l.log(ERROR, msg, args...)
}

func (l *Logger) Fatal(msg string, args ...interface{}) {
	l.log(FATAL, msg, args...)
}

func (l *Logger) DebugStruct(msg string, data interface{}) {
	l.logStruct(DEBUG, msg, data)
}

func (l *Logger) InfoStruct(msg string, data interface{}) {
	l.logStruct(INFO, msg, data)
}

func (l *Logger) WarnStruct(msg string, data interface{}) {
	l.logStruct(WARN, msg, data)
}

func (l *Logger) ErrorStruct(msg string, data interface{}) {
	l.logStruct(ERROR, msg, data)
}

func (l *Logger) FatalStruct(msg string, data interface{}) {
	l.logStruct(FATAL, msg, data)
}

func (l *Logger) Close() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.out.Close()
}

func getLevelString(level LogLevel) string {
	switch level {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	case FATAL:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
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

// Convenience functions that use the default logger
func Debug(msg string, args ...interface{}) {
	GetDefaultLogger().Debug(msg, args...)
}

func Info(msg string, args ...interface{}) {
	GetDefaultLogger().Info(msg, args...)
}

func Warn(msg string, args ...interface{}) {
	GetDefaultLogger().Warn(msg, args...)
}

func Error(msg string, args ...interface{}) {
	GetDefaultLogger().Error(msg, args...)
}

func Fatal(msg string, args ...interface{}) {
	GetDefaultLogger().Fatal(msg, args...)
}

func DebugStruct(msg string, data interface{}) {
	GetDefaultLogger().DebugStruct(msg, data)
}

func InfoStruct(msg string, data interface{}) {
	GetDefaultLogger().InfoStruct(msg, data)
}

func WarnStruct(msg string, data interface{}) {
	GetDefaultLogger().WarnStruct(msg, data)
}

func ErrorStruct(msg string, data interface{}) {
	GetDefaultLogger().ErrorStruct(msg, data)
}

func FatalStruct(msg string, data interface{}) {
	GetDefaultLogger().FatalStruct(msg, data)
}

func CloseDefaultLogger() {
	if defaultLogger != nil {
		defaultLogger.Close()
	}
}
