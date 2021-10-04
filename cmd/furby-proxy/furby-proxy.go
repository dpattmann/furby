package main

import (
	"log"

	"github.com/dpattmann/furby/internal/auth"
	"github.com/dpattmann/furby/internal/config"
	"github.com/dpattmann/furby/internal/server"
	"github.com/dpattmann/furby/internal/store"

	flag "github.com/spf13/pflag"
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

	clientCredentialsConfig := store.NewClientCredentialsConfig(
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

	forwarder := server.NewForwarder(memoryStore, authorizer)

	if c.Server.Tls {
		if err := server.ServeTls(forwarder, c.Server.Cert, c.Server.Key); err != nil {
			log.Fatal("Error running server")
		}
	}

	if err := server.Serve(forwarder); err != nil {
		log.Fatal("Error running server")
	}
}
