package main

import (
	"log"

	flag "github.com/spf13/pflag"

	"github.com/dpattmann/furby/auth"
	"github.com/dpattmann/furby/config"
	"github.com/dpattmann/furby/oauth2"
	"github.com/dpattmann/furby/server"
	"github.com/dpattmann/furby/store"
)

var (
	authorizer auth.Authorizer
)

func main() {
	path := flag.StringP("path", "p", "./furby_config.json", "parameter file")
	flag.Parse()

	if flag.NFlag() == 0 {
		flag.PrintDefaults()
		log.Fatal("Please pass parameter(s)")
	}

	c, err := config.NewConfig(*path)

	if err != nil {
		log.Fatalf("Can't read config: %v", err)
	}

	clientCredentialsConfig := oauth2.NewClientCredentialsConfig(
		c.ClientCredentials.Id,
		c.ClientCredentials.Secret,
		c.ClientCredentials.Url,
		c.ClientCredentials.Scopes,
	)

	switch c.Auth.Type {
	case "user-agent":
		authorizer = auth.NewUserAgentAuthorizer(c.Auth.UserAgents)
	default:
		authorizer = auth.NewNoOpAuthorizer()
	}

	memoryStore := store.NewMemoryStore(clientCredentialsConfig)

	storeHandler := server.NewStoreHandler(memoryStore, authorizer)

	if c.Server.Tls {
		if err := server.ServeTls(storeHandler, c.Server.Cert, c.Server.Key); err != nil {
			log.Fatal("Error running server")
		}
	}

	if err := server.Serve(storeHandler); err != nil {
		log.Fatal("Error running server")
	}
}
