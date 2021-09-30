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

## Testing

```bash
go test ./...
```
### Generate mocks

```bash
mockery --all
```