package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/erupshis/zero_agency_test/internal/config"
	"github.com/erupshis/zero_agency_test/internal/controller"
	"github.com/erupshis/zero_agency_test/internal/helpers"
	"github.com/erupshis/zero_agency_test/internal/logger"
	"github.com/erupshis/zero_agency_test/internal/storage"
	postgresql "github.com/erupshis/zero_agency_test/internal/storage/manager/postgres"
	"github.com/gofiber/fiber/v2"
)

func main() {
	//log.
	log, err := logger.CreateZapLogger("info")
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to create logger: %v", err)
	}

	//config.
	cfg := config.Parse()

	ctxWithCancel, cancel := context.WithCancel(context.Background())
	defer cancel()

	//storage.
	postgreManager, err := postgresql.CreatePostgreDB(ctxWithCancel, cfg, log)
	if err != nil {
		log.Info("failed to open postgre database: %v", err)
		return
	}
	defer helpers.ExecuteWithLogError(postgreManager.Close, log)

	mainStorage := storage.Create(postgreManager, log)
	mainController := controller.Create(mainStorage, log)

	server := fiber.New()
	server.Use(log.LogHandler)
	server.Mount("/", mainController.Route())

	go func(log logger.BaseLogger) {
		log.Info("server is launching with host '%s'", cfg.Host)
		if err = server.Listen(cfg.Host); err != nil {
			log.Info("failed to launch server: %v", err)
		}

		log.Info("server has been stopped")
	}(log)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, os.Kill, syscall.SIGTERM)
	<-sigCh

	if err = server.ShutdownWithContext(ctxWithCancel); err != nil {
		log.Info("failed to stop server by context cancel: %v", err)
	} else {
		log.Info("server gracefully stopped")
	}
}
