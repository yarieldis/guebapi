# guebapi — Manual

This document is the operator/developer manual for **guebapi**, a RESTful API for user
authentication and profile management built with Go, Gin, and JWT.

It runs in two modes:

1. **Local** — a plain HTTP server backed by an in-memory user store (for development).
2. **Serverless** — AWS API Gateway + Lambda, backed by DynamoDB (for production).

---

## 1. Overview

- **Language**: Go 1.26
- **Web framework**: Gin v1.10.0
- **Auth**: JWT (HS256), configurable expiry (default 72h)
- **Storage**: in-memory (local) or DynamoDB (production)
- **Infrastructure**: Terraform (AWS provider)

### Architecture (layered)

```
HTTP request
   → Router (Gin, route groups)
      → Middleware (JWT auth on protected routes)
         → Handlers (auth, profile)
            → Service (auth.JWTService)
               → Repository (memory | dynamodb)
```

Shared wiring lives in `internal/app.BuildRouter(cfg, repo)`, reused by both entry points:

- `main.go` — local HTTP server
- `cmd/lambda/main.go` — AWS Lambda (wraps the Gin engine with `ginadapter.NewV2`)

---

## 2. Endpoints

Base URL: `http://localhost:8080` (local) or your API Gateway URL (production).

### 2.1 Health

| Method | Path | Auth |
|--------|------|------|
| GET | `/health` | No |

Returns `200 OK` with `{"status":"ok"}` when the app is running.

```sh
curl -i http://localhost:8080/health
```

### 2.2 Register

| Method | Path | Auth |
|--------|------|------|
| POST | `/api/register` | No |

Request body (JSON):

```json
{ "username": "alice", "password": "s3cret" }
```

Responses:

| Status | Body |
|--------|------|
| 200 | `{"message":"user created successfully"}` |
| 400 | `{"error":"username already exists"}` |
| 400 | `{"error":"invalid request"}` |

```sh
curl -i -X POST http://localhost:8080/api/register \
  -H "Content-Type: application/json" \
  -d '{"username":"alice","password":"s3cret"}'
```

### 2.3 Login

| Method | Path | Auth |
|--------|------|------|
| POST | `/api/login` | No |

Request body (JSON):

```json
{ "username": "alice", "password": "s3cret" }
```

Responses:

| Status | Body |
|--------|------|
| 200 | `{"token":"<jwt>"}` |
| 401 | `{"error":"invalid credentials"}` |
| 400 | `{"error":"invalid request"}` |

```sh
curl -i -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"username":"alice","password":"s3cret"}'
```

Copy the returned token for the protected routes below.

### 2.4 Get profile (protected)

| Method | Path | Auth |
|--------|------|------|
| GET | `/api/protected/profile` | Yes |

Header: `Authorization: <token>` (raw token, **no** `Bearer ` prefix).

Responses:

| Status | Body |
|--------|------|
| 200 | `{"username":"alice","message":"Welcome to your profile!"}` |
| 401 | `{"error":"missing token"}` or `{"error":"invalid token"}` |

```sh
curl -i http://localhost:8080/api/protected/profile \
  -H "Authorization: <your-jwt-token>"
```

### 2.5 Update password (protected)

| Method | Path | Auth |
|--------|------|------|
| POST | `/api/protected/update` | Yes |

Header: `Authorization: <token>`.

Request body (JSON):

```json
{ "password": "new-password" }
```

Responses:

| Status | Body |
|--------|------|
| 200 | `{"message":"profile updated successfully"}` |
| 400 | `{"error":"invalid request"}` |
| 401 | `{"error":"missing token"}` or `{"error":"invalid token"}` |
| 500 | `{"error":"internal server error"}` |

```sh
curl -i -X POST http://localhost:8080/api/protected/update \
  -H "Authorization: <your-jwt-token>" \
  -H "Content-Type: application/json" \
  -d '{"password":"new-password"}'
```

> **Important**: the JWT middleware reads the token directly from the `Authorization`
> header and does **not** strip a `Bearer ` prefix. Pass the raw token.

---

## 3. Using locally

### Prerequisites

- Go 1.26+
- (optional) `make` for convenience targets

### Start

```sh
git clone https://github.com/yarieldis/guebapi.git
cd guebapi
go mod tidy
go run main.go
```

The server listens on `http://localhost:8080` and seeds two users:

| username | password |
|----------|----------|
| `john` | `password123` |
| `jane` | `password456` |

### Quick smoke test

