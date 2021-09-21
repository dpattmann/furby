package oauth2

type ClientCredentials struct {
	Id     string   `koanf:"id" validate:"required"`
	Scopes []string `koanf:"scopes" validate:"required"`
	Secret string   `koanf:"secret" validate:"required"`
	Url    string   `koanf:"url" validate:"required,url"`
}
