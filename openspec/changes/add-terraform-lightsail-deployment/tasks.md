# Tasks: Configure AWS Deployment with Terraform

## 1. Terraform Project Setup
- [ ] 1.1 Create `terraform/` directory with `hashicorp/aws` provider and version constraint
- [ ] 1.2 Add `required_providers` block and an `aws_region` variable
- [ ] 1.3 Add `.gitignore` entries for `.terraform/`, `*.tfstate`, `*.tfstate.backup`, and `*.tfvars`

## 2. Lightsail Container Service Resource
- [ ] 2.1 Define `aws_lightsail_container_service` with nano power and scale 1
- [ ] 2.2 Parameterize the service name via a variable

## 3. Deployment Version Resource
- [ ] 3.1 Define `aws_lightsail_container_service_deployment_version` referencing the pushed image
- [ ] 3.2 Parameterize the image label/version via a variable
- [ ] 3.3 Configure container port 8080 (HTTP) and environment variables (`PORT`, `JWT_SECRET_KEY`, `JWT_TOKEN_DURATION`)
- [ ] 3.4 Configure the public endpoint with a `/health` health check

## 4. Variables and Outputs
- [ ] 4.1 Define variables for `aws_region`, `service_name`, `image_version`, `jwt_secret_key`, and `jwt_token_duration`
- [ ] 4.2 Mark `jwt_secret_key` as sensitive
- [ ] 4.3 Define an output for the public endpoint URL

## 5. Documentation
- [ ] 5.1 Replace the manual AWS CLI deployment steps in `README.md` with the Terraform workflow (init, plan, apply, destroy)
- [ ] 5.2 Document the image push step (`aws lightsail push-container-image`) that must precede `terraform apply`

## 6. Validation
- [ ] 6.1 Run `terraform fmt` and `terraform validate`
- [ ] 6.2 Run `terraform plan` against a clean environment
- [ ] 6.3 Verify `terraform apply` provisions the service and `GET /health` returns 200
- [ ] 6.4 Verify existing API endpoints still work after deployment
- [ ] 6.5 Verify `terraform destroy` removes the service
