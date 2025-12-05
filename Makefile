# ===========================================
# Ezra API - Makefile (Windows Compatible)
# ===========================================
# This version works better on Windows Command Prompt and PowerShell
# For colored output, use Git Bash or Windows Terminal with the regular Makefile
#
# Quick Start:
#   make help          - Show all available commands
#   make deps          - Download Go dependencies
#   make run           - Run the application locally
#   make docker-up     - Start all services with Docker Compose
#
# ===========================================

.PHONY: help build run test clean docker-build docker-up docker-down migrate-up migrate-down lint format swagger deps

# ===========================================
# Variables
# ===========================================

APP_NAME=ezra-api
VERSION?=1.0.0
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
GIT_COMMIT=$(shell git rev-parse --short HEAD)
DOCKER_IMAGE=$(APP_NAME):$(VERSION)

# ===========================================
# Default Target
# ===========================================

.DEFAULT_GOAL := help

# ===========================================
# Help & Documentation
# ===========================================

## help: Show this help message with all available commands
help:
	@echo Ezra API - Available Commands:
	@echo.
	@echo   help               - Show this help message
	@echo   version            - Display version information
	@echo.
	@echo Dependencies:
	@echo   deps               - Download Go dependencies
	@echo   install-tools      - Install development tools
	@echo.
	@echo Building:
	@echo   build              - Build application binary
	@echo   prod-build         - Production build (clean + deps + test + build)
	@echo.
	@echo Running:
	@echo   run                - Run application locally
	@echo   dev                - Run with hot reload
	@echo.
	@echo Testing:
	@echo   test               - Run all tests
	@echo   test-coverage      - Run tests with HTML coverage report
	@echo   test-unit          - Run unit tests only
	@echo   test-integration   - Run integration tests only
	@echo   bench              - Run benchmarks
	@echo.
	@echo Code Quality:
	@echo   lint               - Run linter
	@echo   format             - Format code
	@echo   vet                - Run go vet
	@echo   security-scan      - Run security scanner
	@echo   check              - Run all quality checks
	@echo.
	@echo Docker:
	@echo   docker-build       - Build Docker image
	@echo   docker-up          - Start services
	@echo   docker-up-full     - Start all services (including optional)
	@echo   docker-down        - Stop services
	@echo   docker-down-volumes - Stop services and delete volumes
	@echo   docker-logs        - Follow all service logs
	@echo   docker-logs-api    - Follow API logs only
	@echo   docker-restart     - Restart services
	@echo.
	@echo Database:
	@echo   migrate-up         - Apply database migrations
	@echo   migrate-down       - Rollback migrations
	@echo   db-connect         - Connect to database
	@echo   db-backup          - Create database backup
	@echo   db-restore         - Restore from backup (requires BACKUP_FILE=...)
	@echo.
	@echo Other:
	@echo   swagger            - Generate API documentation
	@echo   clean              - Remove build artifacts
	@echo   health-check       - Check service health
	@echo   ci                 - Run CI pipeline locally
	@echo.

## version: Display version information
version:
	@echo Version:    $(VERSION)
	@echo Build Time: $(BUILD_TIME)
	@echo Git Commit: $(GIT_COMMIT)

# ===========================================
# Dependency Management
# ===========================================

## deps: Download and verify Go module dependencies
deps:
	@echo [+] Downloading dependencies...
	go mod download
	go mod verify
	@echo [OK] Dependencies downloaded

## install-tools: Install development tools
install-tools:
	@echo [+] Installing development tools...
	@echo [+] Installing golangci-lint (linter)...
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo [+] Installing swag (Swagger docs generator)...
	go install github.com/swaggo/swag/cmd/swag@latest
	@echo [+] Installing goimports (code formatter)...
	go install golang.org/x/tools/cmd/goimports@latest
	@echo [+] Installing air (hot reload)...
	go install github.com/cosmtrek/air@latest
	@echo [OK] Tools installed

# ===========================================
# Building
# ===========================================

## build: Build the application binary for Linux
build:
	@echo [+] Building $(APP_NAME)...
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s -X 'main.Version=$(VERSION)' -X 'main.BuildTime=$(BUILD_TIME)' -X 'main.GitCommit=$(GIT_COMMIT)'" -o bin/$(APP_NAME) ./cmd/main.go
	@echo [OK] Build complete: bin/$(APP_NAME)

## prod-build: Complete production build
prod-build: clean deps test build
	@echo [OK] Production build complete

# ===========================================
# Running Locally
# ===========================================

## run: Run the application locally
run:
	@echo [+] Running $(APP_NAME)...
	go run ./cmd/main.go

## dev: Run with hot reload (requires air)
dev:
	@echo [+] Running in development mode...
	air -c .air.toml

# ===========================================
# Testing
# ===========================================

## test: Run all tests
test:
	@echo [+] Running tests...
	go test -v -race -coverprofile=coverage.out ./...
	@echo [OK] Tests complete

