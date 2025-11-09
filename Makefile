# ===========================================
# Ezra API - Makefile
# ===========================================

.PHONY: help build run test clean docker-build docker-up docker-down migrate-up migrate-down lint format swagger deps

# Variables
APP_NAME=ezra-api
VERSION?=1.0.0
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
GIT_COMMIT=$(shell git rev-parse --short HEAD)
DOCKER_IMAGE=$(APP_NAME):$(VERSION)

# Colors for output
RED=\033[0;31m
GREEN=\033[0;32m
YELLOW=\033[1;33m
NC=\033[0m # No Color

# Default target
.DEFAULT_GOAL := help

## help: Show this help message
help:
	@echo '$(YELLOW)Ezra API - Available Commands:$(NC)'
	@echo ''
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'
	@echo ''

## deps: Download Go dependencies
deps:
	@echo '$(GREEN)Downloading dependencies...$(NC)'
	go mod download
	go mod verify
	@echo '$(GREEN)✓ Dependencies downloaded$(NC)'

## build: Build the application binary
build:
	@echo '$(GREEN)Building $(APP_NAME)...$(NC)'
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
		-ldflags="-w -s \
		-X 'main.Version=$(VERSION)' \
		-X 'main.BuildTime=$(BUILD_TIME)' \
		-X 'main.GitCommit=$(GIT_COMMIT)'" \
		-o bin/$(APP_NAME) \
		./cmd/main.go
	@echo '$(GREEN)✓ Build complete: bin/$(APP_NAME)$(NC)'

## run: Run the application locally
run:
	@echo '$(GREEN)Running $(APP_NAME)...$(NC)'
	go run ./cmd/main.go

## dev: Run the application in development mode with hot reload
dev:
	@echo '$(GREEN)Running in development mode...$(NC)'
	air -c .air.toml

## test: Run all tests
test:
	@echo '$(GREEN)Running tests...$(NC)'
	go test -v -race -coverprofile=coverage.out ./...
	@echo '$(GREEN)✓ Tests complete$(NC)'

## test-coverage: Run tests with coverage report
test-coverage: test
	@echo '$(GREEN)Generating coverage report...$(NC)'
	go tool cover -html=coverage.out -o coverage.html
	@echo '$(GREEN)✓ Coverage report: coverage.html$(NC)'

## test-unit: Run unit tests only
test-unit:
	@echo '$(GREEN)Running unit tests...$(NC)'
	go test -v -short ./...

## test-integration: Run integration tests only
test-integration:
	@echo '$(GREEN)Running integration tests...$(NC)'
	go test -v -run Integration ./...

## bench: Run benchmarks
bench:
	@echo '$(GREEN)Running benchmarks...$(NC)'
	go test -bench=. -benchmem ./...

## lint: Run linter
lint:
	@echo '$(GREEN)Running linter...$(NC)'
	golangci-lint run ./...
	@echo '$(GREEN)✓ Linting complete$(NC)'

## format: Format code
format:
	@echo '$(GREEN)Formatting code...$(NC)'
	go fmt ./...
	goimports -w .
	@echo '$(GREEN)✓ Code formatted$(NC)'

## vet: Run go vet
vet:
	@echo '$(GREEN)Running go vet...$(NC)'
	go vet ./...
	@echo '$(GREEN)✓ Vet complete$(NC)'

## swagger: Generate Swagger documentation
swagger:
	@echo '$(GREEN)Generating Swagger docs...$(NC)'
	swag init -g cmd/main.go --output docs
	@echo '$(GREEN)✓ Swagger docs generated$(NC)'

## clean: Clean build artifacts
clean:
	@echo '$(YELLOW)Cleaning build artifacts...$(NC)'
	rm -rf bin/
	rm -rf tmp/
	rm -f coverage.out coverage.html
	@echo '$(GREEN)✓ Clean complete$(NC)'

## docker-build: Build Docker image
docker-build:
	@echo '$(GREEN)Building Docker image: $(DOCKER_IMAGE)...$(NC)'
	docker build \
		--build-arg VERSION=$(VERSION) \
		--build-arg BUILD_TIME=$(BUILD_TIME) \
		--build-arg GIT_COMMIT=$(GIT_COMMIT) \
		-t $(DOCKER_IMAGE) \
		-t $(APP_NAME):latest \
		.
	@echo '$(GREEN)✓ Docker image built: $(DOCKER_IMAGE)$(NC)'

## docker-up: Start all services with Docker Compose
docker-up:
	@echo '$(GREEN)Starting services...$(NC)'
	docker-compose up -d
	@echo '$(GREEN)✓ Services started$(NC)'

## docker-up-full: Start all services including optional ones
docker-up-full:
	@echo '$(GREEN)Starting all services (including optional)...$(NC)'
	docker-compose --profile full --profile nginx up -d
	@echo '$(GREEN)✓ All services started$(NC)'

## docker-down: Stop all services
docker-down:
	@echo '$(YELLOW)Stopping services...$(NC)'
	docker-compose down
	@echo '$(GREEN)✓ Services stopped$(NC)'

## docker-down-volumes: Stop all services and remove volumes
docker-down-volumes:
	@echo '$(RED)Stopping services and removing volumes...$(NC)'
	docker-compose down -v
	@echo '$(GREEN)✓ Services stopped and volumes removed$(NC)'

## docker-logs: View logs from all services
docker-logs:
	docker-compose logs -f

## docker-logs-api: View logs from API service only
docker-logs-api:
	docker-compose logs -f api

## docker-restart: Restart all services
docker-restart: docker-down docker-up

## migrate-up: Run database migrations
migrate-up:
	@echo '$(GREEN)Running database migrations...$(NC)'
	docker exec -i ezra-postgres psql -U postgres -d ezradb < migrate/000000_postgres.up.sql
	@echo '$(GREEN)✓ Migrations applied$(NC)'

## migrate-down: Rollback database migrations
migrate-down:
	@echo '$(RED)Rolling back database migrations...$(NC)'
	docker exec -i ezra-postgres psql -U postgres -d ezradb < migrate/000000_postgres.down.sql
	@echo '$(GREEN)✓ Migrations rolled back$(NC)'

## db-connect: Connect to the database
db-connect:
	docker exec -it ezra-postgres psql -U postgres -d ezradb

## db-backup: Create database backup
db-backup:
	@echo '$(GREEN)Creating database backup...$(NC)'
	mkdir -p backups
	docker exec ezra-postgres pg_dump -U postgres ezradb > backups/backup_$(shell date +%Y%m%d_%H%M%S).sql
	@echo '$(GREEN)✓ Backup created in backups/$(NC)'

## db-restore: Restore database from backup (requires BACKUP_FILE variable)
db-restore:
	@echo '$(YELLOW)Restoring database from $(BACKUP_FILE)...$(NC)'
	docker exec -i ezra-postgres psql -U postgres -d ezradb < $(BACKUP_FILE)
	@echo '$(GREEN)✓ Database restored$(NC)'

## install-tools: Install development tools
install-tools:
	@echo '$(GREEN)Installing development tools...$(NC)'
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/swaggo/swag/cmd/swag@latest
	go install golang.org/x/tools/cmd/goimports@latest
	go install github.com/cosmtrek/air@latest
	@echo '$(GREEN)✓ Tools installed$(NC)'

## security-scan: Run security vulnerability scan
security-scan:
	@echo '$(GREEN)Running security scan...$(NC)'
	gosec ./...
	@echo '$(GREEN)✓ Security scan complete$(NC)'

## check: Run all checks (format, vet, lint, test)
check: format vet lint test
	@echo '$(GREEN)✓ All checks passed$(NC)'

## ci: Run CI pipeline locally
ci: deps check test-coverage
	@echo '$(GREEN)✓ CI pipeline complete$(NC)'

## prod-build: Build for production
prod-build: clean deps test build
	@echo '$(GREEN)✓ Production build complete$(NC)'

## health-check: Check if services are healthy
health-check:
	@echo '$(GREEN)Checking service health...$(NC)'
	@curl -f http://localhost:8080/health || echo '$(RED)API is not healthy$(NC)'
	@docker exec ezra-postgres pg_isready -U postgres || echo '$(RED)Database is not healthy$(NC)'

## version: Display version information
version:
	@echo 'Version:    $(VERSION)'
	@echo 'Build Time: $(BUILD_TIME)'
	@echo 'Git Commit: $(GIT_COMMIT)'