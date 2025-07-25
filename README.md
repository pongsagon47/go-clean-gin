# Go Clean Gin API

A production-ready RESTful API built with Go, Gin framework following Clean Architecture principles with advanced error handling and standardized response system.

## 🚀 Features

- 🏗️ **Clean Architecture** (Entity, Repository, Usecase, Handler)
- 🔐 **JWT Authentication** with secure token validation
- 🐘 **PostgreSQL Database** with GORM and connection pooling
- 📝 **Advanced Request Validation** with custom error messages
- 🐳 **Docker Support** with docker-compose
- 📊 **Structured Logging** with Zap (JSON/Development formats)
- ⚡ **Hot Reload** with Air/CompileDaemon
- 🧪 **Unit Testing** ready with comprehensive test structure
- 📋 **Standardized API Responses** with pagination support
- ⚠️ **Professional Error Handling** with custom error codes
- 🔍 **Enhanced Input Validation** with field-specific messages
- 🔧 **Database Connection Pooling** with custom configuration

## 📁 Project Structure

```
go-clean-gin/
├── cmd/
│   └── main.go                 # Application entrypoint
├── config/
│   └── config.go              # Configuration management
├── internal/
│   ├── entity/                # Domain entities (User, Product)
│   │   ├── user.go
│   │   └── product.go
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
│   │   └── postgres.go       # PostgreSQL setup with pooling
│   ├── errors/                # Custom error system
│   │   └── errors.go         # Application-specific errors
│   ├── logger/                # Logging utilities
│   │   └── logger.go         # Zap logger with levels
│   ├── response/              # Response system
│   │   └── response.go       # Standardized API responses
│   └── validator/             # Input validation
│       └── validator.go      # Custom validation messages
├── scripts/
│   └── test.sh               # Enhanced test runner
├── tmp/                      # Temporary files (hot reload)
├── .env                      # Environment variables
├── .env.example             # Example environment file
├── .gitignore               # Git ignore rules
├── .golangci.yml            # Linter configuration
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
make run
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

### Handler Architecture

#### Response System (`pkg/response/`)

- **Unified Response Structure**: All endpoints return consistent JSON format
- **Success Responses**: `Success()` and `SuccessWithMeta()` functions
- **Error Responses**: `Error()` and `ValidationError()` functions
- **Pagination Support**: Built-in pagination metadata generation

#### Error System (`pkg/errors/`)

- **Application Errors**: Custom `AppError` type with codes and HTTP status
- **Error Wrapping**: Wrap underlying errors with application context
- **Predefined Errors**: Common error instances ready to use
- **Error Details**: Support for additional error information

#### Validation System (`pkg/validator/`)

- **Enhanced Validation**: Better error messages with field names
- **Custom Messages**: Tailored validation messages for each rule
- **Field Mapping**: JSON field names in error messages
- **Struct Validation**: Comprehensive request validation

## 🛠️ Development Commands

```bash
# Setup & Development
make setup              # First-time project setup
make dev                # Run with hot reload
make dev-force          # Kill port conflicts and run
make run                # Run without hot reload

# Building & Testing
make build              # Build application
make test               # Run unit tests
make test-coverage      # Run tests with coverage

# Code Quality
make fmt                # Format code
make tidy               # Tidy dependencies
make clean              # Clean build artifacts

# Docker
make docker-build       # Build Docker image
make docker-run         # Start containers
make docker-stop        # Stop containers
make docker-logs        # View container logs

# Utilities
make check-port         # Check if port is available
make kill-port          # Kill processes on port
make health             # Check application health
make help               # Show all commands
```

## 🧪 Testing

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage
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
DB_LOG_LEVEL=warn
DB_MAX_IDLE_CONNS=10
DB_MAX_OPEN_CONNS=100
DB_CONN_MAX_LIFETIME=60

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
- **Standardized Error Handling** - Centralized error management
- **Response Standardization** - Consistent API responses

### Advanced Features

- **Connection Pooling** - Optimized database performance
- **Structured Logging** - Production-ready logging with Zap
- **Custom Validation** - Enhanced validation with detailed messages
- **Error Wrapping** - Comprehensive error tracking and debugging

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Run tests (`make test`)
4. Commit your changes (`git commit -m 'Add amazing feature'`)
5. Push to the branch (`git push origin feature/amazing-feature`)
6. Open a Pull Request

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
make run           # Run without hot reload
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
