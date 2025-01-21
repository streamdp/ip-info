FROM golang:1.23.5-alpine as build

ARG VERSION
ENV CGO_ENABLED=0

RUN apk update && apk upgrade

WORKDIR /build
ADD . /build
RUN go mod download
RUN go build -mod=readonly -o app -ldflags="-s -X github.com/streamdp/ip-info/config.Version=$VERSION" ./cmd

# Adding statically compiled wget binaries to makes docker healthcheck possible when using a distroless base image.
ADD https://busybox.net/downloads/binaries/1.31.0-i686-uclibc/busybox_WGET ./wget
RUN chmod a+x ./wget

FROM gcr.io/distroless/static-debian12

WORKDIR /srv
COPY --from=build /build/wget /usr/bin/wget
COPY --from=build /build/app /srv/app

CMD ["/srv/app"]