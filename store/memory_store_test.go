package store

import (
	oauth22 "github.com/dpattmann/furby/oauth2"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
)

var (
	mockClientConfig = oauth22.NewClientCredentialsConfig("ClientIdValue", "ClientSecretValue", "http://localhost:8081", []string{})
)

func TestStore_GetToken(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", "http://localhost:8081", httpmock.NewJsonResponderOrPanic(200, oauth2.Token{
		AccessToken: "AccessTokenValue",
		TokenType:   "TokenTypeValue",
	}))

	tokenStore := NewMemoryStore(mockClientConfig)

	testToken := &oauth2.Token{
		AccessToken: "AccessTokenValue",
		TokenType:   "TokenTypeValue",
	}

	token, err := tokenStore.GetToken()

	assert.NoError(t, err)
	assert.Equal(t, testToken.AccessToken, token.AccessToken)
	assert.Equal(t, testToken.TokenType, token.TokenType)
}

func TestStore_ErrorGettingToken(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", "http://localhost:8081", httpmock.NewStringResponder(400, "Invalid Token Request"))

	tokenStore := NewMemoryStore(mockClientConfig)

	_ = &oauth2.Token{
		AccessToken: "AccessTokenValue",
		TokenType:   "TokenTypeValue",
	}

	_, err := tokenStore.GetToken()

	assert.Error(t, err)
}

func TestStore_ReturnValidToken(t *testing.T) {
	tokenStore := NewMemoryStore(mockClientConfig)

	date := time.Now().Add(time.Minute * 5)

	testToken := &oauth2.Token{
		AccessToken: "AccessTokenValue",
		TokenType:   "TokenTypeValue",
		Expiry:      date,
	}

	tokenStore.setToken(testToken)
	token, err := tokenStore.GetToken()

	assert.NoError(t, err)
	assert.Equal(t, testToken.AccessToken, token.AccessToken)
	assert.Equal(t, testToken.TokenType, token.TokenType)
}
