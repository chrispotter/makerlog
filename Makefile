.PHONY: help install build up down logs clean migrate-up migrate-down migrate-create api-dev web-dev

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-20s %s\n", $$1, $$2}'

install: ## Install dependencies for all services
	@echo "Installing API dependencies..."
	cd services/api && go mod download
	@echo "Installing Web dependencies..."
	cd apps/web && npm install
	@echo "All dependencies installed!"

build: ## Build all Docker containers
	docker-compose build

up: ## Start all services with Docker Compose
	docker-compose up -d
	@echo "Services started! API: http://localhost:8080, Web: http://localhost:3000"

down: ## Stop all services
	docker-compose down

logs: ## Show logs from all services
	docker-compose logs -f

logs-api: ## Show logs from API service
	docker-compose logs -f api

logs-web: ## Show logs from web service
	docker-compose logs -f web

logs-db: ## Show logs from database service
	docker-compose logs -f postgres

clean: ## Stop services and remove volumes
	docker-compose down -v
	@echo "All services stopped and volumes removed!"

# Database migrations
migrate-up: ## Run database migrations up
	cd services/api && goose -dir migrations postgres "postgres://postgres:postgres@localhost:5432/makerlog?sslmode=disable" up

migrate-down: ## Run database migrations down
	cd services/api && goose -dir migrations postgres "postgres://postgres:postgres@localhost:5432/makerlog?sslmode=disable" down

migrate-status: ## Check migration status
	cd services/api && goose -dir migrations postgres "postgres://postgres:postgres@localhost:5432/makerlog?sslmode=disable" status

migrate-create: ## Create a new migration (usage: make migrate-create name=my_migration)
	cd services/api && goose -dir migrations create $(name) sql

# Development commands (without Docker)
api-dev: ## Run API in development mode
	cd services/api && go run cmd/api/main.go

web-dev: ## Run web app in development mode
	cd apps/web && npm run dev

api-build: ## Build API binary
	cd services/api && go build -o bin/makerlog cmd/api/main.go

# Database commands
db-start: ## Start only the database
	docker-compose up -d postgres

db-stop: ## Stop the database
	docker-compose stop postgres

db-shell: ## Open PostgreSQL shell
	docker-compose exec postgres psql -U postgres -d makerlog

# Testing
test-api: ## Run API tests
	cd services/api && go test ./...

test-web: ## Run web tests
	cd apps/web && npm test

# Linting
lint-api: ## Lint API code
	cd services/api && go fmt ./...
	cd services/api && go vet ./...

lint-web: ## Lint web code
	cd apps/web && npm run lint
