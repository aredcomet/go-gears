package sqlstore

import (
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectToDb(dsn string, config *gorm.Config, logger *logrus.Logger) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), config)
	if err != nil {
		logger.Error("Failed to open DB connection: ", err)
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		logger.Error("Failed to get SQL DB: ", err)
		return nil, err
	}

	// ping
	if err = sqlDB.Ping(); err != nil {
		logger.Error("Failed to ping DB: ", err)
		return nil, err
	}

	return db, nil
}
