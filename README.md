# Kong Connect Services API

A RESTful API for managing and retrieving organizational services with versioning support.

## Overview

This API provides endpoints to:
- List services with pagination, filtering, and sorting
- Retrieve individual services with their versions
- Support for searching services by name or description

## Architecture & Design Decisions

### Technology Stack
- **Go 1.24**: Primary language for performance and simplicity
- **Gorilla Mux**: HTTP router for clean URL patterns and middleware support
- **SQLite**: Lightweight database perfect for this use case, easy to setup and deploy
- **Standard Library**: Minimal dependencies for better maintainability

### Architecture Pattern
The application follows a layered architecture:

```
cmd/server/          # Application entry point
internal/
├── handlers/        # HTTP handlers (Presentation layer)
├── service/         # Business logic (Service layer)
├── repository/      # Data access (Repository layer)
└── models/          # Data structures
pkg/
└── database/        # Database connection and setup
```

### Design Considerations

1. **Separation of Concerns**: Each layer has a single responsibility
    - Handlers: HTTP request/response handling
    - Service: Business logic and validation
    - Repository: Data access operations

2. **Database Choice**: SQLite was chosen for:
    - Zero configuration required
    - Perfect for read-heavy workloads
    - Easy deployment and testing
    - Sufficient for the scope of this assignment

3. **API Design**: RESTful endpoints following standard conventions
    - GET /api/v1/services - List all services
    - GET /api/v1/services/{id} - Get specific service

4. **Pagination**: Implemented to handle large datasets efficiently
    - Default page size: 12 (matching UI grid)
    - Maximum page size: 100 (prevent abuse)

## API Endpoints

### GET /api/v1/services

Retrieve a paginated list of services with optional filtering and sorting.

**Query Parameters:**
- `search` (string): Search in service name or description
- `sort_by` (string): Sort field (name, created_at, updated_at)
- `sort_dir` (string): Sort direction (asc, desc)
- `page` (int): Page number (default: 1)
- `page_size` (int): Items per page (default: 12, max: 100)

**Example Request:**
```bash
curl "http://localhost:8080/api/v1/services?search=contact&sort_by=name&sort_dir=asc&page=1&page_size=10"
```

**Response:**
```json
{
  "services": [
    {
      "id": 1,
      "name": "Contact Us",
      "description": "Lorem ipsum dolor sit amet...",
      "created_at": "2023-01-01T00:00:00Z",
      "updated_at": "2023-01-01T00:00:00Z",
      "versions": [
        {
          "id": 1,
          "service_id": 1,
          "version": "2.0.0",
          "created_at": "2023-01-01T00:00:00Z"
        }
      ]
    }
  ],
  "total": 1,
  "page": 1,
  "page_size": 10,
  "total_pages": 1
}
```

### GET /api/v1/services/{id}

Retrieve a specific service by ID with all its versions.

**Example Request:**
```bash
curl "http://localhost:8080/api/v1/services/1"
```

**Response:**
```json
{
  "id": 1,
  "name": "Contact Us",
  "description": "Lorem ipsum dolor sit amet...",
  "created_at": "2023-01-01T00:00:00Z",
  "updated_at": "2023-01-01T00:00:00Z",
  "versions": [
    {
      "id": 1,
      "service_id": 1,
      "version": "2.0.0",
      "created_at": "2023-01-01T00:00:00Z"
    }
  ]
}
```

### GET /health

Health check endpoint.

**Response:** `OK` (200 status)

## Getting Started

### Prerequisites
- Go 1.21 or higher
- Git

### Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd services-api
```

2. Install dependencies:
```bash
go mod tidy
```

3. Run the server:
```bash
go run cmd/server/main.go
```

The server will start on port 8080 by default.

### Environment Variables

- `PORT`: Server port (default: 8080)
- `DB_PATH`: Database file path (default: ./services.db)

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific test
go test ./internal/service/
```

## Development

### Project Structure

```
services-api/
├── cmd/server/                 # Application entry point
│   └── main.go
├── internal/                   # Private application code
│   ├── handlers/              # HTTP handlers
│   │   └── service_handler.go
│   ├── models/                # Data models
│   │   └── service.go
│   ├── repository/            # Data access layer
│   │   └── service_repository.go
│   └── service/               # Business logic layer
│       ├── service_service.go
│       └── service_service_test.go
├── pkg/                       # Public packages
│   └── database/              # Database utilities
│       └── connection.go
├── tests/                     # Integration tests
├── docs/                      # Documentation
├── go.mod                     # Go modules
├── go.sum                     # Go modules checksum
└── README.md
```

### Adding New Features

1. **Models**: Define data structures in `internal/models/`
2. **Repository**: Add data access methods in `internal/repository/`
3. **Service**: Implement business logic in `internal/service/`
4. **Handlers**: Add HTTP endpoints in `internal/handlers/`
5. **Tests**: Write tests alongside your code

## Trade-offs and Assumptions

### Trade-offs Made

1. **SQLite vs PostgreSQL**: Chose SQLite for simplicity and zero configuration
    -  Easy setup and deployment
    -  Perfect for read-heavy workloads
    -  Limited concurrent write performance
    -  Less suitable for high-scale production

2. **In-memory vs Persistent Storage**: Used file-based SQLite
    -  Data persists between restarts
    -  Can be backed up easily
    -  Slightly slower than in-memory

3. **Custom vs Framework**: Used minimal dependencies
    -  Smaller binary size
    -  Better performance
    -  More boilerplate code

### Assumptions Made

1. **Read-Heavy Workload**: API is primarily for displaying services
2. **Moderate Scale**: Hundreds to thousands of services, not millions
3. **Simple Search**: Basic text search is sufficient for MVP
4. **Version Ordering**: Newer versions should appear first
5. **Default Pagination**: 12 items per page matches the UI grid

## Future Enhancements

If given more time, the following features could be added:

### Authentication & Authorization
- JWT-based authentication
- Role-based access control
- API key authentication

### Advanced Features
- Full-text search with indexing
- Service categories/tags
- Caching layer (Redis)
- Rate limiting
- Metrics and monitoring

### CRUD Operations
- POST /api/v1/services - Create service
- PUT /api/v1/services/{id} - Update service
- DELETE /api/v1/services/{id} - Delete service
- POST /api/v1/services/{id}/versions - Add version

### Production Readiness
- Configuration management
- Structured logging
- Graceful shutdown
- Database migrations
- Docker containerization
- CI/CD pipeline

## Performance Considerations

- Database indexes on searchable fields
- Connection pooling for concurrent requests
- Pagination to limit memory usage
- Efficient SQL queries with proper joins
- HTTP middleware for common concerns (CORS, logging)

## Testing Strategy

- **Unit Tests**: Business logic validation
- **Integration Tests**: Database operations
- **API Tests**: HTTP endpoint functionality
- **Load Tests**: Performance under stress

Current test coverage focuses on the service layer as it contains the core business logic.