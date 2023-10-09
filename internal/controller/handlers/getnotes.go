package handlers

import (
	"encoding/json"

	"github.com/erupshis/zero_agency_test/internal/logger"
	"github.com/erupshis/zero_agency_test/internal/storage"
	"github.com/gofiber/fiber/v2"
)

func GetNotes(storage storage.BaseStorage, log logger.BaseLogger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		notes, err := storage.GetNotes(c.Context())
		if err != nil {
			log.Info("[Controller:getNotes] failed to get notes from storage: %v", err)
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		if len(notes) == 0 {
			c.Status(fiber.StatusNoContent)
			log.Info("[Controller:getNotes] couldn't find any note")
			return nil
		}

		resp, err := json.Marshal(notes)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			log.Info("[Controller:getNotes] failed to marshal response body: %v", err)
			return nil
		}

		c.Set("Content-Type", "application/json")
		_, err = c.Write(resp)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			log.Info("[Controller:getNotes] failed to marshal response body: %v", err)
			return nil
		}

		log.Info("[Controller:getNotes] request successfully handled")
		return nil
	}
}
