package config

import (
	"github.com/knadh/koanf/v2"
)

var k = koanf.New(".")

type ConfigOption func(*koanf.Koanf) error

type Config struct {
	config *koanf.Koanf
}

func New(opts ...ConfigOption) (*Config, error) {
	c := koanf.New(".")

	for _, o := range opts {
		err := o(c)
		if err != nil {
			return nil, err
		}
	}

	return &Config{
		config: c,
	}, nil
}

func (c *Config) String(key string) (string, error) {
	return c.config.String(key), nil
}
