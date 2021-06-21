package logging

import (
	"strings"

	"github.com/rs/zerolog"
)

func Parse(level string) zerolog.Level {
	level = strings.ToLower(level)

	switch level {
	case "panic":
		return zerolog.PanicLevel
	case "fatal":
		return zerolog.FatalLevel
	case "error":
		return zerolog.ErrorLevel
	case "warn":
		return zerolog.WarnLevel
	case "debug":
		return zerolog.DebugLevel
	case "trace":
		return zerolog.TraceLevel
	case "info":
		return zerolog.InfoLevel
	}

	return zerolog.Disabled
}
