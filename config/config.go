package config

import (
	"github.com/dpattmann/furby/oauth2"

	"github.com/go-playground/validator/v10"
	"github.com/knadh/koanf"
)

type Config struct {
	ClientCredentials oauth2.ClientCredentials `koanf:"clientcredentials" validate:"required"`
}

func (c *Config) validate() (err error) {
	validate := validator.New()
	err = validate.Struct(c.ClientCredentials)

	return
}

func NewConfig() (config *Config, err error) {
	var k = koanf.New(".")

	builder := NewBuilder(k)

	err = builder.loadConfig()
	if err != nil {
		return
	}

	config = new(Config)
	err = builder.unmarshalConfigToStruct(config)

	if err != nil {
		return
	}

	err = config.validate()

	return
}
