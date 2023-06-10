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

func (repository *PostgreRepository) RidesByFareType(date string) ([]RidesByFareTypeResponse, error) {
	response := []RidesByFareTypeResponse{}

	err := repository.db.Select(&response, fmt.Sprintf(`
		SELECT rates.description, COUNT(vendor_id) as num_trips
		FROM rides
		INNER JOIN rates ON rides.rate_code = rates.rate_code
		WHERE pickup_datetime < '%v'
		GROUP BY rates.description
		ORDER BY LOWER(rates.description);
	`, date))
	if err != nil {
		repository.logger.Error(fmt.Sprintf("Error on query RidesByFareType %v", date), zap.Error(err))
		return nil, err
	}
	return response, nil
}

func (repository *PostgreRepository) RidesByAirportsCodeResponse(date string, airportCodes string) ([]RidesByAirportsCodeResponse, error) {
	response := []RidesByAirportsCodeResponse{}

	query := fmt.Sprintf(`
	SELECT rates.description,
		COUNT(vendor_id) AS num_trips,
		AVG(dropoff_datetime - pickup_datetime) AS avg_trip_duration,
		AVG(total_amount) AS avg_total,
		AVG(passenger_count) AS avg_passengers
	FROM rides
	INNER JOIN rates on rates.rate_code = rides.rate_code
	WHERE rides.rate_code IN(%v) AND pickup_datetime < '%v'
	GROUP BY rates.description
	ORDER BY rates.description;
	`, airportCodes, date)
	repository.logger.Info("Query", zap.String("query", query))

	err := repository.db.Select(&response, query)
	if err != nil {
		repository.logger.Error(fmt.Sprintf("Error on query RidesByAirportsCodeResponse date => %v, airportCodes => %v", airportCodes, date), zap.Error(err))
		return nil, err
	}
	return response, nil
}

func (repository *PostgreRepository) Close() error {
	return nil
}
