package config

import (
	"github.com/knadh/koanf/providers/confmap"
	"strings"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/providers/env"
)

type Builder struct {
	k    *koanf.Koanf
	args Args
}

func NewBuilder(k *koanf.Koanf, args Args) *Builder {
	return &Builder{k: k, args: args}
}

func (c *Builder) unmarshalConfigToStruct(config *Config) error {
	return c.k.Unmarshal("", &config)
}

func (c *Builder) loadConfig() (err error) {
	err = c.loadConfigFromArgs()
	if err != nil {
		return
	}

	err = c.loadConfigFromEnvironment()
	if err != nil {
		return
	}

	return nil
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

func (c *Builder) loadConfigFromArgs() error {
	return c.k.Load(confmap.Provider(map[string]interface{}{
		"auth.type":                *c.args.authType,
		"clientcredentials.id":     *c.args.clientId,
		"clientcredentials.url":    *c.args.tokenUrl,
		"clientcredentials.secret": *c.args.clientSecret,
		"clientcredentials.scopes": *c.args.scopes,
		"server.Cert":              *c.args.cert,
		"server.Key":               *c.args.key,
		"server.Tls":               *c.args.tls,
	}, "."), nil)
}
