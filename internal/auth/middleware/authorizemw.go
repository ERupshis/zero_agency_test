package middleware

import (
	"context"
	"fmt"
	"strings"

	"github.com/erupshis/zero_agency_test/internal/auth/jwtgenerator"
	"github.com/erupshis/zero_agency_test/internal/auth/users/managers"
	"github.com/erupshis/zero_agency_test/internal/auth/users/userdata"
	"github.com/erupshis/zero_agency_test/internal/logger"
	"github.com/gofiber/fiber/v2"
)

func AuthorizeUser(usersManager managers.BaseUsersManager, jwt jwtgenerator.JwtGenerator, userRoleRequirement int, log logger.BaseLogger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if c.Path() == "/login" || c.Path() == "/register" {
			return c.Next()
		}

		authHeader := c.Get("Authorization")
		if authHeader == "" {
			log.Info("[AuthorizeUser] invalid request without authentication token")
			c.Status(fiber.StatusUnauthorized)
			return nil
		}

		token := strings.Split(authHeader, " ")
		if len(token) != 2 || token[0] != "Bearer" {
			log.Info("[AuthorizeUser] invalid invalid token")
			c.Status(fiber.StatusUnauthorized)
			return nil
		}

		userID := jwt.GetUserId(token[1])
		userRole, err := usersManager.GetUserRole(userID)
		if err != nil {
			log.Info("[AuthorizeUser] failed to search user in system: %v", err)
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		if userRole == -1 {
			log.Info("[AuthorizeUser] user is not registered in system")
			c.Status(fiber.StatusUnauthorized)
			return nil
		}

		if userRole < userRoleRequirement {
			log.Info("[AuthorizeUser] user doesn't have permission to resource: %s", c.Path())
			c.Status(fiber.StatusForbidden)
			return nil
		}

		ctxWithValue := context.WithValue(c.Context(), userdata.UserID, fmt.Sprintf("%d", userID))
		c.SetUserContext(ctxWithValue)
		if err = c.Next(); err != nil {
			log.Info("[AuthorizeUser] failed to handle request: %v", err)
			return nil
		}

		return nil
	}
}
