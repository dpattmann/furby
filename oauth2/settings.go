package oauth2

type ClientCredentialSettings struct {
	Id     string   `koanf:"id"`
	Scopes []string `koanf:"scopes"`
	Secret string   `koanf:"secret"`
	Url    string   `koanf:"url"`
}
