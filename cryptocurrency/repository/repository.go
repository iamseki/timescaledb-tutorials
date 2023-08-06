package repository

type config struct {
	URL                    string `envconfig:"DATABASE_URI" default:"postgres://postgres:password@localhost:5432/postgres"`
	MaxPoolSize            int    `envconfig:"MAX_POOL_SIZE" default:"100"`
	MinPoolSize            int    `envconfig:"MIN_POOL_SIZE" default:"2"`
	IdleMaxLifeTimeSeconds int    `envconfig:"IDLE_MAX_LIFE_TIME_SECONDS" default:"3600"`
}

type Repository interface {
	ListRecentBlocks(limit int) ([]RecentBlock, error)
	TxVolumeXFeesByHour(days int) ([]TxVolume, error)
	TxVolumeXUSDByHour(days int) ([]TxVolume, error)
	BlockVolumeTxXMiningFee(days int) ([]BlockVolume, error)
	BlockVolumeXMiningFee(days int) ([]BlockVolume, error)
	MinerReveueByHour(days int) ([]MinerRevenue, error)
	Close() error
}
