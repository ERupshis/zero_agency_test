package controller

import (
	"github.com/erupshis/zero_agency_test/internal/controller/handlers"
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

	app.Post("/edit/:ID", handlers.EditNode(ctrl.strg, ctrl.log))
	app.Get("/list", handlers.GetNotes(ctrl.strg, ctrl.log))

	return app
}
