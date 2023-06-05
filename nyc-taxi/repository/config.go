package repository

type Config struct {
	Database DatabaseConfig
}

type DatabaseConfig struct {
	URL                    string `envconfig:"DATABASE_URI" default:"postgres://postgres:password@localhost:5432/nyc_taxi_cab"`
	MaxPoolSize            int    `envconfig:"MAX_POOL_SIZE" default:"100"`
	MinPoolSize            int    `envconfig:"MIN_POOL_SIZE" default:"2"`
	IdleMaxLifeTimeSeconds int    `envconfig:"IDLE_MAX_LIFE_TIME_SECONDS" default:"3600"`
}
