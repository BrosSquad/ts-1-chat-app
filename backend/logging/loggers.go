package logging

import (
	"github.com/rs/zerolog"
)

type (
	Error struct {
		DefaultLogger
		zerolog.Logger
	}

	Info struct {
		DefaultLogger
		zerolog.Logger
	}

	Debug struct {
		DefaultLogger
		zerolog.Logger
	}
)

func NewErrorLogger(rootDir string, toConsole bool) *Error {
	path, output, logger := createZerologLogger("error", rootDir, "error", toConsole)
	return &Error{
		Logger: logger,
		DefaultLogger: DefaultLogger{
			path: path,
			file: output,
		},
	}
}

func NewInfoLogger(rootDir string, toConsole bool) *Info {
	path, output, logger := createZerologLogger("info", rootDir, "info", toConsole)
	return &Info{
		Logger: logger,
		DefaultLogger: DefaultLogger{
			path: path,
			file: output,
		},
	}
}

func NewDebugLogger(rootDir string) *Debug {
	path, output, logger := createZerologLogger("debug", rootDir, "debug", true)
	return &Debug{
		Logger: logger,
		DefaultLogger: DefaultLogger{
			path: path,
			file: output,
		},
	}
}
