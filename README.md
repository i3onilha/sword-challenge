# Task Management API

A RESTful API for managing maintenance tasks performed by technicians, with manager oversight and notifications.

## Features

- User roles: Manager and Technician
- Task management (CRUD operations)
- Role-based access control with JWT authentication
- Real-time notifications for managers using RabbitMQ
- MySQL database for data persistence
- Unit tests for core functionality
- Swagger API documentation
- CORS support
- Input validation and sanitization
- Error handling and logging

## Prerequisites

- Go 1.23.1 or later
- MySQL 8.0 or later
- RabbitMQ 3.8 or later
- Docker and Docker Compose (for containerized development)
- Kubernetes (for containerized deployment)

## Getting Started

1. Clone the repository:
```bash
git clone https://github.com/i3onilha/sword-challenge
cd sword-challenge
```

2. Set up environment variables:
```bash
cp env-example .env
# Edit .env with your configuration:
# - Database settings (DB_USER, DB_PASSWORD, DB_HOST, DB_PORT, DB_NAME)
# - JWT secret (JWT_SECRET)
# - RabbitMQ URL (RABBITMQ_URL)
```

3. Start the development environment:
```bash
docker-compose up -d app-dev
```

4. Get into container development:
```bash
make dev
```

5. Run the application inside container development:
```bash
make air
```

The API will be available at `http://localhost:3000`.

## API Endpoints

### Tasks

- `POST /api/tasks` - Create a new task (Technician only)
  - Required fields: title, summary, performed_at
  - Title max length: 255 characters
  - Summary max length: 2500 characters
  - Performed_at must be between 1900-01-01 and 2100-12-31

- `GET /api/tasks` - List tasks (Technicians see their own, Managers see all)
- `GET /api/tasks/:id` - Get task details
- `PUT /api/tasks/:id` - Update task (Technician can update own tasks)
- `DELETE /api/tasks/:id` - Delete task (Manager only)

### Notifications

- `GET /api/notifications` - Get unread notifications (Manager only)
- `PUT /api/notifications/:id/read` - Mark notification as read (Manager only)

## Authentication

The API uses JWT (JSON Web Token) authentication:
- Include the JWT token in the Authorization header: `Bearer <token>`
- The token must contain user_id and roles claims
- JWT secret must be at least 32 characters long

## Database Schema

### Users
- id (BIGINT, PRIMARY KEY)
- name (VARCHAR)
- email (VARCHAR, UNIQUE)
- password_hash (VARCHAR)
- role (ENUM: 'manager', 'technician')
- created_at (TIMESTAMP)
- updated_at (TIMESTAMP)

### Tasks
- id (BIGINT, PRIMARY KEY)
- technician_id (BIGINT, FOREIGN KEY)
- title (VARCHAR(255))
- summary (TEXT)
- performed_at (TIMESTAMP)
- created_at (TIMESTAMP)
- updated_at (TIMESTAMP)

### Notifications
- id (BIGINT, PRIMARY KEY)
- task_id (BIGINT, FOREIGN KEY)
- message (TEXT)
- is_read (BOOLEAN)
- created_at (TIMESTAMP)

## Testing

Run the test suite:
```bash
go test ./...
```

View test coverage:
```bash
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## Project Structure

```
.
├── cmd/
│   └── server/
│       └── server.go
├── internal/
│   ├── controllers/
│   ├── middleware/
│   ├── models/
│   ├── repository/
│   │   └── mysql/
│   └── service/
├── pkg/
│   └── messaging/
├── test/
├── docker-compose.yml
├── Dockerfile
├── go.mod
├── go.sum
└── README.md
```

## Security Considerations

- JWT-based authentication
- Role-based access control
- Input validation and sanitization
- SQL injection prevention through prepared statements
- XSS protection through proper content type headers and HTML escaping
- CORS configuration for API access
- Secure error handling without exposing sensitive information
- Environment variable configuration
- Password hashing for user authentication

## Development Tools

- GoDoc server available at `http://localhost:6464`
- Test coverage server available at `http://localhost:6767`

## Future Improvements

- Add rate limiting
- Add request logging
- Add more comprehensive test coverage
- Implement user management endpoints
- Add pagination for task and notification lists
- Add filtering and sorting options
- Implement WebSocket for real-time notifications
