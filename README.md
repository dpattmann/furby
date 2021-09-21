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

mockgen -source store/store.go -destination store/mock/mock_store.go -package mock -mock_names Store=Store
mockgen -source auth/auth.go -destination auth/mock/mock_auth.go -package mock -mock_names Authorization=Authorization
```