# Change: Configure AWS Deployment with Terraform

## Why

The `deploy-aws-lightsail` change deploys the application to AWS Lightsail Container Service using manual AWS CLI commands and explicitly lists "Infrastructure as code" as a non-goal. Manual commands are error-prone, hard to review, and not reproducible. This change replaces that manual workflow: managing the Lightsail infrastructure with Terraform (Infrastructure as Code) makes the deployment declarative, version-controlled, and repeatable.

## What Changes

- Add a `terraform/` directory with HashiCorp Terraform configuration using the `hashicorp/aws` provider
- Declare the Lightsail Container Service (`aws_lightsail_container_service`, nano power, scale 1) as code
- Declare the container deployment version (`aws_lightsail_container_service_deployment_version`) — image, port, environment variables, and `/health` public endpoint — as code
- Parameterize region, service name, image version, and secrets via Terraform variables (`JWT_SECRET_KEY` marked sensitive)
- Add `.gitignore` entries for Terraform state and variable files
- Replace the manual AWS CLI deployment steps in `README.md` with the Terraform-based workflow
- **Note**: The container image push remains an AWS CLI step (`aws lightsail push-container-image`) because the AWS provider does not support pushing images to the Lightsail registry; Terraform references the already-pushed image version

## Impact

- Affected specs: `deployment`
- Affected code:
  - `terraform/` - New directory (Terraform configuration)
  - `README.md` - Deployment instructions
  - `.gitignore` - Terraform state/secret exclusions
- Breaking changes: None (local development unchanged)
- Relationship: Replaces the manual deployment steps of `deploy-aws-lightsail`, which still defines the Docker image, `/health` endpoint, and `PORT` prerequisites that this change depends on
