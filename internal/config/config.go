package config

import (
	"fmt"
	"sync"

	"github.com/spf13/viper"
)

type Provider interface {
	GetString(key string) string
	GetInt(key string) int
}

type Config struct {
	*viper.Viper
}

var (
	_             Provider = (*Config)(nil)
	once          sync.Once
	defaultConfig *Config
)

func defaultLoadConfig() *Config {
	v := viper.New()

	v.SetDefault("PORT", 8080)
	v.SetDefault("LOG_LEVEL", "debug")
	v.SetDefault("GIN_MODE", "release")

	v.SetDefault("DB_READ_USER", "user")     // non production approach
	v.SetDefault("DB_READ_PASSWORD", "pass") // non production approach
	v.SetDefault("DB_READ_HOST", "mysql")
	v.SetDefault("DB_READ_PORT", 3306)
	v.SetDefault("DB_READ_NAME", "users")
	v.SetDefault("DB_READ_MAX_OPEN_CONN", 100)
	v.SetDefault("DB_READ_MAX_IDLE_CONN", 50)

	v.SetDefault("DB_WRITE_USER", "user")     // non production approach
	v.SetDefault("DB_WRITE_PASSWORD", "pass") // non production approach
	v.SetDefault("DB_WRITE_HOST", "mysql")
	v.SetDefault("DB_WRITE_PORT", 3306)
	v.SetDefault("DB_WRITE_NAME", "users")
	v.SetDefault("DB_WRITE_MAX_OPEN_CONN", 150)
	v.SetDefault("DB_WRITE_MAX_IDLE_CONN", 75)

	v.SetConfigName("config")
	v.SetConfigType("yaml")

	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil || len(v.AllSettings()) == 0 {
		v.AddConfigPath(".")
		v.SetConfigName("config")
		v.SetConfigType("yaml")

		if err = v.ReadInConfig(); err != nil {
			fmt.Printf("\n=> no configuration file (.config.yaml) found, using defaults only\n\n")
		}
	}

	return &Config{v}
}

// Load returns application config created only once
func Load() *Config {
	once.Do(func() {
		defaultConfig = defaultLoadConfig()
	})

	return defaultConfig
}
