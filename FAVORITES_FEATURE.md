# Favorites Feature Documentation

## Overview
This document describes the implementation of the favorites feature, which allows users to mark music tracks as favorites.

## Database Schema

### Favorites Table
```sql
CREATE TABLE IF NOT EXISTS favorites (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    music_id INTEGER NOT NULL REFERENCES musics(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ DEFAULT now(),
    UNIQUE(user_id, music_id)
);
```

**Indexes:**
- `idx_favorites_user_id` - Fast lookups by user
- `idx_favorites_music_id` - Fast lookups by music
- `idx_favorites_created_at` - Ordering by creation time

**Key Features:**
- Foreign keys ensure referential integrity
- Unique constraint prevents duplicate favorites
- Cascade delete removes favorites when user or music is deleted

## API Endpoints

All endpoints require authentication via JWT token.

### 1. Add Music to Favorites
**Endpoint:** `POST /api/favorites/music/:id`

**Description:** Adds a music track to the authenticated user's favorites.

**Parameters:**
- `id` (URL parameter) - The ID of the music track

**Response:**
- **201 Created** - Music added successfully
- **404 Not Found** - Music doesn't exist
- **409 Conflict** - Music already in favorites
- **401 Unauthorized** - User not authenticated

**Example:**
```bash
curl -X POST http://localhost:8080/api/favorites/music/1 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**Success Response:**
```json
{
  "message": "music added to favorites"
}
```

### 2. Remove Music from Favorites
**Endpoint:** `DELETE /api/favorites/music/:id`

**Description:** Removes a music track from the authenticated user's favorites.

**Parameters:**
- `id` (URL parameter) - The ID of the music track

**Response:**
- **200 OK** - Music removed successfully
- **404 Not Found** - Favorite doesn't exist
- **401 Unauthorized** - User not authenticated

**Example:**
```bash
curl -X DELETE http://localhost:8080/api/favorites/music/1 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**Success Response:**
```json
{
  "message": "music removed from favorites"
}
```

### 3. Get User's Favorite Music
**Endpoint:** `GET /api/favorites/`

**Description:** Retrieves all music tracks favorited by the authenticated user.

**Response:**
- **200 OK** - Returns array of music tracks
- **401 Unauthorized** - User not authenticated

**Example:**
```bash
curl -X GET http://localhost:8080/api/favorites/ \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**Success Response:**
```json
[
  {
    "id": 1,
    "title": "Song Title",
    "cover": "cover.jpg",
    "audio": "audio.mp3",
    "user_id": 5,
    "created_at": "2025-10-24T10:00:00Z",
    "updated_at": "2025-10-24T10:00:00Z"
  }
]
```

### 4. Check if Music is Favorited
**Endpoint:** `GET /api/favorites/music/:id/check`

**Description:** Checks if a specific music track is in the authenticated user's favorites.

**Parameters:**
- `id` (URL parameter) - The ID of the music track

**Response:**
- **200 OK** - Returns favorite status
- **401 Unauthorized** - User not authenticated

**Example:**
```bash
curl -X GET http://localhost:8080/api/favorites/music/1/check \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**Success Response:**
```json
{
  "is_favorite": true
}
```

### 5. Get Favorite Count for Music
**Endpoint:** `GET /api/favorites/music/:id/count`

**Description:** Gets the total number of users who favorited a specific music track.

**Parameters:**
- `id` (URL parameter) - The ID of the music track

**Response:**
- **200 OK** - Returns favorite count
- **401 Unauthorized** - User not authenticated

**Example:**
```bash
curl -X GET http://localhost:8080/api/favorites/music/1/count \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**Success Response:**
```json
{
  "count": 42
}
```

## Architecture

The implementation follows the Clean Architecture pattern:

### 1. Entity Layer (`domain/entity/favorite.go`)
- Defines the core `Favorite` business entity
- Contains validation logic

### 2. Repository Interface (`domain/repository/favorite_repository.go`)
- Defines data access methods
- Keeps domain layer independent of infrastructure

### 3. Repository Implementation (`infrastructure/persistence/favorite_repository_impl.go`)
- GORM-based PostgreSQL implementation
- Uses optimized queries with joins
- Handles database interactions

### 4. Use Case Layer (`usecase/favorite_usecase.go`)
- Contains business logic
- Validates operations
- Coordinates between repositories

### 5. Handler Layer (`interface/http/handler/favorite_handler.go`)
- HTTP request/response handling
- Input validation
- Error mapping to HTTP status codes

### 6. Router Configuration (`interface/http/router/router.go`)
- Defines API routes
- Applies authentication middleware

## Database Migration

The favorites table is included in the main migration file:

**Up Migration:** `migrate/000000_postgres.up.sql`
**Down Migration:** `migrate/000000_postgres.down.sql`

To apply migrations, use your migration tool (e.g., golang-migrate):

```bash
# Apply migrations
migrate -path migrate -database "postgres://user:pass@localhost:5432/ezradb?sslmode=disable" up

# Rollback migrations
migrate -path migrate -database "postgres://user:pass@localhost:5432/ezradb?sslmode=disable" down
```

## Testing

You can test the favorites feature using curl or any API client:

```bash
# 1. Login to get JWT token
TOKEN=$(curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"username":"your_username","password":"your_password"}' \
  | jq -r '.token')

# 2. Add a favorite
curl -X POST http://localhost:8080/api/favorites/music/1 \
  -H "Authorization: Bearer $TOKEN"

# 3. Get your favorites
curl -X GET http://localhost:8080/api/favorites/ \
  -H "Authorization: Bearer $TOKEN"

# 4. Check if music is favorited
curl -X GET http://localhost:8080/api/favorites/music/1/check \
  -H "Authorization: Bearer $TOKEN"

# 5. Get favorite count
curl -X GET http://localhost:8080/api/favorites/music/1/count \
  -H "Authorization: Bearer $TOKEN"

# 6. Remove favorite
curl -X DELETE http://localhost:8080/api/favorites/music/1 \
  -H "Authorization: Bearer $TOKEN"
```

## Error Handling

The API returns appropriate HTTP status codes and error messages:

- **400 Bad Request** - Invalid input (e.g., invalid music ID)
- **401 Unauthorized** - Missing or invalid JWT token
- **404 Not Found** - Resource doesn't exist
- **409 Conflict** - Duplicate favorite
- **500 Internal Server Error** - Server error

Error response format:
```json
{
  "error": "error message here"
}
```

## Performance Considerations

1. **Indexes:** All foreign keys and frequently queried columns are indexed
2. **Unique Constraint:** Prevents duplicate favorites at database level
3. **Cascade Deletes:** Automatic cleanup when users or music are deleted
4. **Optimized Queries:** Uses JOIN for efficient favorite music retrieval
5. **Context Support:** All methods support context for timeout/cancellation

## Future Enhancements

Possible improvements:
1. Add pagination for user favorites
2. Add sorting options (by date, title, etc.)
3. Add favorite statistics/analytics
4. Add notifications when music is favorited
5. Add favorite playlists
6. Export favorites functionality

