package logging

import (
	"fmt"
	"io"
	"os"
	"path"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/BrosSquad/ts-1-chat-app/backend/utils"
)

const Ext = "jsonl"

type DefaultLogger struct {
	path string
	file *os.File
}

func (f DefaultLogger) GetFileName() string {
	return f.path
}

func (f *DefaultLogger) Close() error {
	return f.file.Close()
}

func createZerologLogger(fileName, rootDir, level string, toConsole bool) (string, *os.File, zerolog.Logger) {
	p := path.Join(rootDir, fmt.Sprintf("%s.%s", fileName, Ext))
	output, err := utils.CreateLogFile(p, 0744)

	if err != nil {
		log.Fatal().
			Err(err).
			Msgf("Cannot open %s file for logging", p)
	}

	var logger zerolog.Logger

	if toConsole {
		writers := make([]io.Writer, 0, 2)

		writers = append(writers, output)

		writers = append(writers, zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		})

		logger = zerolog.New(zerolog.MultiLevelWriter(writers...)).
			With().
			Timestamp().
			Logger().
			Level(Parse(level))
	} else {
		logger = zerolog.New(output).
			With().
			Logger().
			Level(zerolog.ErrorLevel)
	}

	return p, output, logger
}
