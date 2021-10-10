package store

import (
	"context"
	"errors"
	"sync"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
	"golang.org/x/sync/singleflight"
)

var (
	requestGroup singleflight.Group
)

const prefetchTime = 3590

func NewMemoryStore(c *clientcredentials.Config, wg *sync.WaitGroup) *MemoryStore {
	return &MemoryStore{
		client:    c,
		waitGroup: wg,
	}
}

type MemoryStore struct {
	mu        sync.RWMutex
	token     *oauth2.Token
	client    *clientcredentials.Config
	waitGroup *sync.WaitGroup
}

func (s *MemoryStore) GetToken() (*oauth2.Token, error) {
	if s.token != nil && s.token.Expiry.Sub(time.Now()).Seconds() > prefetchTime {
		return s.token, nil
	} else {
		if s.token == nil {
			return s.updateToken()
		} else {
			s.waitGroup.Add(1)
			go s.fetchAndUpdateToken()

			if s.token.Valid() {
				return s.token, nil
			}
		}
	}

	return nil, errors.New("couldn't get a token")
}

func (s *MemoryStore) fetchAndUpdateToken() error {
	_, err, _ := requestGroup.Do("fetchAndToken", func() (token interface{}, err error) {
		return s.updateToken()
	})

	if err != nil {
		return err
	}

	s.waitGroup.Done()

	return nil
}

func (s *MemoryStore) updateToken() (token *oauth2.Token, err error) {
	token, err = s.client.Token(context.Background())

	if err != nil {
		return
	}

	s.setToken(token)

	return
}

func (s *MemoryStore) setToken(token *oauth2.Token) {
	s.mu.Lock()
	s.token = token
	s.mu.Unlock()
}
