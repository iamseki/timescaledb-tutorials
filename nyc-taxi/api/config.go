package api

type Config struct {
	App AppConfig
}

type AppConfig struct {
	Name            string `envconfig:"APP_NAME" default:"nyc_taxi_api"`
	LogLevel        string `envconfig:"APP_LOG_LEVEL" default:"info"`
	Port            string `envconfig:"APP_PORT" default:"5000"`
	HTTPReadTimeout int    `envconfig:"HTTP_READ_TIMEOUT" default:"120"`
}
