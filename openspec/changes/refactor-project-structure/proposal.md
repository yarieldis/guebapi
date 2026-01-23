# Change: Refactor Project Structure to Standard Go Layout

## Why

The current codebase has all application logic in a single `main.go` file (~157 lines), mixing concerns like HTTP handlers, middleware, models, and in-memory storage. As the project grows, this monolithic structure will become difficult to maintain, test, and extend. Reorganizing into a standard Go project layout improves separation of concerns, testability, and team collaboration.

## What Changes

- **Move models** to `internal/models/` package
- **Move HTTP handlers** to `internal/api/handlers/` package
- **Move middleware** to `internal/api/middleware/` package
- **Move authentication logic** to `internal/service/auth/` package
- **Move user storage** to `internal/repository/user/` package
- **Create configuration management** in `config/` package
- **Refactor `main.go`** to wire dependencies and start server
- **Add interfaces** for repository layer to enable testing and future database integration

## Impact

- **Affected specs**: None (first structural spec)
- **Affected code**:
  - `main.go` - Will be reduced to dependency wiring and server startup
  - `internal/` - New packages will be created
  - `config/` - New configuration package
- **Breaking changes**: None (API endpoints remain unchanged)
- **Benefits**:
  - Clear separation of concerns
  - Easier unit testing with mockable interfaces
  - Ready for database integration
  - Follows standard Go project layout conventions
