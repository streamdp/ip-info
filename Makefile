export CGO_ENABLED=0

all: generate test build

build:
	go build -mod=readonly -o bin/app -ldflags="-X github.com/streamdp/ip-info/config.Version=$$IP_INFO_VERSION" ./cmd

generate:
	go generate ./...

test:
	go clean -testcache
	go test -cover ./...

clean:
	rm -rf bin/