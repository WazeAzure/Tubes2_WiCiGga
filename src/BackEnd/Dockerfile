FROM golang:latest

WORKDIR /usr/local/go/src

RUN mkdir -p ./backend
COPY go.mod ./
COPY go.sum ./
COPY server.go  ./
COPY algorithm  ./backend
COPY caching  ./backend
COPY scraper  ./backend
COPY util  ./backend

RUN go mod tidy
RUN go mod vendor

RUN go run server.go