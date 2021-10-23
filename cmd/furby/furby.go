package main

import (
	"github.com/dpattmann/furby/internal/auth"
	"github.com/dpattmann/furby/internal/config"
	"github.com/dpattmann/furby/internal/server"
	"github.com/dpattmann/furby/internal/store"
	flag "github.com/spf13/pflag"
	"log"
	"net/http"
	"strconv"
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

	clientCredentialsConfig := store.NewClientCredentialsConfig(c.ClientCredentials)
	memoryStore := store.NewMemoryStore(clientCredentialsConfig)

	go memoryStore.BackgroundUpdate(300)

	switch c.Auth.Type {
	case "user-agent":
		authorizer = auth.NewUserAgentAuthorizer(c.Auth.UserAgents)
	default:
		authorizer = auth.NewNoOpAuthorizer()
	}

	handler := server.NewHandler(memoryStore, authorizer)

	if c.Server.Tls {
		if err := http.ListenAndServeTLS(":"+strconv.Itoa(c.Server.Port), c.Server.Cert, c.Server.Key, handler); err != nil {
			log.Fatal("Error running server")
		}
	}

	if err := http.ListenAndServe(":"+strconv.Itoa(c.Server.Port), handler); err != nil {
		log.Fatal("Error running server")
	}
}
