package config

import (
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/spf13/viper"
)

type Config struct {
	Port string `mapstructure:"port" env:"PORT"`

	SQL    DB           `mapstructure:"sql" env-prefix:"POSTGRES_"`
	Search SearchConfig `mapstructure:"search" env-prefix:"SEARCH_"`
}

type DB struct {
	Host     string `mapstructure:"host" env:"HOST"`
	Port     int    `mapstructure:"port" env:"PORT"`
	User     string `mapstructure:"user" env:"USER"`
	Password string `mapstructure:"password" env:"PASSWORD"`
	Name     string `mapstructure:"name" env:"NAME"`
}

type SearchConfig struct {
	BaseURL             string        `mapstructure:"base_url" env:"BASE_URL"`
	IndexTimeout        time.Duration `mapstructure:"index_timeout" env:"INDEX_TIMEOUT"`
	SearchTimeout       time.Duration `mapstructure:"search_timeout" env:"SEARCH_TIMEOUT"`
	SchedulerInterval   time.Duration `mapstructure:"scheduler_interval" env:"SCHEDULER_INTERVAL"`
	SchedulerBatchSize  int           `mapstructure:"scheduler_batch_size" env:"SCHEDULER_BATCH_SIZE"`
	SchedulerMaxRetries int           `mapstructure:"scheduler_max_retries" env:"SCHEDULER_MAX_RETRIES"`
}

func LoadConfigFromFile(path string) (*Config, error) {
	config := new(Config)
	viper.SetConfigFile(path)
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	err = viper.Unmarshal(config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func MustLoadConfigFromFile(path string) *Config {
	config, err := LoadConfigFromFile(path)
	if err != nil {
		panic(err)
	}

	return config
}

func LoadConfigFromEnv() (*Config, error) {
	config := new(Config)
	err := cleanenv.ReadEnv(config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func MustLoadConfigFromEnv() *Config {
	config, err := LoadConfigFromEnv()
	if err != nil {
		panic(err)
	}

	return config
}
