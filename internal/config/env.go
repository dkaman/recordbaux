package config

import (
	"strings"

	"github.com/knadh/koanf/v2"

	env "github.com/knadh/koanf/providers/env/v2"
)

func normalizeEnvPath(s string, prefix string) string {
	return strings.Replace(strings.ToLower(strings.TrimPrefix(s, prefix)), "_", ".", -1)
}

func WithEnv(prefix string) ConfigOption {
	return func(k *koanf.Koanf) error {
		pre := strings.ToUpper(prefix + "_")

		k.Load(env.Provider(".", env.Opt{
			Prefix: pre,
			TransformFunc: func(k, v string) (string, any) {
				k = strings.ReplaceAll(strings.ToLower(strings.TrimPrefix(k, pre)), "_", ".")
				if strings.Contains(v, " ") {
					return k, strings.Split(v, " ")
				}
				return k, v
			},
		}), nil)

		return nil
	}
}
