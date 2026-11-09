## Context

guebapi is a Go + Gin REST API that currently stores users in an in-memory `map`. The prior Lightsail deployment plans are abandoned in favor of AWS API Gateway + Lambda, managed with Terraform. Because Lambda functions are stateless and invocations may run on different instances, the in-memory store cannot be used in production; DynamoDB replaces it.

## Goals / Non-Goals

**Goals:**
- Run the existing Gin application unchanged (routes, handlers, JWT logic) inside a single Lambda function
- Proxy all routes through an API Gateway HTTP API
- Persist users in DynamoDB behind the existing `user.Repository` interface
- Provision everything with Terraform (declarative, version-controlled)
- Keep local development unchanged (run `main.go` with the in-memory repository)

**Non-Goals:**
- Password hashing (out of scope; plain-text passwords remain a documented limitation)
- API Gateway authorizers / Cognito (the app's JWT middleware keeps handling auth)
- Custom domain or TLS certificates (API Gateway default endpoint)
- CI/CD automation
- Multi-environment workspaces (single production environment)

## Decisions

### 1. Lambda adaptation: `aws-lambda-go-api-proxy`
**Decision**: Wrap the Gin engine with `github.com/awslabs/aws-lambda-go-api-proxy/gin` (`ginadapter`), using the v2 proxy for API Gateway HTTP API payloads.

**Rationale**:
- Minimal code change: the existing `*gin.Engine` is passed to the adapter unchanged
- No container registry required (plain Go binary zip)

**Alternatives considered**:
- AWS Lambda Web Adapter (container image + adapter layer) — zero code changes but requires ECR and a Docker build/push step
- Hand-written Lambda handler per route — more code, loses the existing Gin routing/middleware

### 2. Packaging: zip binary
**Decision**: Build the Go binary for `linux/arm64`, zip it, and point Terraform's `aws_lambda_function` at the archive.

**Rationale**: Lightweight, no ECR/Docker, fast cold starts.

**Trade-off**: A build step (`make build-lambda` or a Terraform `null_resource`) must run before `terraform apply`.

### 3. API Gateway: HTTP API (v2)
**Decision**: Use `aws_apigatewayv2_api` with a `$default` route and `AWS_PROXY` integration to Lambda.

**Rationale**: Cheaper and simpler than REST API (v1); a single `$default` route forwards every request path to the Gin router, which handles its own routing and JWT auth.

### 4. DynamoDB storage
**Decision**: A `users` table with `username` (string) as the partition key, `PAY_PER_REQUEST` billing, accessed via `aws-sdk-go-v2/service/dynamodb`.

**Rationale**: Matches the existing `Repository` interface (`Create`, `GetByUsername`, `UpdatePassword`, `Exists`). `Create` uses a conditional `PutItem` to detect duplicate usernames.

### 5. Terraform state and secrets
**Decision**: Local Terraform state (S3 backend documented as a future option). `JWT_SECRET_KEY` is a Terraform variable marked `sensitive`, supplied via git-ignored `terraform.tfvars`, and injected as a Lambda environment variable.

## Risks / Trade-offs

| Risk | Impact | Mitigation |
|------|--------|------------|
| Cold starts | Slight latency on first requests | Small Go binary; provisioned concurrency optional later |
| Plain-text passwords | Credential exposure if DB leaked | Document; add bcrypt in a follow-up |
| Secret in Terraform state | Low | Mark sensitive; git-ignore tfstate/tfvars; S3 backend later |
| DynamoDB eventual consistency | Stale reads | Strongly-consistent reads for auth lookups |
| Vendor lock-in | Low | Repository interface makes storage swappable |

## Migration Plan

1. Keep local dev on the in-memory repository (`STORAGE_BACKEND=memory`, default).
2. Build and deploy the Lambda with `STORAGE_BACKEND=dynamodb` and `DYNAMODB_TABLE_NAME=users`.
3. Terraform provisions the table, function, and API Gateway; verify `GET /health` and the existing endpoints via the API Gateway URL.
4. Rollback: `terraform destroy`.

## Open Questions

None — using local Terraform state for now; an S3 backend will be adopted when CI/CD or team collaboration is added.
