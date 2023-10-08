// Package postgresql postgresql handling PostgreSQL database.
package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sync"

	"github.com/erupshis/zero_agency_test/db/models"
	"github.com/erupshis/zero_agency_test/internal/config"
	"github.com/erupshis/zero_agency_test/internal/helpers"
	"github.com/erupshis/zero_agency_test/internal/logger"
	"github.com/erupshis/zero_agency_test/internal/retryer"
	"github.com/erupshis/zero_agency_test/internal/storage/manager"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgerrcode"
	_ "github.com/jackc/pgx/v4/stdlib"
	"gopkg.in/reform.v1"
	"gopkg.in/reform.v1/dialects/postgresql"
)

// databaseErrorsToRetry errors to retry request to database.
var databaseErrorsToRetry = []error{
	errors.New(pgerrcode.UniqueViolation),
	errors.New(pgerrcode.ConnectionException),
	errors.New(pgerrcode.ConnectionDoesNotExist),
	errors.New(pgerrcode.ConnectionFailure),
	errors.New(pgerrcode.SQLClientUnableToEstablishSQLConnection),
	errors.New(pgerrcode.SQLServerRejectedEstablishmentOfSQLConnection),
	errors.New(pgerrcode.TransactionResolutionUnknown),
	errors.New(pgerrcode.ProtocolViolation),
}

// postgresDB storageManager implementation for PostgreSQL. Consist of database and QueriesHandler.
// Request to database are synchronized by sync.RWMutex. All requests is done on united transaction. Multi insert/update/delete is not supported at the moment.
type postgresDB struct {
	database *sql.DB
	reformDB *reform.DB

	log logger.BaseLogger
	mu  sync.RWMutex
}

// CreatePostgreDB creates manager implementation. Supports migrations and check connection to database.
func CreatePostgreDB(ctx context.Context, cfg config.Config, log logger.BaseLogger) (manager.BaseStorageManager, error) {
	log.Info("[CreatePostgreDB] open database with settings: '%s'", cfg.DatabaseDSN)
	createDatabaseError := "create db: %w"
	database, err := sql.Open("pgx", cfg.DatabaseDSN)
	if err != nil {
		return nil, fmt.Errorf(createDatabaseError, err)
	}

	driver, err := postgres.WithInstance(database, &postgres.Config{})
	if err != nil {
		return nil, fmt.Errorf(createDatabaseError, err)
	}

	m, err := migrate.NewWithDatabaseInstance("file://db/migrations", "postgres", driver)
	if err != nil {
		return nil, fmt.Errorf(createDatabaseError, err)
	}

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return nil, fmt.Errorf(createDatabaseError, err)
	}

	dbManager := &postgresDB{
		database: database,
		reformDB: nil,
		log:      log,
	}

	if _, err = dbManager.CheckConnection(ctx); err != nil {
		return nil, fmt.Errorf(createDatabaseError, err)
	}

	database.SetMaxIdleConns(cfg.DatabaseIdleConn)
	database.SetMaxOpenConns(cfg.DatabaseOpenConn)
	dbManager.reformDB = reform.NewDB(database, postgresql.Dialect, reform.NewPrintfLogger(log.Printf))

	log.Info("[CreatePostgreDB] successful")
	return dbManager, nil
}

// CheckConnection checks connection to database.
func (p *postgresDB) CheckConnection(ctx context.Context) (bool, error) {
	exec := func(context context.Context) (int64, []byte, error) {
		return 0, []byte{}, p.database.PingContext(context)
	}
	_, _, err := retryer.RetryCallWithTimeout(ctx, p.log, nil, databaseErrorsToRetry, exec)
	if err != nil {
		return false, fmt.Errorf("check connection: %w", err)
	}
	return true, nil
}

// Close closes database.
func (p *postgresDB) Close() error {
	return p.database.Close()
}
func (p *postgresDB) getNote(ctx context.Context, ID int64) (*models.News, error) {
	note, err := p.reformDB.WithContext(ctx).FindByPrimaryKeyFrom(models.NewsTable, ID)
	if err != nil {
		return nil, fmt.Errorf("get note from db: %w", err)
	}

	return note.(*models.News), nil
}

func (p *postgresDB) EditNote(ctx context.Context, note *models.News) error {
	if err := p.reformDB.WithContext(ctx).Save(note); err != nil {
		return fmt.Errorf("save/update news: %w", err)
	}
	return nil
}

func (p *postgresDB) GetNotes(ctx context.Context) ([]models.News, error) {
	tx, err := p.reformDB.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("create transaction: %w", err)
	}
	defer helpers.ExecuteWithLogError(tx.Rollback, p.log)

	notes, err := tx.SelectAllFrom(models.NewsTable, "")
	if err != nil {
		return nil, fmt.Errorf("get notes from db: %w", err)
	}

	res := make([]models.News, 0, len(notes))
	for _, noteStruct := range notes {
		note := *noteStruct.(*models.News)

		categories, err := tx.SelectAllFrom(models.NewsCategoriesTable, "WHERE news_id = $1", note.ID)
		if err != nil {
			return nil, fmt.Errorf("create transaction: %w", err)
		}

		note.Categories = []int64{}
		for _, category := range categories {
			note.Categories = append(note.Categories, category.(*models.NewsCategories).CategoryID)
		}

		res = append(res, note)
	}

	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("commit transaction: %w", err)
	}

	return res, nil
}