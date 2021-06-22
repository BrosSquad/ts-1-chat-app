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

func NewErrorLogger(rootDir string, enabled, toConsole bool) *Error {
	var level string

	if enabled {
		level = "error"
	} else {
		level = "disabled"
	}

	path, output, logger := createZerologLogger("error", rootDir, level, toConsole)
	return &Error{
		Logger: logger,
		DefaultLogger: DefaultLogger{
			path: path,
			file: output,
		},
	}
}

func NewInfoLogger(rootDir string, enabled, toConsole bool) *Info {
	var level string

	if enabled {
		level = "info"
	} else {
		level = "disabled"
	}

	path, output, logger := createZerologLogger("info", rootDir, level, toConsole)
	return &Info{
		Logger: logger,
		DefaultLogger: DefaultLogger{
			path: path,
			file: output,
		},
	}
}

func NewDebugLogger(rootDir string, enabled bool) *Debug {
	var level string

	if enabled {
		level = "debug"
	} else {
		level = "disabled"
	}

	path, output, logger := createZerologLogger("debug", rootDir, level, true)
	return &Debug{
		Logger: logger,
		DefaultLogger: DefaultLogger{
			path: path,
			file: output,
		},
	}
}
