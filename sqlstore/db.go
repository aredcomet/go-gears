package sqlstore

import (
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// ConnectToDb opens and returns a database connection using the
// provided `dsn` (data source name) according to the rules defined in `config`.
// The provided `logger` is used to log any error during the connection process.
//
// The function uses the "gorm" package to open a connection to
// a Postgres database, then tries to get the underlying SQL database
// instance and sends a ping to ensure the connection is working.
//
// If everything is alright, it returns the `gorm.DB` instance
// which represents a session with the database. In case of any error,
// it logs the error message and returns `nil` for the `gorm.DB`
// instance along with the error.
func ConnectToDb(dsn string, config *gorm.Config, logger *logrus.Logger) (*gorm.DB, error) {
	// Attempt to open a new db connection
	db, err := gorm.Open(postgres.Open(dsn), config)
	if err != nil {
		// Log the error and return if db connection fails
		logger.Error("Failed to open DB connection: ", err)
		return nil, err
	}

	// Attempt to get the SQL db instance
	sqlDB, err := db.DB()
	if err != nil {
		// Log the error and return if getting SQL db instance fails
		logger.Error("Failed to get SQL DB: ", err)
		return nil, err
	}

	// Ping the database to ensure connection is working
	if err = sqlDB.Ping(); err != nil {
		// Log the error and return if ping fails
		logger.Error("Failed to ping DB: ", err)
		return nil, err
	}

	// If no errors, return the db instance
	return db, nil
}
