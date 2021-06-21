package logging

import (
	"context"
	"io"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func ConfigureDefaultLogger(ctx context.Context, level string, writer io.Writer, toJson bool) {
	if writer == nil {
		writer = os.Stderr
	}

	zerolog.SetGlobalLevel(Parse(level))
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	if !toJson {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: writer})
	} else {
		log.Logger = log.Output(writer)
	}
}
