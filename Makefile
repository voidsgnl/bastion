.PHONY: build test lint run-dev run-nmap-lab run-nmap-lab-postgres db-up db-down run-postgres

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

run-nmap-lab-postgres:
	cd server && go run ./cmd/bastion-server -ui-dir ../ui -engine nmap-lab -store postgres

db-up:
	docker compose -f deploy/compose/docker-compose.yml up -d db

db-down:
	docker compose -f deploy/compose/docker-compose.yml stop db

run-postgres:
	cd server && go run ./cmd/bastion-server -ui-dir ../ui -store postgres
