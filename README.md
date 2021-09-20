# Furby

Furby is an easy and simple to use OAuth2 Token cache.

## Testing

```bash
go test ./...
```

```bash
# Generate mock interfaces

mockgen -source store/store.go -destination store/mock/mock_store.go -package mock -mock_names Store=Store
mockgen -source auth/auth.go -destination auth/mock/mock_auth.go -package mock -mock_names Authorization=Authorization
```