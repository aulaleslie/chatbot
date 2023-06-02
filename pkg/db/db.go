package db

import (
	"chatbot/config"
	"errors"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func ConnectDB(databaseConfig *config.DatabaseConf) (*gorm.DB, IDatabaseAdapter, error) {
	log.Debug().Msg("[ConnectDB] connecting to DB...")

	var db *gorm.DB
	var handler IDatabaseAdapter
	var err error
	switch databaseConfig.Type {
	case "postgres":
		db, handler, err = connectPostgresDb(databaseConfig.Postgres.ConnectionString)
	default:
		err = errors.New("no database type provided")
	}

	if err != nil {
		log.Warn().Err(err).Msgf("[ConnectDB] there is an error while connecting to DB: %d", err.Error())
		return nil, nil, err
	}

	log.Debug().Msg("[ConnectDB] connected...")
	return db, handler, nil
}
