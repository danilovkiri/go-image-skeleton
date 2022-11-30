// Package logger provides logging functionality.

package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

// InitLog initializes a logger.
func InitLog() *zerolog.Logger {
	zerolog.TimeFieldFormat = time.RFC3339
	consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout}
	Logger := zerolog.New(consoleWriter).With().Timestamp().Logger()
	return &Logger
}
