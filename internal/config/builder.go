package config

import (
	"errors"
	"path/filepath"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/providers/file"
)

type Builder struct {
	k    *koanf.Koanf
	path string
}

func NewBuilder(k *koanf.Koanf, path string) *Builder {
	return &Builder{k: k, path: path}
}

func (c *Builder) unmarshalConfigToStruct(config *Config) error {
	return c.k.Unmarshal("", &config)
}

func (c *Builder) loadConfig() (err error) {
	err = c.loadConfigMap()
	if err != nil {
		return
	}

	return c.loadConfigFile()
}

// LoadConfigMap is used to set default values
func (c *Builder) loadConfigMap() error {
	return c.k.Load(confmap.Provider(map[string]interface{}{
		"auth.type":      "noop",
		"server.addr":    ":8443",
		"server.tls":     false,
		"store.interval": 300,
	}, "."), nil)
}

func (c *Builder) loadConfigFile() error {
	switch filepath.Ext(c.path) {
	case ".yaml", ".yml":
		return c.k.Load(file.Provider(c.path), yaml.Parser())
	case ".json":
		return c.k.Load(file.Provider(c.path), json.Parser())
	}

	return errors.New("error: unsupported file type")
}
