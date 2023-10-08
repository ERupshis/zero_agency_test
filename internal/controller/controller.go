package controller

import (
	"github.com/erupshis/zero_agency_test/internal/config"
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

func (c *Controller) LaunchServer(cfg config.Config) {
	go func(host string) {
		app := fiber.New()
		c.route(app)
		c.log.Info("[Controller:LaunchServer] server is launching")
		if err := app.Listen(host); err != nil {
			c.log.Info("[Controller:LaunchServer] failed to launch server: %v", err)
		}
	}(cfg.Host)
}

func (c *Controller) route(app *fiber.App) {
	app.Post("/edit/:Id", func(c *fiber.Ctx) error {
		return c.SendString("Add")
	})

	app.Get("/list", func(c *fiber.Ctx) error {
		return c.SendString("Read")
	})
}
