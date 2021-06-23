package logging

import (
	"context"
	"io"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func ConfigureDefaultLogger(ctx context.Context, level string, writer io.Writer, toConsole, toJson bool) {

	zerolog.SetGlobalLevel(Parse(level))
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	var w []io.Writer

	if toConsole {
		if !toJson {
			w = append(w, zerolog.ConsoleWriter{Out: os.Stdout})
		} else {
			w = append(w, os.Stdout)
		}
	}

	if writer != nil {
		w = append(w, writer)
	}

	log.Logger = log.Output(io.MultiWriter(w...))
}
