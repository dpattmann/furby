package handler

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/dpattmann/furby/internal/auth"
	"github.com/dpattmann/furby/internal/metrics"
	"github.com/dpattmann/furby/internal/store"

	"github.com/prometheus/client_golang/prometheus"
)

type ProxyHandler struct {
	target *url.URL
	store  store.Store
	auth   auth.Authorizer
}

func NewProxyHandler(target *url.URL, store store.Store, auth auth.Authorizer) ProxyHandler {
	return ProxyHandler{
		target: target,
		store:  store,
		auth:   auth,
	}
}

func (t ProxyHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	metrics.ReceivedRequests.Inc()

	timer := prometheus.NewTimer(metrics.RequestTime)
	defer timer.ObserveDuration()

	if req.Method != http.MethodGet {
		http.Error(w, TeapotMessage, http.StatusTeapot)
		return
	}

	if !t.auth.IsAuthorized(req) {
		http.Error(w, Unauthorized, http.StatusUnauthorized)
		return
	}

	token, err := t.store.GetToken()

	if err != nil {
		metrics.Http500Errors.Inc()
		http.Error(w, TokenStoreError, http.StatusInternalServerError)
		return
	}

	targetQuery := t.target.RawQuery

	proxy := &httputil.ReverseProxy{
		Director: func(proxyRequest *http.Request) {
			proxyRequest.Header.Add("X-Forwarded-Host", proxyRequest.Host)
			proxyRequest.Header.Add("X-Origin-Host", t.target.Host)
			proxyRequest.Header.Add("Bearer", token.AccessToken)

			proxyRequest.URL.Scheme = t.target.Scheme
			proxyRequest.URL.Host = t.target.Host
			proxyRequest.Host = t.target.Host

			if targetQuery == "" || proxyRequest.URL.RawQuery == "" {
				proxyRequest.URL.RawQuery = targetQuery + proxyRequest.URL.RawQuery
			} else {
				proxyRequest.URL.RawQuery = targetQuery + "&" + proxyRequest.URL.RawQuery
			}

		},
		ErrorHandler: func(rw http.ResponseWriter, r *http.Request, err error) {
			http.Error(rw, ProxyError, http.StatusInternalServerError)
		},
	}

	proxy.ServeHTTP(w, req)

	return
}
