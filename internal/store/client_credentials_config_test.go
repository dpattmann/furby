package store

import (
	"github.com/dpattmann/furby/internal/config"
	"golang.org/x/oauth2/clientcredentials"
	"reflect"
	"testing"
)

func TestNewClientCredentialsConfig(t *testing.T) {
	type args struct {
		config config.Credentials
	}
	tests := []struct {
		name string
		args args
		want *clientcredentials.Config
	}{
		{
			name: "Create Client Credentials Config",
			args: args{
				config: config.Credentials{
					Id:     "ClientIdString",
					Scopes: []string{"scopeA"},
					Secret: "ClientSecretString",
					Url:    "https://localhost/oauth2/token",
				}},
			want: &clientcredentials.Config{
				ClientID:     "ClientIdString",
				ClientSecret: "ClientSecretString",
				TokenURL:     "https://localhost/oauth2/token",
				Scopes:       []string{"scopeA"},
			},
		},
		{
			name: "Create Client Credentials Config with two scopes",
			args: args{
				config: config.Credentials{
					Id:     "ClientIdString",
					Scopes: []string{"scopeA", "scopeB"},
					Secret: "ClientSecretString",
					Url:    "https://localhost/oauth2/token",
				}},
			want: &clientcredentials.Config{
				ClientID:     "ClientIdString",
				ClientSecret: "ClientSecretString",
				TokenURL:     "https://localhost/oauth2/token",
				Scopes:       []string{"scopeA", "scopeB"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewClientCredentialsConfig(tt.args.config); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewClientCredentialsConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}
