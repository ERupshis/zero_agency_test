package auth

import (
	"github.com/erupshis/zero_agency_test/internal/auth/handlers"
	"github.com/erupshis/zero_agency_test/internal/auth/jwtgenerator"
	"github.com/erupshis/zero_agency_test/internal/auth/middleware"
	"github.com/erupshis/zero_agency_test/internal/auth/users/managers"
	"github.com/erupshis/zero_agency_test/internal/logger"
	"github.com/gofiber/fiber/v2"
)

type Controller struct {
	usersStrg managers.BaseUsersManager
	jwt       jwtgenerator.JwtGenerator

	log logger.BaseLogger
}

func CreateAuthenticator(usersStorage managers.BaseUsersManager, jwt jwtgenerator.JwtGenerator, baseLogger logger.BaseLogger) *Controller {
	return &Controller{
		usersStrg: usersStorage,
		jwt:       jwt,
		log:       baseLogger,
	}
}

func (c *Controller) Route() *fiber.App {
	app := fiber.New()

	app.Post("/register", handlers.RegisterHandler(c.usersStrg, c.jwt, c.log))
	app.Post("/login", handlers.LoginHandler(c.usersStrg, c.jwt, c.log))

	return app
}

func (c *Controller) Authorize(userRole int) fiber.Handler {
	return middleware.AuthorizeUser(c.usersStrg, c.jwt, userRole, c.log)
}
