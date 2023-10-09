package storage

import (
	"context"

	"github.com/erupshis/zero_agency_test/db/models"
)

type BaseStorage interface {
	EditNote(ctx context.Context, note *models.News) error
	GetNotes(ctx context.Context, page int64, perPage int64) ([]models.News, error)
}
