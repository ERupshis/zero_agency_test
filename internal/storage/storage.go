package storage

import (
	"context"
	"fmt"
	"sync"

	"github.com/erupshis/zero_agency_test/db/models"
	"github.com/erupshis/zero_agency_test/internal/logger"
	"github.com/erupshis/zero_agency_test/internal/storage/manager"
)

type storage struct {
	mu   sync.RWMutex
	mngr manager.BaseStorageManager

	log logger.BaseLogger
}

func Create(storageManager manager.BaseStorageManager, baseLogger logger.BaseLogger) BaseStorage {
	return &storage{
		mngr: storageManager,
		log:  baseLogger,
	}
}

func (s *storage) EditNote(ctx context.Context, note *models.News) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	err := s.mngr.EditNote(ctx, note)
	if err != nil {
		return fmt.Errorf("storage: %w", err)
	}
	return nil
}

func (s *storage) GetNotes(ctx context.Context, page int64, perPage int64) ([]models.News, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	notes, err := s.mngr.GetNotes(ctx, page, perPage)
	if err != nil {
		return nil, fmt.Errorf("storage: %w", err)
	}
	return notes, nil
}
