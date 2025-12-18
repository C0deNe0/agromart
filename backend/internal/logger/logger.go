package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

type Logger struct {
	zerolog.Logger
}

func New(env string) *Logger {
	zerolog.TimeFieldFormat = time.RFC3339

	var writer zerolog.ConsoleWriter

	//this is for the local devs
	writer.Out = os.Stdout
	writer.TimeFormat = time.RFC3339 // consistent time format

	logger := zerolog.New(writer).With().Timestamp().Logger()

	if env == "local" || env == "development" {
		logger = logger.With().Caller().Logger()
	}

	return &Logger{Logger: logger}
}

func (l *Logger) Info(msg string) {
	l.Logger.Info().Msg(msg)
}

func (l *Logger) Error(err error, msg string) {
	l.Logger.Error().Err(err).Msg(msg)
}

func (l *Logger) Debug(msg string) {
	l.Logger.Debug().Msg(msg)
}
