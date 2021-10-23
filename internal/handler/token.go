package handler

import (
	"encoding/json"
	"github.com/dpattmann/furby/internal/metrics"
	"github.com/prometheus/client_golang/prometheus"
	"net/http"

	"github.com/dpattmann/furby/internal/auth"
	"github.com/dpattmann/furby/internal/store"
)

type TokenHandler struct {
	store store.Store
	auth  auth.Authorizer
}

func NewTokenHandler(store store.Store, auth auth.Authorizer) TokenHandler {
	return TokenHandler{
		store: store,
		auth:  auth,
	}
}

func (t TokenHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
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

	jsonToken, err := json.Marshal(token)

	if err != nil {
		metrics.Http500Errors.Inc()
		http.Error(w, JsonParseError, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonToken)

	return
}
