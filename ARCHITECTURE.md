# Ezra API - Clean Architecture Documentation

## 🎯 Architecture Overview

This document explains the Clean Architecture implementation in the Ezra API project.

## 📐 Layer Diagram

```
┌─────────────────────────────────────────────────────────────────┐
│                        cmd/main.go                              │
│                  (Dependency Injection)                         │
│                                                                 │
│  Wires together all layers:                                    │
│  Infrastructure → Interface → UseCase → Domain                 │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                    INTERFACE LAYER                              │
│              (Interface Adapters)                               │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  interface/http/handler/         (Controllers)                 │
│  ├── auth_handler.go             HTTP Request/Response         │
│  ├── music_handler.go            JSON parsing                  │
│  ├── event_handler.go            Validation                    │
│  ├── booking_handler.go          Error handling                │
│  ├── band_handler.go                                           │
│  ├── favorite_handler.go                                       │
│  ├── setting_handler.go                                        │
│  ├── notification_handler.go                                   │
│  └── device_token_handler.go                                   │
│                                                                 │
│  interface/http/middleware/      (Middleware)                  │
│  └── auth_middleware.go          JWT verification              │
│                                                                 │
│  interface/http/router/          (Routing)                     │
│  └── router.go                   Route configuration           │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                     USECASE LAYER                               │
│              (Application Business Rules)                       │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  usecase/                        (Orchestration)               │
│  ├── auth_usecase.go             Business workflows            │
│  ├── music_usecase.go            Application logic             │
│  ├── event_usecase.go            Coordination                  │
│  ├── booking_usecase.go          Transaction management        │
│  ├── band_usecase.go                                           │
│  ├── favorite_usecase.go                                       │
│  ├── setting_usecase.go                                        │
│  └── notification_usecase.go                                   │
│                                                                 │
│  Dependencies: Only domain layer interfaces                    │
│  No framework dependencies                                     │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                     DOMAIN LAYER                                │
│              (Enterprise Business Rules)                        │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  domain/entity/                  (Business Objects)            │
│  ├── user.go                     Pure business logic           │
│  ├── music.go                    No external dependencies      │
│  ├── event.go                    Validation rules              │
│  ├── booking.go                  Domain rules                  │
│  ├── band.go                                                   │
│  ├── favorite.go                                               │
│  ├── setting.go                                                │
│  ├── notification.go                                           │
│  └── device_token.go                                           │
│                                                                 │
│  domain/repository/              (Port Interfaces)             │
│  ├── user_repository.go          Data access contracts         │
│  ├── music_repository.go         Technology-agnostic           │
│  ├── event_repository.go                                       │
│  ├── booking_repository.go                                     │
│  ├── band_repository.go                                        │
│  ├── favorite_repository.go                                    │
│  ├── setting_repository.go                                     │
│  ├── notification_repository.go                                │
│  └── device_token_repository.go                                │
│                                                                 │
│  domain/dto/                     (Data Transfer)               │
│  ├── request.go                  Input models                  │
│  └── response.go                 Output models                 │
│                                                                 │
│  Zero external dependencies                                    │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
                              ▲
                              │
                    (Implements interfaces)
                              │
┌─────────────────────────────────────────────────────────────────┐
│                  INFRASTRUCTURE LAYER                           │
│              (Frameworks & Drivers)                             │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  infrastructure/persistence/     (Database Adapters)           │
│  ├── user_repository_impl.go    GORM implementations           │
│  ├── music_repository_impl.go   Database queries              │
│  ├── event_repository_impl.go   ORM operations                │
│  ├── booking_repository_impl.go                               │
│  ├── band_repository_impl.go                                  │
│  ├── favorite_repository_impl.go                              │
│  ├── setting_repository_impl.go                               │
│  ├── notification_repository_impl.go                          │
│  └── device_token_repository_impl.go                          │
│                                                                 │
│  infrastructure/firebase/        (External Services)           │
│  └── fcm_service.go             Firebase integration           │
│                                                                 │
│  Framework-specific implementations                            │
│  Easily replaceable                                            │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                    EXTERNAL SYSTEMS                             │
│                                                                 │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐        │
│  │  PostgreSQL  │  │   Firebase   │  │ Google OAuth │        │
│  │   Database   │  │     FCM      │  │              │        │
│  └──────────────┘  └──────────────┘  └──────────────┘        │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

## 🔄 Request Flow Example

### Creating a Music Track

```
1. Client sends HTTP POST /api/musics
   {
     "title": "Amazing Grace",
     "artist": "John Newton"
   }
                │
                ▼
