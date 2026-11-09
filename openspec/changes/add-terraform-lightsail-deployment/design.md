## Context

The `deploy-aws-lightsail` change deploys guebapi to AWS Lightsail Container Service via manual AWS CLI commands and lists "Infrastructure as code" as a non-goal. This change replaces the manual workflow: Terraform becomes the canonical way to manage the Lightsail infrastructure, while the existing Docker, `/health`, and `PORT` prerequisites remain.

## Goals / Non-Goals

**Goals:**
- Provision the Lightsail Container Service with Terraform
- Declare the deployment version (image, port, env vars, health check) as code
- Make the deployment reproducible, reviewable, and version-controlled
- Replace the manual AWS CLI deployment workflow with Terraform

**Non-Goals:**
- Replacing the AWS CLI image push (AWS provider limitation)
- CI/CD automation
- Multi-environment / workspaces (single production environment for now)
- Custom domain or SSL (Lightsail default HTTPS endpoint)

## Decisions

### 1. Terraform with the hashicorp/aws provider
**Decision**: Use HashiCorp Terraform and the official AWS provider.

**Alternatives considered**:
- AWS CDK / CloudFormation — more AWS-native, but Terraform is the most widely adopted IaC tool and has explicit Lightsail resources.
- Pulumi — similar capability, but Terraform's declarative HCL is simpler for a single-service deployment.

### 2. Resources used
```hcl
resource "aws_lightsail_container_service" "guebapi" {
  name  = var.service_name
  power = "nano"
  scale = 1
}

resource "aws_lightsail_container_service_deployment_version" "guebapi" {
  service_name = aws_lightsail_container_service.guebapi.name

  container {
    container_name = "guebapi"
    image          = ":${var.service_name}.${var.image_label}.${var.image_version}"
    environment {
      key   = "PORT"
      value = "8080"
    }
    ports {
      port     = 8080
      protocol = "HTTP"
    }
  }

  public_endpoint {
    container_name = "guebapi"
    container_port = 8080
    health_check {
      path = "/health"
    }
  }
}
```

### 3. Image push stays on the AWS CLI
**Decision**: Keep `aws lightsail push-container-image` as a manual pre-step.

**Rationale**: The AWS provider cannot push container images to the Lightsail registry. Terraform references the image version that was already pushed.

**Trade-off**: Two-step workflow (push then apply). Mitigated by documenting the step clearly and capturing the image version from the push output as a variable.

### 4. State management
**Decision**: Use local Terraform state for simplicity, with an S3 remote backend documented as a future option.

**Rationale**: Single developer / single environment. Local state keeps setup minimal; adopt S3 when team collaboration or CI/CD arrives.

### 5. Secrets
**Decision**: Pass `JWT_SECRET_KEY` via a Terraform variable marked `sensitive`, supplied through a git-ignored `terraform.tfvars`.

**Rationale**: Avoids hardcoding secrets in HCL and matches the existing environment-variable approach.

## Risks / Trade-offs

| Risk | Impact | Mitigation |
|------|--------|------------|
| Image push not automated | Manual step remains | Document clearly; script later if needed |
| Secret stored in state file | Medium | Mark variable sensitive; git-ignore tfstate and tfvars |
| Provider lag for Lightsail features | Low | Pin provider version; run `terraform validate`/`plan` |
| Drift between manual and IaC | Low | Adopt Terraform as the source of truth after this change |

## Migration Plan

1. Ensure Docker image, `/health` endpoint, and `PORT` support (from `deploy-aws-lightsail`) are in place.
2. Build and push the image to the Lightsail registry, capturing the image version.
3. Run `terraform init`, `terraform plan`, then `terraform apply`.
4. Verify `GET /health` returns 200 and existing API endpoints work.
5. Rollback: `terraform destroy`.

## Open Questions

- Should Terraform adopt an S3 remote backend now, or stay local until team collaboration is needed?
