package auth

import (
	"net/http"
	"strings"
)

type UserAgentAuthorizer struct {
	userAgents []string
}

func NewUserAgentAuthorizer(userAgents []string) *UserAgentAuthorizer {
	return &UserAgentAuthorizer{
		userAgents: userAgents,
	}
}

func (u UserAgentAuthorizer) IsAuthorized(r *http.Request) bool {
	userAgentHeader := strings.ToLower(r.UserAgent())

	if userAgentHeader == "" {
		return false
	}

	for _, ua := range u.userAgents {
		if strings.Contains(userAgentHeader, strings.ToLower(ua)) {
			return true
		}
	}

	return false
}
