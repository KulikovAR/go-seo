FROM golang:1.24.3-alpine AS builder

RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main cmd/server/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

RUN adduser -D -s /bin/sh appuser

WORKDIR /root/

COPY --from=builder /app/main .

COPY --from=builder /app/docs ./docs

RUN chown -R appuser:appuser /root

USER appuser

EXPOSE 8080

CMD ["./main"]
