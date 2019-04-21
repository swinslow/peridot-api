FROM golang:1.12

RUN mkdir -p /obsidian-api
WORKDIR /obsidian-api

ADD . /obsidian-api

RUN go get -v ./...
RUN go build
RUN go install github.com/swinslow/obsidian-api
