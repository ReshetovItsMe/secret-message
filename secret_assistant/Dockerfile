# syntax=docker/dockerfile:1

FROM golang:1.18-alpine

WORKDIR /app/secret-assistant

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY ./**/*.go ./

RUN echo $(ls)

RUN go build -o main

EXPOSE 8080

RUN echo $(ls)

RUN ["chmod", "+x", "./main"]
