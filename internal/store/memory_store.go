package store

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/dpattmann/furby/internal/metrics"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
	"golang.org/x/sync/singleflight"
)

var (
	requestGroup singleflight.Group
)

func NewMemoryStore(c *clientcredentials.Config) *MemoryStore {
	return &MemoryStore{
		client: c,
	}
}

type MemoryStore struct {
	mu     sync.RWMutex
	token  *oauth2.Token
	client *clientcredentials.Config
}

func (s *MemoryStore) BackgroundUpdate(interval int) {
	for {
		<-time.After(time.Second * time.Duration(interval))

		fmt.Println("Updating token")
		if _, err := s.updateToken(); err != nil {
			fmt.Printf("Error updating token: %v", err)
		}
		fmt.Println("Updated token successfully")
	}
}

func (s *MemoryStore) GetToken() (*oauth2.Token, error) {
	if s.token.Valid() {
		return s.token, nil
	}

	token, err, _ := requestGroup.Do("GetToken", func() (token interface{}, err error) {
		return s.updateToken()
	})

	if err != nil {
		return nil, err
	}

	return token.(*oauth2.Token), err
}

func (s *MemoryStore) updateToken() (token *oauth2.Token, err error) {
	metrics.BackendRequests.Inc()
	token, err = s.client.Token(context.Background())

	if err != nil {
		metrics.BackendRequestsFailed.Inc()
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
