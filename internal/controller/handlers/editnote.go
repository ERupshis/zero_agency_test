package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/erupshis/zero_agency_test/db/models"
	"github.com/erupshis/zero_agency_test/internal/constants"
	"github.com/erupshis/zero_agency_test/internal/helpers"
	"github.com/erupshis/zero_agency_test/internal/logger"
	"github.com/erupshis/zero_agency_test/internal/storage"
	errStrg "github.com/erupshis/zero_agency_test/internal/storage/errors"
	"github.com/gofiber/fiber/v2"
)

func EditNode(storage storage.BaseStorage, log logger.BaseLogger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		rawNoteID := c.Params("ID")
		if rawNoteID == "" {
			log.Info("[Controller:editNote] missing id in request")
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		noteID, err := strconv.ParseInt(rawNoteID, 10, 64)
		if err != nil {
			log.Info("[Controller:editNote] invalid id: %v", err)
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		note := &models.News{
			Title:      constants.MissingStringFlag,
			Content:    constants.MissingStringFlag,
			Categories: make([]int64, len(constants.MissingInt64ArrayFlag)),
		}
		copy(note.Categories, constants.MissingInt64ArrayFlag)

		c.Body()
		if err = json.Unmarshal(c.Body(), note); err != nil {
			log.Info("[Controller:editNote] failed to parse request body")
			c.Status(fiber.StatusBadRequest)
			return nil
		}

		helpers.NormalizeInt64(&noteID)
		note.ID = noteID
		note.Categories = helpers.RemoveDuplicatesInt(note.Categories)

		if err = storage.EditNote(c.Context(), note); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				log.Info("[Controller:editNote] couldn't find note with Id '%d': %v", note.ID, err)
				c.Status(fiber.StatusBadRequest)
				return nil
			}

			if errors.Is(err, errStrg.ErrIncorrectNewNote) {
				log.Info("[Controller:editNote] failed to create new note: %v", err)
				c.Status(fiber.StatusBadRequest)
				return nil
			}

			log.Info("[Controller:editNote] failed to edit note: %v", err)
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		if _, err = c.Write([]byte(fmt.Sprintf("Id: %d", note.ID))); err != nil {
			log.Info("[Controller:editNote] failed to write response: %v", err)
			c.Status(fiber.StatusInternalServerError)
			return nil
		}

		c.Status(fiber.StatusOK)
		return nil
	}
}
