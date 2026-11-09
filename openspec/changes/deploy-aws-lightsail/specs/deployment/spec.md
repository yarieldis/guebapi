## ADDED Requirements

### Requirement: Health Check Endpoint
The system SHALL provide a health check endpoint at `GET /health` that returns HTTP 200 when the application is running.

#### Scenario: Health check success
- **WHEN** a GET request is made to `/health`
- **THEN** the response status code is 200
- **AND** the response body contains `{"status": "ok"}`

### Requirement: Docker Containerization
The application SHALL be packaged as a Docker container using a multi-stage build to minimize image size.

#### Scenario: Docker build succeeds
- **WHEN** `docker build -t guebapi .` is executed in the project root
- **THEN** a Docker image is created successfully
- **AND** the image size is under 30MB

#### Scenario: Docker container runs correctly
- **WHEN** the container is started with required environment variables
- **THEN** the application starts and listens on the configured port
- **AND** all API endpoints are accessible

### Requirement: Lightsail Container Service Compatibility
The application SHALL be deployable to AWS Lightsail Container Service.

#### Scenario: Container deployment succeeds
- **WHEN** the Docker image is pushed to Lightsail and deployed
- **THEN** the container starts successfully
- **AND** the health check passes
- **AND** the public endpoint is accessible

#### Scenario: Environment variable configuration
- **WHEN** `JWT_SECRET_KEY`, `JWT_TOKEN_DURATION`, and `PORT` are set as environment variables
- **THEN** the application uses these values for configuration
- **AND** the application starts successfully

### Requirement: PORT Environment Variable Support
The application SHALL use the `PORT` environment variable when present, defaulting to 8080 if not set.

#### Scenario: Custom port via environment
- **WHEN** the `PORT` environment variable is set to `5000`
- **THEN** the application listens on port 5000

#### Scenario: Default port when not specified
- **WHEN** the `PORT` environment variable is not set
- **THEN** the application listens on port 8080
