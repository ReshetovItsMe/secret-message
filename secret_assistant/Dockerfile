# syntax=docker/dockerfile:1

FROM golang:1.27.0-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN CGO_ENABLED=0 go build -o /secret-assistant ./cmd

FROM alpine:3.21
ENV RSA_PRIVATE_KEY_PATH=/run/secrets/secret_assistant_private_key
COPY --from=build /secret-assistant /usr/local/bin/secret-assistant
EXPOSE 50052
CMD ["/usr/local/bin/secret-assistant"]
