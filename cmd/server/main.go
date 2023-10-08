package main

import (
	"fmt"
	"os"

	"github.com/erupshis/zero_agency_test/internal/config"
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
	mainStorage := storage.Create(log)

	log.Info(cfg.Host)

}
