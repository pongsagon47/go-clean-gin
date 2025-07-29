# Go Clean Gin API

A production-ready RESTful API built with Go, Gin framework following Clean Architecture principles with **Laravel-style development experience**, advanced error handling, standardized response system, and complete database management tools.

## 🚀 Features

- 🏗️ **Clean Architecture** (Entity, Repository, Usecase, Handler)
- 🔐 **JWT Authentication** with secure token validation
- 🐘 **PostgreSQL Database** with GORM and connection pooling
- 🎨 **Laravel-style Migrations** with file-based versioning and rollback support
- 🌱 **Database Seeders** for development and testing data
- 🛠️ **Artisan CLI Tool** for generating migrations, seeders, entities, and packages
- 📦 **Package Generator** for complete Clean Architecture modules
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
│   ├── main.go                 # Application entrypoint
│   └── artisan/                # 🆕 Laravel-style CLI tool
│       └── main.go             # Migration & package generator
├── config/
│   └── config.go              # Configuration management
├── internal/
│   ├── entity/                # Domain entities (User, Product)
│   │   ├── user.go
│   │   └── product.go
│   ├── migrations/            # 🆕 Laravel-style migration files
│   │   ├── manager.go         # Migration manager
│   │   ├── 2024_01_15_120000_create_users_table.go
│   │   └── 2024_01_15_130000_create_products_table.go
│   ├── seeders/               # 🆕 Laravel-style seeder files
│   │   ├── manager.go         # Seeder manager
│   │   ├── user_seeder.go
│   │   └── product_seeder.go
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
├── bin/                      # 🆕 Built binaries
│   ├── app                   # Main application
│   └── artisan               # CLI tool
├── tmp/                      # Temporary files (hot reload)
├── .env                      # Environment variables
├── .env.example             # Example environment file
├── .gitignore               # Git ignore rules
├── .golangci.yml            # Linter configuration
├── docker-compose.yml       # Docker composition
├── Dockerfile               # Docker image definition
├── go.mod                   # Go module definition
├── go.sum                   # Go dependencies checksum
├── Makefile                 # 🆕 Laravel-style development commands
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

3. **Build Laravel-style CLI tool**

```bash
make build-artisan
```

4. **Start PostgreSQL**

```bash
# Using Docker
docker run --name postgres \
  -e POSTGRES_PASSWORD=password \
  -e POSTGRES_DB=go_clean_gin \
  -p 5432:5432 -d postgres:15-alpine

# Or use docker-compose
make docker-run
```

5. **Setup database**

```bash
# Create database (if needed)
make db-create

# Run migrations
make migrate

# Seed database with sample data
make db-seed
```

6. **Run the application**

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

## 🎨 Laravel-style Development Experience

This project brings the beloved Laravel development experience to Go with powerful generators and migration tools!

### 🏗️ Create Complete Features in Minutes

#### **Option 1: Complete Model Stack**

```bash
# Creates entity + migration + seeder in one command
make make-model NAME=Post FIELDS="title:string,content:text,author_id:uuid,status:string"
```

**Creates:**

- `internal/entity/post.go` - Complete entity with GORM tags, request/response structs
- `internal/migrations/2024_xx_xx_create_posts_table.go` - Migration with proper SQL
- `internal/seeders/post_seeder.go` - Seeder template

#### **Option 2: Step by Step**

```bash
# Create individual components
make make-entity NAME=Post FIELDS="title:string,content:text"
make make-migration NAME=create_posts_table CREATE=true TABLE=posts
make make-seeder NAME=PostSeeder TABLE=posts
```

#### **Option 3: Add Package Structure**

```bash
# Create Clean Architecture package structure
make make-package NAME=Post
```

**Creates:**

- `internal/post/handler.go` - HTTP handlers with proper error handling
- `internal/post/port.go` - Usecase and repository interfaces
- `internal/post/repository.go` - Database operations with GORM
- `internal/post/usecase.go` - Business logic layer

### ⚡ Quick Database Operations

```bash
# Add columns to existing tables
make add-column TABLE=users COLUMN=phone TYPE=string
make add-column TABLE=products COLUMN=sku TYPE=string

# Create indexes for performance
make add-index TABLE=products COLUMNS="category,price"
make add-index TABLE=posts COLUMNS="author_id,created_at"

# Drop columns when no longer needed
make drop-column TABLE=users COLUMN=old_field
```

### 🗄️ Migration Management

