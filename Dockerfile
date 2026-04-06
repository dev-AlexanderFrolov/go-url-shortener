# Этап 1: Собираем бинарник (Builder)
FROM golang:1.26.1-alpine AS builder
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN go build -o shortener .

# Этап 2: Финальный образ (Runner)
FROM alpine:latest
WORKDIR /app

# Забираем готовый бинарник
COPY --from=builder /app/shortener .

# НОВОЕ: Забираем папку с миграциями!
COPY --from=builder /app/migrations ./migrations

EXPOSE 8080
CMD ["./shortener"]