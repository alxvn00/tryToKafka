package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	ConsumerConfig
}

type ConsumerConfig struct {
	Addresses []string `env:"KAFKA_BOOTSTRAP_SERVERS, required"`
}

func NewConfig() (*Config, error) {
	var config Config
	err := cleanenv.ReadEnv(&config)
	if err != nil {
		return nil, fmt.Errorf("failed to read env: %w", err)
	}
	return &config, nil
}
