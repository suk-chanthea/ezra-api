# Ezra API - Clean Architecture

A RESTful API built with Go (Golang) following Clean Architecture principles for managing music, events, bookings, bands, and notifications with Firebase Cloud Messaging (FCM) integration.

## 🏗️ Clean Architecture Structure

This project follows **Clean Architecture** (also known as Hexagonal Architecture or Ports & Adapters) with clear separation of concerns:

```
ezra-api/
├── cmd/                                    # Application Entry Point
│   └── main.go                            # Dependency injection & initialization
│
├── domain/                                 # Enterprise Business Layer (Innermost)
│   ├── entity/                            # Business Entities
│   │   ├── user.go                        # User domain model
│   │   ├── music.go                       # Music domain model
│   │   ├── event.go                       # Event domain model
│   │   ├── booking.go                     # Booking domain model
│   │   ├── band.go                        # Band domain model
│   │   ├── favorite.go                    # Favorite domain model
│   │   ├── setting.go                     # Setting domain model
│   │   ├── notification.go                # Notification domain model
│   │   └── device_token.go                # Device token domain model
│   │
│   ├── repository/                        # Repository Interfaces (Ports)
│   │   ├── user_repository.go             # User data access interface
│   │   ├── music_repository.go            # Music data access interface
│   │   ├── event_repository.go            # Event data access interface
│   │   ├── booking_repository.go          # Booking data access interface
│   │   ├── band_repository.go             # Band data access interface
│   │   ├── favorite_repository.go         # Favorite data access interface
│   │   ├── setting_repository.go          # Setting data access interface
│   │   ├── notification_repository.go     # Notification data access interface
│   │   └── device_token_repository.go     # Device token data access interface
│   │
│   └── dto/                               # Data Transfer Objects
│       ├── request.go                     # Request DTOs (input models)
│       └── response.go                    # Response DTOs (output models)
│
├── usecase/                                # Application Business Layer
│   ├── auth_usecase.go                    # Authentication business logic
│   ├── music_usecase.go                   # Music business logic
│   ├── event_usecase.go                   # Event business logic
│   ├── booking_usecase.go                 # Booking business logic
│   ├── band_usecase.go                    # Band business logic
│   ├── favorite_usecase.go                # Favorite business logic
│   ├── setting_usecase.go                 # Setting business logic
│   └── notification_usecase.go            # Notification business logic (with FCM)
│
├── infrastructure/                         # Infrastructure Layer (Outermost)
│   ├── persistence/                       # Database Implementations (Adapters)
│   │   ├── user_repository_impl.go        # User repository with GORM
│   │   ├── music_repository_impl.go       # Music repository with GORM
│   │   ├── event_repository_impl.go       # Event repository with GORM
│   │   ├── booking_repository_impl.go     # Booking repository with GORM
│   │   ├── band_repository_impl.go        # Band repository with GORM
│   │   ├── favorite_repository_impl.go    # Favorite repository with GORM
│   │   ├── setting_repository_impl.go     # Setting repository with GORM
│   │   ├── notification_repository_impl.go # Notification repository with GORM
│   │   └── device_token_repository_impl.go # Device token repository with GORM
│   │
│   └── firebase/                          # Firebase Services
│       └── fcm_service.go                 # Firebase Cloud Messaging service
│
├── interface/                              # Interface Adapters Layer
│   └── http/                              # HTTP Delivery Mechanism
│       ├── handler/                       # HTTP Handlers (Controllers)
│       │   ├── auth_handler.go            # Auth endpoints
│       │   ├── music_handler.go           # Music endpoints
│       │   ├── event_handler.go           # Event endpoints
│       │   ├── booking_handler.go         # Booking endpoints
│       │   ├── band_handler.go            # Band endpoints
│       │   ├── favorite_handler.go        # Favorite endpoints
│       │   ├── setting_handler.go         # Setting endpoints
│       │   ├── notification_handler.go    # Notification endpoints
│       │   └── device_token_handler.go    # Device token endpoints
│       │
│       ├── middleware/                    # HTTP Middleware
│       │   └── auth_middleware.go         # JWT authentication
│       │
│       └── router/                        # Route Configuration
│           └── router.go                  # Gin router setup
│
├── migrate/                                # Database Migrations
│   ├── 000000_postgres.up.sql            # Schema creation
│   └── 000000_postgres.down.sql          # Schema rollback
│
├── docker-compose.yml                      # Docker services configuration
├── Dockerfile                              # Application container
├── go.mod                                  # Go dependencies
├── go.sum                                  # Go dependency checksums
├── ROUTEMAP.md                            # Complete API documentation
└── README.md                              # This file
```

