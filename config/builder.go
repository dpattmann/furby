package config

import (
	"strings"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/providers/env"
)

type Builder struct {
	k *koanf.Koanf
}

func NewBuilder(k *koanf.Koanf) *Builder {
	return &Builder{k: k}
}

func (c *Builder) unmarshalConfigToStruct(config *Config) error {
	return c.k.Unmarshal("", &config)
}

func (c *Builder) loadConfig() (err error) {
	return c.loadConfigFromEnvironment()
}

func (c *Builder) loadConfigFromEnvironment() error {
	return c.k.Load(env.ProviderWithValue("FURBY_", ".", func(s string, v string) (string, interface{}) {
		key := strings.Replace(
			strings.ToLower(strings.TrimPrefix(s, "FURBY_")), "_", ".", -1)

		if strings.Contains(v, " ") {
			return key, strings.Split(v, " ")
		}

		return key, v
	}), nil)
}
