FROM golang:1.24.2-alpine as build

ARG VERSION
ENV CGO_ENABLED=0

RUN apk update && apk upgrade

WORKDIR /build
ADD . /build
RUN go mod download
RUN go build -mod=readonly -o app -ldflags="-s -X github.com/streamdp/ip-info/config.Version=$VERSION" ./cmd

FROM gcr.io/distroless/static-debian12

WORKDIR /srv
COPY --from=build /build/app /srv/app

CMD ["/srv/app"]