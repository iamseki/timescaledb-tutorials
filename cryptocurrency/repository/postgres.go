package repository

import (
	"time"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/kelseyhightower/envconfig"
	"go.uber.org/zap"
)

type PGSQLRepository struct {
	db     *sqlx.DB
	logger *zap.Logger
}

func NewPostgreSQL(logger *zap.Logger) (Repository, error) {
	dbConfig := &config{}
	if err := envconfig.Process("cryptocurrency-apy", dbConfig); err != nil {
		logger.Fatal("Error loading DATABASE envconfig", zap.Any("dbConfig", dbConfig))
	}

	db, err := sqlx.Connect("pgx", dbConfig.URL)
	if err != nil {
		logger.Error("sqlx Connect error", zap.Error(err), zap.Any("config", dbConfig))
	}

	db.SetMaxIdleConns(dbConfig.MinPoolSize)
	db.SetMaxOpenConns(dbConfig.MaxPoolSize)
	db.SetConnMaxLifetime(time.Duration(time.Duration(dbConfig.IdleMaxLifeTimeSeconds).Seconds()))

	logger.Info("PostgreSQL Connected")
	return &PGSQLRepository{db, logger}, nil
}

func (repository *PGSQLRepository) Close() error {
	return nil
}
