package noop

import "net/http"

type Authorizer struct {

}

func NewAuthorizer() *Authorizer {
	return &Authorizer{}
}

func (a Authorizer) IsAuthorized(_ *http.Request) bool {
	return true
}

