package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func Init() {
	zerolog.TimeFieldFormat = time.RFC3339
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}).With().Logger()
}
func Info(msg string, fields map[string]interface{}) {
	event := log.Info()
	for k, v := range fields {
		event = event.Interface(k, v)
	}

	event.Msg(msg)
}
func Debug(msg string, fields map[string]interface{}) {
	event := log.Debug()

	for k, v := range fields {
		event = event.Interface(k, v)
	}

	event.Msg(msg)
}
func Error(msg string, err error, fields map[string]interface{}) {
	event := log.Error().Err(err).Caller(1)
	for k, v := range fields {
		event = event.Interface(k, v)
	}
	event.Msg(msg)
}
