# Golang JWT Auth API

A simple RESTful API for user authentication and profile management using Go, Gin, and JWT.

## Features

- **User Registration:** Create new users.
- **User Login:** Authenticate users and receive a JWT token.
- **JWT Authentication:** Protect routes using JWT middleware.
- **Profile Management:** View and update user profile (password).

## Endpoints

| Method | Endpoint                   | Description                | Auth Required |
|--------|----------------------------|----------------------------|--------------|
| POST   | `/api/register`            | Register a new user        | No           |
| POST   | `/api/login`               | Login and get JWT token    | No           |
| GET    | `/api/protected/profile`   | Get user profile           | Yes          |
| POST   | `/api/protected/update`    | Update user password       | Yes          |

## Getting Started

### Prerequisites

- Go 1.18+
- [Gin](https://github.com/gin-gonic/gin)
- [JWT](https://github.com/golang-jwt/jwt)

### Installation

1. Clone the repository:
    ```sh
    git clone https://github.com/yarieldis/guebapi.git
    cd guebapi
    ```

2. Install dependencies:
    ```sh
    go mod tidy
    ```

### Running the Server

```sh
go run main.go
```
