package logger

import (
	"os"

	"github.com/rs/zerolog"
)

var Logger zerolog.Logger

func InitLogger(jsonbool, debug bool) {
	logLevel := zerolog.InfoLevel
	zerolog.SetGlobalLevel(logLevel)

	if debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	Logger = zerolog.New(os.Stdout).With().Timestamp().Logger()

	if !jsonbool {
		output := zerolog.ConsoleWriter{Out: os.Stdout}
		Logger = zerolog.New(output).With().Timestamp().Logger()
	}
}
