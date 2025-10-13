# make run // use
CONTAINER_ENGINE = podman#docker
COMPOSE = $(CONTAINER_ENGINE)-compose
BINARY = app
GOPATH = $(HOME)/go
LINT_PATH = $(GOPATH)/bin/golangci-lint

up:
	$(COMPOSE) down
	$(COMPOSE) up -d --build
	$(CONTAINER_ENGINE) ps -a

build:
	go build -o bin/main ./cmd/app

run:
	go fmt ./...
	clear
	go run cmd/app/main.go || echo "exit code: $?"

test:
	go test -v -cover ./...


lint: $(LINT_TARGET)
	@echo "==> Linting Go code..."
	@$(LINT_PATH) run --config ./config/.golangci.yml ./internal/... ./cmd/...

installLint:
	@echo "==> Installing golangci-lint..."
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

installDebugging:
	go install github.com/go-delve/delve/cmd/dlv@latest




# $HOME/go/bin/golangci-lint --version
# or $GOPATH/bin/golangci-lint --version
# or $(GOPATH)/bin/golangci-lint

# # Libs
# go mod init app
# Postgres | go get -u gorm.io/driver/postgres"
# Kafka by confluentinc | go get -u github.com/confluentinc/confluent-kafka-go/kafka && go get -u github.com/confluentinc/confluent-kafka-go/v2/kafka"


# # All Libs
# Gin Server | go get -u github.com/gin-gonic/gin"
# MySQL | go get -u github.com/go-sql-driver/mysql"
# BD gorm | go get -u gorm.io/gorm"
# Token jwt | go get -u github.com/golang-jwt/jwt"
# Hash bcrypt | go get -u golang.org/x/crypto/bcrypt"
