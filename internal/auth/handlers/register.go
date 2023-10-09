package handlers

import (
	"github.com/erupshis/zero_agency_test/internal/auth/jwtgenerator"
	"github.com/erupshis/zero_agency_test/internal/auth/users/managers"
	"github.com/erupshis/zero_agency_test/internal/auth/users/userdata"
	"github.com/erupshis/zero_agency_test/internal/helpers"
	"github.com/erupshis/zero_agency_test/internal/logger"
	"github.com/gofiber/fiber/v2"
)

func RegisterHandler(usersManager managers.BaseUsersManager, jwt jwtgenerator.JwtGenerator, log logger.BaseLogger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		buf := c.Body()
		var user userdata.User
		if err := helpers.UnmarshalData(buf, &user); err != nil {
			c.Status(fiber.StatusBadRequest)
			log.Info("[RegisterHandler] bad new user input userdata")
			return nil
		}

		userID, err := usersManager.GetUserId(user.Login)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			log.Info("[RegisterHandler] failed to check user in database")
			return nil
		}

		if userID != -1 {
			c.Status(fiber.StatusConflict)
			log.Info("[RegisterHandler] login already exists")
			return nil
		}

		userID, err = usersManager.AddUser(user.Login, user.Password)
		if err != nil || userID == -1 {
			c.Status(fiber.StatusInternalServerError)
			log.Info("[RegisterHandler] failed to add new user '%s'", user.Login)
			return nil
		}

		token, err := jwt.BuildJWTString(userID)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			log.Info("[RegisterHandler] new token generation failed: %w", err)
			return nil
		}

		c.Set("Authorization", "Bearer "+token)
		c.Status(fiber.StatusCreated)

		log.Info("[RegisterHandler] user '%s' registered successfully", user.Login)

		return nil
	}
}
