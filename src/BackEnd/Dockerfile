FROM golang:latest

RUN mkdir -p /app
WORKDIR /app

# RUN mkdir -p ./src
# RUN mkdir -p ./src/backend

COPY Dockerfile ./
COPY go.mod go.sum ./
RUN go mod download

COPY server.go  ./

WORKDIR /usr/local/go

COPY algorithm  ./src/backend/
COPY caching  ./src/backend/
COPY scraper  ./src/backend/
COPY util  ./src/backend/

WORKDIR /app
RUN go build -o server .

# FROM alpine:3.19
