# What is Furby?

Furby is an easy and simple to use OAuth2 Token cache. Run Furby in your microservice
infrastructure to move the OAuth2.0 Token management out of your services.

## Build

```bash
go build cmd/furby/furby.go
```

## Configuration

### Configuration by configuration file

Configuration can be passed with json or yaml file by command line argument "--path" or "-p". See example configs.

## Config

## Example Config

```yaml
---
auth:
    type: "noop|user-agent|header"
    user-agents: [] # required if type is user-agent
    header-name: "" # required if type is header
    header-value: [] # required if type is header

credentials:
    id: "ClientId"
    scopes: []
    secret: "ClientSecret"
    url: "https://oauth.server/oauth2/token"

store:
    interval: 300 # default

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
