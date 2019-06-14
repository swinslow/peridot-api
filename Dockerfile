FROM golang:1.12

RUN mkdir -p /peridot-api
WORKDIR /peridot-api

ADD . /peridot-api

RUN go get -v ./...
RUN go build
RUN go install github.com/swinslow/peridot-api
