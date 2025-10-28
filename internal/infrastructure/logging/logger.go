package logging

import (
	"log"
	"os"
)

// Logger provides logging functionality
type Logger struct {
	infoLog  *log.Logger
	errorLog *log.Logger
}

// NewLogger creates a new logger
func NewLogger() *Logger {
	return &Logger{
		infoLog:  log.New(os.Stdout, "INFO: ", log.LstdFlags),
		errorLog: log.New(os.Stderr, "ERROR: ", log.LstdFlags),
	}
}

// Info logs an info message
func (l *Logger) Info(message string) {
	l.infoLog.Println(message)
}

// Error logs an error message
func (l *Logger) Error(message string) {
	l.errorLog.Println(message)
}

// Infof logs a formatted info message
func (l *Logger) Infof(format string, args ...interface{}) {
	l.infoLog.Printf(format, args...)
}

// Errorf logs a formatted error message
func (l *Logger) Errorf(format string, args ...interface{}) {
	l.errorLog.Printf(format, args...)
}