## 🎯 Clean Architecture Principles

### 1. **Dependency Rule**
Dependencies flow **inward** only. Inner layers know nothing about outer layers:

```
Infrastructure ──▶ Interface ──▶ UseCase ──▶ Domain
   (Outermost)                              (Innermost)
```

### 2. **Layer Responsibilities**

#### **Domain Layer** (Innermost - No dependencies)
- **Entities**: Pure business objects with business rules
- **Repository Interfaces**: Define contracts for data access
- **DTOs**: Data transfer between layers

**Key Principle**: This layer has **zero external dependencies**. It's pure Go code with business logic only.

```go
// Example: domain/entity/music.go
type Music struct {
    ID          uint
    Title       string
    Artist      string
    // ... business logic methods
}

func (m *Music) IsValid() bool {
    return m.Title != "" && m.Artist != ""
}
```

#### **UseCase Layer** (Application Business Logic)
- Orchestrates the flow of data to and from entities
- Implements application-specific business rules
- Depends only on Domain layer interfaces

**Key Principle**: Technology-agnostic. Doesn't know about HTTP, databases, or frameworks.

```go
// Example: usecase/music_usecase.go
type MusicUseCase interface {
    CreateMusic(req *dto.CreateMusicRequest, userID uint) error
    GetAllMusic() ([]*dto.MusicResponse, error)
}

type musicUseCase struct {
    musicRepo repository.MusicRepository  // Interface, not implementation
}
```

#### **Interface Layer** (Interface Adapters)
- Converts data between use cases and external world
- HTTP handlers, presenters, controllers
- Depends on UseCase layer

**Key Principle**: Handles all input/output conversions and HTTP concerns.

```go
// Example: interface/http/handler/music_handler.go
type MusicHandler struct {
    musicUseCase usecase.MusicUseCase  // Uses interface
}

func (h *MusicHandler) Create(c *gin.Context) {
    var req dto.CreateMusicRequest
    c.ShouldBindJSON(&req)
    err := h.musicUseCase.CreateMusic(&req, userID)
    // ... handle response
}
```

#### **Infrastructure Layer** (Outermost - Frameworks & Drivers)
- Database implementations (GORM)
- External services (Firebase)
- Framework details

**Key Principle**: All technical implementations. Easily replaceable without affecting business logic.

```go
// Example: infrastructure/persistence/music_repository_impl.go
type musicRepositoryImpl struct {
    db *gorm.DB  // GORM specific
}

func (r *musicRepositoryImpl) Save(music *entity.Music) error {
    // GORM implementation
    return r.db.Create(music).Error
}
```

## 🔄 Data Flow Example

**Creating a Music Track:**

```
1. HTTP Request (JSON)
   ↓
2. Handler (interface/http/handler/music_handler.go)
   - Validates request
   - Extracts user ID from JWT
   ↓
3. UseCase (usecase/music_usecase.go)
   - Business logic validation
   - Orchestrates operation
   ↓
4. Repository Interface (domain/repository/music_repository.go)
   - Abstract data access
   ↓
5. Repository Implementation (infrastructure/persistence/music_repository_impl.go)
   - GORM database operations
   ↓
6. Database (PostgreSQL)
```

**Response flows back up the same path.**

## ✅ Benefits of This Architecture

### 1. **Testability**
```go
// Easy to test use cases with mock repositories
mockRepo := &MockMusicRepository{}
useCase := usecase.NewMusicUseCase(mockRepo)
```

### 2. **Independence**
- Business logic independent of frameworks
- Database can be swapped (GORM → SQL → MongoDB)
- HTTP framework can be changed (Gin → Echo → Chi)

### 3. **Maintainability**
- Clear separation of concerns
- Easy to locate and fix bugs
- Changes in one layer don't affect others

### 4. **Scalability**
- Easy to add new features
- Can split into microservices along use case boundaries

## 🚀 Getting Started

### Prerequisites
- Docker & Docker Compose
- Go 1.21+ (for local development)

### Quick Start

1. **Clone and Setup:**
```bash
git clone <repository>
cd ezra-api
```

2. **Configure Environment:**
```bash
# Copy and edit configuration
cp .env.example .env

# Edit .env with your settings:
# - Database connection
# - JWT secret
# - Google OAuth credentials (optional)
# - Firebase credentials path (optional)
```

