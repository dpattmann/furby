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
go test -v ./...
go test -v ./... -bench=.
```
### Generate mocks

```bash
mockery --all
```

## Development environment

### active development
Use it in case you need some token creating software as docker compose environment but without furby in it.
```bash
docker-compose up -d hydra
```

### passive development
Use it in case you need furby and some token creating software as docker compose environment
```bash
docker-compose up --build
```
