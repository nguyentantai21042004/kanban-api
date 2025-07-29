# Kanban API - Personal Project Management System

> **From Tan Tai API V1 With Love** ❤️

## 📋 Overview

**Kanban API** là một hệ thống quản lý dự án cá nhân được thiết kế cho mục đích internal và nhóm dự án nhỏ. Dự án này được phát triển với kiến trúc Clean Architecture, sử dụng Go 1.23.8 và các công nghệ hiện đại để tạo ra một API mạnh mẽ, dễ bảo trì và mở rộng.

## 🎯 Mục đích

- **Internal Use**: Phục vụ cho nhu cầu quản lý dự án nội bộ
- **Small Team**: Tối ưu cho nhóm dự án nhỏ (2-10 người)
- **Personal Project**: Hỗ trợ quản lý dự án cá nhân hiệu quả
- **Learning Purpose**: Thực hành và học tập các công nghệ mới

## 🏗️ Kiến trúc

### Clean Architecture Pattern
```
├── cmd/                    # Application Entry Points
│   ├── api/               # HTTP API Server
│   └── consumer/          # Message Queue Consumer
├── internal/              # Private Application Code
│   ├── auth/             # Authentication & Authorization
│   ├── boards/           # Board Management
│   ├── cards/            # Card Management
│   ├── lists/            # List Management
│   ├── labels/           # Label Management
│   ├── user/             # User Management
│   ├── role/             # Role Management
│   ├── upload/           # File Upload Management
│   ├── websocket/        # Real-time Communication
│   ├── httpserver/       # HTTP Server Configuration
│   ├── middleware/       # HTTP Middleware
│   ├── models/           # Domain Models
│   ├── dbmodels/         # Database Models (SQLBoiler)
│   └── appconfig/        # Application Configuration
├── pkg/                  # Public Libraries
│   ├── log/             # Logging
│   ├── response/        # HTTP Response Helpers
│   ├── errors/          # Error Handling
│   ├── encrypter/       # Encryption Utilities
│   ├── minio/           # MinIO Client
│   ├── postgres/        # PostgreSQL Utilities
│   ├── rabbitmq/        # RabbitMQ Client
│   ├── discord/         # Discord Webhook
│   ├── websocket/       # WebSocket Utilities
│   └── util/            # General Utilities
├── config/              # Configuration Management
├── migrations/          # Database Migrations
└── docs/               # API Documentation
```

## 🚀 Tính năng chính

### 📊 Kanban Board Management
- **Boards**: Tạo và quản lý bảng kanban
- **Lists**: Quản lý các cột trong bảng (To Do, In Progress, Done)
- **Cards**: Quản lý công việc với metadata phong phú
- **Labels**: Phân loại và tag công việc
- **Real-time Updates**: WebSocket cho cập nhật real-time

### 👥 User Management
- **Authentication**: JWT-based authentication
- **Authorization**: Role-based access control
- **User Profiles**: Quản lý thông tin người dùng
- **Team Collaboration**: Hỗ trợ làm việc nhóm

### 📁 File Management
- **File Upload**: Tải lên file đính kèm
- **MinIO Integration**: Object storage cho file
- **Image Processing**: Xử lý hình ảnh

### 🔔 Notifications
- **Discord Integration**: Webhook notifications
- **Real-time Alerts**: Thông báo real-time
- **Email Notifications**: Email alerts (planned)

### 🌐 API Features
- **RESTful API**: REST API đầy đủ
- **Swagger Documentation**: Auto-generated API docs
- **Health Checks**: `/health`, `/ready`, `/live` endpoints
- **Internationalization**: Multi-language support
- **Error Handling**: Comprehensive error management

## 🛠️ Technology Stack

### Backend
- **Language**: Go 1.23.8
- **Framework**: Gin (HTTP framework)
- **Database**: PostgreSQL
- **ORM**: SQLBoiler (code generation)
- **Message Queue**: RabbitMQ
- **Cache**: Redis
- **Object Storage**: MinIO
- **Documentation**: Swagger/OpenAPI

### DevOps & Deployment
- **Containerization**: Docker
- **Orchestration**: Kubernetes
- **CI/CD**: Jenkins
- **Registry**: Harbor
- **Monitoring**: Discord Webhooks

### Development Tools
- **Code Generation**: SQLBoiler, Swag
- **Logging**: Zap Logger
- **Validation**: Environment-based config
- **Testing**: Go testing framework

## 📦 Installation & Setup

### Prerequisites
- Go 1.23.8+
- PostgreSQL 12+
- Redis 6+
- MinIO
- RabbitMQ (optional)

### Quick Start

1. **Clone Repository**
```bash
git clone https://gitlab.com/tantai-kanban/kanban-api.git
cd kanban-api
```

2. **Install Dependencies**
```bash
go mod download
go mod vendor
```

3. **Setup Environment**
```bash
cp env.template .env
# Edit .env with your configuration
```

4. **Generate Code**
```bash
# Generate database models
make models

# Generate Swagger docs
make swagger
```

