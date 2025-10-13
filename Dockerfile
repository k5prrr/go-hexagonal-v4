# Dockerfile
FROM docker.io/library/golang:1.23-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o bin/main ./cmd/app

# --- Финальный образ ---
FROM docker.io/library/alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/bin/main .
CMD ["./main"]
