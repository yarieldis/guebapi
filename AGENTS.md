# guebapi

This file provides guidance to AI coding assistants when working with the **guebapi** repository.

## Repository Overview

**Purpose**: A RESTful API for user authentication and profile management using Go, Gin, and JWT. Provides endpoints for user registration, login, JWT-based authentication, and profile management. Runs locally as an HTTP server and deploys to AWS as a serverless API (API Gateway + Lambda + DynamoDB).

**Status**: Active development

**Technology Stack**:
- Go 1.26
- Gin Framework v1.10.0
- golang-jwt/jwt v4.5.1
- aws-lambda-go v1.55.1
- aws-lambda-go-api-proxy v0.16.2
- aws-sdk-go-v2 (DynamoDB)
- Terraform (hashicorp/aws ~> 5.0)

## Project Structure

```
guebapi/
├── main.go                           # Local HTTP server entry point (in-memory repo)
├── Makefile                          # Build/test/lambda-package targets
├── go.mod / go.sum                   # Go module and dependency checksums
├── README.md                         # Project documentation
├── .gitignore                        # Git ignore rules
├── .env.example                      # Environment variables template
├── config/                           # Configuration management
│   ├── config.go                     # Configuration structs (server, JWT, storage, AWS)
│   └── loader.go                     # Environment variable loader
├── cmd/
│   └── lambda/
│       └── main.go                   # AWS Lambda entry point (DynamoDB repo)
├── internal/                         # Private application code
│   ├── app/
│   │   └── app.go                    # Shared dependency wiring (BuildRouter)
│   ├── api/                          # API layer
│   │   ├── handlers/                 # HTTP request handlers
│   │   │   ├── auth.go               # Login/register handlers
│   │   │   ├── profile.go            # Profile handlers
│   │   │   └── handlers_test.go      # Handler tests
│   │   ├── middleware/               # HTTP middleware
│   │   │   ├── auth.go               # JWT authentication
│   │   │   └── auth_test.go          # Middleware tests
│   │   └── router/                   # Route configuration
│   │       └── router.go             # Route setup (incl. /health)
│   ├── models/                       # Domain models
│   │   ├── user.go                   # User struct
│   │   └── claims.go                 # JWT claims struct
│   ├── repository/                   # Data access layer
│   │   └── user/                     # User storage
│   │       ├── repository.go         # Repository interface
│   │       ├── memory.go             # In-memory implementation (local dev)
│   │       ├── dynamodb.go           # DynamoDB implementation (production)
│   │       ├── memory_test.go        # Repository tests
│   │       └── dynamodb_test.go      # DynamoDB repository tests
│   └── service/                      # Business logic layer
│       └── auth/                     # Authentication service
│           ├── service.go            # Service interface
│           ├── jwt.go                # JWT implementation
│           └── jwt_test.go           # Service tests
└── terraform/                        # Infrastructure as code (AWS)
    ├── versions.tf                   # Terraform + AWS provider requirements
    ├── variables.tf                  # Input variables
    ├── main.tf                       # DynamoDB, Lambda, API Gateway, IAM
    ├── outputs.tf                    # API URL and table name outputs
    └── terraform.tfvars.example      # Variable values template
```

## Quick Start

### Installation

```bash
git clone https://github.com/yarieldis/guebapi.git
cd guebapi
go mod tidy
```

### Development Server

```bash
go run main.go
```

The server starts on `http://localhost:8080` and uses the in-memory repository (seeded with `john`/`password123` and `jane`/`password456`).

### Build

```bash
go build ./...
```

### Tests

```bash
go test ./...
go test -cover ./...
```

### Linting

```bash
go vet ./...
gofmt -l .
```

## API Endpoints

| Method | Endpoint                   | Description              | Auth Required |
|--------|----------------------------|--------------------------|---------------|
| GET    | `/health`                  | Health check             | No            |
| POST   | `/api/register`            | Register a new user      | No            |
| POST   | `/api/login`               | Login and get JWT token  | No            |
| GET    | `/api/protected/profile`   | Get user profile         | Yes           |
| POST   | `/api/protected/update`    | Update user password     | Yes           |

