package repository

type RidesByDaySinceResponse struct {
	Day   string
	Count int
}

type Repository interface {
	RidesByDaySince(date string) ([]RidesByDaySinceResponse, error)
	Close() error
}
