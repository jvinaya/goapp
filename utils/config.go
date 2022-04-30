package utils

import (
	"fmt"
	"time"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

//config store all configuration of application
//the value are read by viper from config file or evn variable
type Config struct {
	DBDriver            string        `env:"DB_DRIVER"`
	DBSource            string        `env:"DB_SOURCE"`
	ServerAddress       string        `env:"SERVER_ADDRESS"`
	TokenSymmetricKey   string        `env:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration time.Duration `env:"ACCESS_TOKEN_DURATION"`
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig(configPath string) (config Config, err error) {
	err = godotenv.Load(configPath)
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		fmt.Printf("%+v\n", err)
	}

	fmt.Printf("%+v\n", cfg)
	return cfg, err
}
