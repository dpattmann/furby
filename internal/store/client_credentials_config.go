package store

import (
	"github.com/dpattmann/furby/internal/config"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

func NewClientCredentialsConfig(config config.Credentials) *clientcredentials.Config {
	return &clientcredentials.Config{
		ClientID:     config.Id,
		ClientSecret: config.Secret,
		TokenURL:     config.Url,
		Scopes:       config.Scopes,
		AuthStyle:    oauth2.AuthStyleAutoDetect,
	}
}
