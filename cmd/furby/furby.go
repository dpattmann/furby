package main

import (
	"log"
	"net/http"

	"github.com/dpattmann/furby/internal/auth"
	"github.com/dpattmann/furby/internal/config"
	"github.com/dpattmann/furby/internal/handler"
	"github.com/dpattmann/furby/internal/metrics"
	"github.com/dpattmann/furby/internal/store"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	flag "github.com/spf13/pflag"
)

var (
	authorizer auth.Authorizer
)

func main() {
	path := flag.StringP("path", "p", "/etc/furby/config.yaml", "parameter file")
	flag.Parse()

	c, err := config.NewConfig(*path)

	if err != nil {
		log.Fatalf("Can't read config: %v", err)
	}

	clientCredentialsConfig := store.NewClientCredentialsConfig(c.Credentials)
	memoryStore := store.NewMemoryStore(clientCredentialsConfig)

	if c.Store.Interval > 0 {
		go memoryStore.BackgroundUpdate(c.Store.Interval)
	}

	switch c.Auth.Type {
	case "user-agent":
		authorizer = auth.NewUserAgentAuthorizer(c.Auth.UserAgents)
	case "header":
		authorizer = auth.NewHeaderAuthorizer(c.Auth.HeaderName, c.Auth.HeaderValues)
	default:
		authorizer = auth.NewNoOpAuthorizer()
	}

	tokenHandler := handler.NewTokenHandler(memoryStore, authorizer)

	http.Handle("/metrics", promhttp.HandlerFor(metrics.PrometheusRegister, promhttp.HandlerOpts{}))
	http.Handle("/health", handler.HealthHandler())
	http.Handle("/", tokenHandler)

	if c.Server.Tls {
		if err := http.ListenAndServeTLS(c.Server.Addr, c.Server.Cert, c.Server.Key, nil); err != nil {
			log.Fatal("Error running server")
		}
	}

	if err := http.ListenAndServe(c.Server.Addr, nil); err != nil {
		log.Fatal("Error running server")
	}
}
