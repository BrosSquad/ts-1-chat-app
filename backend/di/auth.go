package di

import "github.com/BrosSquad/ts-1-chat-app/backend/services/auth"

func (c *container) GetTokenService() auth.TokenService {
	if c.tokenService == nil {
		c.tokenService = auth.NewTokenService(c.GetTokenRepository())
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
		c.loginService = auth.NewLoginService(c.GetUserRepository(), c.GetTokenService(), c.GetPasswordHasher())
	}

	return c.loginService
}
