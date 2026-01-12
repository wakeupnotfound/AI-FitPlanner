# AI Fitness Planner - Root Makefile
# This Makefile provides convenient commands for building and running the entire application

.PHONY: help build up down clean test lint deps

# Default target
help:
	@echo "AI Fitness Planner - Available commands:"
	@echo ""
	@echo "  build        - Build all services (backend and frontend)"
	@echo "  up           - Start all services with Docker Compose"
	@echo "  down         - Stop all services with Docker Compose"
	@echo "  test         - Run tests for both backend and frontend"
	@echo "  lint         - Run linting for both backend and frontend"
	@echo "  clean        - Clean build artifacts"
	@echo "  deps         - Install dependencies for all services"
	@echo ""
	@echo "Development commands:"
	@echo "  dev          - Start development environment (backend + frontend)"
	@echo "  dev-backend  - Start only backend in development mode"
	@echo "  dev-frontend - Start only frontend in development mode"
	@echo ""
	@echo "Production commands:"
	@echo "  prod         - Start production environment"
	@echo "  deploy       - Deploy to production (build + up)"

# Build all services
build:
	@echo "Building backend..."
	$(MAKE) -C backend build
	@echo "Building frontend..."
	$(MAKE) -C frontend build
	@echo "Build complete!"

# Start all services
up:
	docker-compose up -d

# Stop all services
down:
	docker-compose down

# Clean up
clean:
	@echo "Cleaning backend..."
	$(MAKE) -C backend clean
	@echo "Cleaning frontend..."
	$(MAKE) -C frontend clean
	docker-compose down -v --remove-orphans
	docker system prune -f

# Run all tests
test:
	@echo "Running backend tests..."
	$(MAKE) -C backend test
	@echo "Running frontend tests..."
	$(MAKE) -C frontend test:unit
	@echo "All tests complete!"

# Run linting
lint:
	@echo "Linting backend..."
	$(MAKE) -C backend lint
	@echo "Linting frontend..."
	$(MAKE) -C frontend lint

# Install dependencies
deps:
	@echo "Installing backend dependencies..."
	$(MAKE) -C backend deps
	@echo "Installing frontend dependencies..."
	$(MAKE) -C frontend deps

# Development environment
dev:
	@echo "Starting development environment..."
	docker-compose -f docker-compose.yml up -d

# Backend only development
dev-backend:
	@echo "Starting backend in development mode..."
	$(MAKE) -C backend run

# Frontend only development  
dev-frontend:
	@echo "Starting frontend in development mode..."
	$(MAKE) -C frontend dev

# Production environment
prod:
	@echo "Starting production environment..."
	docker-compose -f docker-compose.yml up -d --build

# Deploy to production
deploy: build
	@echo "Deploying to production..."
	docker-compose -f docker-compose.yml up -d --build