package config

import (
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"

	"errors"
	"path/filepath"
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
	return c.loadConfigFile()
}

func (c *Builder) loadConfigFile() (err error) {
	switch {
	case filepath.Ext(c.path) == ".yaml" || filepath.Ext(c.path) == ".yml":
		err = c.k.Load(file.Provider(c.path), yaml.Parser())
	default:
		err = c.k.Load(file.Provider(c.path), json.Parser())
	}

	if err != nil {
		err = errors.New("error loading config: " + err.Error())
		return
	}

	return nil
}
