package oauth2

import (
	"reflect"
	"testing"

	"golang.org/x/oauth2/clientcredentials"
)

func TestNewClientCredentialsConfig(t *testing.T) {
	type args struct {
		id     string
		secret string
		url    string
		scopes []string
	}
	tests := []struct {
		name string
		args args
		want *clientcredentials.Config
	}{
		{
			name: "Create Client Credentials Config",
			args: args{
				id:     "ClientIdString",
				secret: "ClientSecretString",
				url:    "https://localhost/oauth2/token",
				scopes: []string{"scopeA"},
			},
			want: NewClientCredentialsConfig("ClientIdString", "ClientSecretString", "https://localhost/oauth2/token", []string{"scopeA"}),
		},
		{
			name: "Create Client Credentials Config with two scopes",
			args: args{
				id:     "ClientIdString",
				secret: "ClientSecretString",
				url:    "https://localhost/oauth2/token",
				scopes: []string{"scopeA", "scopeB"},
			},
			want: NewClientCredentialsConfig("ClientIdString", "ClientSecretString", "https://localhost/oauth2/token", []string{"scopeA", "scopeB"}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewClientCredentialsConfig(tt.args.id, tt.args.secret, tt.args.url, tt.args.scopes); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewClientCredentialsConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}
