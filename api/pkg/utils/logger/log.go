package logger

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// SurpressLog disable any logging
func SurpressLog() {
	zerolog.SetGlobalLevel(zerolog.Disabled)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
}

func Info(msg string, data map[string]interface{}) {
	l := log.Info()
	for k, v := range data {
		l.Interface(k, v)
	}
	l.Msg(msg)
}

func Error(msg string, data map[string]interface{}, err error) {
	l := log.Error().Err(err)
	for k, v := range data {
		l.Interface(k, v)
	}
	l.Msg(msg)
}
