package handlers

import (
	"github.com/erupshis/zero_agency_test/internal/auth/jwtgenerator"
	"github.com/erupshis/zero_agency_test/internal/auth/users/managers"
	"github.com/erupshis/zero_agency_test/internal/auth/users/userdata"
	"github.com/erupshis/zero_agency_test/internal/helpers"
	"github.com/erupshis/zero_agency_test/internal/logger"
	"github.com/gofiber/fiber/v2"
)

func LoginHandler(usersManager managers.BaseUsersManager, jwt jwtgenerator.JwtGenerator, log logger.BaseLogger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		buf := c.Body()
		var user userdata.User
		if err := helpers.UnmarshalData(buf, &user); err != nil {
			c.Status(fiber.StatusBadRequest)
			log.Info("[LoginHandler] bad new user input userdata")
			return nil
		}

		userID, err := usersManager.GetUserId(user.Login)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			log.Info("[LoginHandler] failed to get userID from user's database: %w", err)
			return nil
		}

		if userID == -1 {
			c.Status(fiber.StatusUnauthorized)
			log.Info("[LoginHandler] failed to get userID from user's database: %w", err)
			return nil
		}

		authorized, err := usersManager.ValidateUser(user.Login, user.Password)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			log.Info("[LoginHandler] failed to check user's login/password in database")
			return nil
		}

		if !authorized {
			c.Status(fiber.StatusUnauthorized)
			log.Info("[LoginHandler] failed to authorize user")
			return nil
		}

		token, err := jwt.BuildJWTString(userID)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			log.Info("[LoginHandler] new token generation failed: %w", err)
			return nil
		}

		c.Set("Authorization", "Bearer "+token)
		c.Status(fiber.StatusOK)

		log.Info("[LoginHandler] user '%s' authenticated successfully", user.Login)

		return nil
	}
}
