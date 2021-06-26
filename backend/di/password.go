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
			c.passwordHasher = c.getBcryptDriver()
		case "argon":
			c.passwordHasher = c.getArgon2Driver()
			fallthrough
		default:
			log.Fatal().
				Str("driver", driver).
				Msg("Unknown Password Hashing driver")
		}

	}

	return c.passwordHasher
}

func (c *container) getBcryptDriver() password.Hasher {

	cost := c.viper.GetInt("password.bcrypt.cost")

	log.Trace().Str("driver", "bcrypt").
		Int("cost", cost).
		Msg("Creating BCrypt Driver")

	return password.NewBCryptHasher(c.GetErrorLogger(), cost)
}

func (c *container) getArgon2Driver() password.Hasher {
	return nil
}
