## ADDED Requirements

### Requirement: Terraform Infrastructure Provisioning
The system SHALL provide Terraform configuration that provisions the AWS Lightsail Container Service.

#### Scenario: Terraform applies successfully
- **WHEN** `terraform apply` is executed with valid AWS credentials
- **THEN** a Lightsail Container Service is created
- **AND** the service uses the nano power tier with a single node

### Requirement: Declarative Container Deployment
The system SHALL define the container deployment version (image, port, environment variables, and health check) as Terraform configuration.

#### Scenario: Deployment version references a pushed image
- **WHEN** the Terraform deployment version references a valid image version already pushed to the Lightsail registry
- **THEN** the container is deployed with the configured port and environment variables
- **AND** the public endpoint health check uses `/health`

### Requirement: Terraform Variables and Secrets
The system SHALL parameterize deployment configuration via Terraform variables, with `JWT_SECRET_KEY` marked as sensitive.

#### Scenario: Secrets are parameterized
- **WHEN** a git-ignored `terraform.tfvars` file supplies the JWT secret
- **THEN** the secret is passed to the container without being committed to version control

### Requirement: Terraform State Exclusion
The system SHALL exclude Terraform state files and variable files from version control.

#### Scenario: State files are not committed
- **WHEN** Terraform runs in the project directory
- **THEN** `*.tfstate`, `*.tfstate.backup`, `*.tfvars`, and `.terraform/` are ignored by git
