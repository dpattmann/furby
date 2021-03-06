# What is Furby?

[![made-with-Go](https://img.shields.io/badge/Made%20with-Go-1f425f)](http://golang.org)
[![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/dpattmann/furby)](https://github.com/dpattman/furby)

Furby is an easy and simple to use OAuth2 Token cache. Run Furby in your microservice
infrastructure to move the OAuth2.0 Token management out of your services.

## Build

```bash
go build cmd/furby/furby.go
```

## Configuration

### Configuration by configuration file

Configuration can be passed with json or yaml file by command line argument "--config" or "-c". See example configs.

## Config

## Example Config

```yaml
---
stores:
  - interval: 5 # time in minutes
    path: /token # Handler path
    credentials:
      id: "ClientId"
      scopes: []
      secret: "ClientSecret"
      url: "https://oauth.server/oauth2/token"
    auth:
      type: "noop|user-agent|header"
      user-agents: [] # required if type is user-agent
      header-name: "" # required if type is header
      header-value: [] # required if type is header

server:
    addr: ":8080" # default
    cert: ""
    key: ""
    tls: false

```

## Authorizer

| Name | Value | Description | 
| --- | --- | --- |
| noop | --- | Default authorizer |
| user-agent | User-Agents (case insensitive) | Restrict access to specified user-agents |
| header | Header name and values | Restrict access by specifying own header and values  |

## Testing

```bash
go test -v ./...
go test -v ./... -bench=.
```
### Generate mocks

```bash
mockery --all
```

## Development environment

### Active development
Use it in case you need some token creating software as docker compose environment but without furby in it.
```bash
docker-compose up -d hydra
```
