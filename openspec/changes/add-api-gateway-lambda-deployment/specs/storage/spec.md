## ADDED Requirements

### Requirement: DynamoDB User Repository
The system SHALL provide a DynamoDB implementation of the `user.Repository` interface that persists users keyed by username.

#### Scenario: Create, read, and update users
- **WHEN** `Create`, `GetByUsername`, `UpdatePassword`, and `Exists` are called
- **THEN** the operations read from and write to a DynamoDB `users` table using `username` as the partition key

#### Scenario: Duplicate username is rejected
- **WHEN** `Create` is called with an existing username
- **THEN** `ErrUserAlreadyExists` is returned without overwriting the existing user

### Requirement: Persistent User Storage
The system SHALL persist registered users across Lambda invocations and instances.

#### Scenario: Users survive Lambda invocations
- **WHEN** a user is registered via one Lambda invocation
- **THEN** a subsequent login invocation retrieves the user from DynamoDB
