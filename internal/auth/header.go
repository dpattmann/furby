package auth

import (
	"net/http"
)

type HeaderAuthorizer struct {
	headerName   string
	headerValues []string
}

func NewHeaderAuthorizer(headerName string, headerValues []string) *HeaderAuthorizer {
	return &HeaderAuthorizer{
		headerName:   headerName,
		headerValues: headerValues,
	}
}

func (h HeaderAuthorizer) IsAuthorized(r *http.Request) bool {
	header := r.Header.Get(h.headerName)

	if header == "" {
		return false
	}

	for _, ua := range h.headerValues {
		if ua == header {
			return true
		}
	}

	return false
}