### Authentication

Protected routes require a JWT token in the `Authorization` header:

```
Authorization: <jwt_token>
```

## Architecture Patterns

### Layered Architecture

The application follows a standard Go project layout with clear separation of concerns:

- **Models** (`internal/models/`): Domain models (`User`, `Claims`)
- **Repository** (`internal/repository/`): Data access layer with interface-based design
  - `user.Repository` interface for user storage operations
  - `user.MemoryRepository` in-memory implementation (local dev and tests)
  - `user.DynamoDBRepository` persistent implementation (production on Lambda)
- **Service** (`internal/service/`): Business logic layer
  - `auth.Service` interface for authentication operations
  - `auth.JWTService` implementation with JWT token handling
- **Handlers** (`internal/api/handlers/`): HTTP request handlers
  - `AuthHandler` for login/register endpoints
  - `ProfileHandler` for profile management
- **Middleware** (`internal/api/middleware/`): Cross-cutting concerns
  - JWT authentication middleware
- **Router** (`internal/api/router/`): Route configuration and setup
- **App** (`internal/app/`): Shared dependency wiring used by both entry points
- **Config** (`config/`): Externalized configuration via environment variables

### Dependency Flow

```
Local server (main.go)                  Lambda (cmd/lambda/main.go)
  └── config.Load()                       └── config.Load()
  └── memory repository                   └── AWS config + DynamoDB client
  └── app.BuildRouter(cfg, repo)          └── DynamoDB repository
  └── r.Run(addr)                         └── app.BuildRouter(cfg, repo)
                                          └── lambda.Start(ginadapter.NewV2(r).ProxyWithContext)
```

`app.BuildRouter` wires the inner layers:

```
app.BuildRouter(cfg, userRepo)
  └── auth.NewJWTService(repo, cfg.JWT.SecretKey, cfg.JWT.TokenDuration)
  └── handlers.NewAuthHandler(authService)
  └── handlers.NewProfileHandler(userRepo)
  └── router.SetupRouter(authHandler, profileHandler, authService)
```

### Storage Backends

The `user.Repository` interface allows swapping storage backends:

- **`memory`** (default): in-memory map with `sync.RWMutex`; local dev and tests only; data lost on restart
- **`dynamodb`**: persistent `users` table keyed on `username`; used in the Lambda/API Gateway deployment

Select the backend via the `STORAGE_BACKEND` environment variable.

### Gin Router Pattern

Routes are organized into groups:

- `/health` - Health check (no auth)
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
- Configurable token expiration (default 72 hours)
- Custom claims with username

### AWS Lambda + API Gateway

The app runs unchanged on Lambda via the `aws-lambda-go-api-proxy` Gin adapter:
- `ginadapter.NewV2(r)` adapts API Gateway HTTP API (payload v2) events to the Gin engine
- The whole router is wrapped, so routes and middleware are identical to local dev
- `$default` route on an HTTP API proxies every path to the single Lambda function

### DynamoDB

Persistent user storage for the serverless deployment:
- `users` table, `username` string partition key, `PAY_PER_REQUEST` billing
- `Create` uses a conditional `PutItem` (`attribute_not_exists(username)`) to reject duplicates

## Common Development Workflows

### Adding a New Public Endpoint

1. Create or modify a handler in `internal/api/handlers/`:
   ```go
   func (h *MyHandler) NewEndpoint(c *gin.Context) {
       // Handler logic
   }
   ```
2. Register the route in `internal/api/router/router.go`:
   ```go
   public.POST("/new-endpoint", myHandler.NewEndpoint)
   ```
3. Add tests in `internal/api/handlers/handlers_test.go`

### Adding a New Protected Endpoint

1. Create a handler method in the appropriate handler struct
2. Register the route in the `protected` group in `router.go` (middleware applied automatically)
3. Add tests with an authentication token

### Adding a New Storage Backend

1. Implement the `user.Repository` interface in `internal/repository/user/`
2. Wire it into `main.go` (local) or `cmd/lambda/main.go` (Lambda) as appropriate
3. If configurable, add an option to `config.StorageConfig`

