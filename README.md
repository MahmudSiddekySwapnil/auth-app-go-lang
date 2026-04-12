# auth-app-go-lang

Authentication App built with Go

## Setup & Installation

```bash
# Initialize Go module
go mod init auth-app

# Install Dependencies

# Gin Web Framework - Fast HTTP web framework for building REST APIs and routing
# Why: Provides routing, middleware support, and JSON binding for HTTP requests
go get github.com/gin-gonic/gin

# PostgreSQL Driver (pq - pure Go) - Database driver for PostgreSQL
# Why: Allows Go application to connect, query, and execute commands on PostgreSQL database
go get github.com/lib/pq

# JWT (JSON Web Token) v5 - Library for creating and verifying JWT tokens
# Why: Used for user authentication - generates tokens on login and validates them for protected routes
go get github.com/golang-jwt/jwt/v5

# Bcrypt - Password hashing library for secure password storage
# Why: Encrypts user passwords before storing in database for security; prevents plain-text password storage
go get golang.org/x/crypto/bcrypt

# Godotenv - Loads environment variables from .env file
# Why: Reads DB_URL, JWT_SECRET, PORT from .env file without hardcoding sensitive credentials
go get github.com/joho/godotenv
```

## Running the Application

```bash
# From project root directory
go run ./cmd/main.go
```

Server runs on `http://localhost:8080`

## Environment Variables (.env)

```
PORT=8080
DB_URL=postgres://user:password@localhost:5432/dbname?sslmode=disable
JWT_SECRET=your-secret-key
```