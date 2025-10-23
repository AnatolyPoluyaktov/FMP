# FMP Makefile
# Автоматизация сборки и генерации кода

.PHONY: help generate build test clean deps install-tools

# Default target
help:
	@echo "FMP - Financial Manager Platform"
	@echo ""
	@echo "Available targets:"
	@echo "  generate     - Generate code from OpenAPI spec"
	@echo "  build        - Build all applications"
	@echo "  test         - Run tests for all applications"
	@echo "  clean        - Clean build artifacts"
	@echo "  deps         - Install dependencies"
	@echo "  install-tools - Install required tools"
	@echo "  dev          - Start development environment"
	@echo "  deploy       - Deploy to production"

# Install required tools
install-tools:
	@echo "Installing required tools..."
	go install github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen@latest
	go install github.com/swaggo/swag/cmd/swag@latest
	go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Install dependencies
deps:
	@echo "Installing dependencies..."
	cd fmp-core && go mod tidy
	cd minapp/backend && go mod tidy
	cd minapp/frontend && npm install
	cd fmp-analytics && npm install

# Generate code from OpenAPI spec
generate:
	@echo "Generating code from OpenAPI specification..."
	./scripts/generate.sh

# Build all applications
build: generate
	@echo "Building applications..."
	cd fmp-core && go build -o bin/fmp-core .
	cd minapp/backend && go build -o bin/fmp-minapp-backend .
	cd minapp/frontend && npm run build
	cd fmp-analytics && npm run build

# Run tests
test:
	@echo "Running tests..."
	cd fmp-core && go test ./...
	cd minapp/backend && go test ./...
	cd minapp/frontend && npm test -- --watchAll=false
	cd fmp-analytics && npm test -- --watchAll=false

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf fmp-core/bin/
	rm -rf minapp/backend/bin/
	rm -rf minapp/frontend/build/
	rm -rf fmp-analytics/build/
	rm -rf generated/
	rm -rf fmp-core/docs/

# Development environment
dev:
	@echo "Starting development environment..."
	@echo "Make sure PostgreSQL is running and database is created"
	@echo "Run 'make dev-core' in one terminal and 'make dev-minapp' in another"

dev-core:
	cd fmp-core && go run main.go

dev-minapp-backend:
	cd minapp/backend && go run main.go

dev-minapp-frontend:
	cd minapp/frontend && npm start

dev-analytics:
	cd fmp-analytics && npm start

# Docker commands
docker-up:
	./docker.sh up

docker-down:
	./docker.sh down

docker-restart:
	./docker.sh restart

docker-logs:
	./docker.sh logs

docker-build:
	./docker.sh build

docker-clean:
	./docker.sh clean

docker-status:
	./docker.sh status

docker-shell:
	./docker.sh shell

docker-db:
	./docker.sh db

docker-redis:
	./docker.sh redis

# Database operations
db-migrate:
	cd fmp-core && migrate -path migrations -database "postgres://user:password@localhost/fmp?sslmode=disable" up

db-rollback:
	cd fmp-core && migrate -path migrations -database "postgres://user:password@localhost/fmp?sslmode=disable" down

db-create:
	createdb fmp

# Docker operations
docker-build:
	docker-compose build

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

# Deployment
deploy:
	cd deploy && ./deploy.sh all

# Linting
lint:
	cd fmp-core && golangci-lint run
	cd minapp/backend && golangci-lint run
	cd minapp/frontend && npm run lint
	cd fmp-analytics && npm run lint

# Format code
fmt:
	cd fmp-core && go fmt ./...
	cd minapp/backend && go fmt ./...
	cd minapp/frontend && npm run format
	cd fmp-analytics && npm run format
