package config

import (
	"github.com/go-playground/validator/v10"
	"github.com/knadh/koanf"
)

var (
	defaultConfig = map[string]interface{}{
		"auth.type":   "noop",
		"server.addr": ":8443",
		"server.tls":  false,
	}
)

type Config struct {
	Server Server  `koanf:"server" validate:"required"`
	Stores []Store `koanf:"stores" validate:"required,gt=0,dive,required"`
}

type Store struct {
	Path        string      `koanf:"path" validate:"required"`
	Interval    int         `koanf:"interval" validate:"required"`
	Auth        Auth        `koanf:"auth" validate:"required"`
	Credentials Credentials `koanf:"credentials" validate:"required"`
}

type Auth struct {
	Type         string   `koanf:"type" validate:"oneof=noop user-agent header"`
	UserAgents   []string `koanf:"user_agents" validate:"required_if=Type user-agent"`
	HeaderValues []string `koanf:"header_values" validate:"required_if=Type header"`
	HeaderName   string   `koanf:"header_name" validate:"required_if=Type header"`
}

type Credentials struct {
	Id     string   `koanf:"id" validate:"required"`
	Scopes []string `koanf:"scopes"`
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

	builder := NewBuilder(k, path, defaultConfig)

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
