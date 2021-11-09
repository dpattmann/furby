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

	for _, s := range c.Stores {
		clientCredentialsConfig := store.NewClientCredentialsConfig(s.Credentials)
		memoryStore := store.NewMemoryStore(clientCredentialsConfig)

		for _, s := range c.Stores {
			if s.Interval > 0 {
				go memoryStore.BackgroundUpdate(s.Interval)
			}
		}

		switch s.Auth.Type {
		case "user-agent":
			authorizer = auth.NewUserAgentAuthorizer(s.Auth.UserAgents)
		case "header":
			authorizer = auth.NewHeaderAuthorizer(s.Auth.HeaderName, s.Auth.HeaderValues)
		default:
			authorizer = auth.NewNoOpAuthorizer()
		}

		tokenHandler := handler.NewTokenHandler(memoryStore, authorizer)

		http.Handle(s.Path, tokenHandler)
	}

	http.Handle("/metrics", promhttp.HandlerFor(metrics.PrometheusRegister, promhttp.HandlerOpts{}))
	http.Handle("/health", handler.HealthHandler())

	if c.Server.Tls {
		if err := http.ListenAndServeTLS(c.Server.Addr, c.Server.Cert, c.Server.Key, nil); err != nil {
			log.Fatal("Error running server")
		}
	}

	if err := http.ListenAndServe(c.Server.Addr, nil); err != nil {
		log.Fatal("Error running server")
	}
}
