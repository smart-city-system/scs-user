# SCS User Service API

A robust user management service built with Go, Echo framework, and PostgreSQL. This service provides comprehensive user authentication, management, and account verification capabilities with full Swagger API documentation.

## 🚀 Features

- **User Management**: Create, retrieve, and manage user accounts
- **Authentication**: JWT-based authentication system
- **Account Verification**: Email-based account verification
- **API Documentation**: Complete Swagger/OpenAPI documentation
- **Database**: PostgreSQL with GORM ORM
- **Message Queue**: Kafka integration for event-driven architecture
- **Logging**: Structured logging with Zap
- **Validation**: Request validation with custom validators
- **Middleware**: CORS, request logging, error handling, and response standardization

## 📋 Prerequisites

Before running this project, make sure you have the following installed:

- **Go 1.23.3+**
- **PostgreSQL 12+**
- **Apache Kafka** (optional, for message queue features)
- **Git**

## 🛠️ Installation & Setup

### 1. Clone the Repository

```bash
git clone <repository-url>
cd scs-user
```

### 2. Install Dependencies

```bash
go mod download
```

### 3. Environment Configuration

Create a `.env` file in the root directory with the following variables:

```env
# Server Configuration
PORT=1326
MODE=development
READ_TIMEOUT=5s
WRITE_TIMEOUT=5s

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=your_db_user
DB_PASSWORD=your_db_password
DB_NAME=smart_city

# Logging
LOG_LEVEL=debug

# Kafka Configuration (optional)
KAFKA_BROKERS=localhost:9093
```

### 4. Database Setup

1. **Create PostgreSQL Database:**
   ```sql
   CREATE DATABASE smart_city;
   ```

2. **Database Migration:**
   The application automatically runs migrations on startup, creating the necessary tables.

### 5. Run the Application

```bash
# Development mode
go run cmd/server/main.go

# Or build and run
go build -o scs-user cmd/server/main.go
./scs-user
```

The server will start on `http://localhost:1326` (or the port specified in your .env file).

## 📚 API Documentation

### Swagger UI

Once the server is running, access the interactive API documentation at:

**🌐 http://localhost:1326/swagger/index.html**

### API Endpoints

#### Authentication (`/api/v1/auth`)
- `POST /auth/login` - User login
- `POST /auth/validate-token` - Validate JWT token

#### User Management (`/api/v1/users`)
- `POST /users` - Create new user (🔒 Protected)
- `GET /users` - Get paginated list of users (🔒 Protected)
- `GET /users/me` - Get current user profile (🔒 Protected)
- `POST /users/verify` - Verify user account

#### Health Check (`/api/v1/health`)
- `GET /health` - Service health status

### Authentication

Protected endpoints require a Bearer token in the Authorization header:

```bash
Authorization: Bearer <your-jwt-token>
```

**Example API Usage:**

1. **Login to get token:**
   ```bash
   curl -X POST http://localhost:1326/api/v1/auth/login \
     -H "Content-Type: application/json" \
     -d '{"email": "user@example.com", "password": "password123"}'
   ```

2. **Use token for protected endpoints:**
   ```bash
   curl -X GET http://localhost:1326/api/v1/users/me \
     -H "Authorization: Bearer <your-token>"
   ```

## 🐳 Docker Deployment

### Build Docker Image

```bash
docker build -t scs-user:latest .
```

### Run with Docker

```bash
docker run -p 1326:1326 \
  -e DB_HOST=your_db_host \
  -e DB_USER=your_db_user \
  -e DB_PASSWORD=your_db_password \
  -e DB_NAME=smart_city \
  scs-user:latest
```

## 🏗️ Project Structure

```
scs-user/
├── cmd/server/          # Application entry point
├── config/              # Configuration management
├── docs/                # Generated Swagger documentation
├── internal/
│   ├── controllers/     # HTTP handlers
│   ├── dto/            # Data Transfer Objects
│   ├── middlewares/    # Custom middleware
│   ├── models/         # Database models
│   ├── repositories/   # Data access layer
│   ├── server/         # Server setup and routing
│   ├── services/       # Business logic
│   └── types/          # Custom types
├── pkg/
│   ├── db/             # Database connection
│   ├── errors/         # Error handling
│   ├── kafka/          # Kafka client
│   ├── logger/         # Logging utilities
│   ├── utils/          # Utility functions
│   └── validation/     # Validation helpers
├── .env                # Environment variables
├── Dockerfile          # Docker configuration
├── go.mod              # Go module dependencies
└── README.md           # This file
```

## 🔧 Development

### Regenerate Swagger Documentation

After modifying API endpoints:

```bash
go run github.com/swaggo/swag/cmd/swag init -g cmd/server/main.go
```

### Run Tests

```bash
go test ./...
```

### Code Formatting

```bash
go fmt ./...
```

## 🌟 Key Technologies

- **Framework**: Echo v4
- **Database**: PostgreSQL with GORM
- **Authentication**: JWT tokens
- **Documentation**: Swagger/OpenAPI
- **Logging**: Zap logger
- **Validation**: Go Playground Validator
- **Message Queue**: Kafka
- **Configuration**: Environment variables with caarlos0/env

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🆘 Troubleshooting

### Common Issues

1. **Database Connection Error**: Ensure PostgreSQL is running and credentials are correct
2. **Port Already in Use**: Change the PORT in your .env file
3. **Kafka Connection Error**: Kafka is optional; the service will log errors but continue running

### Support

For support and questions, please open an issue in the repository.
