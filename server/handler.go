package server

import (
	"encoding/json"
	"fmt"
	"github.com/dpattmann/furby/auth"
	"net/http"

	"github.com/dpattmann/furby/store"
)

type StoreHandler struct {
	store store.Store
	auth  auth.Authorization
}

func NewStoreHandler(store store.Store, auth auth.Authorization) StoreHandler {
	return StoreHandler{
		store: store,
		auth:  auth,
	}
}

func (t StoreHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		w.WriteHeader(http.StatusTeapot)
		w.Write([]byte("I'm a teapot"))
		return
	}

	if !t.auth.IsAuthorized(req) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Not authorized"))
	}

	token, err := t.store.GetToken()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error getting token from store"))
		return
	}

	jsonToken, err := json.Marshal(token)

	w.WriteHeader(200)
	w.Write(jsonToken)

	return
}

func Serve(tokenEndpointHandler StoreHandler) error {
	fmt.Println("Server is running on port *:8080")
	return http.ListenAndServe(":8080", tokenEndpointHandler)
}
