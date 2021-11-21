package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

const (
	// Name of the application
	AppName = "todo-go"
)

type Config struct {
	*viper.Viper
}

func New() *Config {
	config := &Config{
		Viper: viper.New(),
	}

	// Set default configurations
	config.setDefaults()

	// SetConfigName sets name for the config file.
	// Does not include extension
	// Select the .env file
	config.SetConfigName(config.GetString("APP_CONFIG_NAME"))
	config.SetConfigType("dotenv")
	// AddConfigPath adds a path for Viper to search for the config file in.
	config.AddConfigPath(config.GetString("APP_CONFIG_PATH"))

	// Automatically refresh environment variables
	config.AutomaticEnv()

	// Read configuration
	if err := config.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			fmt.Println("failed to read configuration:", err.Error())
			os.Exit(1)
		}
	}
	return config
}

func (c *Config) setDefaults() {
	// Set default App configuration
	c.SetDefault("APP_ADDR", ":8080")
	c.SetDefault("APP_CONFIG_NAME", ".env")
	c.SetDefault("APP_CONFIG_PATH", ".")

	// Set default database options
	c.SetDefault("DATABASE_TYPE", "memory") // Availables: "memory"
	c.SetDefault("DATABASE_HOST", "")
	c.SetDefault("DATABASE_PORT", "")
	c.SetDefault("DATABASE_USERNAME", "")
	c.SetDefault("DATABASE_PASSWORD", "")
	c.SetDefault("DATABASE_NAME", "")
}
