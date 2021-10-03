package store

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

func NewClientCredentialsConfig(id, secret, url string, scopes []string) *clientcredentials.Config {
	return &clientcredentials.Config{
		ClientID:     id,
		ClientSecret: secret,
		TokenURL:     url,
		Scopes:       scopes,
		AuthStyle:    oauth2.AuthStyleAutoDetect,
	}
}
