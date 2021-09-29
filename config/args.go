package config

import (
	flag "github.com/spf13/pflag"
)

type Args struct {
	authType     *string
	userAgents   *[]string
	clientId     *string
	clientSecret *string
	tokenUrl     *string
	scopes       *[]string
	cert         *string
	key          *string
	tls          *bool
}

func NewArgs() *Args {
	return &Args{}
}

func (a *Args) Parse() (parsedArgs *Args) {
	a.authType = flag.StringP("auth", "a", "noop", "type of authorizer")
	a.userAgents = flag.StringSliceP("user_agents", "u", []string{}, "user agents as authorization strings")
	a.clientId = flag.StringP("client_id", "i", "", "OAuth2 client id")
	a.clientSecret = flag.StringP("client_secret", "s", "", "OAuth2 client secret")
	a.tokenUrl = flag.StringP("token_url", "o", "", "OAuth2 token server url")
	a.scopes = flag.StringSliceP("scopes", "p", []string{}, "OAuth2 token scopes")
	a.cert = flag.StringP("cert", "c", "", "TLS cert file")
	a.key = flag.StringP("key", "k", "", "TLS key file")
	a.tls = flag.BoolP("isTLS", "t", false, "Is TLS active")

	flag.Parse()

	return a
}
