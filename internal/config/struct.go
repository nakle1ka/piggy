package config

type Config struct {
	App      App      `mapstructure:",squash"`
	Database Database `mapstructure:",squash"`
	Redis    Redis    `mapstructure:",squash"`
	JWT      JWT      `mapstructure:",squash"`
}

type App struct {
	Host string `mapstructure:"APP_HOST"`
	Port int    `mapstructure:"APP_PORT"`
	Env  string `mapstructure:"APP_ENV"`
}

type Database struct {
	Host     string `mapstructure:"DB_HOST"`
	Port     int    `mapstructure:"DB_PORT"`
	Username string `mapstructure:"DB_USERNAME"`
	Password string `mapstructure:"DB_PASSWORD"`
	Name     string `mapstructure:"DB_NAME"`
	SSLMode  bool   `mapstructure:"DB_SSL_MODE"`
}

type Redis struct {
	Host     string `mapstructure:"REDIS_HOST"`
	Port     int    `mapstructure:"REDIS_PORT"`
	Password string `mapstructure:"REDIS_PASSWORD"`
	DB       int    `mapstructure:"REDIS_DB"`
}

type JWT struct {
	SecretKey            string `mapstructure:"JWT_SECRET_KEY"`
	RefreshTokenLifeTime int    `mapstructure:"JWT_REFRESH_TOKEN_LIFE_TIME"`
	AccessTokenLifeTime  int    `mapstructure:"JWT_ACCESS_TOKEN_LIFE_TIME"`
}