5. **Run Database Migrations**
```bash
# Apply migrations to PostgreSQL
psql -h localhost -U postgres -d kanban_db -f migrations/01_init_user.sql
psql -h localhost -U postgres -d kanban_db -f migrations/02_init_role.sql
psql -h localhost -U postgres -d kanban_db -f migrations/03_kanban_init.sql
psql -h localhost -U postgres -d kanban_db -f migrations/04_init_data.sql
psql -h localhost -U postgres -d kanban_db -f migrations/05_update_model.sql
psql -h localhost -U postgres -d kanban_db -f migrations/06_upload_model.sql
```

6. **Run Application**
```bash
# Run API server
make run-api

# Run consumer (optional)
make run-consumer
```

### Docker Deployment

```bash
# Build and run with Docker Compose
make build-docker-compose
```

### Kubernetes Deployment

```bash
# Apply Kubernetes manifests
kubectl apply -f deployment.yaml
```

## 🔧 Configuration

### Environment Variables

```bash
# Server Configuration
HOST=0.0.0.0
APP_PORT=8080
API_MODE=debug

# Database Configuration
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_USER=postgres
POSTGRES_PASSWORD=your_password
POSTGRES_DB=kanban_db

# Storage Configuration
MINIO_ENDPOINT=localhost:9000
MINIO_ACCESS_KEY=minioadmin
MINIO_SECRET_KEY=minioadmin
MINIO_USE_SSL=false
MINIO_REGION=us-east-1
MINIO_BUCKET=kanban-files

# Security Configuration
JWT_SECRET=your_jwt_secret
ENCRYPT_KEY=your_encryption_key
INTERNAL_KEY=your_internal_key

# Monitoring Configuration
DISCORD_REPORT_BUG_ID=your_discord_webhook_id
DISCORD_REPORT_BUG_TOKEN=your_discord_webhook_token
```

## 📚 API Documentation

### Health Check Endpoints
- `GET /health` - Basic health check
- `GET /ready` - Readiness check (with DB connectivity)
- `GET /live` - Liveness check

### Main API Endpoints
- `GET /swagger/*` - Swagger UI documentation
- `GET /api/v1/boards` - Board management
- `GET /api/v1/cards` - Card management
- `GET /api/v1/lists` - List management
- `GET /api/v1/labels` - Label management
- `GET /api/v1/users` - User management
- `GET /api/v1/auth` - Authentication
- `GET /api/v1/uploads` - File upload
- `GET /api/v1/websocket/ws/{board_id}` - WebSocket connection

## 🧪 Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific test
go test ./internal/boards/...
```

## 📊 Monitoring & Logging

### Health Checks
```bash
# Test health endpoints
curl http://localhost:8080/health
curl http://localhost:8080/ready
curl http://localhost:8080/live
```

### Logs
- **Structured Logging**: JSON format in production
- **Log Levels**: Debug, Info, Warn, Error
- **Log Rotation**: Automatic log management

## 🔒 Security Features

- **JWT Authentication**: Secure token-based auth
- **Role-based Access Control**: Fine-grained permissions
- **Data Encryption**: Sensitive data encryption
- **Input Validation**: Comprehensive validation
- **CORS Support**: Cross-origin resource sharing
- **Rate Limiting**: API rate limiting (planned)

## 🚀 Deployment

### Docker
```bash
# Build image
docker build -f cmd/api/Dockerfile -t kanban-api .

# Run container
docker run -p 8080:8080 kanban-api
```

### Kubernetes
```bash
# Apply deployment
kubectl apply -f deployment.yaml

# Check status
kubectl get pods -n kanban-api
kubectl logs deployment/kanban-api -n kanban-api
```

### Jenkins CI/CD
- **Automated Build**: Docker image building
- **Automated Testing**: Unit and integration tests
- **Automated Deployment**: Kubernetes deployment
- **Discord Notifications**: Build status notifications

## 🤝 Contributing

### Development Guidelines
1. **Code Style**: Follow Go conventions
2. **Testing**: Write tests for new features
3. **Documentation**: Update docs for changes
4. **Git Flow**: Feature branches with PRs

### Code Generation
```bash
# Generate models after DB changes
make models

# Generate Swagger docs
make swagger
```

## 📈 Performance

### Optimizations
- **Database Indexing**: Optimized queries
- **Connection Pooling**: Efficient DB connections
- **Caching**: Redis for frequently accessed data
- **Async Processing**: Background job processing

### Scalability
- **Horizontal Scaling**: Kubernetes deployment
- **Load Balancing**: Multiple pod instances
- **Database Scaling**: Read replicas support
- **Microservices Ready**: Modular architecture

## 🔮 Roadmap

### Planned Features
- [ ] **Email Notifications**: SMTP integration
- [ ] **Advanced Analytics**: Project metrics
- [ ] **Mobile API**: Mobile-optimized endpoints
- [ ] **API Rate Limiting**: Request throttling
- [ ] **Advanced Search**: Full-text search
- [ ] **Export/Import**: Data portability

### Technical Improvements
- [ ] **GraphQL API**: Alternative to REST
- [ ] **Event Sourcing**: Audit trail
- [ ] **CQRS Pattern**: Command Query Separation
- [ ] **Distributed Tracing**: OpenTelemetry
- [ ] **Metrics Collection**: Prometheus integration

## 📄 License

This project is developed for personal and internal use. All rights reserved.

## 👨‍💻 Author

**Nguyen Tan Tai** - Personal Kanban API Project

---

> **"From Tan Tai API V1 With Love"** - A personal project built with passion and care for efficient project management. 🚀
