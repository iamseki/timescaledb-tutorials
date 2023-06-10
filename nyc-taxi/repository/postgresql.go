package repository

import (
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type PostgreRepository struct {
	db     *sqlx.DB
	logger *zap.Logger
}

func NewPostgreSQL(config *Config, logger *zap.Logger) (Repository, error) {
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

func (repository *PostgreRepository) RidesByDay(date string) ([]RidesByDayResponse, error) {
	response := []RidesByDayResponse{}
	err := repository.db.Select(&response, fmt.Sprintf(`
		SELECT date_trunc('day', pickup_datetime) as day,
		COUNT(*)
		FROM rides
		WHERE pickup_datetime < '%v'
		GROUP BY day
		ORDER BY day;
	`, date))
	if err != nil {
		repository.logger.Error(fmt.Sprintf("Error on query RidesByDay %v", date), zap.Error(err))
		return nil, err
	}
	return response, nil
}

func (repository *PostgreRepository) AverageFareByDay(date string) ([]AverageFareByDayResponse, error) {
	response := []AverageFareByDayResponse{}

	err := repository.db.Select(&response, fmt.Sprintf(`
		SELECT date_trunc('day', pickup_datetime)
		AS day, avg(fare_amount)
		FROM rides
		WHERE pickup_datetime < '%v'
		GROUP BY day
		ORDER BY day;
	`, date))
	if err != nil {
		repository.logger.Error(fmt.Sprintf("Error on query AverageFare %v", date), zap.Error(err))
		return nil, err
	}
	return response, nil
}

func (repository *PostgreRepository) Close() error {
	return nil
}
