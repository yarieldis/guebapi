## ADDED Requirements

### Requirement: Lambda Compatibility
The application SHALL run as an AWS Lambda function behind API Gateway using the `aws-lambda-go-api-proxy` Gin adapter, exposing the same routes, handlers, and JWT middleware as the local server.

#### Scenario: Lambda handler serves existing routes
- **WHEN** a request is sent through the API Gateway to the Lambda function
- **THEN** the Gin router handles the request using the same routing and middleware as local development

### Requirement: API Gateway HTTP API
The system SHALL be exposed through an AWS API Gateway HTTP API that proxies all routes to the Lambda function.

#### Scenario: Requests reach the application
- **WHEN** a client calls any route on the API Gateway URL
- **THEN** the request is forwarded to the Lambda function and the response is returned to the client

### Requirement: Health Check Endpoint
The system SHALL provide a health check endpoint at `GET /health` that returns HTTP 200 when the application is running.

#### Scenario: Health check success
- **WHEN** a GET request is made to `/health`
- **THEN** the response status code is 200
- **AND** the response body contains `{"status": "ok"}`

### Requirement: Terraform Infrastructure Provisioning
The system SHALL provide Terraform configuration that provisions the DynamoDB table, the Lambda function with an IAM role and policies, and the API Gateway HTTP API.

#### Scenario: Terraform applies successfully
- **WHEN** `terraform apply` is executed with valid AWS credentials
- **THEN** the DynamoDB table, Lambda function, IAM role, and API Gateway are created

### Requirement: Secrets and State Exclusion
The system SHALL parameterize secrets via Terraform variables marked `sensitive`, and SHALL exclude Terraform state files, variable files, and build artifacts from version control.

#### Scenario: Secrets are not committed
- **WHEN** a git-ignored `terraform.tfvars` supplies the JWT secret
- **THEN** the secret is injected as a Lambda environment variable without being committed to version control

#### Scenario: Terraform state files are ignored
- **WHEN** Terraform runs in the project directory
- **THEN** `.terraform/`, `*.tfstate`, `*.tfstate.backup`, and `*.tfvars` are ignored by git
