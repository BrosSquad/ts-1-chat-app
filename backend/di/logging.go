package di

import "github.com/BrosSquad/ts-1-chat-app/backend/logging"

func (c *container) GetDebugLogger() *logging.Debug {
	if c.debugLogger == nil {
		c.debugLogger = logging.NewDebugLogger(c.LogsRoot, c.viper.GetBool("logging.debug.enabled"))
	}

	return c.debugLogger
}

func (c *container) GetErrorLogger() *logging.Error {
	if c.errorLogger == nil {
		c.errorLogger = logging.NewErrorLogger(c.LogsRoot, c.viper.GetBool("logging.error.enabled"), c.LogToConsole)
	}

	return c.errorLogger

}

func (c *container) GetInfoLogger() *logging.Info {
	if c.infoLogger == nil {
		c.infoLogger = logging.NewInfoLogger(c.LogsRoot, c.viper.GetBool("logging.info.enabled"), c.LogToConsole)
	}

	return c.infoLogger
}
