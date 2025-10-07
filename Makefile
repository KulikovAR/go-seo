.PHONY: migrate run build clean swagger test test-unit test-integration test-coverage

# Запуск миграций
migrate:
	go run cmd/migrate/main.go

# Запуск сервера
run:
	go run cmd/server/main.go

# Сборка приложения
build:
	go build -o bin/go-seo cmd/server/main.go

# Очистка
clean:
	rm -rf bin/

# Установка зависимостей
deps:
	go mod tidy
	go mod download

# Генерация Swagger документации
swagger:
	swag init -g cmd/server/main.go -o docs/

# Запуск всех тестов
test:
	go test ./...

# Запуск unit тестов
test-unit:
	go test ./tests/unit/...

# Запуск интеграционных тестов
test-integration:
	go test ./tests/integration/...

# Запуск тестов с покрытием
test-coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Создание .env файла из примера
env:
	cp .env.example .env
