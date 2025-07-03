package config

import (
	"github.com/knadh/koanf/v2"
)

type ConfigOption func(*koanf.Koanf) error

type Config struct {
	*koanf.Koanf
}

func New(opts ...ConfigOption) (*Config, error) {
	c := koanf.New(".")

	for _, o := range opts {
		err := o(c)
		if err != nil {
			return nil, err
		}
	}

	return &Config{c}, nil
}
