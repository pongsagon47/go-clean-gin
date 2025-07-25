# Go Clean Gin API

A RESTful API built with Go, Gin framework following Clean Architecture principles with standardized response and error handling system.

## 🚀 Features

- 🏗️ **Clean Architecture** (Entity, Repository, Usecase, Handler)
- 🔐 **JWT Authentication** with secure token validation
- 🐘 **PostgreSQL Database** with GORM
- 🔄 **Database Migrations** with version control
- 📝 **Request Validation** with detailed error messages
- 🐳 **Docker Support** with docker-compose
- 📊 **Structured Logging** with Zap
- ⚡ **Hot Reload** with Air/CompileDaemon
- 🧪 **Testing Ready** with unit and integration tests
- 📋 **Standardized API Responses** with pagination
- ⚠️ **Advanced Error Handling** with custom error codes
- 🔍 **Input Validation** with field-specific messages

## 📁 Project Structure

```
go-clean-gin/
├── cmd/
│   └── main.go                 # Application entrypoint
├── config/
│   └── config.go              # Configuration management
├── internal/
│   ├── entity/                # Domain entities (User, Product)
│   ├── auth/                  # Authentication module
│   │   ├── handler.go         # HTTP handlers
│   │   ├── usecase.go         # Business logic
│   │   ├── repository.go      # Data access
│   │   └── port.go           # Interfaces
│   ├── product/               # Product module
│   │   ├── handler.go         # HTTP handlers
│   │   ├── usecase.go         # Business logic
│   │   ├── repository.go      # Data access
│   │   └── port.go           # Interfaces
│   ├── middleware/            # HTTP middlewares
│   │   ├── auth.go           # JWT authentication
│   │   ├── cors.go           # CORS configuration
│   │   ├── error.go          # Error handling
│   │   ├── logging.go        # Request logging
│   │   └── recovery.go       # Panic recovery
│   ├── router/                # Route definitions
│   │   └── router.go         # API routes setup
│   └── container/             # Dependency injection
│       └── container.go      # DI container
├── pkg/
│   ├── database/              # Database connection
│   │   └── postgres.go       # PostgreSQL setup
│   ├── errors/                # Custom error system
│   │   └── errors.go         # Application errors
│   ├── logger/                # Logging utilities
│   │   └── logger.go         # Zap logger setup
│   ├── response/              # Response system
│   │   └── response.go       # Standardized responses
│   └── validator/             # Input validation
│       └── validator.go      # Custom validators
├── scripts/
│   └── test.sh               # Test runner script
├── test/
│   └── integration_test.go   # Integration tests
├── tmp/                      # Temporary files (hot reload)
├── .env                      # Environment variables
├── .env.example             # Example environment file
├── .gitignore               # Git ignore rules
├── docker-compose.yml       # Docker composition
├── Dockerfile               # Docker image definition
├── go.mod                   # Go module definition
├── go.sum                   # Go dependencies checksum
├── Makefile                 # Development commands
└── README.md               # Project documentation
```

## 🚀 Quick Start

### Prerequisites

- Go 1.21+
- PostgreSQL 13+
- Docker & Docker Compose (optional)

### Local Development

1. **Clone and setup**

```bash
git clone <repository-url>
cd go-clean-gin
make setup
```

2. **Configure environment**

```bash
# Edit .env file with your settings
cp .env.example .env
vim .env
```

3. **Start PostgreSQL**

```bash
# Using Docker
docker run --name postgres \
  -e POSTGRES_PASSWORD=password \
  -e POSTGRES_DB=go_clean_gin \
  -p 5432:5432 -d postgres:15-alpine

# Or use docker-compose
make docker-run
```

4. **Run the application**

```bash
# Development with hot reload
make dev

# Or simple run
make dev-simple
```

### Using Docker

```bash
# Start all services
make docker-run

# View logs
make docker-logs

# Stop services
make docker-stop
```

## 📡 API Endpoints

### Base URL

```
http://localhost:8080/api/v1
```

### Authentication

```http
# Register
POST /auth/register
{
  "email": "user@example.com",
  "username": "username",
  "password": "password123",
  "first_name": "John",
  "last_name": "Doe"
}

# Login
POST /auth/login
{
  "email": "user@example.com",
  "password": "password123"
}

# Get Profile (Protected)
GET /auth/profile
Authorization: Bearer <token>
```

### Products

```http
# Get Products (with filters & pagination)
GET /products?page=1&limit=10&category=electronics&search=phone

# Get Product by ID
GET /products/{id}

# Create Product (Protected)
POST /products
Authorization: Bearer <token>
{
  "name": "iPhone 15",
  "description": "Latest iPhone model",
  "price": 999.99,
  "stock": 10,
  "category": "electronics"
}

# Update Product (Protected)
PUT /products/{id}
Authorization: Bearer <token>
{
  "name": "iPhone 15 Pro",
  "price": 1099.99
}

# Delete Product (Protected)
DELETE /products/{id}
Authorization: Bearer <token>
```

### Health Check

```http
GET /health
```

## 📋 Response & Error Handling System

### Standardized Response Format

#### Success Response

```json
{
  "success": true,
  "message": "Operation completed successfully",
  "data": {
    // Response data here
  },
  "timestamp": "2024-01-15T10:30:00Z"
}
```

#### Success Response with Pagination

