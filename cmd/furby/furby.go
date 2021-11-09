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

	// The defaultMuxRouter exposes /debug/vars as a default route. As we don't want this,
	// we create a new router
	m := http.NewServeMux()

	for _, s := range c.Stores {
		clientCredentialsConfig := store.NewClientCredentialsConfig(s.Credentials)
		tokenStore := store.NewMemoryStore(clientCredentialsConfig)

		for _, s := range c.Stores {
			if s.Interval > 0 {
				go tokenStore.BackgroundUpdate(s.Interval)
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

		tokenHandler := handler.NewStoreHandler(tokenStore, authorizer)

		m.Handle(s.Path, tokenHandler)
	}

	m.Handle("/metrics", promhttp.HandlerFor(metrics.PrometheusRegister, promhttp.HandlerOpts{}))
	m.Handle("/health", handler.HealthHandler())

	s := http.Server{
		Addr:    ":8443",
		Handler: m,
	}

	if c.Server.Tls {
		if err := s.ListenAndServeTLS(c.Server.Cert, c.Server.Key); err != nil {
			log.Fatal("Error running server")
		}
	}

	if err := s.ListenAndServe(); err != nil {
		log.Fatal("Error running server")
	}
}
