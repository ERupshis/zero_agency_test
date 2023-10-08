package storage

import (
	"context"

	"github.com/erupshis/zero_agency_test/db/models"
)

type BaseStorage interface {
	EditNote(ctx context.Context, note *models.News) error
	GetNotes(ctx context.Context) ([]models.News, error)
}
