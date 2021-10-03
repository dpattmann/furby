package server

import (
	"github.com/dpattmann/furby/internal/auth"
	"github.com/dpattmann/furby/internal/store"
	"io"
	"net/http"
)

type Forwarder struct {
	store store.Store
	auth  auth.Authorizer
}

func NewForwarder(store store.Store, auth auth.Authorizer) Forwarder {
	return Forwarder{
		store: store,
		auth:  auth,
	}
}

func (f Forwarder) addAccessTokenToRequest(req *http.Request) (err error) {
	token, err := f.store.GetToken()

	if err != nil {
		return
	}

	req.Header.Add("Authorization", "Bearer "+token.AccessToken)

	return
}

func (f Forwarder) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if !f.auth.IsAuthorized(req) {
		http.Error(w, Unauthorized, http.StatusUnauthorized)
		return
	}

	err := f.addAccessTokenToRequest(req)

	if err != nil {
		http.Error(w, TokenStoreError, http.StatusInternalServerError)
		return
	}

	resp, err := http.DefaultTransport.RoundTrip(req)

	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	defer resp.Body.Close()

	f.copyHeader(w.Header(), resp.Header)

	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func (f Forwarder) copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}
