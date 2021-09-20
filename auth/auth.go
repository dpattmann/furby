package auth

import "net/http"

type Authorization interface {
	IsAuthorized (r *http.Request) bool
}
