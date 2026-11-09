## 1. Application Changes

- [ ] 1.1 Add health check endpoint (`GET /health`) to router
- [ ] 1.2 Update config loader to support `PORT` environment variable

## 2. Docker Configuration

- [ ] 2.1 Create `Dockerfile` with multi-stage build for Go application
- [ ] 2.2 Create `.dockerignore` to exclude unnecessary files

## 3. Documentation

- [ ] 3.1 Add deployment section to README.md with:
  - Prerequisites (AWS CLI, credentials)
  - Lightsail container service creation
  - Image build and push commands
  - Deployment commands
  - Environment variable configuration
  - How to get public URL

## 4. Local Validation

- [ ] 4.1 Test Docker build locally: `docker build -t guebapi .`
- [ ] 4.2 Test Docker run locally: `docker run -p 8080:8080 -e JWT_SECRET_KEY=test -e PORT=8080 guebapi`
- [ ] 4.3 Verify health endpoint responds: `curl http://localhost:8080/health`
- [ ] 4.4 Verify existing API endpoints still work
