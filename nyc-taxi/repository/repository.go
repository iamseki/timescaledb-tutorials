package repository

type RidesByDayResponse struct {
	Day   string `db:"day" json:"day"`
	Count int    `db:"count" json:"count"`
}

type AverageFareByDayResponse struct {
	Day     string  `db:"day" json:"day"`
	Average float64 `db:"avg" json:"average"`
}

type RidesByFareTypeResponse struct {
	Description string `db:"description" json:"description"`
	TotalTrips  int    `db:"num_trips" json:"totalTrips"`
}

type RidesByAirportsCodeResponse struct {
	Description       string  `db:"description" json:"description"`
	TotalTrips        string  `db:"num_trips" json:"totalTrips"`
	AverageDuration   string  `db:"avg_trip_duration" json:"averageDuration"`
	AverageTotal      float64 `db:"avg_total" json:"averageTotal"`
	AveragePassengers string  `db:"avg_passengers" json:"averagePassengers"`
}

type RidesByTimeBucketResponse struct {
	TimeBucket string `db:"time_bucket" json:"timeBucket"`
	Rides      int    `db:"count" json:"rides"`
}

type Repository interface {
	RidesByDay(date string) ([]RidesByDayResponse, error)
	AverageFareByDay(date string) ([]AverageFareByDayResponse, error)
	RidesByFareType(date string) ([]RidesByFareTypeResponse, error)
	RidesByAirportsCodeResponse(date string, airportCodes string) ([]RidesByAirportsCodeResponse, error)
	RidesByTimeBucket(date string, interval string) ([]RidesByTimeBucketResponse, error)
	Close() error
}
