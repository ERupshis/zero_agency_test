package ram

import (
	"github.com/erupshis/zero_agency_test/internal/logger"
	"github.com/erupshis/zero_agency_test/internal/storage/manager"
)

type ram struct {
	log logger.BaseLogger
}

func Create(baseLogger logger.BaseLogger) manager.BaseStorageManager {
	return &ram{
		log: baseLogger,
	}
}

func (r *ram) AddNote() error {
	return nil
}

func (r *ram) EditNote() error {
	return nil
}

func (r *ram) GetNotes() error {
	return nil
}
