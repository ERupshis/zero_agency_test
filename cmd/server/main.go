package main

import (
	"fmt"
	"os"

	"github.com/erupshis/zero_agency_test/internal/config"
	"github.com/erupshis/zero_agency_test/internal/logger"
)

func main() {
	//log.
	log, err := logger.CreateZapLogger("info")
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to create logger: %v", err)
	}

	//config.
	cfg := config.Parse()

	log.Info(cfg.Host)

}
