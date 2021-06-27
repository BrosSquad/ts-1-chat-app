package di

import (
	"github.com/BrosSquad/ts-1-chat-app/backend/repositories/token"
	"github.com/BrosSquad/ts-1-chat-app/backend/repositories/user"
	"github.com/BrosSquad/ts-1-chat-app/backend/services/auth"
	"github.com/BrosSquad/ts-1-chat-app/backend/validators"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"gorm.io/gorm"

	"github.com/BrosSquad/ts-1-chat-app/backend/logging"
	"github.com/BrosSquad/ts-1-chat-app/backend/services/password"
)

type (
	Config struct {
		LogsRoot     string
		LogToConsole bool
	}

	Container interface {
		GetDatabase() *gorm.DB

		GetDebugLogger() *logging.Debug
		GetErrorLogger() *logging.Error
		GetInfoLogger() *logging.Info

		GetPasswordHasher() password.Hasher
		GetValidator() validators.Validator
		GetTokenService() auth.TokenService

		GetRegisterService() auth.RegisterService
		GetLoginService() auth.LoginService
		GetLogoutService() auth.LogoutService

		GetChatBuffer() uint16
		GetConfig() *viper.Viper
	}

	container struct {
		viper *viper.Viper
		Config

		db          *gorm.DB
		errorLogger *logging.Error
		debugLogger *logging.Debug
		infoLogger  *logging.Info

		// Services
		passwordHasher  password.Hasher
		validator       validators.Validator
		tokenService    auth.TokenService
		logoutService   auth.LogoutService
		registerService auth.RegisterService
		loginService    auth.LoginService

		// Repositories
		userRepository  user.Repository
		tokenRepository token.Repository
	}
)

func (c *container) GetChatBuffer() uint16 {
	return uint16(c.viper.GetUint32("chat.buffer"))
}

func (c *container) GetConfig() *viper.Viper {
	return c.viper
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
