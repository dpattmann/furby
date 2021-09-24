# What is Furby?

Furby is an easy and simple to use OAuth2 Token cache.

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

```bash
# Generate mock interfaces
mockery --all
```