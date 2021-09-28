package auth

import (
	"net/http"
	"strings"
)

type UserAgentAuthorizer struct {
	userAgents []string
}

func NewUserAgentAuthorizer(userAgents []string) *UserAgentAuthorizer {
	var lowerCaseUserAgens []string

	for _, ua := range userAgents {
		lowerCaseUserAgens = append(lowerCaseUserAgens, strings.ToLower(ua))
	}

	return &UserAgentAuthorizer{
		userAgents: lowerCaseUserAgens,
	}
}

func (u UserAgentAuthorizer) IsAuthorized(r *http.Request) bool {
	userAgentHeader := strings.ToLower(r.UserAgent())

	if userAgentHeader == "" {
		return false
	}

	for _, ua := range u.userAgents {
		if strings.Contains(userAgentHeader, ua) {
			return true
		}
	}

	return false
}
