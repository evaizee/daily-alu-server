# DailyAlu Server

A robust REST API server built with Go, implementing clean architecture and Domain-Driven Design principles.

## Prerequisites

- Go 1.22 or newer
- PostgreSQL 12 or newer
- [golang-migrate/migrate](https://github.com/golang-migrate/migrate) for database migrations

## Project Structure

```
.
├── cmd/                    # Command line interfaces
├── config/                 # Configuration files
├── database/              # Database migrations and queries
├── internal/              # Private application code
│   ├── container/         # Dependency injection container
│   ├── handler/           # HTTP handlers
│   ├── middleware/        # HTTP middleware
│   ├── module/            # Business logic modules
│   │   └── user/         # User module
│   │       ├── domain/   # Domain models and interfaces
│   │       ├── repository/# Data access layer
│   │       └── usecase/  # Business logic
│   ├── router/           # HTTP routing
│   ├── security/         # Security utilities
│   └── validator/        # Request validation
└── pkg/                   # Public libraries
    ├── app_log/          # Logging utilities
    ├── db/               # Database utilities
    └── cache/            # Caching utilities
```

## Configuration

1. Copy the example configuration:
   ```bash
   cp config/config.yaml config/config.local.yaml
   ```

2. Update the configuration in `config/config.local.yaml`:
   ```yaml
   database:
     host: localhost
     port: 5432
     user: your_user
     password: your_password
     name: dailyalu
     sslmode: disable

   jwt:
     secret: your-super-secret-key-change-this
     expiry: 24  # hours
   ```

## Database Setup

1. Create the database:
   ```sql
   CREATE DATABASE dailyalu;
   ```

2. Run migrations:
   ```bash
   make migrate-up
   ```

   To rollback migrations:
   ```bash
   make migrate-down
   ```

## Running the Server

1. Install dependencies:
   ```bash
   make deps
   ```

2. Run the server:
   ```bash
   make run
   ```

   Or with a specific config:
   ```bash
   make run CONFIG=config/config.local.yaml
   ```

## API Endpoints

### Authentication

#### Register
```http
POST /api/auth/register
Content-Type: application/json

{
    "email": "user@example.com",
    "password": "your-password",
    "name": "John Doe"
}
```

#### Login
```http
POST /api/auth/login
Content-Type: application/json

{
    "email": "user@example.com",
    "password": "your-password"
}
```

### Users

#### Get User
```http
GET /api/users/:id
Authorization: Bearer your-jwt-token
```

#### Update User
```http
PUT /api/users/:id
Authorization: Bearer your-jwt-token
Content-Type: application/json

{
    "email": "new.email@example.com",
    "name": "New Name"
}
```

#### Delete User (Admin only)
```http
DELETE /api/users/:id
Authorization: Bearer your-jwt-token
```

## Development

### Running Tests
```bash
make test
```

### Code Quality
```bash
make lint
```

## Security

- All endpoints except `/api/auth/register` and `/api/auth/login` require JWT authentication
- Passwords are hashed using bcrypt
- Rate limiting is enabled by default
- CORS is configured for security