### Working with Request Data

```go
var data MyStruct
if err := c.BindJSON(&data); err != nil {
    c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
    return
}

username, exists := c.Get("username")
```

### Sending Responses

```go
c.JSON(http.StatusOK, gin.H{"message": "success"})
c.JSON(http.StatusBadRequest, gin.H{"error": "error message"})
```

## Configuration

### Environment Variables

Configuration is loaded from environment variables with sensible defaults:

| Variable              | Description                             | Default            |
|-----------------------|-----------------------------------------|--------------------|
| `SERVER_HOST`         | Server host                             | `` (empty)         |
| `SERVER_PORT`         | Server port                             | `8080`             |
| `JWT_SECRET_KEY`      | Secret key for JWT signing              | `secret_key_here`  |
| `JWT_TOKEN_DURATION`  | JWT token duration                      | `72h`              |
| `STORAGE_BACKEND`     | Storage backend (`memory`/`dynamodb`)   | `memory`           |
| `DYNAMODB_TABLE_NAME` | DynamoDB table name                     | `users`            |
| `AWS_REGION`          | AWS region                              | `us-east-1`        |

Copy `.env.example` to `.env` and modify as needed for local development.

## Deployment (AWS API Gateway + Lambda)

See the `## Deployment` section in `README.md` for the full workflow. Summary:

1. Build the Lambda package: `GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -o bootstrap ./cmd/lambda && zip lambda.zip bootstrap`
2. Configure `terraform/terraform.tfvars` (from `terraform.tfvars.example`)
3. Run `terraform init`, `terraform plan`, `terraform apply`
4. Read the URL: `terraform output api_url`
5. Tear down with `terraform destroy`

Terraform provisions the DynamoDB table, Lambda function (with IAM role and policies), and API Gateway HTTP API.

## Testing

### Running Tests

```bash
go test ./...
go test ./internal/repository/...   # specific package
go test -cover ./...
```

### Test Patterns

Use table-driven tests for Go:

```go
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
```

Repository and service layers are unit-tested with mocks/fakes (e.g., the DynamoDB repository is tested against an in-memory fake of the `dynamoAPI` interface).

## Important Notes

- **Storage**: in-memory storage is for local dev/tests only; production uses DynamoDB
- **Default secrets**: JWT key has a default value; set `JWT_SECRET_KEY` in production
- **No password hashing**: passwords are stored in plain text (implement bcrypt for production)
- **No input validation**: beyond JSON binding; add validation for production use
- **Terraform state**: local for now; adopt an S3 backend before team collaboration or CI/CD

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
- Use environment variables and Terraform variables (`sensitive`) for configuration
- Implement password hashing (bcrypt)
- Validate and sanitize all inputs

## Dependencies

| Package | Version | Purpose |
|---------|---------|---------|
| gin-gonic/gin | v1.10.0 | HTTP web framework |
| golang-jwt/jwt | v4.5.1 | JWT implementation |
| aws-lambda-go | v1.55.1 | AWS Lambda runtime |
| aws-lambda-go-api-proxy | v0.16.2 | Gin → Lambda adapter |
| aws-sdk-go-v2 | v1.47.8 | AWS SDK (DynamoDB) |
| go-playground/validator | v10.20.0 | Input validation |

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
- **Import errors**: Check import paths match the module name in `go.mod`

### Runtime Errors

- **Port already in use**: Change `SERVER_PORT` or kill the process using port 8080
- **Invalid token**: Ensure the token is passed without the "Bearer " prefix

### Lambda / DynamoDB

- **Lambda timeouts or errors**: Check CloudWatch Logs for the function
- **DynamoDB access denied**: Verify the Lambda IAM role has `dynamodb:GetItem`, `dynamodb:PutItem`, and `dynamodb:UpdateItem` on the table
- **API Gateway 500**: Confirm the `$default` route targets the Lambda integration and the `AWS_PROXY` payload format version is `2.0`

### Common Issues

- **Authentication fails**: Check JWT token format and expiration
- **User not found**: Users live in DynamoDB in production; local dev uses in-memory storage and loses data on restart
