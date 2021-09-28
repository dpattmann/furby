# What is Furby?

Furby is an easy and simple to use OAuth2 Token cache. Run Furby in your microservice
infrastructure to move the OAuth2.0 Token management out of your services.

## Build

```bash
go build cmd/furby/furby.go
```

## Configuration

### Configuration by environment variables

| Name | Description |
|--- |--- |
|FURBY_CLIENTCREDENTIALS_ID | OAuth2 Client Id |
|FURBY_CLIENTCREDENTIALS_SECRET | OAuth2 Client Secret |
|FURBY_CLIENTCREDENTIALS_URL | Oauth2 Server Token Url |
|FURBY_CLIENTCREDENTIALS_SCOPES | OAauth2 Token Scopes |

## Testing

```bash
go test ./...
```
### Generate mocks

```bash
mockery --all
```