package utils

import (
	"fmt"
	"os"
	"time"
)

type Logger struct {
	out *os.File
}

func NewLogger(filename string) (*Logger, error) {
	out, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	return &Logger{out: out}, nil
}

func (l *Logger) Info(msg string) {
	l.write("INFO", msg)
}

func (l *Logger) Error(msg string) {
	l.write("ERROR", msg)
}

func (l *Logger) write(level, msg string) {
	now := time.Now().Format("2006-01-02 15:04:05")
	line := fmt.Sprintf("[%s] [%s] %s\n", now, level, msg)
	l.out.WriteString(line)
}

func (l *Logger) Close() {
	l.out.Close()
}
