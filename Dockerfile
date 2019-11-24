# SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

FROM golang:1.13

RUN mkdir -p /peridot-api
WORKDIR /peridot-api

ADD . /peridot-api

RUN go get -v ./...
RUN go build
RUN go install github.com/swinslow/peridot-api