```sh
# register
curl -s -X POST localhost:8080/api/register -H "Content-Type: application/json" \
  -d '{"username":"bob","password":"pw"}'

# login (capture the token)
TOKEN=$(curl -s -X POST localhost:8080/api/login -H "Content-Type: application/json" \
  -d '{"username":"bob","password":"pw"}' | sed -E 's/.*"token":"([^"]+)".*/\1/')

# profile
curl -s localhost:8080/api/protected/profile -H "Authorization: $TOKEN"

# update password
curl -s -X POST localhost:8080/api/protected/update \
  -H "Authorization: $TOKEN" -H "Content-Type: application/json" \
  -d '{"password":"newpw"}'
```

---

## 4. Configuration

Configuration is read from environment variables. Copy `.env.example` to `.env` for local
development (values are read from the process environment, so export them or use a dotenv loader).

| Variable | Description | Default |
|----------|-------------|---------|
| `SERVER_HOST` | Server host | (empty) |
| `SERVER_PORT` | Server port | `8080` |
| `JWT_SECRET_KEY` | JWT signing key | `secret_key_here` |
| `JWT_TOKEN_DURATION` | JWT expiry duration | `72h` |
| `STORAGE_BACKEND` | `memory` or `dynamodb` | `memory` |
| `DYNAMODB_TABLE_NAME` | DynamoDB table | `users` |
| `AWS_REGION` | AWS region | `us-east-1` |

Example:

```sh
SERVER_PORT=9000 JWT_SECRET_KEY=dev-secret go run main.go
```

> `STORAGE_BACKEND` selects the repository. The local server (`main.go`) always uses the
> in-memory store; the Lambda entry point (`cmd/lambda/main.go`) always uses DynamoDB.

---

## 5. Testing

### Run tests

```sh
go test ./...                  # all packages
go test ./internal/repository/user/...   # single package
go test -v ./...               # verbose
go test -race ./...            # race detector
```

### Coverage

```sh
go test -cover ./...
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out   # open in browser
```

### Static checks / formatting

```sh
go vet ./...
gofmt -l .        # lists files needing formatting (empty = clean)
```

### What is covered

- `internal/repository/user` — memory and DynamoDB repositories (DynamoDB uses an in-memory fake of the SDK client)
- `internal/service/auth` — JWT login, register, token generation/validation
- `internal/api/handlers` — HTTP-level route tests via `httptest`
- `internal/api/middleware` — missing/invalid/valid token

### Makefile targets

```sh
make build        # go build ./...
make test         # go test ./...
make build-lambda # build + zip the Lambda binary
make clean        # remove lambda.zip and bootstrap
```

---

## 6. Deployment (AWS API Gateway + Lambda + DynamoDB)

The production topology is a serverless API:

```
Client → API Gateway (HTTP API, $default route) → Lambda (Gin app) → DynamoDB (users)
```

### Prerequisites

- Terraform 1.5+
- AWS CLI with credentials configured
- Go 1.26+
- `zip` (Linux/macOS/Git Bash) or PowerShell `Compress-Archive` (Windows)

### 6.1 Build the Lambda package

The Lambda uses a Go binary (custom runtime `provided.al2023`, arm64) named `bootstrap`,
zipped as `lambda.zip`.

```sh
# macOS / Linux / Git Bash
GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -o bootstrap ./cmd/lambda
zip lambda.zip bootstrap
rm -f bootstrap
```

Windows (PowerShell):

```powershell
$env:GOOS="linux"; $env:GOARCH="arm64"; $env:CGO_ENABLED="0"
go build -o bootstrap ./cmd/lambda
Compress-Archive -Path bootstrap -DestinationPath lambda.zip
```

(Or use `make build-lambda` where `make` is available.)

### 6.2 Configure Terraform

```sh
cd terraform
cp terraform.tfvars.example terraform.tfvars
```

Edit `terraform.tfvars`:

```hcl
aws_region         = "us-east-1"
service_name       = "guebapi"
lambda_zip_path    = "../lambda.zip"   # path to the zip built above
dynamodb_table_name = "users"
jwt_secret_key     = "a-long-random-secret"
jwt_token_duration = "72h"
```

`terraform.tfvars` is git-ignored and must not be committed.

### 6.3 Apply

```sh
cd terraform
terraform init
terraform plan
terraform apply
```

This provisions:

- DynamoDB `users` table (partition key `username`, PAY_PER_REQUEST)
- Lambda function + IAM role/policy (DynamoDB + CloudWatch Logs)
- API Gateway HTTP API with a `$default` route proxying to Lambda

### 6.4 Verify

```sh
terraform output api_url
curl <api_url>/health                    # expect {"status":"ok"}
```

