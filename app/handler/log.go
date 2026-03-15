package handler

import (
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// ConfigureLogger sets up the global zerolog logger according to the requested
// output mode and minimum log level.
//
// mode:
//   - "STACKDRIVER" – structured JSON to stdout with field names and severity
//     values that match the Google Cloud Logging / Stackdriver schema.
//   - anything else  – human-readable coloured output to stderr (dev default).
//
// level: any value accepted by zerolog.ParseLevel ("debug", "info", "warn",
// "error", "fatal", "panic").  Defaults to "debug" when the value is empty or
// unrecognised.
func ConfigureLogger(mode string, level string) {
	zerolog.TimeFieldFormat = time.RFC3339Nano

	if strings.ToUpper(mode) == "STACKDRIVER" {
		log.Logger = zerolog.New(os.Stdout).With().Timestamp().Logger()

		zerolog.LevelFieldName = "severity"
		zerolog.TimestampFieldName = "time"

		zerolog.LevelFieldMarshalFunc = func(l zerolog.Level) string {
			severity := map[zerolog.Level]string{
				zerolog.DebugLevel: "DEBUG",
				zerolog.InfoLevel:  "INFO",
				zerolog.WarnLevel:  "WARNING",
				zerolog.ErrorLevel: "ERROR",
				zerolog.FatalLevel: "CRITICAL",
				zerolog.PanicLevel: "EMERGENCY",
			}[l]
			return severity
		}
	} else {
		log.Logger = log.Output(
			zerolog.ConsoleWriter{
				Out:        os.Stderr,
				TimeFormat: time.RFC3339Nano,
			},
		)
	}

	logLevel, err := zerolog.ParseLevel(level)
	if err != nil {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(logLevel)
	}
}

// NewLogger returns a child of the global logger. Call ConfigureLogger first
// so that the global logger is fully initialised before any component requests
// its own logger instance.
func NewLogger() zerolog.Logger {
	return log.Logger
}
