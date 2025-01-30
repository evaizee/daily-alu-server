.PHONY: all deps clean build run test lint migrate-up migrate-down

# Default config file
CONFIG ?= config/config.yaml

# Database connection string (update according to your config)
DB_URL = postgres://$(shell go run main.go config get database.user):$(shell go run main.go config get database.password)@$(shell go run main.go config get database.host):$(shell go run main.go config get database.port)/$(shell go run main.go config get database.name)?sslmode=$(shell go run main.go config get database.sslmode)

all: deps build

# Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod download
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	go clean

# Build the application
build:
	@echo "Building application..."
	go build -o bin/server main.go

# Run the application
run:
	@echo "Running server..."
	go run main.go serve --config $(CONFIG)

# Run the application in debug mode
debug:
	@echo "Running server in debug mode..."
	dlv debug main.go -- serve --config $(CONFIG)

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Run linter
lint:
	@echo "Running linter..."
	go vet ./...
	test -z $(shell gofmt -l .)

# Database migrations
migrate-up:
	@echo "Running database migrations..."
	$(HOME)/go/bin/migrate -path database/migrations -database "$(DB_URL)" up

migrate-down:
	@echo "Rolling back database migrations..."
	$(HOME)/go/bin/migrate -path database/migrations -database "$(DB_URL)" down

migrate-create:
	@read -p "Enter migration name: " name; \
	$(HOME)/go/bin/migrate create -ext sql -dir database/migrations -seq $$name

migrate-status:
	@echo "Migration status:"
	$(HOME)/go/bin/migrate -path database/migrations -database "$(DB_URL)" version

# Help command
help:
	@echo "Available commands:"
	@echo "  make deps          - Install project dependencies"
	@echo "  make build         - Build the application"
	@echo "  make run          - Run the application"
	@echo "  make debug        - Run the application in debug mode"
	@echo "  make test         - Run tests"
	@echo "  make lint         - Run linter"
	@echo "  make migrate-up    - Run database migrations"
	@echo "  make migrate-down  - Rollback database migrations"
	@echo "  make migrate-create - Create a new migration file"
	@echo "  make migrate-status - Show migration status"
	@echo "  make clean         - Clean build artifacts"
	@echo "  make help          - Show this help message"
