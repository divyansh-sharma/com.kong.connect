# Kong Connect Services API

A RESTful API for managing and retrieving organizational services with versioning support.

## Overview

This API provides endpoints to:

* List services with pagination, filtering, and sorting
* Retrieve individual services with their versions
* Support for searching services by name or description
* Token-based authentication and role-based access control

## Architecture & Design Decisions

### Technology Stack

* **Go 1.24**: Primary language for performance and simplicity
* **Gorilla Mux**: HTTP router for clean URL patterns and middleware support
* **SQLite**: Lightweight database perfect for this use case, easy to setup and deploy
* **Standard Library**: Minimal dependencies for better maintainability

### Architecture Pattern

The application follows a layered architecture:

```
cmd/server/          # Application entry point
internal/
├── handlers/        # HTTP handlers (Presentation layer)
├── service/         # Business logic (Service layer)
├── repository/      # Data access (Repository layer)
├── middleware/      # Authentication & Authorization
└── models/          # Data structures
pkg/
└── database/        # Database connection and setup
```

### Design Considerations

1. **Separation of Concerns**: Each layer has a single responsibility
2. **Database Choice**: SQLite for its simplicity and zero config
3. **API Design**: RESTful endpoints
4. **Pagination**: Supports large datasets efficiently

##  Authentication & Authorization

This API supports **token-based authentication** and **role-based authorization** for all endpoints under `/api/v1`.

### Authentication

All requests must include a valid **Bearer Token** in the `Authorization` header:

```
Authorization: Bearer <token>
```

#### Supported Tokens (for development/testing):

| Token           | Role     | Access Level       |
| --------------- | -------- | ------------------ |
| `admin-token`   | `admin`  | Full access        |
| `viewer-token`  | `viewer` | Read-only access   |
| *Invalid token* | -        | `401 Unauthorized` |

> In production, replace this with proper JWT validation.

### Authorization

Role-based access control is enforced via middleware:

* `admin` and `viewer` roles can **read services**
* (Planned) Only `admin` will be allowed to **create/update/delete**

### Authenticated Request Examples

```bash
# List services with viewer role
curl -H "Authorization: Bearer viewer-token" \
     "http://localhost:8080/api/v1/services?search=chat&page=1&page_size=5"

# Get specific service with admin role
curl -H "Authorization: Bearer admin-token" \
     "http://localhost:8080/api/v1/services/1"

# Missing or invalid token
curl "http://localhost:8080/api/v1/services"
# Response: 401 Unauthorized
```

### Auth Internals

* **middleware/auth.go**

   * `AuthMiddleware`: Validates the token and injects user context
   * `RoleAuthorization`: Ensures user has required role(s)
* **Route Protection (in `main.go`)**

  ```go
  router.Use(AuthMiddleware)
  api.Use(RoleAuthorization("admin", "viewer"))
  ```

---

## API Endpoints

### GET /api/v1/services

Retrieve a paginated list of services with optional filtering and sorting.

**Query Parameters:**

* `search` (string): Search in service name or description
* `sort_by` (string): Sort field (name, created\_at, updated\_at)
* `sort_dir` (string): Sort direction (asc, desc)
* `page` (int): Page number (default: 1)
* `page_size` (int): Items per page (default: 12, max: 100)

**Example Request:**

```bash
curl -H "Authorization: Bearer viewer-token" \
     "http://localhost:8080/api/v1/services?search=contact&sort_by=name&sort_dir=asc&page=1&page_size=10"
```

### GET /api/v1/services/{id}

Retrieve a specific service by ID with all its versions.

**Example Request:**

```bash
curl -H "Authorization: Bearer admin-token" \
     "http://localhost:8080/api/v1/services/1"
```

### GET /health

Health check endpoint.

**Response:** `OK` (200 status)

---

## Getting Started

### Prerequisites

* Go 1.24 or higher
* Git

### Installation

```bash
git clone <repository-url>
cd services-api
go mod tidy
go run cmd/server/main.go
```

The server will start on port 8080 by default.

### Environment Variables

* `PORT`: Server port (default: 8080)
* `DB_PATH`: Database file path (default: ./services.db)

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific test
go test ./internal/service/
```

---

## Development

### Project Structure

```
services-api/
├── cmd/server/
│   └── main.go
├── internal/
│   ├── handlers/
│   ├── models/
│   ├── repository/
│   ├── service/
│   └── middleware/          # Auth middleware lives here
├── pkg/
│   └── database/
├── tests/
├── docs/
├── go.mod
├── go.sum
└── README.md
```

### Adding New Features

1. **Models**: Define data structures in `internal/models/`
2. **Repository**: Add data access methods in `internal/repository/`
3. **Service**: Implement business logic in `internal/service/`
4. **Handlers**: Add HTTP endpoints in `internal/handlers/`
5. **Tests**: Write unit and integration tests

---

## Trade-offs and Assumptions

### Trade-offs Made

1. **SQLite vs PostgreSQL**: Chose SQLite for simplicity and zero configuration
2. **In-memory vs Persistent**: Used file-based SQLite for data persistence
3. **Custom vs Framework**: Used minimal dependencies for better control

### Assumptions Made

1. Read-heavy workload
2. Moderate dataset (1000s of records)
3. Basic search is sufficient
4. Default pagination matches UI grid

---

## Future Enhancements

### Auth & Security

* Full JWT validation (with secret/key rotation)
* Role-specific access control per route
* Token expiry, refresh tokens

### Advanced Features

* CRUD for services and versions
* Caching layer (Redis)
* Service tags/categories
* Full-text search
* Docker and CI/CD integration

---

## Performance Considerations

* Database indexes for search/sort
* Pagination to limit memory usage
* HTTP middleware for CORS, logging, and auth
* Efficient query design

---

## Testing Strategy

* **Unit Tests**: Service layer logic
* **Integration Tests**: DB interactions
* **API Tests**: Endpoint behavior
* **Load Tests**: Scalability under pressure

Current coverage is strongest in the service layer.

---
