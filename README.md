A RESTful API built with Go, Gin framework following Clean Architecture principles.

## Features

- 🏗️ Clean Architecture (Entity, Repository, Usecase, Handler)
- 🔐 JWT Authentication
- 🐘 PostgreSQL Database with GORM
- 🔄 Database Migrations
- 📝 Request Validation
- 🚀 Docker Support
- 📊 Structured Logging with Zap
- ⚡ Hot Reload with Air
- 🧪 Testing Ready

## Project Structure

```
go-clean-gin/
├── cmd/
│   └── main.go                 # Application entrypoint
├── config/
│   └── config.go              # Configuration management
├── internal/
│   ├── entity/                # Domain entities
│   ├── auth/                  # Authentication module
│   ├── product/               # Product module
│   ├── middleware/            # HTTP middlewares
│   ├── router/                # Route definitions
│   └── container/             # Dependency injection
├── pkg/
│   ├── database/              # Database connection
│   └── logger/                # Logging utilities
├── .env                       # Environment variables
├── docker-compose.yml         # Docker composition
├── Dockerfile                 # Docker image definition
├── Makefile                   # Development commands
└── README.md
```

## Quick Start

### Prerequisites

- Go 1.21+
- PostgreSQL 13+
- Docker (optional)

### Local Development

1. **Clone the repository**

```bash
git clone <repository-url>
cd go-clean-gin
```

2. **Install dependencies**

```bash
make install
```

3. **Setup environment**

```bash
cp .env.example .env
# Edit .env with your database credentials
```

4. **Start PostgreSQL** (using Docker)

```bash
docker run --name postgres -e POSTGRES_PASSWORD=password -e POSTGRES_DB=go_clean_gin -p 5432:5432 -d postgres:15-alpine
```

5. **Run the application**

```bash
make run
# or for hot reload
make dev
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

## API Endpoints

### Authentication

- `POST /api/v1/auth/register` - Register new user
- `POST /api/v1/auth/login` - Login user
- `GET /api/v1/auth/profile` - Get user profile (protected)

### Products

- `GET /api/v1/products` - Get products with filters
- `GET /api/v1/products/:id` - Get product by ID
- `POST /api/v1/products` - Create product (protected)
- `PUT /api/v1/products/:id` - Update product (protected)
- `DELETE /api/v1/products/:id` - Delete product (protected)

### Health Check

- `GET /health` - Health check endpoint

## Environment Variables

```bash
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=go_clean_gin
DB_SSLMODE=disable

# Server
SERVER_PORT=8080
SERVER_HOST=0.0.0.0

# JWT
JWT_SECRET=your-secret-key
JWT_EXPIRATION_HOURS=24

# Logging
LOG_LEVEL=info
LOG_FORMAT=json

# Environment
ENV=development
```

## Development Commands

```bash
make build          # Build the application
make run             # Run the application
make dev             # Run with hot reload
make test            # Run tests
make test-coverage   # Run tests with coverage
make lint            # Run linter
make fmt             # Format code
make clean           # Clean build artifacts
```

## Testing

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License.
