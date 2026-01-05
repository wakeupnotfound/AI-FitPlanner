# AI Fitness Planner Backend

Backend API service for the AI Fitness Planning System built with Go, Gin framework, and clean architecture.

## Features

- 🔐 JWT-based authentication with refresh tokens
- 🤖 AI-powered fitness and nutrition plan generation
- 📊 Training and nutrition tracking
- 📈 Progress analytics and statistics
- 🔒 AES-256 encryption for sensitive data
- 🚦 Rate limiting and security middleware
- 📝 Structured logging with Zap
- 🗄️ MySQL database with GORM
- ⚡ Redis for caching and session management

## Prerequisites

- Go 1.21 or higher
- MySQL 8.0 or higher
- Redis 7.0 or higher
- Docker and Docker Compose (optional)

## Project Structure

```
backend/
├── cmd/
│   └── api/
│       └── main.go           # Application entry point
├── configs/
│   └── config.yaml           # Configuration file
├── internal/
│   ├── api/
│   │   ├── request/          # Request DTOs
│   │   └── response/         # Response DTOs
│   ├── config/               # Configuration management
│   ├── errors/               # Error definitions
│   ├── model/                # Data models
│   └── pkg/
│       ├── crypto/           # Encryption utilities
│       ├── database/         # Database connection
│       ├── jwt/              # JWT utilities
│       ├── logger/           # Logging utilities
│       └── redis/            # Redis client
├── Dockerfile
├── docker-compose.yml
├── Makefile
└── README.md
```

## Getting Started

### Quick Start (Docker)

The fastest way to get started is using Docker Compose:

```bash
# 1. Clone the repository
cd backend

# 2. Copy environment file
cp .env.example .env

# 3. Start all services (MySQL, Redis, API)
docker-compose up -d

# 4. Check service health
curl http://localhost:8080/health

# 5. View API documentation
open http://localhost:8080/swagger/index.html
```

### Manual Setup

#### 1. Clone and Setup

```bash
cd backend
cp .env.example .env
# Edit .env with your configuration
```

#### 2. Install Dependencies

```bash
make deps
```

#### 3. Start Infrastructure

**Option A: Using Docker**
```bash
make docker-up
```

**Option B: Manual Installation**

Install and start MySQL:
```bash
# macOS
brew install mysql
brew services start mysql

# Create database and user
mysql -u root -p
CREATE DATABASE fitness_planner;
CREATE USER 'fitness_user'@'localhost' IDENTIFIED BY 'your_password';
GRANT ALL PRIVILEGES ON fitness_planner.* TO 'fitness_user'@'localhost';
FLUSH PRIVILEGES;
```

Install and start Redis:
```bash
# macOS
brew install redis
brew services start redis

# Linux
sudo apt-get install redis-server
sudo systemctl start redis
```

#### 4. Run Database Migrations

```bash
make migrate
```

Or manually:

```bash
mysql -h localhost -u fitness_user -p fitness_planner < ../database/schema.sql
```

#### 5. Configure Environment

Edit `.env` file with your settings:
```bash
FITNESS_APP_PORT=8080
FITNESS_DATABASE_MYSQL_HOST=localhost
FITNESS_DATABASE_MYSQL_PORT=3306
FITNESS_DATABASE_MYSQL_USER=fitness_user
FITNESS_DATABASE_MYSQL_PASSWORD=your_password
FITNESS_DATABASE_MYSQL_DBNAME=fitness_planner
FITNESS_DATABASE_REDIS_HOST=localhost:6379
FITNESS_JWT_SECRET=your-jwt-secret-key-min-32-chars
FITNESS_APP_SECRET_KEY=your-encryption-key-32-chars
```

#### 6. Run the Application

```bash
make run
```

The API will be available at `http://localhost:8080`

### Verify Installation

```bash
# Check health
curl http://localhost:8080/health

# Expected response:
# {"status":"healthy","timestamp":1234567890,"services":{"database":"healthy","redis":"healthy"}}
```

## Configuration

Configuration can be set via:
1. `configs/config.yaml` file
2. Environment variables with `FITNESS_` prefix

Example environment variable:
```bash
export FITNESS_APP_PORT=8080
export FITNESS_DATABASE_MYSQL_HOST=localhost
```

## Development

### Build

```bash
make build
```

### Run Tests

Run all tests:
```bash
make test
```

Run tests with verbose output:
```bash
go test -v ./...
```

Run tests for a specific package:
```bash
go test -v ./internal/service/...
```

Run a specific test:
```bash
go test -v -run TestAuthService_Register ./internal/service/
```

### Test Types

