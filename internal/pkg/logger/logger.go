package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var logger zerolog.Logger

// Init initializes the logger
func Init(level string) {
	zerolog.TimeFieldFormat = time.RFC3339

	// Parse log level
	logLevel, err := zerolog.ParseLevel(level)
	if err != nil {
		logLevel = zerolog.InfoLevel
	}

	// Create logger with console writer for better readability
	logger = zerolog.New(os.Stdout).
		Level(logLevel).
		With().
		Timestamp().
		Caller().
		Logger()

	log.Logger = logger
}

// Get returns the global logger instance
func Get() *zerolog.Logger {
	return &logger
}

// Debug returns a debug level logger
func Debug() *zerolog.Event {
	return logger.Debug()
}

// Info returns an info level logger
func Info() *zerolog.Event {
	return logger.Info()
}

// Warn returns a warn level logger
func Warn() *zerolog.Event {
	return logger.Warn()
}

// Error returns an error level logger
func Error() *zerolog.Event {
	return logger.Error()
}

// Fatal returns a fatal level logger
func Fatal() *zerolog.Event {
	return logger.Fatal()
}
