package main

import (
	"log"

	"github.com/dpattmann/furby/auth/noop"
	"github.com/dpattmann/furby/config"
	"github.com/dpattmann/furby/oauth2"
	"github.com/dpattmann/furby/server"
	"github.com/dpattmann/furby/store/memory"
)

func main() {
	c, err := config.NewConfig()

	if err != nil {
		log.Fatalf("Can't read config: %v", err)
	}

	clientCredentialsConfig := oauth2.NewClientCredentialsConfig(
		c.ClientCredentials.Id,
		c.ClientCredentials.Secret,
		c.ClientCredentials.Url,
		c.ClientCredentials.Scopes,
	)

	noopAuthorizer := noop.NewAuthorizer()
	memoryStore := memory.NewMemoryStore(clientCredentialsConfig)
	storeHandler := server.NewStoreHandler(memoryStore, noopAuthorizer)

	if err := server.Serve(storeHandler); err != nil {
		log.Fatal("Error running server")
	}
}
