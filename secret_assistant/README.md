# Secret Assistant

The Secret Assistant service is a data encryption and decryption service built in Golang. It is designed to provide secure encryption and decryption capabilities for sensitive data.

## gRPC Endpoints

The Secret Assistant service provides the following API endpoints:

- `/encrypt`: Accepts data as input and encrypts it using the specified encryption key.
- `/decrypt`: Accepts encrypted data as input and decrypts it using the specified encryption key.

## Lint

Use Go 1.27.0 and golangci-lint v2.15.0 on your `PATH`. From `secret_assistant/`, run:

```bash
make lint
# Equivalent: golangci-lint run
```

Both commands use the shared `../.golangci.yml` configuration.

## Generate proto run (if changed)

```bash
export GO_PATH=~/go
export PATH=$PATH:/$GO_PATH/bin
```

```bash
protoc --proto_path=../proto \
 --go_out=./proto --go_opt=paths=source_relative \
 --go-grpc_out=./proto --go-grpc_opt=paths=source_relative \
 message.proto
```
