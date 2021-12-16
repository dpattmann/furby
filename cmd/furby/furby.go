package main

import (
	"log"
	"net/http"
	"net/url"

	"github.com/dpattmann/furby/internal/auth"
	"github.com/dpattmann/furby/internal/config"
	"github.com/dpattmann/furby/internal/handler"
	"github.com/dpattmann/furby/internal/metrics"
	"github.com/dpattmann/furby/internal/store"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	flag "github.com/spf13/pflag"
)

func setupStoreAndReturnHandler(s config.Store) http.Handler {
	clientCredentialsConfig := store.NewClientCredentialsConfig(s.Credentials)
	tokenStore := store.NewMemoryStore(clientCredentialsConfig)

	if s.Interval > 0 {
		go tokenStore.BackgroundUpdate(s.Interval)
	}

	var authorizer auth.Authorizer
	switch s.Auth.Type {
	case "user-agent":
		authorizer = auth.NewUserAgentAuthorizer(s.Auth.UserAgents)
	case "header":
		authorizer = auth.NewHeaderAuthorizer(s.Auth.HeaderName, s.Auth.HeaderValues)
	default:
		authorizer = auth.NewNoOpAuthorizer()
	}

	if s.Target != "" {
		target, err := url.Parse(s.Target)
		if err != nil {
			log.Fatalf("Can't parse url. Error: %v", err)
		}
		return handler.NewProxyHandler(target, tokenStore, authorizer)

	}

	return handler.NewStoreHandler(tokenStore, authorizer)
}

func main() {
	path := flag.StringP("config", "c", "/etc/furby/config.yaml", "config file")
	flag.Parse()

	c, err := config.NewConfig(*path)

	if err != nil {
		log.Fatalf("Can't read config: %v", err)
	}

	// The defaultMuxRouter exposes /debug/vars as a default route. As we don't want this,
	// we create a new router
	serveMux := http.NewServeMux()

	for _, s := range c.Stores {
		h := setupStoreAndReturnHandler(s)
		serveMux.Handle(s.Path, h)
	}

	serveMux.Handle("/metrics", promhttp.HandlerFor(metrics.PrometheusRegister, promhttp.HandlerOpts{}))
	serveMux.Handle("/health", handler.HealthHandler())

	server := http.Server{
		Addr:    ":8443",
		Handler: serveMux,
	}

	if c.Server.Tls {
		if err := server.ListenAndServeTLS(c.Server.Cert, c.Server.Key); err != nil {
			log.Fatal("Error running server")
		}
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatal("Error running server")
	}
}
