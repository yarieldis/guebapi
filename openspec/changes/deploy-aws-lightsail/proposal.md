# Change: Deploy Application to AWS Lightsail Container Service

## Why

The application currently runs only locally. Deploying to AWS Lightsail Container Service provides a simple, cost-effective production environment with predictable pricing (~$7/month), automatic HTTPS, and managed container orchestration.

## What Changes

- **Add Dockerfile** for containerized deployment
- **Add `.dockerignore`** to exclude unnecessary files from image
- **Add health check endpoint** (`GET /health`) for container health monitoring
- **Update configuration** to support `PORT` environment variable
- **Add deployment documentation** with step-by-step Lightsail instructions

## Impact

- **Affected specs**: None (new deployment capability)
- **Affected code**:
  - `Dockerfile` - New file for container build
  - `.dockerignore` - New file to optimize image
  - `internal/api/router/router.go` - Add health endpoint
  - `config/loader.go` - Support PORT env var
  - `README.md` - Deployment instructions
- **Breaking changes**: None (local development unchanged)
- **Deployment model**: Single production environment, manual deployments, in-memory storage (ephemeral)
