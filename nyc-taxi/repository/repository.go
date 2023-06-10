package repository

type RidesByDaySinceResponse struct {
	Day   string `db:"day"`
	Count int    `db:"count"`
}

type AverageFareSinceResponse struct {
	Day     string  `db:"day"`
	Average float64 `db:"avg"`
}

type Repository interface {
	RidesByDaySince(date string) ([]RidesByDaySinceResponse, error)
	AverageFareSince(date string) ([]AverageFareSinceResponse, error)
	Close() error
}
