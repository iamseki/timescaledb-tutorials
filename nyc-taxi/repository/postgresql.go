package repository

import (
	"time"

	"github.com/iamseki/timescaledb-tutorials/nyc-taxi/api"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type PostgreRepository struct {
	db     *sqlx.DB
	logger *zap.Logger
}

func NewPostgreSQL(config *api.Config, logger *zap.Logger) (Repository, error) {
	db, err := sqlx.Connect("pgx", config.Database.URL)
	if err != nil {
		logger.Error("sqlx Connect error", zap.Error(err), zap.Any("config", config.Database))
	}

	db.SetMaxIdleConns(config.Database.MinPoolSize)
	db.SetMaxOpenConns(config.Database.MaxPoolSize)
	db.SetConnMaxLifetime(time.Duration(time.Duration(config.Database.IdleMaxLifeTimeSeconds).Seconds()))

	logger.Info("PostgreSQL Connected")
	return &PostgreRepository{db, logger}, nil
}

func (r *PostgreRepository) Close() error {
	return nil
}
