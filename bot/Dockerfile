# Stage 1: Сборка бинарника
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o bot_binary main.go

# Stage 2: Запуск
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/bot_binary .
CMD ["./bot_binary"]