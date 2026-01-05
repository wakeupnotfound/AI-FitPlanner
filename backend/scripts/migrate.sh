#!/bin/bash

# Database Migration Script
# This script runs the database migration to create tables and insert initial data

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}=== AI Fitness Planner Database Migration ===${NC}"
echo ""

# Check if config file exists
if [ ! -f "configs/config.yaml" ]; then
    echo -e "${RED}Error: configs/config.yaml not found${NC}"
    echo "Please create the configuration file before running migration"
    exit 1
fi

# Check if schema file exists
if [ ! -f "database/schema.sql" ]; then
    echo -e "${RED}Error: database/schema.sql not found${NC}"
    exit 1
fi

echo -e "${YELLOW}Starting database migration...${NC}"
echo ""

# Run the migration
cd "$(dirname "$0")/.."
go run scripts/migrate.go

if [ $? -eq 0 ]; then
    echo ""
    echo -e "${GREEN}✓ Database migration completed successfully${NC}"
else
    echo ""
    echo -e "${RED}✗ Database migration failed${NC}"
    exit 1
fi
