package config

import (
	"github.com/go-playground/validator/v10"
	"github.com/knadh/koanf"
)

type Config struct {
	Auth        Auth        `koanf:"auth" validate:"required"`
	Credentials Credentials `koanf:"credentials" validate:"required"`
	Server      Server      `koanf:"server" validate:"required"`
	Store       Store       `koanf:"store" validate:"required"`
}

type Store struct {
	Interval int `koanf:"interval" validate:"required"`
}

type Auth struct {
	Type         string   `koanf:"type" validate:"oneof=noop user-agent header"`
	UserAgents   []string `koanf:"user_agents" validate:"required_if=Type user-agent"`
	HeaderValues []string `koanf:"header_values" validate:"required_if=Type header"`
	HeaderName   string   `koanf:"header_name" validate:"required_if=Type header"`
}

type Credentials struct {
	Id     string   `koanf:"id" validate:"required"`
	Scopes []string `koanf:"scopes" validate:"required"`
	Secret string   `koanf:"secret" validate:"required"`
	Url    string   `koanf:"url" validate:"required,url"`
}

type Server struct {
	Addr string `koanf:"addr" validate:"required"`
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
