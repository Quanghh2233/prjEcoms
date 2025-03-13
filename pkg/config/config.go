package config

import (
	"fmt"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	Environment  string        `mapstructure:"ENVIRONMENT"`
	DBDriver     string        `mapstructure:"DB_DRIVER"`
	DBSource     string        `mapstructure:"DB_SOURCE"`
	ServerPort   string        `mapstructure:"PORT"`
	JWTSecret    string        `mapstructure:"JWT_SECRET"`
	JWTExpiresIn time.Duration `mapstructure:"JWT_EXPIRES_IN"`
}

// LoadConfig reads configuration from environment variables
func LoadConfig() (Config, error) {
	// Try to load .env file if it exists
	_ = godotenv.Load()

	viper.AutomaticEnv()

	// Set default values
	viper.SetDefault("ENVIRONMENT", "development")
	viper.SetDefault("DB_DRIVER", "postgres")
	viper.SetDefault("DB_SOURCE", "postgresql://qhh:2203@localhost:5432/ecommerce?sslmode=disable")
	viper.SetDefault("PORT", "8080")
	viper.SetDefault("JWT_SECRET", "your_secret_key")
	viper.SetDefault("JWT_EXPIRES_IN", time.Hour*24)

	var config Config

	config = Config{
		Environment:  viper.GetString("ENVIRONMENT"),
		DBDriver:     viper.GetString("DB_DRIVER"),
		DBSource:     viper.GetString("DB_SOURCE"),
		ServerPort:   viper.GetString("PORT"),
		JWTSecret:    viper.GetString("JWT_SECRET"),
		JWTExpiresIn: viper.GetDuration("JWT_EXPIRES_IN"),
	}

	// Validate required configurations
	if config.JWTSecret == "your_secret_key" {
		fmt.Println("Warning: Using default JWT secret. This is not secure for production.")
	}

	return config, nil
}
