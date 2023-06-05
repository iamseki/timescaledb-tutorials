package api

type Config struct {
	App      AppConfig
	Database DatabaseConfig
}

type AppConfig struct {
	Name            string `envconfig:"APP_NAME" default:"nyc_taxi_api"`
	LogLevel        string `envconfig:"APP_LOG_LEVEL" default:"info"`
	Port            string `envconfig:"APP_PORT" default:"5000"`
	HTTPReadTimeout int    `envconfig:"HTTP_READ_TIMEOUT" default:"120"`
}

type DatabaseConfig struct {
	URL                    string `envconfig:"DATABASE_URI" default:"postgres://postgres:password@localhost:5432/nyc_taxi_cab"`
	MaxPoolSize            int    `envconfig:"MAX_POOL_SIZE" default:"100"`
	MinPoolSize            int    `envconfig:"MIN_POOL_SIZE" default:"2"`
	IdleMaxLifeTimeSeconds int    `envconfig:"IDLE_MAX_LIFE_TIME_SECONDS" default:"3600"`
}
