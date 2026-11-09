## Context

The guebapi application is a Go REST API using Gin framework. It needs to be deployed to AWS with minimal operational complexity and predictable costs. User preferences:
- Lightsail Container Service
- In-memory storage (no database)
- Single environment (production)
- Manual deployments
- No infrastructure as code
- Default AWS URL (no custom domain)

## Goals / Non-Goals

**Goals:**
- Deploy the application to AWS Lightsail Container Service
- Keep deployment process simple and repeatable
- Minimize cost (~$7/month for nano tier)
- Externalize secrets via environment variables

**Non-Goals:**
- Persistent storage (data will be lost on container restart)
- CI/CD automation
- Infrastructure as code (Terraform/CDK) — superseded by `add-terraform-lightsail-deployment`
- Custom domain or SSL certificates (Lightsail provides default HTTPS endpoint)
- Multi-environment setup
- Auto-scaling (single container instance)

## Decisions

### 1. Docker-based deployment
**Decision**: Use multi-stage Docker build for minimal image size.

**Rationale**:
- Lightsail Container Service requires Docker images
- Multi-stage build produces small images (~15-20MB)
- Scratch or distroless base for minimal attack surface

### 2. Lightsail Container Service configuration
**Decision**: Use Nano tier (512MB RAM, 0.25 vCPU) with 1 node.

**Rationale**:
- Sufficient for a lightweight Go API
- Lowest cost tier (~$7/month)
- Can scale up later if needed

### 3. Image registry
**Decision**: Use Lightsail's built-in container image registry.

**Rationale**:
- No need for ECR or Docker Hub
- Simpler workflow with AWS CLI
- Images pushed directly to Lightsail

### 4. Environment configuration
**Decision**: Use Lightsail environment variables for configuration.

**Rationale**:
- Built-in secret management
- Easy to modify via AWS Console or CLI

**Configuration mapping**:
| App Config | Environment Variable |
|------------|---------------------|
| SERVER_PORT | PORT (set to container port) |
| JWT_SECRET_KEY | JWT_SECRET_KEY |
| JWT_TOKEN_DURATION | JWT_TOKEN_DURATION |

### 5. Health check endpoint
**Decision**: Add a `/health` endpoint for container health monitoring.

**Rationale**:
- Lightsail uses health checks to determine container status
- Enables automatic container restart on failure
- Standard practice for containerized applications

## Risks / Trade-offs

| Risk | Impact | Mitigation |
|------|--------|------------|
| Data loss on restart | High | Document clearly; user accepted in-memory storage |
| Single container SPOF | Medium | Lightsail auto-restarts unhealthy containers |
| Limited scaling | Low | Can upgrade tier or add nodes later |
| Vendor lock-in | Low | Standard Docker image works anywhere |

## Migration Plan

No migration needed - this is a new deployment capability. Existing local development workflow remains unchanged.

## Deployment Steps (Manual)

> **Note**: These manual steps are superseded by the `add-terraform-lightsail-deployment` change, which replaces them with Terraform. Retained here for reference only.

1. Install AWS CLI and configure credentials
2. Build Docker image: `docker build -t guebapi .`
3. Create Lightsail container service:
   ```bash
   aws lightsail create-container-service \
     --service-name guebapi \
     --power nano \
     --scale 1
   ```
4. Push image to Lightsail:
   ```bash
   aws lightsail push-container-image \
     --service-name guebapi \
     --label latest \
     --image guebapi:latest
   ```
5. Deploy container with environment variables:
   ```bash
   aws lightsail create-container-service-deployment \
     --service-name guebapi \
     --containers '{
       "guebapi": {
         "image": ":guebapi.latest.X",
         "ports": {"8080": "HTTP"},
         "environment": {
           "JWT_SECRET_KEY": "<your-secret>",
           "JWT_TOKEN_DURATION": "72h",
           "PORT": "8080"
         }
       }
     }' \
     --public-endpoint '{
       "containerName": "guebapi",
       "containerPort": 8080,
       "healthCheck": {
         "path": "/health"
       }
     }'
   ```
6. Get public URL:
   ```bash
   aws lightsail get-container-services --service-name guebapi
   ```

## Open Questions

None - all requirements clarified.
