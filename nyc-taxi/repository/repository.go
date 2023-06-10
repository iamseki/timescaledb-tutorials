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

type Repository interface {
	RidesByDay(date string) ([]RidesByDayResponse, error)
	AverageFareByDay(date string) ([]AverageFareByDayResponse, error)
	RidesByFareType(date string) ([]RidesByFareTypeResponse, error)
	Close() error
}
