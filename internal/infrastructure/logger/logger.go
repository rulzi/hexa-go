package logger

import (
	"github.com/sirupsen/logrus"
)

// Logger is a simple logger interface
type Logger interface {
	Info(msg string)
	Error(msg string)
	Fatal(msg string)
}

// LogrusLogger is a logrus implementation of Logger
type LogrusLogger struct {
	logger *logrus.Logger
}

// NewSimpleLogger creates a new LogrusLogger
func NewSimpleLogger() *LogrusLogger {
	log := logrus.New()
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
		ForceColors:   true,
	})
	log.SetLevel(logrus.InfoLevel)

	return &LogrusLogger{
		logger: log,
	}
}

// Info logs an info message
func (l *LogrusLogger) Info(msg string) {
	l.logger.Info(msg)
}

// Error logs an error message
func (l *LogrusLogger) Error(msg string) {
	l.logger.Error(msg)
}

// Fatal logs a fatal message and exits
func (l *LogrusLogger) Fatal(msg string) {
	l.logger.Fatal(msg)
}
