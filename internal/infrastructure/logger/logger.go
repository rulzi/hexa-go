package logger

import (
	"io"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

// Logger is an enhanced logger interface with structured logging support
type Logger interface {
	Debug(msg string)
	DebugWithFields(msg string, fields map[string]interface{})
	Info(msg string)
	InfoWithFields(msg string, fields map[string]interface{})
	Warn(msg string)
	WarnWithFields(msg string, fields map[string]interface{})
	Error(msg string)
	ErrorWithFields(msg string, fields map[string]interface{})
	Fatal(msg string)
	FatalWithFields(msg string, fields map[string]interface{})
	WithFields(fields map[string]interface{}) Logger
}

// LogrusLogger is a logrus implementation of Logger
type LogrusLogger struct {
	logger *logrus.Logger
	entry  *logrus.Entry
	file   *os.File
}

// LoggerConfig holds logger configuration
type LoggerConfig struct {
	Level         string // debug, info, warn, error
	Format        string // text, json
	FilePath      string // path to log file
	ReportCaller  bool
	EnableFile    bool
	EnableConsole bool
}

// NewLogger creates a new LogrusLogger with configuration
func NewLogger(config LoggerConfig) (*LogrusLogger, error) {
	log := logrus.New()

	// Set log level
	level, err := logrus.ParseLevel(config.Level)
	if err != nil {
		level = logrus.InfoLevel
	}
	log.SetLevel(level)

	// Set formatter
	if config.Format == "json" {
		log.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
		})
	} else {
		log.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
			ForceColors:     config.EnableConsole && !config.EnableFile,
		})
	}

	// Set report caller
	log.SetReportCaller(config.ReportCaller)

	var logFile *os.File
	var writers []io.Writer

	// Setup file output
	if config.EnableFile && config.FilePath != "" {
		// Create log directory if it doesn't exist
		logDir := filepath.Dir(config.FilePath)
		if err := os.MkdirAll(logDir, 0755); err != nil {
			return nil, err
		}

		// Open or create log file
		logFile, err = os.OpenFile(config.FilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return nil, err
		}
		writers = append(writers, logFile)
	}

	// Setup console output
	if config.EnableConsole {
		writers = append(writers, os.Stdout)
	}

	// Set output to multiple writers or default to stdout
	if len(writers) > 0 {
		log.SetOutput(io.MultiWriter(writers...))
	} else {
		log.SetOutput(os.Stdout)
	}

	return &LogrusLogger{
		logger: log,
		entry:  logrus.NewEntry(log),
		file:   logFile,
	}, nil
}

// NewSimpleLogger creates a new LogrusLogger with default configuration
func NewSimpleLogger() *LogrusLogger {
	logger, _ := NewLogger(LoggerConfig{
		Level:         "info",
		Format:        "text",
		EnableConsole: true,
		EnableFile:    false,
		ReportCaller:  false,
	})
	return logger
}

// Close closes the log file if it was opened
func (l *LogrusLogger) Close() error {
	if l.file != nil {
		return l.file.Close()
	}
	return nil
}

// Debug logs a debug message
func (l *LogrusLogger) Debug(msg string) {
	l.entry.Debug(msg)
}

// DebugWithFields logs a debug message with structured fields
func (l *LogrusLogger) DebugWithFields(msg string, fields map[string]interface{}) {
	l.entry.WithFields(fields).Debug(msg)
}

// Info logs an info message
func (l *LogrusLogger) Info(msg string) {
	l.entry.Info(msg)
}

// InfoWithFields logs an info message with structured fields
func (l *LogrusLogger) InfoWithFields(msg string, fields map[string]interface{}) {
	l.entry.WithFields(fields).Info(msg)
}

// Warn logs a warning message
func (l *LogrusLogger) Warn(msg string) {
	l.entry.Warn(msg)
}

// WarnWithFields logs a warning message with structured fields
func (l *LogrusLogger) WarnWithFields(msg string, fields map[string]interface{}) {
	l.entry.WithFields(fields).Warn(msg)
}

// Error logs an error message
func (l *LogrusLogger) Error(msg string) {
	l.entry.Error(msg)
}

// ErrorWithFields logs an error message with structured fields
func (l *LogrusLogger) ErrorWithFields(msg string, fields map[string]interface{}) {
	l.entry.WithFields(fields).Error(msg)
}

// Fatal logs a fatal message and exits
func (l *LogrusLogger) Fatal(msg string) {
	l.entry.Fatal(msg)
}

// FatalWithFields logs a fatal message with structured fields and exits
func (l *LogrusLogger) FatalWithFields(msg string, fields map[string]interface{}) {
	l.entry.WithFields(fields).Fatal(msg)
}

// WithFields returns a new logger with pre-set fields
func (l *LogrusLogger) WithFields(fields map[string]interface{}) Logger {
	return &LogrusLogger{
		logger: l.logger,
		entry:  l.entry.WithFields(fields),
		file:   l.file,
	}
}
