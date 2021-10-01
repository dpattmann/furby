package auth

import "net/http"

type Authorizer interface {
	IsAuthorized(r *http.Request) bool
}
