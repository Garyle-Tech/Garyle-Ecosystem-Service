# Garyle Ecosystem Service API

A Golang API service built with the Gin framework following clean architecture principles.

## Project Structure

```
garyle-ecosystem-service/
├── cmd/                  # Application entry points
│   └── api/              # Main API server entry point
├── internal/             # Private application code
│   ├── app/              # Application layer
│   │   ├── api/          # API handlers and routes
│   │   └── config/       # Application configuration
│   ├── domain/           # Domain layer (business logic)
│   │   ├── model/        # Domain models/entities
│   │   ├── repository/   # Repository interfaces
│   │   └── service/      # Business logic services
│   └── infrastructure/   # Infrastructure layer
│       ├── database/     # Database implementations
│       └── middleware/   # HTTP middleware
├── pkg/                  # Reusable public packages
│   ├── logger/           # Logging utilities
│   └── utils/            # Common utility functions
└── docs/                 # API documentation
```

## Getting Started

1. Create your main entry point in `cmd/api/main.go`
2. Define your domain models in `internal/domain/model/`
3. Implement repositories in `internal/infrastructure/database/`
4. Create services in `internal/domain/service/`
5. Set up API routes and handlers in `internal/app/api/`

## Clean Architecture Layers

- **Domain Layer**: Contains business logic and entities
- **Application Layer**: Orchestrates the domain layer and handles use cases
- **Infrastructure Layer**: Implements interfaces defined in the domain layer

## Example Use

For a feature like "User Management":
1. Define User model in `internal/domain/model/user.go`
2. Create UserRepository interface in `internal/domain/repository/user_repository.go`
3. Implement UserRepository in `internal/infrastructure/database/user_repository.go`
4. Create UserService in `internal/domain/service/user_service.go`
5. Set up User handlers in `internal/app/api/user_handler.go` 