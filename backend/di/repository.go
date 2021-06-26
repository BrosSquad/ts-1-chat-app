package di

import (
	"github.com/BrosSquad/ts-1-chat-app/backend/repositories/token"
	"github.com/BrosSquad/ts-1-chat-app/backend/repositories/user"
)

func (c *container) GetUserRepository() user.Repository {
	if c.userRepository == nil {
		c.userRepository = user.New(c.GetDatabase(), c.GetErrorLogger())
	}

	return c.userRepository
}

func (c *container) GetTokenRepository() token.Repository {
	if c.tokenRepository == nil {
		c.tokenRepository = token.New(c.GetDatabase(), c.GetErrorLogger())
	}

	return c.tokenRepository
}
