package controller

import (
	"github.com/erupshis/zero_agency_test/internal/logger"
)

type Controller struct {
	strg storage.BaseStorage

	log logger.BaseLogger
}
