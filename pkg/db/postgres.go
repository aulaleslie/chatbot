package db

import (
	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func connectPostgresDb(dsn string) (*gorm.DB, IDatabaseAdapter, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, nil, err
	}

	handler := newPostgresDbHandler(db)
	return db, handler, nil
}

func newPostgresDbHandler(postgresDb *gorm.DB) IDatabaseAdapter {
	return &PostgresDbHandler{
		DB: postgresDb,
	}
}

func (handler *PostgresDbHandler) Any(model interface{}, key, value string) (bool, error) {
	log.Debug().Msg("[PostgresDBHandler] query any")
	result := handler.DB.First(&model, key+" = ?", value)

	if result.Error != nil {
		return false, result.Error
	}

	return true, nil
}

func (handler *PostgresDbHandler) Insert(model interface{}) (int, error) {
	log.Debug().Msg("[PostgresDBHandler] query insert")
	result := handler.DB.Create(model)

	if result.Error != nil {
		return 0, result.Error
	}

	return int(result.RowsAffected), nil
}
