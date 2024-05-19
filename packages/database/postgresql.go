package database

import (
	"darkness8129/currency-converter/packages/logging"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgreSQLDatabase struct {
	db     *gorm.DB
	logger *logging.Logger
}

type Options struct {
	User     string
	Password string
	Database string
	Port     string
	Host     string
	Logger   *logging.Logger
}

func NewPostgreSQLDatabase(opt Options) (*PostgreSQLDatabase, error) {
	logger := opt.Logger.Named("PostgreSQLDatabase")

	dsn := fmt.Sprintf(
		"user=%s password=%s dbname=%s port=%s host=%s",
		opt.User, opt.Password, opt.Database, opt.Port, opt.Host,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		PrepareStmt: true,
	})
	if err != nil {
		logger.Error("failed to connect to DB", "err", err)
		return nil, fmt.Errorf("failed to connect to DB: %w", err)
	}

	// needed for automatic creating IDs for new records
	err = db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`).Error
	if err != nil {
		logger.Error("failed to create uuid-ossp extension", "err", err)
		return nil, fmt.Errorf("failed to create uuid-ossp extension: %w", err)
	}

	return &PostgreSQLDatabase{
		db:     db,
		logger: logger,
	}, nil
}

func (p *PostgreSQLDatabase) DB() *gorm.DB {
	return p.db
}

func (p *PostgreSQLDatabase) Close() error {
	logger := p.logger.Named("Close")

	db, err := p.db.DB()
	if err != nil {
		logger.Error("failed to get db", "err", err)
		return fmt.Errorf("failed to get db: %w", err)
	}

	err = db.Close()
	if err != nil {
		logger.Error("failed to close postgresql connection", "err", err)
		return fmt.Errorf("failed to close postgresql connection: %w", err)
	}

	logger.Info("successfully closed connection to DB")
	return nil
}
