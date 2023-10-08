package manager

import (
	"context"

	"github.com/erupshis/zero_agency_test/db/models"
)

type BaseStorageManager interface {
	EditNote(ctx context.Context, note *models.News) error
	GetNotes(ctx context.Context) ([]models.News, error)

	Close() error
	CheckConnection(ctx context.Context) (bool, error)
}
