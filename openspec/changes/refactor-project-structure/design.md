# Design: Project Structure Reorganization

## Context

The guebapi project is a JWT authentication API built with Go and Gin. Currently, all code resides in `main.go`, which includes:
- Model definitions (`User`, `Claims`)
- HTTP handlers (login, register, profile, update)
- JWT middleware
- In-memory user storage
- Router configuration

This single-file approach works for prototypes but limits maintainability and testability as the project evolves.

## Goals / Non-Goals

**Goals:**
- Establish clear package boundaries following Go conventions
- Enable unit testing of individual components
- Prepare architecture for database integration
- Maintain backward compatibility (same API behavior)

**Non-Goals:**
- Database migration (future change)
- Adding new features or endpoints
- Changing authentication mechanism
- Performance optimization

## Decisions

### Decision 1: Follow Standard Go Project Layout

Adopt the community-standard layout with `cmd/`, `internal/`, and `config/` directories.

**Rationale:** Well-documented, widely understood pattern that new contributors can navigate immediately.

**Alternatives considered:**
- Flat structure with packages at root level - Rejected: pollutes import namespace
- Domain-driven design folders - Rejected: over-engineered for current scope

### Decision 2: Use `internal/` for Private Packages

All application code goes under `internal/` to prevent external imports.

**Structure:**
```
internal/
├── api/
│   ├── handlers/    # HTTP request handlers
│   ├── middleware/  # JWT auth, logging, etc.
│   └── router/      # Route definitions
├── models/          # Domain models (User, Claims)
├── repository/      # Data access interfaces and implementations
│   └── user/        # User storage
└── service/         # Business logic
    └── auth/        # Authentication service
```

### Decision 3: Interface-Based Repository Layer

Define `UserRepository` interface to abstract storage implementation.

```go
type UserRepository interface {
    Create(ctx context.Context, user *models.User) error
    GetByUsername(ctx context.Context, username string) (*models.User, error)
    UpdatePassword(ctx context.Context, username, password string) error
    Exists(ctx context.Context, username string) (bool, error)
}
```

**Rationale:** Enables testing with mocks and future database swap without changing business logic.

### Decision 4: Configuration Package

Create `config/` package for externalized configuration.

```go
type Config struct {
    Server   ServerConfig
    JWT      JWTConfig
}

type ServerConfig struct {
    Port string
    Host string
}

type JWTConfig struct {
    SecretKey     string
    TokenDuration time.Duration
}
```

**Rationale:** Separates configuration from code, supports environment variables, and improves security by not hardcoding secrets.

### Decision 5: Dependency Injection in main.go

Wire all dependencies explicitly in `main.go`:
1. Load configuration
2. Create repository instances
3. Create service instances
4. Create handlers with injected services
5. Configure router
6. Start server

**Rationale:** Explicit wiring makes dependencies visible and testable without frameworks.

## Package Dependencies

```
main.go
  └── config
  └── internal/api/router
        └── internal/api/handlers
              └── internal/service/auth
                    └── internal/repository/user
                          └── internal/models
        └── internal/api/middleware
              └── internal/service/auth
```

## Risks / Trade-offs

| Risk | Mitigation |
|------|------------|
| More files to navigate | Clear naming conventions, IDE navigation |
| Initial refactoring effort | Comprehensive task breakdown, incremental commits |
| Potential for circular imports | Careful package design, interfaces at boundaries |

## Migration Plan

1. Create new package structure (empty packages with interfaces)
2. Move models first (no dependencies)
3. Implement repository with current in-memory storage
4. Extract authentication service
5. Move handlers with service injection
6. Extract middleware
7. Configure router in dedicated package
8. Reduce main.go to wiring only
9. Add configuration loading
10. Run tests to verify behavior unchanged

**Rollback:** Each step is a separate commit; revert to previous commit if issues arise.

## Open Questions

- Should we add structured logging now or defer to separate change?
- Should configuration support both YAML files and environment variables?