The project includes multiple types of tests:

#### Unit Tests
Test individual functions and methods in isolation:
```bash
go test ./internal/service/... -v
go test ./internal/repository/... -v
go test ./internal/handler/... -v
```

#### Property-Based Tests
Test universal properties across randomized inputs (using gopter):
```bash
go test ./tests/property/... -v
```

Property tests validate correctness properties like:
- Encryption/decryption round trips
- Token generation uniqueness
- Session invalidation
- Data persistence consistency

#### Integration Tests
Test complete workflows with real dependencies:
```bash
go test ./tests/integration/... -v
```

### View Test Coverage

Generate and view test coverage report:
```bash
make test-coverage
```

Or manually:
```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

View coverage by package:
```bash
go test -cover ./...
```

### Running Tests with Database

Some tests require a test database. Use the test configuration:

```bash
# Set test environment
export FITNESS_APP_MODE=test
export FITNESS_DATABASE_MYSQL_DBNAME=fitness_planner_test

# Run tests
go test ./...
```

Or use Docker for isolated test database:
```bash
docker-compose -f docker-compose.test.yml up -d
go test ./...
docker-compose -f docker-compose.test.yml down
```

### Format Code

```bash
make fmt
```

### Run Linter

```bash
make lint
```

Install golangci-lint if not already installed:
```bash
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

## Docker Deployment

### Build and Run with Docker Compose

```bash
docker-compose up -d
```

### View Logs

```bash
docker-compose logs -f api
```

### Stop Services

```bash
docker-compose down
```

## API Documentation

### Swagger UI

Interactive API documentation is available via Swagger UI:

**Local Development:**
```
http://localhost:8080/swagger/index.html
```

**Production:**
```
https://your-domain.com/swagger/index.html
```

### Regenerate Swagger Documentation

After modifying API handlers or adding new endpoints:

```bash
make swagger
```

Or manually:

```bash
swag init -g cmd/api/main.go -o docs --parseDependency --parseInternal
```

### API Endpoints Overview

#### Authentication
- `POST /api/v1/auth/register` - Register new user
- `POST /api/v1/auth/login` - User login
- `POST /api/v1/auth/logout` - User logout (requires auth)
- `POST /api/v1/auth/refresh` - Refresh access token

#### User Management
- `GET /api/v1/user/profile` - Get user profile
- `PUT /api/v1/user/profile` - Update user profile
- `POST /api/v1/user/body-data` - Add body measurements
- `GET /api/v1/user/body-data` - Get body data history
- `POST /api/v1/user/fitness-goals` - Set fitness goals

#### AI API Management
- `POST /api/v1/ai-apis` - Add AI API configuration
- `GET /api/v1/ai-apis` - List AI APIs
- `GET /api/v1/ai-apis/:id` - Get AI API details
- `PUT /api/v1/ai-apis/:id` - Update AI API
- `DELETE /api/v1/ai-apis/:id` - Delete AI API
- `POST /api/v1/ai-apis/:id/test` - Test AI API connection
- `POST /api/v1/ai-apis/:id/set-default` - Set as default API

#### Fitness Assessments
- `POST /api/v1/assessments` - Create fitness assessment
- `GET /api/v1/assessments/latest` - Get latest assessment

#### Training Plans
- `POST /api/v1/training-plans/generate` - Generate training plan (AI)
- `GET /api/v1/training-plans/tasks/:taskId` - Get generation task status
- `GET /api/v1/training-plans` - List training plans
- `GET /api/v1/training-plans/:id` - Get plan details
- `GET /api/v1/training-plans/today` - Get today's training

#### Training Records
- `POST /api/v1/training-records` - Record training session
- `GET /api/v1/training-records` - List training records

#### Nutrition Plans
- `POST /api/v1/nutrition-plans/generate` - Generate nutrition plan (AI)
- `GET /api/v1/nutrition-plans` - List nutrition plans
- `GET /api/v1/nutrition-plans/:id` - Get plan details
- `GET /api/v1/nutrition-plans/today` - Get today's meals

#### Nutrition Records
- `POST /api/v1/nutrition-records` - Record meal
- `GET /api/v1/nutrition-records` - List nutrition records
- `GET /api/v1/nutrition-records/daily-summary` - Get daily nutrition summary

#### Statistics
- `GET /api/v1/stats/training` - Get training statistics
- `GET /api/v1/stats/progress` - Get progress report
- `GET /api/v1/stats/trends` - Get trend analysis

#### System
- `GET /health` - Health check endpoint

### Authentication

Most endpoints require authentication. Include the JWT token in the Authorization header:

