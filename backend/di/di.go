package di

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"gorm.io/gorm"

	"github.com/BrosSquad/ts-1-chat-app/backend/logging"
)

type (
	Config struct {
		LogsRoot     string
		DbPath       string
		LogToConsole bool
	}

	Container interface {
		GetDatabase() *gorm.DB

		GetDebugLogger() *logging.Debug
		GetErrorLogger() *logging.Error
		GetInfoLogger() *logging.Info

		GetChatBuffer() uint16
	}

	container struct {
		viper *viper.Viper
		Config

		db          *gorm.DB
		errorLogger *logging.Error
		debugLogger *logging.Debug
		infoLogger  *logging.Info
	}
)



func (c *container) GetChatBuffer() uint16 {
	return uint16(c.viper.GetUint32("chat.buffer"))
}

func New(config Config, pathToConfigFile string) Container {
	v := viper.New()

	v.AddConfigPath(pathToConfigFile)
	v.SetConfigType("yaml")
	v.SetConfigName("config")

	err := v.ReadInConfig()

	if err != nil {
		log.Fatal().
			Err(err).
			Msg("Cannot read config file")
	}

	return &container{
		Config: config,
		viper:  v,
	}
}
