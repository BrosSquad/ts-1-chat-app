package di

import "github.com/BrosSquad/ts-1-chat-app/backend/services/auth"

func (c *container) GetTokenService() auth.TokenService {
	if c.tokenService == nil {
		c.tokenService = auth.NewTokenService()
	}

	return c.tokenService
}
