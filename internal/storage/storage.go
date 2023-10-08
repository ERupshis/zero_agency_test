package storage

import (
	"github.com/erupshis/zero_agency_test/internal/logger"
	"github.com/erupshis/zero_agency_test/internal/storage/manager"
)

type storage struct {
	mngr manager.BaseStorageManager
	log  logger.BaseLogger
}

func Create(storageManager manager.BaseStorageManager, baseLogger logger.BaseLogger) BaseStorage {
	return &storage{
		mngr: storageManager,
		log:  baseLogger,
	}
}

func (s *storage) AddNote() error {
	return nil
}

func (s *storage) EditNote() error {
	return nil
}

func (s *storage) GetNotes() error {
	return nil
}