```bash
# Run pending migrations
make migrate

# Check what's been applied
make migrate-status

# Rollback if needed
make migrate-rollback           # Last migration
make migrate-rollback COUNT=3  # Last 3 migrations

# Fresh start (DANGER: destroys data!)
make migrate-fresh
```

### 🌱 Database Seeding

```bash
# Run all seeders
make db-seed

# Run Specific Seeder
make db-seed NAME=ProductSeeder

# Create new seeders
make make-seeder NAME=ProductSeeder TABLE=products
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

### Basic Development

```bash
# Setup & Development
make setup              # First-time project setup
make dev                # Run with hot reload
make dev-force          # Kill port conflicts and run
make run                # Run without hot reload

# Building & Testing
make build              # Build application
make build-artisan      # Build artisan CLI tool
make test               # Run unit tests
make test-coverage      # Run tests with coverage

# Code Quality
make fmt                # Format code
make tidy               # Tidy dependencies
make clean              # Clean build artifacts
```

### Laravel-style Commands

```bash
# 🎨 Generators
make make-migration     # Create new migration file
make make-seeder        # Create new seeder file
make make-entity        # Create new entity/model file
make make-package       # Create new package structure
make make-model         # Create complete model stack

# ⚡ Quick Actions
make add-column         # Add column to existing table
make drop-column        # Drop column from table
make add-index          # Add index to table

# 🗄️ Migration Management
make migrate            # Run pending migrations
make migrate-status     # Show migration status
make migrate-rollback   # Rollback migrations
make migrate-fresh      # Fresh migration (DANGER!)

# 🌱 Database Seeding
make db-seed            # Run database seeders

# 🏭 Database Management
make db-create          # Create database
make db-drop            # Drop database (DANGER!)
make db-reset           # Reset database completely
make db-info            # Show database information

# 🔍 Utilities
make list-migrations    # List all migration files
make validate-migrations # Validate migration syntax
make examples           # Show usage examples
```

### Docker Commands

```bash
# Docker
make docker-build       # Build Docker image
make docker-run         # Start containers
make docker-stop        # Stop containers
make docker-logs        # View container logs

# Health & Status
make health             # Check application health
make status             # Show application status
```

## 📖 Complete Workflow Examples

### Create a Blog System

```bash
# 1. Setup project
make setup
make build-artisan

# 2. Create blog entities
make make-model NAME=Post FIELDS="title:string,content:text,author_id:uuid,status:string"
make make-model NAME=Comment FIELDS="post_id:uuid,content:text,author_id:uuid"
make make-model NAME=Category FIELDS="name:string,description:text"

# 3. Create package structures
make make-package NAME=Post
make make-package NAME=Comment
make make-package NAME=Category

# 4. Add relationships and indexes
make add-column TABLE=posts COLUMN=category_id TYPE=uuid
make add-index TABLE=posts COLUMNS="author_id,created_at"
make add-index TABLE=posts COLUMNS="category_id,status"
make add-index TABLE=comments COLUMNS="post_id"

# 5. Deploy and seed
make migrate
make db-seed

# 6. Start development
make dev
```

### Add Features to Existing Models

```bash
# Add phone to users
make add-column TABLE=users COLUMN=phone TYPE=string

# Add SKU to products
make add-column TABLE=products COLUMN=sku TYPE=string

# Add search indexes
make add-index TABLE=products COLUMNS="name,category"

# Apply changes
make migrate
```

### Database Management

```bash
# Check current status
make migrate-status
make db-info

# Reset everything (development)
make migrate-fresh

# Rollback problematic migration
make migrate-rollback COUNT=2

