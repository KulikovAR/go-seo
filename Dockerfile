# Многоэтапная сборка для оптимизации размера образа
FROM golang:1.24.3-alpine AS builder

# Установка необходимых пакетов
RUN apk add --no-cache git ca-certificates tzdata

# Создание рабочей директории
WORKDIR /app

# Копирование go.mod и go.sum для кэширования зависимостей
COPY go.mod go.sum ./

# Загрузка зависимостей
RUN go mod download

# Копирование исходного кода
COPY . .

# Обновление зависимостей и создание go.sum
RUN go mod tidy

# Сборка приложения
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main cmd/server/main.go

# Финальный образ
FROM alpine:latest

# Установка ca-certificates для HTTPS запросов
RUN apk --no-cache add ca-certificates tzdata

# Создание пользователя для безопасности
RUN adduser -D -s /bin/sh appuser

WORKDIR /root/

# Копирование бинарного файла из builder
COPY --from=builder /app/main .

# Копирование документации Swagger
COPY --from=builder /app/docs ./docs

# Смена владельца файлов
RUN chown -R appuser:appuser /root

# Переключение на непривилегированного пользователя
USER appuser

# Открытие порта
EXPOSE 8080

# Команда запуска
CMD ["./main"]
