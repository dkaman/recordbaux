package config

import (
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

func WithFile(path string) ConfigOption {
	return func(k *koanf.Koanf) error {
		if err := k.Load(file.Provider(path), yaml.Parser()); err != nil {
			return err
		}
		return nil
	}
}
