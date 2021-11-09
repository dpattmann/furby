package config

import (
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	validConfig = `
		{
			"stores": [{
				"interval": "5",
				"path": "/",
				"credentials" : {
			  		"id": "TestClientId",
			  		"scopes": [
						"scopeA",
						"scopeB"
			  		],
			  	"secret": "TestClientSecret",
			  	"url": "https://localhost"
				},
				"auth": {
					"type": "noop"
				}
			}],
			"server": {
				"addr": ":8443",
				"cert": "foo.cert",
				"key": "foo.key",
				"tls": "true"
		  	}
		}
	`
)

func createTempConfig(content string) (err error) {
	f, err := os.Create("test_temp.json")

	if err != nil {
		return
	}

	defer f.Close()

	_, err = f.WriteString(content)

	return
}

func removeTempConfig() {
	err := os.Remove("test_temp.json")

	if err != nil {
		log.Fatal(err)
	}
}

func TestNewValidConfig(t *testing.T) {
	t.Run("Create valid config from environment", func(t *testing.T) {
		err := createTempConfig(validConfig)
		assert.NoError(t, err)

		got, err := NewConfig("./test_temp.json")

		want := &Config{
			Server: Server{
				Addr: ":8443",
				Cert: "foo.cert",
				Key:  "foo.key",
				Tls:  true,
			},
			Stores: []Store{
				{
					Interval: 5,
					Path:     "/",
					Credentials: Credentials{
						Id:     "TestClientId",
						Scopes: []string{"scopeA", "scopeB"},
						Secret: "TestClientSecret",
						Url:    "https://localhost",
					},
					Auth: Auth{
						Type: "noop",
					},
				},
			},
		}

		assert.NoError(t, err)
		assert.Equal(t, want, got)

		removeTempConfig()
	})
}

func TestConfig_validate(t *testing.T) {
	type fields struct {
		Server Server
		Stores []Store
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "Should return no validation error on valid config",
			fields: fields{
				Stores: []Store{
					{
						Path:     "/",
						Interval: 5,
						Auth: Auth{
							Type: "noop",
						},
						Credentials: Credentials{
							Id:     "ClientId",
							Secret: "ClientSecret",
							Url:    "https://localhost/oauth2/token",
						},
					},
				},
				Server: Server{
					Addr: ":8080",
					Tls:  false,
				},
			},
			wantErr: false,
		},
		{
			name: "Should return validation error if tls is true but no cert and key is specified",
			fields: fields{
				Stores: []Store{
					{
						Path:     "/",
						Interval: 5,
						Auth: Auth{
							Type: "noop",
						},
						Credentials: Credentials{
							Id:     "ClientId",
							Secret: "ClientSecret",
							Url:    "https://localhost/oauth2/token",
						},
					},
				},
				Server: Server{
					Addr: ":8080",
					Tls:  true,
				},
			},
			wantErr: true,
		},
		{
			name: "Should return validation error if no store is specified",
			fields: fields{
				Stores: []Store{},
				Server: Server{
					Addr: ":8080",
					Tls:  false,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				Server: tt.fields.Server,
				Stores: tt.fields.Stores,
			}
			if err := c.validate(); (err != nil) != tt.wantErr {
				t.Errorf("validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
