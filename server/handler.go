package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/dpattmann/furby/auth"
	"github.com/dpattmann/furby/store"
)

type StoreHandler struct {
	store store.Store
	auth  auth.Authorizer
}

func NewStoreHandler(store store.Store, auth auth.Authorizer) StoreHandler {
	return StoreHandler{
		store: store,
		auth:  auth,
	}
}

func (t StoreHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		writeResponse(w, http.StatusTeapot, "I'm a teapot")
		return
	}

	if !t.auth.IsAuthorized(req) {
		writeResponse(w, http.StatusUnauthorized, "Not authorized")
	}

	token, err := t.store.GetToken()

	if err != nil {
		writeResponse(w, http.StatusInternalServerError, "Error getting token from store")
		return
	}

	jsonToken, err := json.Marshal(token)

	writeResponse(w, http.StatusOK, string(jsonToken))

	return
}

func ServeTls(tokenEndpointHandler StoreHandler, cert, key string) error {
	fmt.Println("Server is running on port *:4443")
	return http.ListenAndServeTLS(":4443", cert, key, tokenEndpointHandler)
}

func Serve(tokenEndpointHandler StoreHandler) error {
	fmt.Println("Server is running on port *:8080")
	return http.ListenAndServe(":8080", tokenEndpointHandler)
}

func writeResponse(writer http.ResponseWriter, status int, message string) {
	writer.WriteHeader(status)
	_, err := writer.Write([]byte(message))

	if err != nil {
		log.Printf("error '%v' while writing message response", err.Error())
	}
}
