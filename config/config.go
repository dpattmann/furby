package config

import (
	"strings"

	"github.com/dpattmann/furby/oauth2"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
)

type Config struct {
	ClientCredentialSettings oauth2.ClientCredentialSettings `koanf:"clientcredentialsettings"`
}

type Builder struct {
	k *koanf.Koanf
}

func NewConfig() (config *Config, err error) {
	var k = koanf.New(".")

	builder := NewBuilder(k)

	err = builder.loadConfigFromJsonFile()
	if err != nil {
		return
	}

	err = builder.loadConfigFromEnvironment()
	if err != nil {
		return
	}

	config = new(Config)
	err = builder.unmarshalConfigToStruct(config)

	return
}

func NewBuilder(k *koanf.Koanf) *Builder {
	return &Builder{k: k}
}

func (c *Builder) unmarshalConfigToStruct(config *Config) error {
	return c.k.Unmarshal("", &config)
}

func (c *Builder) loadConfigFromEnvironment() error {
	return c.k.Load(env.Provider("FURBY_", ".", func(s string) string {
		return strings.Replace(
			strings.ToLower(strings.TrimPrefix(s, "FURBY_")), "_", ".", -1)
	}), nil)
}

func (c *Builder) loadConfigFromJsonFile() error {
	return c.k.Load(file.Provider("config/config.json"), json.Parser())
}
