# Project Structure Specification

## ADDED Requirements

### Requirement: Layered Package Architecture
The project SHALL organize code into distinct layers following Go conventions: models, repository, service, and API (handlers/middleware/router).

#### Scenario: Package imports follow dependency direction
- **WHEN** inspecting package imports
- **THEN** handlers import services, services import repositories, repositories import models
- **AND** no circular dependencies exist

#### Scenario: Internal packages are private
- **WHEN** application code is placed under `internal/`
- **THEN** external packages cannot import internal packages

### Requirement: Model Package
The system SHALL define domain models in `internal/models/` package, separate from business logic and data access.

#### Scenario: User model definition
- **WHEN** the `User` model is needed
- **THEN** it is imported from `internal/models`
- **AND** contains `Username` and `Password` fields with JSON tags

#### Scenario: Claims model definition
- **WHEN** JWT claims are needed
- **THEN** the `Claims` struct is imported from `internal/models`
- **AND** embeds `jwt.RegisteredClaims`

### Requirement: Repository Interface Pattern
The system SHALL define repository interfaces to abstract data storage, enabling testability and storage backend swapping.

#### Scenario: UserRepository interface
- **WHEN** user data access is needed
- **THEN** code depends on `UserRepository` interface, not concrete implementation
- **AND** the interface defines `Create`, `GetByUsername`, `UpdatePassword`, and `Exists` methods

#### Scenario: In-memory implementation
- **WHEN** the application runs without a database
- **THEN** an in-memory `UserRepository` implementation is used
- **AND** user data is stored in a thread-safe map

### Requirement: Service Layer
The system SHALL encapsulate business logic in service packages under `internal/service/`.

#### Scenario: AuthService interface
- **WHEN** authentication operations are needed
- **THEN** code depends on `AuthService` interface
- **AND** the interface defines `GenerateToken`, `ValidateToken`, and `HashPassword` methods

#### Scenario: Service depends on repository
- **WHEN** `AuthService` is instantiated
- **THEN** it receives a `UserRepository` via constructor injection
- **AND** does not directly access storage

### Requirement: HTTP Handler Separation
The system SHALL place HTTP handlers in `internal/api/handlers/` with explicit service dependencies.

#### Scenario: Handler constructor injection
- **WHEN** a handler is created
- **THEN** required services are passed via constructor
- **AND** handlers do not instantiate their own dependencies

#### Scenario: Handler grouping
- **WHEN** handlers are organized
- **THEN** authentication handlers (login, register) are in `auth.go`
- **AND** profile handlers (profile, update) are in `profile.go`

### Requirement: Middleware Package
The system SHALL place HTTP middleware in `internal/api/middleware/` package.

#### Scenario: JWT authentication middleware
- **WHEN** protected routes are accessed
- **THEN** JWT middleware from `internal/api/middleware` validates the token
- **AND** sets the username in request context

### Requirement: Router Configuration
The system SHALL configure routes in `internal/api/router/` package with injected handlers.

#### Scenario: Router setup function
- **WHEN** the router is initialized
- **THEN** `SetupRouter` function receives handler instances
- **AND** configures public and protected route groups

### Requirement: Configuration Management
The system SHALL manage configuration through a dedicated `config/` package supporting environment variables.

#### Scenario: Configuration loading
- **WHEN** the application starts
- **THEN** configuration is loaded from environment variables
- **AND** defaults are provided for development

#### Scenario: Sensitive values externalized
- **WHEN** JWT secret key is needed
- **THEN** it is read from `JWT_SECRET` environment variable
- **AND** is never hardcoded in source files

### Requirement: Minimal main.go
The `main.go` file SHALL only perform dependency wiring and server startup, containing no business logic.

#### Scenario: Dependency wiring
- **WHEN** `main()` executes
- **THEN** it loads configuration
- **AND** creates repository, service, and handler instances
- **AND** configures the router
- **AND** starts the HTTP server

#### Scenario: No business logic in main
- **WHEN** `main.go` is reviewed
- **THEN** it contains no HTTP handler implementations
- **AND** contains no authentication logic
- **AND** contains no data access code