2. Gin Router (interface/http/router/router.go)
   Routes to MusicHandler.Create()
                │
                ▼
3. Auth Middleware (interface/http/middleware/auth_middleware.go)
   Validates JWT token
   Extracts user_id
                │
                ▼
4. Music Handler (interface/http/handler/music_handler.go)
   • Parses JSON to CreateMusicRequest DTO
   • Validates input format
   • Calls: musicUseCase.CreateMusic(req, userID)
                │
                ▼
5. Music UseCase (usecase/music_usecase.go)
   • Validates business rules
   • Creates entity: entity.NewMusic(...)
   • Validates entity: music.IsValid()
   • Calls: musicRepo.Save(music)
                │
                ▼
6. Music Repository Interface (domain/repository/music_repository.go)
   Defines: Save(music *entity.Music) error
                │
                ▼
7. Music Repository Implementation (infrastructure/persistence/music_repository_impl.go)
   • Converts entity to GORM model
   • Executes: db.Create(model).Error
   • Saves to PostgreSQL
                │
                ▼
8. PostgreSQL Database
   Music record inserted
                │
                ▼
   Response flows back up:
   Repository → UseCase → Handler → Client
                │
                ▼
9. Client receives: 201 Created
   {
     "id": 42,
     "title": "Amazing Grace",
     "artist": "John Newton",
     "created_at": "2025-10-29T..."
   }
```

## 🎭 Dependency Inversion in Action

### The Problem (Without Clean Architecture)

```go
// ❌ BAD: Handler depends on concrete implementation
type MusicHandler struct {
    db *gorm.DB  // Direct database dependency
}

func (h *MusicHandler) Create(c *gin.Context) {
    // Handler knows about database details
    h.db.Create(&music)
}
```

**Problems:**
- Hard to test (need real database)
- Tightly coupled to GORM
- Can't swap database easily
- Business logic mixed with infrastructure

### The Solution (Clean Architecture)

```go
// ✅ GOOD: Handler depends on abstraction

// Domain Layer - Interface (Port)
type MusicRepository interface {
    Save(music *entity.Music) error
}

// UseCase Layer - Uses interface
type musicUseCase struct {
    musicRepo repository.MusicRepository  // Depends on interface
}

// Infrastructure Layer - Implementation (Adapter)
type musicRepositoryImpl struct {
    db *gorm.DB  // GORM specific
}

func (r *musicRepositoryImpl) Save(music *entity.Music) error {
    return r.db.Create(music).Error
}

// main.go - Dependency Injection
repo := persistence.NewMusicRepository(db)  // Concrete implementation
useCase := usecase.NewMusicUseCase(repo)    // Inject dependency
handler := handler.NewMusicHandler(useCase)
```

**Benefits:**
- Easy to test with mocks
- Database can be swapped
- Clear separation of concerns
- Business logic independent of infrastructure

## 📦 Module Dependencies

```
domain/          (No dependencies - pure Go)
   ↑
   │ (depends on interfaces)
   │
usecase/         (Depends on: domain/)
   ↑
   │ (depends on use cases)
   │
interface/       (Depends on: domain/, usecase/)
   ↑
   │ (implements interfaces)
   │
infrastructure/  (Depends on: domain/)
   │
   │ (external dependencies)
   ↓
External frameworks (GORM, Gin, Firebase, etc.)
```

## 🧪 Testing Strategy

### 1. Domain Layer Tests
```go
// Test pure business logic
func TestMusicIsValid(t *testing.T) {
    music := entity.NewMusic("Title", "Artist", userID)
    assert.True(t, music.IsValid())
}
```

### 2. UseCase Layer Tests (with Mocks)
```go
type MockMusicRepository struct{}

