PREFIX ?= $(HOME)/.local
BIN := bin/repos-sync

.PHONY: build install test vet lint release

build:
	go build -o $(BIN) ./cmd/repos-sync

install: build
	install -Dm755 $(BIN) $(PREFIX)/bin/repos-sync

test:
	go test ./...

vet:
	go vet ./...

lint:
	test -z "$$(gofmt -l .)"
	go vet ./...

release:
	GOOS=linux GOARCH=amd64 go build -o dist/repos-sync-linux-amd64 ./cmd/repos-sync
	GOOS=linux GOARCH=arm64 go build -o dist/repos-sync-linux-arm64 ./cmd/repos-sync
	sha256sum dist/repos-sync-linux-* > dist/SHA256SUMS
