package handlers

import (
	"encoding/json"
	"strconv"

	"github.com/erupshis/zero_agency_test/internal/logger"
	"github.com/erupshis/zero_agency_test/internal/storage"
	"github.com/gofiber/fiber/v2"
)

func GetNotes(storage storage.BaseStorage, log logger.BaseLogger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		page, err := strconv.ParseInt(c.Query("page", "1"), 10, 64)
		if err != nil {
			log.Info("[Controller:getNotes] invalid 'page' in query: %v", err)
			c.Status(fiber.StatusBadGateway)
			return nil
		}
		perPage, err := strconv.ParseInt(c.Query("perPage", "10"), 10, 64)
		if err != nil {
			log.Info("[Controller:getNotes] invalid 'perPage' in query: %v", err)
			c.Status(fiber.StatusBadGateway)
			return nil
		}

		notes, err := storage.GetNotes(c.Context(), page, perPage)
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