```json
{
  "success": true,
  "message": "Products retrieved successfully",
  "data": [
    // Array of items
  ],
  "meta": {
    "page": 1,
    "limit": 10,
    "total": 25,
    "total_pages": 3,
    "has_next": true,
    "has_previous": false
  },
  "timestamp": "2024-01-15T10:30:00Z"
}
```

#### Error Response

```json
{
  "success": false,
  "message": "Request failed",
  "error": {
    "code": "ERROR_CODE",
    "message": "Detailed error message",
    "details": "Additional error information"
  },
  "timestamp": "2024-01-15T10:30:00Z"
}
```

#### Validation Error Response

```json
{
  "success": false,
  "message": "Validation failed",
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Validation failed",
    "fields": {
      "email": "email is required",
      "password": "password must be at least 6 characters"
    }
  },
  "timestamp": "2024-01-15T10:30:00Z"
}
```

### Error Codes

#### General Errors

- `INTERNAL_ERROR` - Internal server error
- `NOT_FOUND` - Resource not found
- `BAD_REQUEST` - Invalid request format
- `UNAUTHORIZED` - Authentication required
- `FORBIDDEN` - Insufficient permissions
- `VALIDATION_ERROR` - Request validation failed

#### Authentication Errors

- `INVALID_CREDENTIALS` - Invalid email or password
- `TOKEN_EXPIRED` - JWT token has expired
- `TOKEN_INVALID` - Invalid JWT token
- `USER_EXISTS` - User already exists
- `USER_NOT_FOUND` - User not found

#### Product Errors

- `PRODUCT_NOT_FOUND` - Product not found
- `PRODUCT_EXISTS` - Product already exists
- `INSUFFICIENT_STOCK` - Not enough stock available
- `INVALID_OWNER` - User can only modify own resources

## 🛠️ Development Commands

```bash
# Setup & Installation
make setup              # First-time project setup
make install            # Install dependencies
make install-tools      # Install development tools

# Development
make dev                # Run with hot reload
make dev-simple         # Run without hot reload
make dev-force          # Kill port conflicts and run

# Building
make build              # Build for current platform
make build-all          # Build for multiple platforms
make run                # Run built binary

# Testing
make test               # Run all tests
make test-coverage      # Run tests with coverage
make test-integration   # Run integration tests
make test-unit          # Run unit tests only

# Code Quality
make fmt                # Format code
make lint               # Run linter
make lint-fix           # Fix linting issues
make security           # Security check

# Docker
make docker-build       # Build Docker image
make docker-run         # Start containers
make docker-stop        # Stop containers
make docker-logs        # View logs

# Utilities
make clean              # Clean build artifacts
make tidy               # Tidy dependencies
make update             # Update dependencies
make health             # Check app health
make logs               # View application logs

# Performance
make benchmark          # Run benchmark tests
make profile-cpu        # CPU profiling
make profile-mem        # Memory profiling

# Help
make help               # Show all commands
```

## 🧪 Testing

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run integration tests
make test-integration

# Run benchmark tests
make benchmark
```

## 🔧 Configuration

### Environment Variables

```bash
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=go_clean_gin
DB_SSLMODE=disable
SERVER_PORT=8080
SERVER_HOST=0.0.0.0
SERVER_READ_TIMEOUT=30s
SERVER_WRITE_TIMEOUT=30s

# Server
SERVER_PORT=8080
SERVER_HOST=0.0.0.0
SERVER_READ_TIMEOUT=30s
SERVER_WRITE_TIMEOUT=30s

# JWT
JWT_SECRET=your-super-secret-jwt-key
JWT_EXPIRATION_HOURS=24

# Logging
LOG_LEVEL=info
LOG_FORMAT=json

# Environment
ENV=development
```

## 🐳 Docker Deployment

```bash
# Build and run with Docker Compose
make docker-run

# Or manually
docker build -t go-clean-gin .
docker run -p 8080:8080 --env-file .env go-clean-gin
```

## 🏗️ Architecture

This project follows **Clean Architecture** principles:

- **Entities** (`internal/entity/`) - Core business models
- **Use Cases** (`internal/*/usecase.go`) - Business logic
- **Interface Adapters** (`internal/*/handler.go`, `internal/*/repository.go`) - External interfaces
- **Frameworks & Drivers** (`pkg/`, `cmd/`) - External frameworks

### Key Patterns

- **Repository Pattern** - Data access abstraction
- **Dependency Injection** - Loose coupling
- **Middleware Pattern** - Cross-cutting concerns
- **Error Handling** - Centralized error management

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Run tests (`make test`)
4. Run linter (`make lint`)
5. Commit your changes (`git commit -m 'Add amazing feature'`)
6. Push to the branch (`git push origin feature/amazing-feature`)
7. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🆘 Troubleshooting

### Common Issues

#### Port Already in Use

```bash
make kill-port  # Kill processes on configured port
make dev-force  # Force run (kills port first)
```

#### Hot Reload Not Working

```bash
make install-tools  # Install Air/CompileDaemon
make dev-simple     # Run without hot reload
```

#### Database Connection Issues

```bash
# Check PostgreSQL is running
docker ps
make docker-run     # Start PostgreSQL with Docker
```

#### Build Issues

```bash
make clean      # Clean build artifacts
make tidy       # Tidy dependencies
go mod download # Re-download dependencies
```

## 📞 Support

If you have any questions or need help, please:

1. Check the [troubleshooting section](#-troubleshooting)
2. Review existing [issues](https://github.com/your-repo/issues)
3. Create a new issue with detailed information

---

**Happy coding! 🚀**
