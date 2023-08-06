package repository

type config struct {
	URL                    string `envconfig:"DATABASE_URI" default:"postgres://postgres:password@localhost:5432/postgres"`
	MaxPoolSize            int    `envconfig:"MAX_POOL_SIZE" default:"100"`
	MinPoolSize            int    `envconfig:"MIN_POOL_SIZE" default:"2"`
	IdleMaxLifeTimeSeconds int    `envconfig:"IDLE_MAX_LIFE_TIME_SECONDS" default:"3600"`
}

type RecentBlock struct {
	Time    string  `db:"time" json:"time"`
	Hash    string  `db:"hash" json:"hash"`
	BlockID int     `db:"block_id" json:"blockID"`
	FeeUSD  float64 `db:"fee_usd" json:"feeUSD"`
}

type Repository interface {
	ListRecentBlocks(limit int) ([]RecentBlock, error)
	Close() error
}
