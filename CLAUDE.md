<!-- OPENSPEC:START -->
# OpenSpec Instructions

These instructions are for AI assistants working in this project.

Always open `@/openspec/AGENTS.md` when the request:
- Mentions planning or proposals (words like proposal, spec, change, plan)
- Introduces new capabilities, breaking changes, architecture shifts, or big performance/security work
- Sounds ambiguous and you need the authoritative spec before coding

Use `@/openspec/AGENTS.md` to learn:
- How to create and apply change proposals
- Spec format and conventions
- Project structure and guidelines

Keep this managed block so 'openspec update' can refresh the instructions.

<!-- OPENSPEC:END -->

# CLAUDE.md - guebapi

This file provides guidance to Claude Code when working with the **guebapi** repository.

## Repository Overview

**Purpose**: A RESTful API for user authentication and profile management using Go, Gin, and JWT. Provides endpoints for user registration, login, JWT-based authentication, and profile management.

**Status**: Active development

**Technology Stack**:
- Go 1.24
- Gin Framework v1.10.0
- golang-jwt/jwt v4.5.1
- go-playground/validator v10.20.0

## Project Structure

```
guebapi/
├── main.go             # Application entry point with all handlers and middleware
├── go.mod              # Go module definition and dependencies
├── go.sum              # Dependency checksums
├── README.md           # Project documentation
├── .gitignore          # Git ignore rules
├── cmd/                # Command-line application entry points (future)
├── config/             # Configuration files (future)
├── internal/           # Private application code
│   ├── api/            # API handlers (future)
│   ├── pkg/            # Internal packages (future)
│   ├── repository/     # Data access layer (future)
│   └── service/        # Business logic layer (future)
├── tests/              # Test files
├── openspec/           # API specifications
│   ├── specs/          # OpenAPI/specification files
│   └── changes/        # Change tracking
└── .vscode/            # VS Code configuration
    ├── launch.json     # Debug configuration
    ├── tasks.json      # Build/run/test tasks
    └── settings.json   # Editor settings
```

## Quick Start

### Installation

```bash
# Clone the repository
git clone https://github.com/yarieldis/guebapi.git
cd guebapi

# Install dependencies
go mod tidy
```

### Development Server

```bash
# Run the application
go run main.go

# Or using VS Code task
# Ctrl+Shift+B -> Run
```

The server starts on `http://localhost:8080`.

### Build

```bash
# Build the application
go build

# Build with output name
go build -o guebapi.exe
```

### Tests

```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run tests with coverage
go test -cover ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Linting

```bash
# Run Go vet
go vet ./...

# Run staticcheck (if installed)
staticcheck ./...

# Format code
go fmt ./...
```

## API Endpoints

| Method | Endpoint                 | Description              | Auth Required |
|--------|--------------------------|--------------------------|---------------|
| POST   | `/api/register`          | Register a new user      | No            |
| POST   | `/api/login`             | Login and get JWT token  | No            |
| GET    | `/api/protected/profile` | Get user profile         | Yes           |
| POST   | `/api/protected/update`  | Update user password     | Yes           |

### Authentication

Protected routes require a JWT token in the `Authorization` header:
```
Authorization: <jwt_token>
```

## Architecture Patterns

### Current Architecture

The application currently follows a monolithic single-file structure in `main.go` with:
- **Models**: `User` and `Claims` structs
- **Handlers**: HTTP request handlers (`loginHandler`, `registerHandler`, `profileHandler`, `updateHandler`)
- **Middleware**: JWT authentication middleware (`authMiddleware`)
- **In-memory storage**: User data stored in a map

### Planned Architecture (internal/ directories)

- **api/**: HTTP handlers and route definitions
- **service/**: Business logic layer
- **repository/**: Data access layer (database interactions)
- **pkg/**: Shared internal packages

### Gin Router Pattern

Routes are organized into groups:
- `/api` - Public routes (login, register)
- `/api/protected` - Protected routes requiring JWT authentication

## Key Technologies

### Gin Framework

Web framework for routing and middleware:
- Route groups for organizing endpoints
- Middleware support for authentication
- JSON binding and response helpers

### JWT Authentication

JSON Web Tokens for stateless authentication:
- HS256 signing method
- 72-hour token expiration
- Custom claims with username

## Common Development Workflows

### Adding a New Public Endpoint

1. Create handler function in `main.go`:
   ```go
   func newHandler(c *gin.Context) {
       // Handler logic
   }
   ```
2. Register route in the `public` group:
   ```go
   public.POST("/new-endpoint", newHandler)
   ```

### Adding a New Protected Endpoint

1. Create handler function
2. Register route in the `protected` group (middleware applied automatically)

### Working with Request Data

```go
// Bind JSON request body
var data MyStruct
if err := c.BindJSON(&data); err != nil {
    c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
    return
}

