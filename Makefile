.PHONY: migrate run build clean swagger test test-unit test-integration test-coverage

migrate:
	go run cmd/migrate/main.go

run:
	go run cmd/server/main.go

build:
	go build -o bin/go-seo cmd/server/main.go

clean:
	rm -rf bin/

deps:
	go mod tidy
	go mod download

swagger:
	swag init -g cmd/server/main.go -o docs/

test:
	go test ./...

test-unit:
	go test ./tests/unit/...

test-integration:
	go test ./tests/integration/...

test-coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

env:
	cp .env.example .env
