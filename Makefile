.PHONY: build test lint run-dev

build:
	cd server && go build ./cmd/bastion-server

test:
	cd server && go test ./...

lint:
	cd server && go vet ./...

run-dev:
	cd server && go run ./cmd/bastion-server -ui-dir ../ui
