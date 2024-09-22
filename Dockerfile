FROM golang:1.23.0-alpine as build
ARG VERSION
WORKDIR /build
ADD . /build
COPY go.* ./
RUN go get ./...

RUN apk update && apk upgrade
RUN go build -mod=readonly -o app -ldflags="-X github.com/streamdp/ip-info/config.Version=$VERSION" ./cmd

FROM alpine:3.20.3

COPY --from=build /build/app 	    /srv/app

WORKDIR /srv

CMD ["/srv/app"]
