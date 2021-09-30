package config

import (
	"github.com/go-playground/validator/v10"
	"github.com/knadh/koanf"
)

type Config struct {
	Auth              Auth              `koanf:"auth" validate:"required"`
	ClientCredentials ClientCredentials `koanf:"client_credentials" validate:"required"`
	Server            Server            `koanf:"server" validate:"required"`
}

type Auth struct {
	Type       string   `koanf:"type"`
	UserAgents []string `koanf:"user_agents" validate:"required_if=Type user-agent"`
}

type ClientCredentials struct {
	Id     string   `koanf:"id" validate:"required"`
	Scopes []string `koanf:"scopes" validate:"required"`
	Secret string   `koanf:"secret" validate:"required"`
	Url    string   `koanf:"url" validate:"required,url"`
}

type Server struct {
	Cert string `koanf:"cert" validate:"required_if=Tls true"`
	Key  string `koanf:"key" validate:"required_if=Tls true"`
	Tls  bool   `koanf:"tls"`
}

func (c *Config) validate() (err error) {
	validate := validator.New()
	err = validate.Struct(c)

	return
}

func NewConfig(path string) (config *Config, err error) {
	var k = koanf.New(".")

	builder := NewBuilder(k, path)

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
