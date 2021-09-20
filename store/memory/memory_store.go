package memory

import (
	"context"
	"sync"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

func NewMemoryStore(c *clientcredentials.Config) *Store {
	return &Store{
		client: c,
	}
}

type Store struct {
	mu     sync.RWMutex
	token  *oauth2.Token
	client *clientcredentials.Config
}

func (s *Store) GetToken() (t *oauth2.Token, err error) {
	if s.token.Valid() {
		return s.token, nil
	}

	t, err = s.getTokenFromServer()

	if err != nil {
		return
	}

	s.setToken(t)

	return
}

func (s *Store) getTokenFromServer() (*oauth2.Token, error) {
	return s.client.Token(context.Background())
}

func (s *Store) setToken(token *oauth2.Token) {
	s.mu.Lock()
	s.token = token
	s.mu.Unlock()
}
