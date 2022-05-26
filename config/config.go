package config

import (
	"github.com/caarlos0/env"
)

type EngineConfig struct {
	Host string `env:"HOST" envDefault:"0.0.0.0"`
	Port string `env:"PORT" envDefault:":8080"`
}

func NewEngineConfig() (*EngineConfig, error) {
	var engineConfig EngineConfig

	if err := env.Parse(&engineConfig); err != nil {
		return nil, err
	}

	return &engineConfig, nil
}
