package server

import (
	"encoding/json"
	"net/http"

	"github.com/dpattmann/furby/internal/auth"
	"github.com/dpattmann/furby/internal/store"
)

type Handler struct {
	store store.Store
	auth  auth.Authorizer
}

func NewHandler(store store.Store, auth auth.Authorizer) Handler {
	return Handler{
		store: store,
		auth:  auth,
	}
}

func (t Handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
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
		http.Error(w, TokenStoreError, http.StatusInternalServerError)
		return
	}

	jsonToken, err := json.Marshal(token)

	w.WriteHeader(http.StatusOK)
	w.Write(jsonToken)

	return
}
