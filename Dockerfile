FROM golang:1.18

RUN mkdir -p /var/www/iot.go
RUN go install github.com/cosmtrek/air@latest

COPY . /var/www/iot.go

WORKDIR /var/www/iot.go

