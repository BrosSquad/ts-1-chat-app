package di

import (
	"github.com/rs/zerolog/log"

	"github.com/BrosSquad/ts-1-chat-app/backend/services/password"
)

func (c *container) GetPasswordHasher() password.Hasher {
	if c.passwordHasher == nil {
		driver := c.viper.GetString("password.driver")

		log.Trace().
			Str("driver", driver).
			Msg("Password hashing driver")

		switch driver {
		case "bcrypt":
			cost := c.viper.GetInt("password.bcrypt.cost")

			log.Trace().Str("driver", driver).
				Int("cost", cost).
				Msg("Creating BCrypt Driver")

			c.passwordHasher = password.NewBCryptHasher(c.GetErrorLogger(), cost)
		default:
			log.Fatal().Str("driver", driver).Msg("Unknown Password Hashing driver")
		}

	}

	return c.passwordHasher
}
