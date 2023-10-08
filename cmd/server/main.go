package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/erupshis/zero_agency_test/internal/config"
	"github.com/erupshis/zero_agency_test/internal/controller"
	"github.com/erupshis/zero_agency_test/internal/logger"
	"github.com/erupshis/zero_agency_test/internal/storage"
	"github.com/erupshis/zero_agency_test/internal/storage/manager/ram"
)

func main() {
	//log.
	log, err := logger.CreateZapLogger("info")
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to create logger: %v", err)
	}

	//config.
	cfg := config.Parse()

	ramManager := ram.Create(log)
	mainStorage := storage.Create(ramManager, log)
	mainController := controller.Create(mainStorage, log)

	mainController.LaunchServer(cfg)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, os.Kill, syscall.SIGTERM)
	<-sigCh
}
