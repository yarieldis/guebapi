# Change: Deploy guebapi to AWS API Gateway + Lambda with DynamoDB

## Why

The Lightsail deployment plans are abandoned. Lightsail requires a managed container service and costs ~$7/month even when idle. API Gateway + Lambda is a serverless, pay-per-request model better suited to a small API. This change replaces the two Lightsail proposals with a serverless deployment managed by Terraform, and moves user storage from in-memory to DynamoDB (required because Lambda functions are stateless and concurrent invocations do not share memory).

## What Changes

- Adapt the Go application to run on AWS Lambda using the `aws-lambda-go-api-proxy` Gin adapter (routes and business logic unchanged)
- Add a DynamoDB implementation of the existing `user.Repository` interface
- Add a `GET /health` endpoint for deployment health checks
- Add Terraform configuration that provisions: DynamoDB `users` table, Lambda function (with IAM role and policies), and an API Gateway HTTP API (`$default` route proxying to Lambda)
- Remove the two Lightsail change proposals (superseded)
- Update `README.md` deployment instructions and `.gitignore`

## Impact

- Affected specs: `deployment`, `storage`
- Affected code:
  - `internal/repository/user/dynamodb.go` (new)
  - `internal/app/` (new; shared router + wiring used by both server and Lambda)
  - `cmd/lambda/main.go` (new Lambda entry point)
  - `main.go` (refactor to use shared wiring; still the local HTTP server)
  - `internal/api/router/router.go` (`/health` endpoint)
  - `config/` (DynamoDB table name, AWS region, storage backend selection)
  - `terraform/` (new)
  - `README.md`, `.gitignore`
- Breaking changes: None (local development unchanged; in-memory repository remains the default for local dev and tests)
