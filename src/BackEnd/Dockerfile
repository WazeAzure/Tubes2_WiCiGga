FROM golang:latest

RUN mkdir -p app/

WORKDIR /app


RUN go mod init backend/web-api-gin
COPY . ./

RUN GOPROXY="https://goproxy.io" go mod tidy

RUN go build -o server server.go