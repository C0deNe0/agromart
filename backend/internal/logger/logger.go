package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

// type Logger struct {
// 	zerolog.Logger
// }

func New(env string) *zerolog.Logger {
	//how the time should be displayed
	zerolog.TimeFieldFormat = time.RFC3339

	//this is for the local devs
	// consistent time format
	writer := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
	}

	log := zerolog.New(writer).With().Timestamp().Logger()

	if env == "local" || env == "development" {
		log = log.With().Caller().Logger()
	}

	return &log
}
