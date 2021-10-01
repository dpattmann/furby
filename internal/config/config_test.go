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
			"client_credentials" : {
        		"id": "TestClientId",
        		"scopes": [
        		    "scopeA",
        		    "scopeB"
        		],
        		"secret": "TestClientSecret",
        		"url": "https://localhost"
			},
    		"server": {
    		    "cert": "foo.cert",
    		    "key": "foo.key",
    		    "tls": "true"
    		}
		}
	`
	configWithInvalidUrl = `
		{    
			"client_credentials" : {
        		"id": "TestClientId",
        		"scopes": [
        		    "scopeA",
        		    "scopeB"
        		],
        		"secret": "TestClientSecret",
        		"url": "localhost"
			}
		}
	`
	configWithDisabledTls = `
		{    
			"client_credentials" : {
        		"id": "TestClientId",
        		"scopes": [
        		    "scopeA",
        		    "scopeB"
        		],
        		"secret": "TestClientSecret",
        		"url": "https://localhost"
			},
    		"server": {
    		    "cert": "",
    		    "key": "",
    		    "tls": "false"
    		}
		}
	`

	configWithTlsButWithoutCert = `
		{    
			"client_credentials" : {
        		"id": "TestClientId",
        		"scopes": [
        		    "scopeA",
        		    "scopeB"
        		],
        		"secret": "TestClientSecret",
        		"url": "https://localhost"
			},
    		"server": {
    		    "cert": "",
    		    "key": "",
    		    "tls": "true"
    		}
		}
	`
)

func createTempConfig(content string) {
	f, err := os.Create("test_temp.json")

	if err != nil {
		log.Fatal(err)
	}

	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Println("couldn't close test temp file")
		}
	}(f)

	_, err2 := f.WriteString(content)

	if err2 != nil {
		log.Fatal(err2)
	}

}

func removeTempConfig() {
	err := os.Remove("test_temp.json")

	if err != nil {
		log.Fatal(err)
	}
}

func TestNewValidConfig(t *testing.T) {
	t.Run("Create valid config from environment", func(t *testing.T) {
		createTempConfig(validConfig)

		got, err := NewConfig("./test_temp.json")

		want := &Config{
			ClientCredentials: ClientCredentials{
				Id:     "TestClientId",
				Scopes: []string{"scopeA", "scopeB"},
				Secret: "TestClientSecret",
				Url:    "https://localhost",
			},
			Server: Server{
				Cert: "foo.cert",
				Key:  "foo.key",
				Tls:  true,
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
			createTempConfig(tt.testCase)
			if _, err := NewConfig("./test_temp.json"); (err != nil) != tt.wantErr {
				t.Errorf("Expected %v, wantErr %v", err, tt.wantErr)
			}
			removeTempConfig()
		})
	}

}
