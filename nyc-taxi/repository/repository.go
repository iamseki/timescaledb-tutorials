package repository

type RidesByDayResponse struct {
	Day   string `db:"day"`
	Count int    `db:"count"`
}

type AverageFareByDayResponse struct {
	Day     string  `db:"day"`
	Average float64 `db:"avg"`
}

type Repository interface {
	RidesByDay(date string) ([]RidesByDayResponse, error)
	AverageFareByDay(date string) ([]AverageFareByDayResponse, error)
	Close() error
}
