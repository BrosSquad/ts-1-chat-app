package logging

import "github.com/rs/zerolog"

type DBLogger struct {
	DefaultLogger
	zerolog.Logger
}

func NewDBLoger() *DBLogger {
	return &DBLogger{}
}