# Complete database reset
make db-reset
```

## 🧪 Testing

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Validate migrations
make validate-migrations
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

This project follows **Clean Architecture** principles with **Laravel-style database management**:

- **Entities** (`internal/entity/`) - Core business models
- **Use Cases** (`internal/*/usecase.go`) - Business logic
- **Interface Adapters** (`internal/*/handler.go`, `internal/*/repository.go`) - External interfaces
- **Frameworks & Drivers** (`pkg/`, `cmd/`) - External frameworks
- **Migrations** (`internal/migrations/`) - Database schema versioning
- **Seeders** (`internal/seeders/`) - Database data seeding

### Key Patterns

- **Repository Pattern** - Data access abstraction
- **Dependency Injection** - Loose coupling
- **Middleware Pattern** - Cross-cutting concerns
- **Migration Pattern** - Database schema versioning (Laravel-style)
- **Seeder Pattern** - Database data management
- **Standardized Error Handling** - Centralized error management
- **Response Standardization** - Consistent API responses

### Advanced Features

- **Laravel-style Migrations** - File-based database versioning with rollback support
- **Artisan CLI Tool** - Command-line interface for generating migrations, entities, and packages
- **Auto-registration** - Automatic migration and seeder discovery
- **Package Generator** - Complete Clean Architecture scaffolding
- **Connection Pooling** - Optimized database performance
- **Structured Logging** - Production-ready logging with Zap
- **Custom Validation** - Enhanced validation with detailed messages
- **Error Wrapping** - Comprehensive error tracking and debugging

## 🎨 Laravel-style Features

### Migration System

- **File-based Migrations**: Each migration is a separate Go file
- **Automatic Versioning**: Timestamp-based migration ordering
- **Rollback Support**: Every migration has an up and down method
- **Transaction Safety**: All migrations run in database transactions
- **Status Tracking**: See which migrations have been applied

### Seeder System

- **Data Seeding**: Populate database with development/test data
- **Environment Awareness**: Different seeds for different environments
- **Dependency Management**: Seeders can depend on each other
- **Idempotent**: Safe to run multiple times

### Artisan CLI

- **Code Generation**: Generate migration, seeder, entity, and package boilerplate
- **Database Management**: Run migrations, rollbacks, and seeders
- **Status Monitoring**: Check current database state
- **Laravel-familiar**: Commands similar to Laravel artisan

### Package Generator

- **Complete Scaffolding**: Generate handler, usecase, repository, and port files
- **Clean Architecture**: Follows established patterns
- **Ready-to-customize**: Basic structure with TODO comments
- **Interface-driven**: Proper dependency injection setup

## 🚀 Type Mapping System

The Laravel-style generator automatically maps field types:

| Field Type  | Go Type           | SQL Type                   | GORM Tag                        | Validation               |
| ----------- | ----------------- | -------------------------- | ------------------------------- | ------------------------ |
| `string`    | `string`          | `VARCHAR(255)`             | `not null`                      | `required,min=1,max=255` |
| `text`      | `string`          | `TEXT`                     | `type:text`                     | `required`               |
| `int`       | `int`             | `INTEGER`                  | `not null`                      | `required,min=0`         |
| `decimal`   | `decimal.Decimal` | `DECIMAL(10,2)`            | `type:decimal(10,2);not null`   | `required,min=0`         |
| `bool`      | `bool`            | `BOOLEAN`                  | `default:false`                 | ``                       |
| `uuid`      | `uuid.UUID`       | `UUID`                     | `type:uuid;not null`            | `required`               |
| `timestamp` | `time.Time`       | `TIMESTAMP WITH TIME ZONE` | `type:timestamp with time zone` | ``                       |

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Run tests (`make test`)
4. Test migrations (`make migrate-status`, `make migrate`, `make migrate-rollback`)
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
make run           # Run without hot reload
```

#### Database Connection Issues

```bash
# Check PostgreSQL is running
docker ps
make docker-run     # Start PostgreSQL with Docker
```

#### Migration Issues

```bash
# Check migration status
make migrate-status

# Validate migration files
make validate-migrations

# List migration files
make list-migrations

# Reset database (DANGER!)
make db-reset
```

#### Artisan CLI Issues

```bash
# Build artisan CLI
make build-artisan

# Check if artisan is working
./bin/artisan -help

# Use go run if binary doesn't work
go run cmd/artisan/main.go -help
```

#### Build Issues

```bash
make clean      # Clean build artifacts
make tidy       # Tidy dependencies
go mod download # Re-download dependencies
```

## 📚 Learn More

### Migration Documentation

```bash
# Show detailed examples
make examples

# Show all available commands
make help
```

### Laravel-style Commands Reference

For complete command reference and examples, run:

```bash
make examples
```

This will show you detailed examples of:

- Creating complete features with entities, migrations, and packages
- Managing database schema with migrations
- Adding columns, indexes, and constraints
- Creating and running seeders
- Rolling back migrations
- Database management and reset operations

## 📞 Support

If you have any questions or need help, please:

1. Check the [troubleshooting section](#-troubleshooting)
2. Run `make examples` for Laravel-style command examples
3. Review existing [issues](https://github.com/your-repo/issues)
4. Create a new issue with detailed information

---

**Happy coding with Laravel-style migrations in Go! 🚀✨🎨**

_Experience the best of both worlds: Laravel's developer experience with Go's performance and type safety._