3. **Run with Docker:**
```bash
docker-compose up -d
```

The API will be available at `http://localhost:8080`

### Firebase Cloud Messaging Setup (Optional)

1. Create a Firebase project at [Firebase Console](https://console.firebase.google.com/)
2. Download service account key (JSON)
3. Set environment variable:
```bash
export FIREBASE_CREDENTIALS_PATH=./firebase-adminsdk.json
```

Without Firebase credentials, the API runs with a dummy FCM service (no actual push notifications).

## 🛠️ Development

### Local Development (without Docker)

```bash
# Install dependencies
go mod download

# Run database (PostgreSQL) separately
docker-compose up -d postgres

# Run application
go run cmd/main.go
```

### Building

```bash
# Build binary
go build -o ezra-api ./cmd/main.go

# Run
./ezra-api
```

### Database Migrations

Migrations are automatically applied when the database container starts for the first time.

**Manual migration:**
```bash
# Apply migration
docker exec -i ezra-postgres psql -U postgres -d ezradb < migrate/000000_postgres.up.sql

# Rollback
docker exec -i ezra-postgres psql -U postgres -d ezradb < migrate/000000_postgres.down.sql
```

## 📚 API Documentation

Complete API documentation with all endpoints, request/response examples, and integration guides is available in:

📖 **[ROUTEMAP.md](./ROUTEMAP.md)**

### Key Features:
- User authentication (JWT + Google OAuth)
- Music management with metadata
- Event scheduling and management
- Booking system
- Band/team management
- Favorites system
- User settings/preferences
- Notifications (in-app + push)
- Firebase Cloud Messaging integration

## 🔧 Technology Stack

### Core
- **Go 1.21+** - Programming language
- **Gin** - HTTP web framework
- **GORM** - ORM for database operations
- **PostgreSQL 16** - Primary database

### External Services
- **Firebase Cloud Messaging** - Push notifications
- **Google OAuth** - Social authentication

### DevOps
- **Docker** - Containerization
- **Docker Compose** - Multi-container orchestration

## 🏛️ Architecture Advantages

### Why Clean Architecture?

1. **Framework Independence**
   - Business rules don't depend on Gin, GORM, or any library
   - Can change frameworks without touching business logic

2. **Testable**
   - Business rules tested without UI, database, or server
   - Mock repositories for unit testing use cases

3. **UI Independence**
   - Could easily add gRPC, GraphQL, or CLI alongside HTTP

4. **Database Independence**
   - Switch from PostgreSQL to MySQL, MongoDB, etc.
   - Business logic remains unchanged

5. **External Agency Independence**
   - Business rules don't know about Firebase, OAuth, etc.
   - Easy to add/remove external services

## 📦 Dependency Injection

All dependencies are injected in `cmd/main.go`:

```go
// 1. Initialize Repositories (Infrastructure)
userRepo := persistence.NewUserRepository(db)
deviceTokenRepo := persistence.NewDeviceTokenRepository(db)

// 2. Initialize Services (Infrastructure)
fcmService := firebase.NewFCMService(credPath, deviceTokenRepo)

// 3. Initialize Use Cases (Application)
authUseCase := usecase.NewAuthUseCase(userRepo, secretKey)
notificationUseCase := usecase.NewNotificationUseCase(notifRepo, fcmService)

// 4. Initialize Handlers (Interface)
authHandler := handler.NewAuthHandler(authUseCase)

// 5. Setup Router (Interface)
router := router.NewRouter(authHandler, ...)
```

This makes the entire application's dependency graph visible in one place.

## 🧪 Testing Strategy

```
Domain Layer     → Unit tests (pure business logic)
UseCase Layer    → Unit tests with mocks
Interface Layer  → Integration tests
Infrastructure   → Integration tests with test database
```

## 📈 Future Enhancements

Following clean architecture makes these additions easy:

- [ ] gRPC API (add new interface layer)
- [ ] WebSocket support (add new interface layer)
- [ ] GraphQL API (add new interface layer)
- [ ] Caching layer (add to infrastructure)
- [ ] Message queue (add to infrastructure)
- [ ] Elasticsearch (add new repository implementation)
- [ ] Multi-tenancy support (modify use cases)

## 📄 License

MIT License - feel free to use this project as a reference or starting point for your own clean architecture projects!

## 🤝 Contributing

This project demonstrates clean architecture principles. Feel free to:
- Use it as a template for your projects
- Submit issues or suggestions
- Create pull requests for improvements

---

**Built with ❤️ using Clean Architecture principles**
