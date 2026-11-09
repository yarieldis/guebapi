# Tasks: Deploy guebapi to AWS API Gateway + Lambda with DynamoDB

## 1. Application Changes
- [ ] 1.1 Extract shared router/wiring into `internal/app` so both server and Lambda entry points reuse it
- [ ] 1.2 Add `GET /health` endpoint to the router
- [ ] 1.3 Add config for storage backend, DynamoDB table name, and AWS region

## 2. DynamoDB Repository
- [ ] 2.1 Add `aws-sdk-go-v2` DynamoDB dependency
- [ ] 2.2 Implement `user.DynamoDBRepository` satisfying `user.Repository` (Create with conditional PutItem, GetByUsername, UpdatePassword, Exists)
- [ ] 2.3 Add unit tests using the DynamoDB local/test doubles

## 3. Lambda Entry Point
- [ ] 3.1 Add `github.com/awslabs/aws-lambda-go-api-proxy` and `github.com/aws/aws-lambda-go` dependencies
- [ ] 3.2 Add `cmd/lambda/main.go` that builds the router, wraps it with `ginadapter` (HTTP API/v2), and calls `lambda.Start`
- [ ] 3.3 Keep `main.go` as the local HTTP server entry point

## 4. Terraform
- [ ] 4.1 Create `terraform/` with the `hashicorp/aws` provider and required version constraints
- [ ] 4.2 Provision `aws_dynamodb_table` (`users`, hash key `username`, PAY_PER_REQUEST)
- [ ] 4.3 Provision Lambda IAM role + policies (DynamoDB access, CloudWatch Logs)
- [ ] 4.4 Provision `aws_lambda_function` from the zipped binary with env vars (`JWT_SECRET_KEY`, `JWT_TOKEN_DURATION`, `STORAGE_BACKEND=dynamodb`, `DYNAMODB_TABLE_NAME`)
- [ ] 4.5 Provision `aws_apigatewayv2_api` + stage + `$default` route + `AWS_PROXY` integration
- [ ] 4.6 Add variables (`aws_region`, `jwt_secret_key` sensitive, `jwt_token_duration`) and an output for the API Gateway URL
- [ ] 4.7 Add `.gitignore` entries for `.terraform/`, `*.tfstate`, `*.tfstate.backup`, `*.tfvars`, and the build output zip

## 5. Documentation
- [ ] 5.1 Replace the Lightsail deployment docs in `README.md` with the API Gateway + Lambda + Terraform workflow
- [ ] 5.2 Document the pre-apply build step and the DynamoDB/Lambda architecture

## 6. Validation
- [ ] 6.1 `go build ./...` and `go test ./...` pass locally
- [ ] 6.2 `terraform fmt` and `terraform validate` pass
- [ ] 6.3 `terraform apply` provisions the stack; `GET /health` returns 200 via the API Gateway URL
- [ ] 6.4 Register → login → profile flow works through API Gateway
- [ ] 6.5 `terraform destroy` tears everything down