Then exercise register/login/profile/update against `<api_url>` exactly as in section 3
(remember: raw token in the `Authorization` header).

### 6.5 Tear down

```sh
terraform destroy
```

### Terraform notes

- State is **local** by default (`.terraform/`, `*.tfstate` are git-ignored). Adopt an S3
  backend before team/CI use.
- `.terraform.lock.hcl` is committed to pin provider versions.
- `JWT_SECRET_KEY` is marked `sensitive` and injected as a Lambda environment variable.

---

## 7. CI/CD (GitHub Actions)

### Current state

`.github/workflows/go.yml` runs on every push and pull request to `master`:

- `go build -v ./...`
- `go test -v ./...` (Go 1.26)

### Recommended CI additions

- `go vet ./...`
- `gofmt -l .` (fail if any file is unformatted)
- `terraform fmt -check -recursive` + `terraform validate` (requires `hashicorp/setup-terraform`)

### What CD (auto-deploy) requires

1. **Remote Terraform state** — GitHub runners are ephemeral, so local state would be lost
   every run. Requires an **S3 backend + DynamoDB lock table**.
2. **AWS credentials for GitHub** — via **OIDC federation** (recommended, no static keys) or
   static `AWS_ACCESS_KEY_ID` / `AWS_SECRET_ACCESS_KEY` repo secrets.
3. **`JWT_SECRET_KEY`** — as a repo secret, injected via `TF_VAR_jwt_secret_key`.

### Proposed workflow

| Trigger | Action |
|---------|--------|
| Pull request to `master` | build + test + vet + `terraform plan` (no changes applied) |
| Push/merge to `master` | build + test + build Lambda zip + `terraform apply` |
| `workflow_dispatch` | manual deploy button |

The deploy job would:

1. Cross-compile the Go binary (`linux/arm64`) and zip it → `lambda.zip`
2. Run `terraform init` (S3 backend), `terraform plan`, `terraform apply -auto-approve`
3. Print `terraform output api_url`

### Backend bootstrap (chicken-and-egg)

The S3 bucket + DynamoDB lock table must exist before Terraform can use them as its backend:

- **Option A**: add a `terraform/bootstrap/` (or script) that creates those two resources once.
- **Option B**: create the S3 bucket + lock table manually, then wire the backend config to use them.

---

## 8. Debugging

### Local server

- Gin logs every request to stdout. Errors return JSON with a short `error` message.
- `go run main.go` uses the in-memory store; registered users are lost on restart.
- Set `SERVER_PORT` to change the port if 8080 is taken.

### VS Code

The repo ships `.vscode/` config:

- **launch.json** — `Launch Go Program` (builds first, runs tests after via tasks).
- **tasks.json** — `Build`, `Run`, `Test` tasks.

### Lambda (CloudWatch)

- Lambda logs go to **CloudWatch Logs** (`/aws/lambda/guebapi`). Check there for panics or
  startup errors (e.g., AWS credential/config failures).
- A cold start is normal; subsequent requests reuse the warm container.

### DynamoDB

- **Access denied** → verify the Lambda IAM policy grants `dynamodb:GetItem`,
  `dynamodb:PutItem`, and `dynamodb:UpdateItem` on the table ARN.
- **"user not found" in production** → users live in DynamoDB; confirm you registered
  against the deployed URL, not localhost.

### API Gateway

- **HTTP 500 / missing responses** → confirm:
  - the `$default` route targets the Lambda integration,
  - the integration is `AWS_PROXY` with `payload_format_version = "2.0"`,
  - the `aws_lambda_permission` allows `apigateway.amazonaws.com` to invoke the function.
- **Binary/garbled response** → the app returns JSON; ensure content-type headers are correct.

### Common issues

| Symptom | Likely cause / fix |
|---------|--------------------|
| `401 invalid token` | Token expired or malformed; ensure you send the raw token without `Bearer `. |
| `401 missing token` | `Authorization` header absent. |
| `401 invalid credentials` | Wrong username/password. |
| `400 username already exists` | Duplicate registration. |
| `dial tcp :8080: connect: connection refused` | Server not running or wrong port. |
| User disappears after restart | In-memory store (local) — expected; use DynamoDB in production. |

---

## 9. Security notes (current limitations)

- Passwords are stored in **plain text** — add bcrypt before real production use.
- Input validation is minimal (JSON binding only) — add stricter validation.
- `JWT_SECRET_KEY` has a default value — always set it in production.
- Terraform state is local — use an S3 backend for shared/CI workflows.
