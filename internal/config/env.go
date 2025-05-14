package config

import (
	"strings"

	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/v2"
)

func normalizeEnvPath(s string, prefix string) string {
	return strings.Replace(strings.ToLower(strings.TrimPrefix(s, prefix)), "_", ".", -1)
}

func WithEnv(prefix string) ConfigOption {
	return func(k *koanf.Koanf) error {
		k.Load(env.Provider("MYVAR_", ".", func(s string) string {
			return normalizeEnvPath(s, prefix)
		}), nil)
		return nil
	}
}