```bash
Authorization: Bearer <your-access-token>
```

Example with curl:

```bash
curl -H "Authorization: Bearer eyJhbGc..." http://localhost:8080/api/v1/user/profile
```

## Health Check

```bash
curl http://localhost:8080/health
```

## Environment Variables

See `.env.example` for all available configuration options.

Key variables:
- `FITNESS_APP_PORT`: API server port (default: 8080)
- `FITNESS_DATABASE_MYSQL_HOST`: MySQL host
- `FITNESS_DATABASE_MYSQL_DBNAME`: Database name
- `FITNESS_DATABASE_REDIS_HOST`: Redis host
- `FITNESS_JWT_SECRET`: JWT signing secret
- `FITNESS_APP_SECRET_KEY`: Application encryption key

## Security

- All passwords are hashed with bcrypt
- API keys are encrypted with AES-256-GCM
- JWT tokens for authentication
- Rate limiting on all endpoints
- Input validation and sanitization
- SQL injection prevention via GORM
- XSS prevention via output escaping

## Troubleshooting

### Database Connection Issues

**Problem:** `Failed to initialize database`

**Solutions:**
1. Verify MySQL is running:
   ```bash
   # macOS
   brew services list | grep mysql
   
   # Linux
   sudo systemctl status mysql
   ```

2. Check database credentials in `.env`
3. Verify database exists:
   ```bash
   mysql -u fitness_user -p -e "SHOW DATABASES;"
   ```

4. Check MySQL logs:
   ```bash
   # macOS
   tail -f /usr/local/var/mysql/*.err
   
   # Linux
   sudo tail -f /var/log/mysql/error.log
   ```

### Redis Connection Issues

**Problem:** `Failed to initialize Redis`

**Solutions:**
1. Verify Redis is running:
   ```bash
   redis-cli ping
   # Should return: PONG
   ```

2. Check Redis configuration in `.env`
3. Test Redis connection:
   ```bash
   redis-cli -h localhost -p 6379
   ```

### Port Already in Use

**Problem:** `bind: address already in use`

**Solutions:**
1. Find process using the port:
   ```bash
   lsof -i :8080
   ```

2. Kill the process or change port in `.env`:
   ```bash
   FITNESS_APP_PORT=8081
   ```

### Migration Failures

**Problem:** Migration script fails

**Solutions:**
1. Check if database exists and is accessible
2. Verify user has proper permissions:
   ```sql
   SHOW GRANTS FOR 'fitness_user'@'localhost';
   ```

3. Run migration manually:
   ```bash
   mysql -h localhost -u fitness_user -p fitness_planner < ../database/schema.sql
   ```

### Swagger Documentation Not Loading

**Problem:** `/swagger/index.html` returns 404

**Solutions:**
1. Regenerate Swagger docs:
   ```bash
   make swagger
   ```

2. Verify docs directory exists:
   ```bash
   ls -la docs/
   ```

3. Check that docs package is imported in `main.go`:
   ```go
   _ "github.com/ai-fitness-planner/backend/docs"
   ```

### Test Failures

**Problem:** Tests fail with database errors

**Solutions:**
1. Use test database:
   ```bash
   export FITNESS_DATABASE_MYSQL_DBNAME=fitness_planner_test
   ```

2. Create test database:
   ```sql
   CREATE DATABASE fitness_planner_test;
   ```

3. Run migrations on test database

### JWT Token Issues

**Problem:** `Invalid token` or `Token expired`

**Solutions:**
1. Verify JWT secret is set in `.env` (minimum 32 characters)
2. Check token expiration settings
3. Use refresh token to get new access token:
   ```bash
   curl -X POST http://localhost:8080/api/v1/auth/refresh \
     -H "Content-Type: application/json" \
     -d '{"refresh_token":"your-refresh-token"}'
   ```

### Rate Limiting Issues

**Problem:** Getting 429 Too Many Requests

**Solutions:**
1. Wait for rate limit window to reset
2. Adjust rate limits in `config.yaml`:
   ```yaml
   rate_limit:
     api_calls_per_minute: 100
     api_calls_per_hour: 2000
   ```

3. Clear Redis rate limit counters:
   ```bash
   redis-cli KEYS "ratelimit:*" | xargs redis-cli DEL
   ```

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Code Style

- Follow Go best practices and idioms
- Run `make fmt` before committing
- Run `make lint` and fix any issues
- Add tests for new features
- Update documentation as needed

## Support

For issues and questions:
- Open an issue on GitHub
- Check existing documentation
- Review troubleshooting section above

## License

Copyright © 2024 AI Fitness Planner
