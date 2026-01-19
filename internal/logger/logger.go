package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var Logger zerolog.Logger

// InitLogger initializes the global logger with configuration
func InitLogger(logLevel, logFormat string) {
	// Set log level
	level := parseLogLevel(logLevel)
	zerolog.SetGlobalLevel(level)

	// Configure output format
	if logFormat == "console" {
		// Pretty console output for development
		output := zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		}
		Logger = zerolog.New(output).With().Timestamp().Caller().Logger()
	} else {
		// JSON output for production
		Logger = zerolog.New(os.Stdout).With().Timestamp().Caller().Logger()
	}

	// Set as global logger
	log.Logger = Logger

	Logger.Info().
		Str("level", level.String()).
		Str("format", logFormat).
		Msg("Logger initialized")
}

// parseLogLevel converts string log level to zerolog.Level
func parseLogLevel(level string) zerolog.Level {
	switch level {
	case "debug":
		return zerolog.DebugLevel
	case "info":
		return zerolog.InfoLevel
	case "warn":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	case "fatal":
		return zerolog.FatalLevel
	default:
		return zerolog.InfoLevel
	}
}

// GetLogger returns the global logger instance
func GetLogger() *zerolog.Logger {
	return &Logger
}

// Debug logs a debug message with fields
func Debug() *zerolog.Event {
	return Logger.Debug()
}

// Info logs an info message with fields
func Info() *zerolog.Event {
	return Logger.Info()
}

// Warn logs a warning message with fields
func Warn() *zerolog.Event {
	return Logger.Warn()
}

// Error logs an error message with fields
func Error() *zerolog.Event {
	return Logger.Error()
}

// Fatal logs a fatal message with fields and exits
func Fatal() *zerolog.Event {
	return Logger.Fatal()
}
