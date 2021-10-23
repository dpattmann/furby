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
			"credentials" : {
        		"id": "TestClientId",
        		"scopes": [
        		    "scopeA",
        		    "scopeB"
        		],
        		"secret": "TestClientSecret",
        		"url": "https://localhost"
			},
			"store": {
				"interval": "300"
			},
    		"server": {
				"addr": ":8443",
    		    "cert": "foo.cert",
    		    "key": "foo.key",
    		    "tls": "true"
    		},
			"auth": {
				"type": "noop"
			}
		}
	`
	configWithInvalidUrl = `
		{    
			"credentials" : {
        		"id": "TestClientId",
        		"scopes": [
        		    "scopeA",
        		    "scopeB"
        		],
        		"secret": "TestClientSecret",
        		"url": "localhost"
			},
			"store": {
				"interval": "300"
			},
			"auth": {
				"type": "noop"
			}
		}
	`
	configWithDisabledTls = `
		{    
			"credentials" : {
        		"id": "TestClientId",
        		"scopes": [
        		    "scopeA",
        		    "scopeB"
        		],
        		"secret": "TestClientSecret",
        		"url": "https://localhost"
			},
    		"server": {
				"addr": ":8443",
    		    "cert": "",
    		    "key": "",
    		    "tls": "false"
    		},
			"store": {
				"interval": "300"
			},
			"auth": {
				"type": "noop"
			}
		}
	`

	configWithTlsButWithoutCert = `
		{    
			"credentials" : {
        		"id": "TestClientId",
        		"scopes": [
        		    "scopeA",
        		    "scopeB"
        		],
        		"secret": "TestClientSecret",
        		"url": "https://localhost"
			},
    		"server": {
				"addr": ":8443",
    		    "cert": "",
    		    "key": "",
    		    "tls": "true"
    		},
			"store": {
				"interval": "300"
			},
			"auth": {
				"type": "noop"
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
			Credentials: Credentials{
				Id:     "TestClientId",
				Scopes: []string{"scopeA", "scopeB"},
				Secret: "TestClientSecret",
				Url:    "https://localhost",
			},
			Server: Server{
				Addr: ":8443",
				Cert: "foo.cert",
				Key:  "foo.key",
				Tls:  true,
			},
			Store: Store{
				Interval: 300,
			},
			Auth: Auth{
				Type: "noop",
			},
		}

		assert.NoError(t, err)
		assert.Equal(t, want, got)

		removeTempConfig()
	})
}

func TestNewFailureConfig(t *testing.T) {
	tests := []struct {
		testCase        string
		testDescription string
		wantErr         bool
	}{
		{
			testCase:        configWithInvalidUrl,
			testDescription: "Create config with invalid url",
			wantErr:         true,
		},
		{
			testCase:        configWithDisabledTls,
			testDescription: "Create config without tls config",
			wantErr:         false,
		},
		{
			testCase:        configWithTlsButWithoutCert,
			testDescription: "Create config with tls but without cert files",
			wantErr:         true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testDescription, func(t *testing.T) {
			err := createTempConfig(tt.testCase)
			assert.NoError(t, err)
			if _, err := NewConfig("./test_temp.json"); (err != nil) != tt.wantErr {
				t.Errorf("Expected %v, wantErr %v", err, tt.wantErr)
			}
			removeTempConfig()
		})
	}

}
