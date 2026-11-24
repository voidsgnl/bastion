.PHONY: build test lint run-dev run-nmap-lab

build:
	cd server && go build ./cmd/bastion-server

test:
	cd server && go test ./...

lint:
	cd server && go vet ./...

run-dev:
	cd server && go run ./cmd/bastion-server -ui-dir ../ui

run-nmap-lab:
	cd server && go run ./cmd/bastion-server -ui-dir ../ui -engine nmap-lab
