# Golang JWT Auth API

A simple RESTful API for user authentication and profile management using Go, Gin, and JWT.

## Features

- **User Registration:** Create new users.
- **User Login:** Authenticate users and receive a JWT token.
- **JWT Authentication:** Protect routes using JWT middleware.
- **Profile Management:** View and update user profile (password).

## Project Structure

```
guebapi/
├── main.go                           # Application entry point (dependency wiring)
├── config/                           # Configuration management
│   ├── config.go                     # Configuration structs
│   └── loader.go                     # Environment variable loader
├── internal/                         # Private application code
│   ├── api/                          # API layer
│   │   ├── handlers/                 # HTTP request handlers
│   │   │   ├── auth.go               # Login/register handlers
│   │   │   └── profile.go            # Profile handlers
│   │   ├── middleware/               # HTTP middleware
│   │   │   └── auth.go               # JWT authentication
│   │   └── router/                   # Route configuration
│   │       └── router.go             # Route setup
│   ├── models/                       # Domain models
│   │   ├── user.go                   # User struct
│   │   └── claims.go                 # JWT claims struct
│   ├── repository/                   # Data access layer
│   │   └── user/                     # User storage
│   │       ├── repository.go         # Repository interface
│   │       └── memory.go             # In-memory implementation
│   └── service/                      # Business logic layer
│       └── auth/                     # Authentication service
│           ├── service.go            # Service interface
│           └── jwt.go                # JWT implementation
├── go.mod                            # Go module definition
└── go.sum                            # Dependency checksums
```

## Endpoints

| Method | Endpoint                   | Description                | Auth Required |
|--------|----------------------------|----------------------------|--------------|
| POST   | `/api/register`            | Register a new user        | No           |
| POST   | `/api/login`               | Login and get JWT token    | No           |
| GET    | `/api/protected/profile`   | Get user profile           | Yes          |
| POST   | `/api/protected/update`    | Update user password       | Yes          |
| GET    | `/health`                  | Health check               | No           |

## Getting Started

### Prerequisites

- Go 1.26+
- [Gin](https://github.com/gin-gonic/gin)
- [JWT](https://github.com/golang-jwt/jwt)

### Installation

1. Clone the repository:
    ```sh
    git clone https://github.com/yarieldis/guebapi.git
    cd guebapi
    ```

2. Install dependencies:
    ```sh
    go mod tidy
    ```

3. (Optional) Configure environment variables:
    ```sh
    cp .env.example .env
    # Edit .env with your settings
    ```

### Configuration

The application can be configured via environment variables:

| Variable            | Description                  | Default            |
|---------------------|------------------------------|--------------------|
| `SERVER_HOST`       | Server host                  | `` (empty)         |
| `SERVER_PORT`       | Server port                  | `8080`             |
| `JWT_SECRET_KEY`      | Secret key for JWT signing | `secret_key_here` |
| `JWT_TOKEN_DURATION`  | JWT token duration         | `72h`             |
| `STORAGE_BACKEND`     | Storage backend (`memory`/`dynamodb`) | `memory` |
| `DYNAMODB_TABLE_NAME` | DynamoDB table name        | `users`           |
| `AWS_REGION`          | AWS region                 | `us-east-1`       |

### Running the Server

```sh
go run main.go
```

The server starts on `http://localhost:8080` by default.

### Running Tests

```sh
go test ./...
```

## Usage Examples

### Register a new user

```sh
curl -X POST http://localhost:8080/api/register \
  -H "Content-Type: application/json" \
  -d '{"username": "newuser", "password": "mypassword"}'
```

### Login

```sh
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"username": "john", "password": "password123"}'
```

### Access protected route

```sh
curl http://localhost:8080/api/protected/profile \
  -H "Authorization: <your-jwt-token>"
```

## Deployment (AWS API Gateway + Lambda)

The application can be deployed to AWS as a serverless API: API Gateway (HTTP API) fronts a Lambda function, with user data persisted in DynamoDB.

### Prerequisites

- [Terraform](https://www.terraform.io/) 1.5+
- [AWS CLI](https://aws.amazon.com/cli/) with credentials configured
- Go 1.26+

### 1. Build the Lambda package

Build the Go binary for Lambda (Linux, arm64) and zip it:

```sh
# macOS / Linux / Git Bash
GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -o bootstrap ./cmd/lambda
zip lambda.zip bootstrap
```

On Windows (PowerShell):

```powershell
$env:GOOS = "linux"; $env:GOARCH = "arm64"; $env:CGO_ENABLED = "0"
go build -o bootstrap ./cmd/lambda
Compress-Archive -Path bootstrap -DestinationPath lambda.zip
```

### 2. Configure Terraform variables

```sh
cd terraform
cp terraform.tfvars.example terraform.tfvars
# Edit terraform.tfvars with your AWS region and a real JWT secret
```

### 3. Deploy

```sh
terraform init
terraform plan
terraform apply
```

### 4. Verify

```sh
terraform output api_url
curl <api_url>/health
```

Then register/login as normal against `<api_url>`.

### 5. Tear down

```sh
terraform destroy
```

## Architecture

The project follows a standard Go layout with clear separation of concerns:

- **Repository Layer**: Abstracts data access with interfaces, enabling easy swapping of storage backends
- **Service Layer**: Contains business logic (authentication, JWT operations)
- **Handler Layer**: HTTP request handling and response formatting
- **Middleware**: Cross-cutting concerns like authentication
- **Configuration**: Externalized configuration via environment variables

This architecture enables:
- Unit testing with mock dependencies
- Easy database integration (replace in-memory repository)
- Clear boundaries between layers
