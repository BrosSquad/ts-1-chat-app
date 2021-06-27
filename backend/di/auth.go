package di

import (
	"lukechampine.com/blake3"

	"github.com/BrosSquad/ts-1-chat-app/backend/services/auth"
)

func (c *container) GetTokenService() auth.TokenService {
	if c.tokenService == nil {
		c.tokenService = auth.NewTokenService(c.GetTokenRepository(), blake3.New(64, nil))
	}

	return c.tokenService
}

func (c *container) GetRegisterService() auth.RegisterService {
	if c.registerService == nil {
		c.registerService = auth.NewRegisterService(c.GetUserRepository(), c.GetPasswordHasher())
	}

	return c.registerService
}

func (c *container) GetLoginService() auth.LoginService {
	if c.loginService == nil {
		c.loginService = auth.NewLoginService(
			c.GetUserRepository(),
			c.GetTokenService(),
			c.GetPasswordHasher(),
			c.GetErrorLogger(),
		)
	}

	return c.loginService
}

func (c *container) GetLogoutService() auth.LogoutService {
	if c.logoutService == nil {
		c.logoutService = auth.NewLogoutService(c.GetTokenService())
	}

	return c.logoutService
}
