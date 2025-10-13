# example https://github.com/zeromicro/go-zero/blob/master/tools/goctl/Dockerfile
# https://habr.com/ru/companies/otus/articles/660301/
# https://habr.com/ru/articles/647255/

FROM golang:1.21-alpine as builder
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/bin/main ./cmd/app

FROM gcr.io/distroless/static-debian11
COPY --from=builder /app/bin/main /main
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ENTRYPOINT ["/main"]
