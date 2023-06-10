package repository

type RidesByDayResponse struct {
	Day   string `db:"day"`
	Count int    `db:"count"`
}

type AverageFareByDayResponse struct {
	Day     string  `db:"day"`
	Average float64 `db:"avg"`
}

type RidesByFareTypeResponse struct {
	Description string `db:"description"`
	TotalTrips  int    `db:"num_trips"`
}

type Repository interface {
	RidesByDay(date string) ([]RidesByDayResponse, error)
	AverageFareByDay(date string) ([]AverageFareByDayResponse, error)
	RidesByFareType(date string) ([]RidesByFareTypeResponse, error)
	Close() error
}