func (m *MockMusicRepository) Save(music *entity.Music) error {
    return nil  // Mock implementation
}

func TestCreateMusic(t *testing.T) {
    mockRepo := &MockMusicRepository{}
    useCase := usecase.NewMusicUseCase(mockRepo)
    
    err := useCase.CreateMusic(req, userID)
    assert.NoError(t, err)
}
```

### 3. Integration Tests
```go
func TestMusicAPI(t *testing.T) {
    // Test full stack with test database
    resp := httptest.Post("/api/musics", body)
    assert.Equal(t, 201, resp.StatusCode)
}
```

## 🔌 Adding a New Feature (Example: Comments)

Following clean architecture makes this systematic:

### Step 1: Domain Layer
```go
// domain/entity/comment.go
type Comment struct {
    ID        uint
    UserID    uint
    MusicID   uint
    Content   string
    CreatedAt time.Time
}

// domain/repository/comment_repository.go
type CommentRepository interface {
    Save(comment *entity.Comment) error
    FindByMusicID(musicID uint) ([]*entity.Comment, error)
}
```

### Step 2: UseCase Layer
```go
// usecase/comment_usecase.go
type CommentUseCase interface {
    CreateComment(req *dto.CreateCommentRequest, userID uint) error
    GetCommentsByMusic(musicID uint) ([]*dto.CommentResponse, error)
}

type commentUseCase struct {
    commentRepo repository.CommentRepository
    musicRepo   repository.MusicRepository
}
```

### Step 3: Infrastructure Layer
```go
// infrastructure/persistence/comment_repository_impl.go
type commentRepositoryImpl struct {
    db *gorm.DB
}

func (r *commentRepositoryImpl) Save(comment *entity.Comment) error {
    return r.db.Create(comment).Error
}
```

### Step 4: Interface Layer
```go
// interface/http/handler/comment_handler.go
type CommentHandler struct {
    commentUseCase usecase.CommentUseCase
}

func (h *CommentHandler) Create(c *gin.Context) {
    // Handle HTTP request
}
```

### Step 5: Wire in main.go
```go
// cmd/main.go
commentRepo := persistence.NewCommentRepository(db)
commentUseCase := usecase.NewCommentUseCase(commentRepo, musicRepo)
commentHandler := handler.NewCommentHandler(commentUseCase)
```

**Notice:** Each layer has a clear responsibility. Adding features is systematic and predictable.

## 🚀 Scalability Patterns

### Horizontal Scalability
```
Load Balancer
      │
      ├─── Instance 1 (ezra-api)
      ├─── Instance 2 (ezra-api)
      └─── Instance 3 (ezra-api)
              │
              └─── Shared PostgreSQL
```

Clean architecture makes this easy:
- Stateless application
- Shared database
- Easy to containerize

### Microservices Migration
```
Monolith (Current)
├── Music UseCase      →  Music Service
├── Event UseCase      →  Event Service
├── Notification UseCase → Notification Service
```

Clean architecture makes splitting easy:
- Use cases are already isolated
- Clear boundaries between features
- Shared domain entities via contracts

## 🎓 Key Takeaways

### 1. **Dependency Rule**
Inner layers don't know about outer layers. Dependencies point inward.

### 2. **Testability**
Each layer can be tested independently with appropriate mocks.

### 3. **Flexibility**
Swap implementations without touching business logic:
- Database: PostgreSQL → MySQL → MongoDB
- Framework: Gin → Echo → Chi
- Delivery: HTTP → gRPC → GraphQL

### 4. **Maintainability**
Clear responsibilities make code easy to understand and modify.

### 5. **Business Logic First**
Domain layer is pure business logic, not polluted with infrastructure concerns.

## 📚 References

- [Clean Architecture by Robert C. Martin](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Hexagonal Architecture (Ports and Adapters)](https://alistair.cockburn.us/hexagonal-architecture/)
- [Domain-Driven Design](https://martinfowler.com/bliki/DomainDrivenDesign.html)

---

**This architecture ensures your codebase remains clean, testable, and maintainable as it grows! 🚀**

