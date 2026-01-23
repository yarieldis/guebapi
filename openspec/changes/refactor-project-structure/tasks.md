# Tasks: Refactor Project Structure

## 1. Setup Package Structure
- [ ] 1.1 Create `internal/models/` directory
- [ ] 1.2 Create `internal/repository/user/` directory
- [ ] 1.3 Create `internal/service/auth/` directory
- [ ] 1.4 Create `internal/api/handlers/` directory
- [ ] 1.5 Create `internal/api/middleware/` directory
- [ ] 1.6 Create `internal/api/router/` directory
- [ ] 1.7 Create `config/` directory

## 2. Extract Models
- [ ] 2.1 Create `internal/models/user.go` with `User` struct
- [ ] 2.2 Create `internal/models/claims.go` with `Claims` struct
- [ ] 2.3 Verify models compile independently

## 3. Implement Repository Layer
- [ ] 3.1 Create `internal/repository/user/repository.go` with `UserRepository` interface
- [ ] 3.2 Create `internal/repository/user/memory.go` implementing in-memory storage
- [ ] 3.3 Write unit tests for memory repository in `internal/repository/user/memory_test.go`

## 4. Implement Authentication Service
- [ ] 4.1 Create `internal/service/auth/service.go` with `AuthService` interface
- [ ] 4.2 Create `internal/service/auth/jwt.go` implementing JWT operations
- [ ] 4.3 Write unit tests for auth service in `internal/service/auth/jwt_test.go`

## 5. Extract Middleware
- [ ] 5.1 Create `internal/api/middleware/auth.go` with JWT middleware
- [ ] 5.2 Write unit tests for middleware in `internal/api/middleware/auth_test.go`

## 6. Extract HTTP Handlers
- [ ] 6.1 Create `internal/api/handlers/auth.go` with login and register handlers
- [ ] 6.2 Create `internal/api/handlers/profile.go` with profile and update handlers
- [ ] 6.3 Write unit tests for handlers using mock services

## 7. Create Router Package
- [ ] 7.1 Create `internal/api/router/router.go` with route configuration
- [ ] 7.2 Define `SetupRouter` function accepting handler dependencies

## 8. Implement Configuration
- [ ] 8.1 Create `config/config.go` with configuration structs
- [ ] 8.2 Create `config/loader.go` for loading from environment variables
- [ ] 8.3 Create `.env.example` documenting required variables
- [ ] 8.4 Update `.gitignore` to exclude `.env` files

## 9. Refactor main.go
- [ ] 9.1 Remove all extracted code from `main.go`
- [ ] 9.2 Add configuration loading
- [ ] 9.3 Wire dependencies (repository → service → handlers → router)
- [ ] 9.4 Start server with configured port

## 10. Validation
- [ ] 10.1 Run `go build` to verify compilation
- [ ] 10.2 Run `go test ./...` to verify all tests pass
- [ ] 10.3 Manual test: POST `/api/register` creates user
- [ ] 10.4 Manual test: POST `/api/login` returns JWT token
- [ ] 10.5 Manual test: GET `/api/protected/profile` with valid token returns profile
- [ ] 10.6 Manual test: POST `/api/protected/update` updates password
- [ ] 10.7 Run `go vet ./...` for static analysis

## 11. Documentation
- [ ] 11.1 Update README.md with new project structure
- [ ] 11.2 Update CLAUDE.md to reflect actual architecture

## Dependencies

- Tasks 2.x must complete before 3.x (models needed for repository)
- Tasks 3.x must complete before 4.x (repository needed for service)
- Tasks 4.x and 5.x can run in parallel
- Tasks 6.x depend on 4.x (handlers need services)
- Task 7.x depends on 5.x and 6.x (router needs middleware and handlers)
- Tasks 8.x can run in parallel with 2.x-7.x
- Task 9.x depends on all previous tasks
- Task 10.x depends on 9.x
- Task 11.x depends on 10.x
