package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"time"
)

var Config *RateLimiterConfig

type RateLimiterConfig struct {
	Redis struct {
		Address  string `yaml:"address"`
		Password string `yaml:"password"`
		DB       int    `yaml:"db"`
	} `yaml:"redis"`

	RateLimiter struct {
		DefaultRateLimit struct {
			Rate          int           `yaml:"rate"`
			WindowSeconds time.Duration `yaml:"window_seconds"`
		} `yaml:"default_rate_limit"`

		UserTypes map[string]struct {
			KeyPrefix string `yaml:"key_prefix"`
			RateLimit struct {
				Rate          int           `yaml:"rate"`
				WindowSeconds time.Duration `yaml:"window_seconds"`
			} `yaml:"rate_limit"`
		} `yaml:"user_types"`
	} `yaml:"rate_limiter"`
}

func LoadConfig(filename string) (*RateLimiterConfig, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Println("error closing the config file")
		}
	}(file)

	var config RateLimiterConfig
	if err := yaml.NewDecoder(file).Decode(&config); err != nil {
		return nil, fmt.Errorf("failed to decode config file: %w", err)
	}

	Config = &config

	return Config, nil
}
