package di

import (
	"github.com/rs/zerolog/log"

	"github.com/BrosSquad/ts-1-chat-app/backend/logging"
)

func (c *container) GetDebugLogger() *logging.Debug {
	if c.debugLogger == nil {
		enabled := c.viper.GetBool("logging.debug.enabled")

		log.Trace().Bool("enabled", enabled).Str("logger", "debug").Msg("Creating Logger")

		c.debugLogger = logging.NewDebugLogger(c.LogsRoot, enabled)
	}

	return c.debugLogger
}

func (c *container) GetErrorLogger() *logging.Error {
	if c.errorLogger == nil {
		enabled := c.viper.GetBool("logging.error.enabled")

		log.Trace().Bool("enabled", enabled).Str("logger", "error").Msg("Creating Logger")

		c.errorLogger = logging.NewErrorLogger(c.LogsRoot, enabled, c.LogToConsole)
	}

	return c.errorLogger
}

func (c *container) GetInfoLogger() *logging.Info {
	if c.infoLogger == nil {
		enabled := c.viper.GetBool("logging.info.enabled")

		log.Trace().Bool("enabled", enabled).Str("logger", "info").Msg("Creating Logger")

		c.infoLogger = logging.NewInfoLogger(c.LogsRoot, enabled, c.LogToConsole)
	}

	return c.infoLogger
}