## test-coverage: Run tests with HTML coverage report
test-coverage: test
	@echo [+] Generating coverage report...
	go tool cover -html=coverage.out -o coverage.html
	@echo [OK] Coverage report: coverage.html

## test-unit: Run unit tests only
test-unit:
	@echo [+] Running unit tests...
	go test -v -short ./...

## test-integration: Run integration tests only
test-integration:
	@echo [+] Running integration tests...
	go test -v -run Integration ./...

## bench: Run benchmarks
bench:
	@echo [+] Running benchmarks...
	go test -bench=. -benchmem ./...

# ===========================================
# Code Quality
# ===========================================

## lint: Run golangci-lint
lint:
	@echo [+] Running linter...
	golangci-lint run ./...
	@echo [OK] Linting complete

## format: Format code
format:
	@echo [+] Formatting code...
	go fmt ./...
	goimports -w .
	@echo [OK] Code formatted

## vet: Run go vet
vet:
	@echo [+] Running go vet...
	go vet ./...
	@echo [OK] Vet complete

## security-scan: Run security scanner
security-scan:
	@echo [+] Running security scan...
	gosec ./...
	@echo [OK] Security scan complete

## check: Run all quality checks
check: format vet lint test
	@echo [OK] All checks passed

# ===========================================
# Documentation
# ===========================================

## swagger: Generate Swagger documentation
swagger:
	@echo [+] Generating Swagger docs...
	swag init -g cmd/main.go --output docs
	@echo [OK] Swagger docs generated

# ===========================================
# Cleanup
# ===========================================

## clean: Remove build artifacts
clean:
	@echo [+] Cleaning build artifacts...
	@if exist bin rmdir /s /q bin
	@if exist tmp rmdir /s /q tmp
	@if exist coverage.out del /q coverage.out
	@if exist coverage.html del /q coverage.html
	@echo [OK] Clean complete

# ===========================================
# Docker Operations
# ===========================================

## docker-build: Build Docker image
docker-build:
	@echo [+] Building Docker image: $(DOCKER_IMAGE)...
	docker build --build-arg VERSION=$(VERSION) --build-arg BUILD_TIME=$(BUILD_TIME) --build-arg GIT_COMMIT=$(GIT_COMMIT) -t $(DOCKER_IMAGE) -t $(APP_NAME):latest .
	@echo [OK] Docker image built: $(DOCKER_IMAGE)

## docker-up: Start services
docker-up:
	@echo [+] Starting services...
	docker compose up -d
	@echo [OK] Services started

## docker-up-full: Start all services including optional ones
docker-up-full:
	@echo [+] Starting all services...
	docker compose --profile full --profile nginx up -d
	@echo [OK] All services started

## docker-down: Stop services
docker-down:
	@echo [+] Stopping services...
	docker compose down
	@echo [OK] Services stopped

## docker-down-volumes: Stop services and remove volumes
docker-down-volumes:
	@echo [WARNING] Stopping services and removing volumes...
	docker compose down -v
	@echo [OK] Services stopped and volumes removed

## docker-logs: Follow all service logs
docker-logs:
	docker compose logs -f

## docker-logs-api: Follow API logs only
docker-logs-api:
	docker compose logs -f api

## docker-restart: Restart services
docker-restart: docker-down docker-up

# ===========================================
# Database Operations
# ===========================================

## migrate-up: Apply database migrations
migrate-up:
	@echo [+] Running database migrations...
	docker exec -i ezra-postgres psql -U postgres -d ezradb < migrate/000000_postgres.up.sql
	@echo [OK] Migrations applied

## migrate-down: Rollback database migrations
migrate-down:
	@echo [WARNING] Rolling back database migrations...
	docker exec -i ezra-postgres psql -U postgres -d ezradb < migrate/000000_postgres.down.sql
	@echo [OK] Migrations rolled back

## db-connect: Connect to database
db-connect:
	docker exec -it ezra-postgres psql -U postgres -d ezradb

## db-backup: Create database backup
db-backup:
	@echo [+] Creating database backup...
	@if not exist backups mkdir backups
	docker exec ezra-postgres pg_dump -U postgres ezradb > backups/backup_$(shell date +%Y%m%d_%H%M%S).sql
	@echo [OK] Backup created in backups/

## db-restore: Restore database from backup
db-restore:
	@echo [+] Restoring database from $(BACKUP_FILE)...
	docker exec -i ezra-postgres psql -U postgres -d ezradb < $(BACKUP_FILE)
	@echo [OK] Database restored

# ===========================================
# Health & Monitoring
# ===========================================

## health-check: Check service health
health-check:
	@echo [+] Checking service health...
	@curl -f http://localhost:8080/health || echo [ERROR] API is not healthy
	@docker exec ezra-postgres pg_isready -U postgres || echo [ERROR] Database is not healthy

# ===========================================
# CI/CD Pipeline
# ===========================================

## ci: Run CI pipeline locally
ci: deps check test-coverage
	@echo [OK] CI pipeline complete