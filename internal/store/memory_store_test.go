package store

import (
	"net/http"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
)

var (
	validMockClientConfig   = NewClientCredentialsConfig("ClientIdValue", "ClientSecretValue", "http://localhost:8080", []string{})
	invalidMockClientConfig = NewClientCredentialsConfig("ClientIdValue", "ClientSecretValue", "http://localhost:8081", []string{})
)

func TestStore_GetToken(t *testing.T) {
	// Setup http mock backend
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	expirationDate := time.Now().Add(time.Minute * 5)
	httpmock.RegisterResponder("POST", "http://localhost:8080", httpmock.NewJsonResponderOrPanic(http.StatusOK, oauth2.Token{
		AccessToken: "AccessTokenValue",
		TokenType:   "TokenTypeValue",
		Expiry:      expirationDate,
	}))

	httpmock.RegisterResponder("POST", "http://localhost:8081", httpmock.NewStringResponder(400, "Invalid Token Request"))

	validTokenStore := NewMemoryStore(validMockClientConfig)
	invalidTokenStore := NewMemoryStore(invalidMockClientConfig)

	wantedToken := &oauth2.Token{
		AccessToken: "AccessTokenValue",
		TokenType:   "TokenTypeValue",
		Expiry:      expirationDate,
	}
	// mock backend ready

	t.Run("Get new token from token server", func(t *testing.T) {
		assert.Equal(t, 0, httpmock.GetTotalCallCount())
		token, err := validTokenStore.GetToken()

		assert.NoError(t, err)
		assert.Equal(t, wantedToken.AccessToken, token.AccessToken)
		assert.Equal(t, wantedToken.TokenType, token.TokenType)
		assert.Equal(t, 1, httpmock.GetTotalCallCount())
	})

	t.Run("Get valid token from store", func(t *testing.T) {
		token, err := validTokenStore.GetToken()

		assert.NoError(t, err)
		assert.Equal(t, wantedToken.AccessToken, token.AccessToken)
		assert.Equal(t, wantedToken.TokenType, token.TokenType)
		assert.Equal(t, 1, httpmock.GetTotalCallCount())
	})

	t.Run("Get error from token server", func(t *testing.T) {
		_, err := invalidTokenStore.GetToken()

		assert.Error(t, err)
	})
}

func BenchmarkGetToken1(b *testing.B)     { benchmarkGetTokenWithOneSecondHttpDelay(b) }
func BenchmarkGetToken10(b *testing.B)    { benchmarkGetTokenWithOneSecondHttpDelay(b) }
func BenchmarkGetToken100(b *testing.B)   { benchmarkGetTokenWithOneSecondHttpDelay(b) }
func BenchmarkGetToken1000(b *testing.B)  { benchmarkGetTokenWithOneSecondHttpDelay(b) }
func BenchmarkGetToken10000(b *testing.B) { benchmarkGetTokenWithOneSecondHttpDelay(b) }

func benchmarkGetTokenWithOneSecondHttpDelay(b *testing.B) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", "http://localhost:8081", func(request *http.Request) (*http.Response, error) {
		time.Sleep(1 * time.Second)

		return httpmock.NewJsonResponse(http.StatusOK, oauth2.Token{
			AccessToken: "AccessTokenValue",
			TokenType:   "TokenTypeValue",
		})
	})

	tokenStore := NewMemoryStore(validMockClientConfig)

	for i := 0; i < b.N; i++ {
		_, _ = tokenStore.GetToken()
	}
}
