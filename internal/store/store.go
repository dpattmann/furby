package store

import "golang.org/x/oauth2"

type Store interface {
	GetToken() (*oauth2.Token, error)
	BackgroundUpdate(interval int)
}
