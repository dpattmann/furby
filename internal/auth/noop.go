package auth

import "net/http"

type NoOp struct {
}

func NewNoOpAuthorizer() *NoOp {
	return &NoOp{}
}

func (a NoOp) IsAuthorized(_ *http.Request) bool {
	return true
}
