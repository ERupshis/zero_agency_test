package controller

import (
	"encoding/json"
	"strconv"

	"github.com/erupshis/zero_agency_test/db/models"
	"github.com/erupshis/zero_agency_test/internal/logger"
	"github.com/erupshis/zero_agency_test/internal/storage"
	"github.com/gofiber/fiber/v2"
)

type Controller struct {
	strg storage.BaseStorage

	log logger.BaseLogger
}

func Create(baseStorage storage.BaseStorage, baseLogger logger.BaseLogger) *Controller {
	return &Controller{
		strg: baseStorage,
		log:  baseLogger,
	}
}

func (ctrl *Controller) Route() *fiber.App {
	app := fiber.New()

	app.Post("/edit/:Id", ctrl.editNote)
	app.Get("/list", ctrl.getNotes)

	return app
}

func (ctrl *Controller) editNote(c *fiber.Ctx) error {
	rawNoteID := c.Params("Id")
	if rawNoteID == "" {
		ctrl.log.Info("[Controller:editNote] missing id in request")
		c.Status(fiber.StatusBadRequest)
		return nil
	}

	noteID, err := strconv.ParseInt(rawNoteID, 10, 64)
	if err != nil {
		ctrl.log.Info("[Controller:editNote] invalid id: %v", err)
		c.Status(fiber.StatusBadRequest)
		return nil
	}

	note := &models.News{}
	note.ID = noteID

	if err = ctrl.strg.EditNote(c.Context(), note); err != nil {
		ctrl.log.Info("[Controller:editNote] failed to edit note: %v", err)
		c.Status(fiber.StatusInternalServerError)
		return nil
	}

	c.Status(fiber.StatusOK)
	return nil
}

func (ctrl *Controller) getNotes(c *fiber.Ctx) error {
	notes, err := ctrl.strg.GetNotes(c.Context())
	if err != nil {
		ctrl.log.Info("[Controller:getNotes] failed to get notes from storage: %v", err)
		c.Status(fiber.StatusInternalServerError)
		return nil
	}

	if len(notes) == 0 {
		c.Status(fiber.StatusNoContent)
		ctrl.log.Info("[Controller:getNotes] couldn't find any note")
		return nil
	}

	resp, err := json.Marshal(notes)
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		ctrl.log.Info("[Controller:getNotes] failed to marshal response body: %v", err)
		return nil
	}

	c.Set("Content-Type", "application/json")
	_, err = c.Write(resp)
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		ctrl.log.Info("[Controller:getNotes] failed to marshal response body: %v", err)
		return nil
	}

	ctrl.log.Info("[Controller:getNotes] request successfully handled")
	return nil
}