// Get data from middleware context
username, exists := c.Get("username")
```

### Sending Responses

```go
// Success response
c.JSON(http.StatusOK, gin.H{"message": "success"})

// Error response
c.JSON(http.StatusBadRequest, gin.H{"error": "error message"})
```

## Configuration

### Environment Variables

Currently hardcoded (to be externalized):
- JWT secret key
- Server port (8080)

### VS Code Configuration

- **launch.json**: Debug configuration with pre-build and post-test tasks
- **tasks.json**: Build, Run, and Test tasks
- **settings.json**: Editor preferences

## Testing

### Running Tests

```bash
# All tests
go test ./...

# Specific package
go test ./internal/service/...

# With coverage
go test -cover ./...
```

### Test Patterns

Use table-driven tests for Go:
```go
func TestHandler(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected int
    }{
        {"valid input", "test", http.StatusOK},
        {"invalid input", "", http.StatusBadRequest},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test logic
        })
    }
}
```

## Important Notes

- **In-memory storage**: User data is stored in memory and will be lost on restart
- **Hardcoded secrets**: JWT key is hardcoded (should use environment variables in production)
- **No password hashing**: Passwords are stored in plain text (implement bcrypt for production)
- **No input validation**: Add validation for production use

## Best Practices

### Code Style

- Follow Go conventions and idioms
- Use `gofmt` for formatting
- Run `go vet` before committing
- Keep functions focused and small

### Error Handling

- Always check and handle errors
- Return meaningful error messages
- Use appropriate HTTP status codes

### Security

- Never commit secrets or credentials
- Use environment variables for configuration
- Implement password hashing (bcrypt)
- Validate and sanitize all inputs

## Dependencies

| Package | Version | Purpose |
|---------|---------|---------|
| gin-gonic/gin | v1.10.0 | HTTP web framework |
| golang-jwt/jwt | v4.5.1 | JWT implementation |
| go-playground/validator | v10.20.0 | Input validation |
| golang.org/x/crypto | v0.23.0 | Cryptographic functions |

## Git Workflow

## Version Control Guidelines

- **NEVER** commit changes without user approval. Ask systematically for approval before committing.
- Commit messages should be clear and follow convention:
  - ai-tooling: AI agents, automation commands, workflows, or other AI-enabled developer tooling
  - feat: New feature
  - fix: Bug fix
  - docs: Documentation
  - style: Formatting
  - refactor: Code restructuring
  - test: Adding tests
  - chore: Maintenance tasks
- **NEVER** mention AI/Claude authorship in commit messages (no "Generated with Claude Code", "AI-assisted", etc.)

## Troubleshooting

### Build Errors

- **Module errors**: Run `go mod tidy` to sync dependencies
- **Import errors**: Check import paths match module name in `go.mod`

### Runtime Errors

- **Port already in use**: Change port or kill process using port 8080
- **Invalid token**: Ensure token is passed without "Bearer " prefix

### Common Issues

- **Authentication fails**: Check JWT token format and expiration
- **User not found**: Users are stored in memory; restart loses data
