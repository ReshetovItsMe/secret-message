# syntax=docker/dockerfile:1

# ---- Build stage ----
FROM golang:1.27-alpine AS build
WORKDIR /app/gateway

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o /gateway ./cmd/gateway

# ---- Runtime stage ----
FROM alpine:3.21
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=build /gateway /usr/local/bin/gateway

EXPOSE 3000

CMD ["gateway"]
