# syntax=docker/dockerfile:1

FROM golang:1.18-alpine
ENV APP_HOME /app/secret-assistant

WORKDIR "$APP_HOME"

COPY go.mod ./
COPY go.sum ./

RUN go mod download
RUN go mod verify

COPY . . 

WORKDIR "./cmd"
RUN go build -o /secret-assistant

EXPOSE 50052

CMD ["/secret-assistant"]
