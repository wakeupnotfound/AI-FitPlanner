.PHONY: help build run test clean deps migrate

help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

deps: ## Download dependencies
	go mod download
	go mod tidy

build: ## Build the application
	go build -o bin/api cmd/api/main.go

run: ## Run the application
	go run cmd/api/main.go

test: ## Run tests
	go test -v -race -coverprofile=coverage.out ./...

test-coverage: test ## Run tests with coverage report
	go tool cover -html=coverage.out -o coverage.html

clean: ## Clean build artifacts
	rm -rf bin/
	rm -rf logs/
	rm -f coverage.out coverage.html

migrate: ## Run database migrations
	@echo "Running database migration..."
	@bash scripts/migrate.sh

migrate-manual: ## Run database migrations manually with MySQL client
	mysql -h localhost -u fitness_user -p fitness_planner < ../database/schema.sql

docker-up: ## Start Docker containers
	docker-compose up -d

docker-down: ## Stop Docker containers
	docker-compose down

docker-logs: ## View Docker logs
	docker-compose logs -f

lint: ## Run linter
	golangci-lint run

fmt: ## Format code
	go fmt ./...
	goimports -w .

swagger: ## Generate Swagger documentation
	swag init -g cmd/api/main.go -o docs --parseDependency --parseInternal

swagger-install: ## Install Swagger CLI tool
	go install github.com/swaggo/swag/cmd/swag@latest

.DEFAULT_GOAL := help
