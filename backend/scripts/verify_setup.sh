#!/bin/bash

# Script to verify the project setup

echo "=== AI Fitness Planner Backend Setup Verification ==="
echo ""

# Check Go installation
echo "1. Checking Go installation..."
if command -v go &> /dev/null; then
    GO_VERSION=$(go version)
    echo "   ✓ Go is installed: $GO_VERSION"
else
    echo "   ✗ Go is not installed"
    exit 1
fi

# Check project structure
echo ""
echo "2. Checking project structure..."
REQUIRED_DIRS=("cmd/api" "internal/config" "internal/pkg/database" "internal/pkg/logger" "internal/pkg/redis" "configs")
for dir in "${REQUIRED_DIRS[@]}"; do
    if [ -d "$dir" ]; then
        echo "   ✓ Directory exists: $dir"
    else
        echo "   ✗ Directory missing: $dir"
    fi
done

# Check required files
echo ""
echo "3. Checking required files..."
REQUIRED_FILES=("go.mod" "cmd/api/main.go" "configs/config.yaml" ".env.example" "Makefile" "docker-compose.yml")
for file in "${REQUIRED_FILES[@]}"; do
    if [ -f "$file" ]; then
        echo "   ✓ File exists: $file"
    else
        echo "   ✗ File missing: $file"
    fi
done

# Check configuration files
echo ""
echo "4. Checking configuration..."
if [ -f "configs/config.yaml" ]; then
    echo "   ✓ Configuration file exists"
else
    echo "   ✗ Configuration file missing"
fi

if [ -f ".env" ]; then
    echo "   ✓ .env file exists"
else
    echo "   ⚠ .env file not found (copy from .env.example)"
fi

# Check Go modules
echo ""
echo "5. Checking Go modules..."
if [ -f "go.mod" ]; then
    echo "   ✓ go.mod exists"
    echo "   Note: Run 'go mod download' to fetch dependencies"
else
    echo "   ✗ go.mod missing"
fi

echo ""
echo "=== Setup Verification Complete ==="
echo ""
echo "Next steps:"
echo "1. Copy .env.example to .env and configure"
echo "2. Run 'make deps' to download dependencies"
echo "3. Start infrastructure with 'make docker-up'"
echo "4. Run migrations with 'make migrate'"
echo "5. Start the application with 'make run'"
