package config

import (
	"os"
	"testing"

	"github.com/dpattmann/furby/oauth2"

	"github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) {
	t.Run("Create valid config from environment", func(t *testing.T) {
		_ = os.Setenv("FURBY_CLIENTCREDENTIALS_ID", "TestClientId")
		_ = os.Setenv("FURBY_CLIENTCREDENTIALS_SECRET", "TestClientSecret")
		_ = os.Setenv("FURBY_CLIENTCREDENTIALS_URL", "https://localhost")
		_ = os.Setenv("FURBY_CLIENTCREDENTIALS_SCOPES", "scopeA scopeB")

		got, err := NewConfig()

		want := &Config{ClientCredentials: oauth2.ClientCredentials{
			Id:     "TestClientId",
			Scopes: []string{"scopeA", "scopeB"},
			Secret: "TestClientSecret",
			Url:    "https://localhost",
		}}

		assert.NoError(t, err)
		assert.Equal(t, want, got)
	})

	t.Run("Create config with invalid url", func(t *testing.T) {
		_ = os.Setenv("FURBY_CLIENTCREDENTIALS_ID", "TestClientId")
		_ = os.Setenv("FURBY_CLIENTCREDENTIALS_SECRET", "TestClientSecret")
		_ = os.Setenv("FURBY_CLIENTCREDENTIALS_URL", "localhost")
		_ = os.Setenv("FURBY_CLIENTCREDENTIALS_SCOPES", "scopeA scopeB")

		_, err := NewConfig()

		assert.Error(t, err)
	})
}
