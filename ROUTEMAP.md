# Ezra API - Route Documentation

Complete API route reference for the Ezra music management system.

## Base URL
```
http://localhost:8080
```

## Table of Contents
- [Authentication](#authentication)
- [User Management](#user-management)
- [Music](#music)
- [Events](#events)
- [Bookings](#bookings)
- [Favorites](#favorites)
- [Bands](#bands)
- [Settings](#settings)

---

## Authentication

### Register
Create a new user account.

**Endpoint:** `POST /register`  
**Authentication:** None (Public)

**Request Body:**
```json
{
  "username": "johndoe",
  "fullname": "John Doe",
  "email": "john@example.com",
  "password": "password123"
}
```

**Success Response (201):**
```json
{
  "message": "user registered successfully",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Error Response (400):**
```json
{
  "error": "username already exists"
}
```

---

### Login
Authenticate and get an access token.

**Endpoint:** `POST /login`  
**Authentication:** None (Public)

**Request Body:**
```json
{
  "username": "johndoe",
  "password": "password123"
}
```

**Success Response (200):**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Error Response (401):**
```json
{
  "error": "invalid credentials"
}
```

---

### Google Login
Authenticate using Google OAuth.

**Endpoint:** `POST /auth/google`  
**Authentication:** None (Public)

**Request Body:**
```json
{
  "id_token": "google_id_token_here"
}
```

**Success Response (200):**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

---

## User Management

### Logout
Invalidate the current user's token.

**Endpoint:** `POST /api/logout`  
**Authentication:** Required (Bearer Token)

**Headers:**
```
Authorization: Bearer <token>
```

**Success Response (200):**
```json
{
  "message": "logged out successfully"
}
```

---

### Delete User
Delete the current user's account.

**Endpoint:** `DELETE /api/user`  
**Authentication:** Required (Bearer Token)

**Headers:**
```
Authorization: Bearer <token>
```

**Success Response (200):**
```json
{
  "message": "user deleted successfully"
}
```

---

## Music

### Create Music
Add a new music track.

**Endpoint:** `POST /api/musics`  
**Authentication:** Required (Bearer Token)

**Request Body:**
```json
{
  "title": "Amazing Grace",
  "cover": "https://example.com/cover.jpg",
  "audio": "https://example.com/audio.mp3"
}
```

**Success Response (201):**
```json
{
  "id": 1,
  "title": "Amazing Grace",
  "cover": "https://example.com/cover.jpg",
  "audio": "https://example.com/audio.mp3",
  "user_id": 5,
  "created_at": "2025-10-25T10:30:00Z",
  "updated_at": "2025-10-25T10:30:00Z"
}
```

---

### Get All Music
Retrieve all music tracks with pagination.

**Endpoint:** `GET /api/musics`  
**Authentication:** Required (Bearer Token)

**Query Parameters:**
- `page` (optional, default: 1) - Page number
- `page_size` (optional, default: 20, max: 100) - Items per page

**Example:** `GET /api/musics?page=1&page_size=20`

**Success Response (200):**
```json
{
  "data": [
    {
      "id": 1,
      "title": "Amazing Grace",
      "cover": "https://example.com/cover.jpg",
      "audio": "https://example.com/audio.mp3",
      "user_id": 5,
      "created_at": "2025-10-25T10:30:00Z",
      "updated_at": "2025-10-25T10:30:00Z"
    }
  ],
  "pagination": {
    "current_page": 1,
    "page_size": 20,
    "total_pages": 5,
    "total_records": 95,
    "has_next_page": true,
    "has_prev_page": false
  }
}
```

---

### Get User's Music
Retrieve music created by the current user.

**Endpoint:** `GET /api/musics/user`  
**Authentication:** Required (Bearer Token)

**Query Parameters:**
- `page` (optional, default: 1)
- `page_size` (optional, default: 20)

**Success Response (200):** Same as "Get All Music"

---

### Get Music by ID
Retrieve a specific music track.

**Endpoint:** `GET /api/musics/:id`  
**Authentication:** Required (Bearer Token)

**Success Response (200):**
```json
{
  "id": 1,
  "title": "Amazing Grace",
  "cover": "https://example.com/cover.jpg",
  "audio": "https://example.com/audio.mp3",
  "user_id": 5,
  "created_at": "2025-10-25T10:30:00Z",
  "updated_at": "2025-10-25T10:30:00Z"
}
```

**Error Response (404):**
```json
{
  "error": "music not found"
}
```

---

### Update Music
Update an existing music track.

**Endpoint:** `PUT /api/musics/:id`  
**Authentication:** Required (Bearer Token, must be owner)

**Request Body:**
```json
{
  "title": "Amazing Grace (Updated)",
  "cover": "https://example.com/new-cover.jpg",
  "audio": "https://example.com/new-audio.mp3"
}
```

**Success Response (200):**
```json
{
  "id": 1,
  "title": "Amazing Grace (Updated)",
  "cover": "https://example.com/new-cover.jpg",
  "audio": "https://example.com/new-audio.mp3",
  "user_id": 5,
  "created_at": "2025-10-25T10:30:00Z",
  "updated_at": "2025-10-25T11:00:00Z"
}
```

---

### Delete Music
Delete a music track.

**Endpoint:** `DELETE /api/musics/:id`  
**Authentication:** Required (Bearer Token, must be owner)

**Success Response (200):**
```json
{
  "message": "music deleted successfully"
}
```

---

## Events

### Create Event
Create a new event.

**Endpoint:** `POST /api/events`  
**Authentication:** Required (Bearer Token)

**Request Body:**
```json
{
  "title": "Sunday Service",
  "content": "Weekly Sunday worship service",
  "cover": "https://example.com/event-cover.jpg",
  "location": "Main Church Auditorium",
  "start_time": "2025-11-01T09:00:00Z",
  "end_time": "2025-11-01T11:00:00Z",
  "music_ids": [1, 2, 3]
}
```

**Success Response (201):**
```json
{
  "id": 1,
  "title": "Sunday Service",
  "content": "Weekly Sunday worship service",
  "cover": "https://example.com/event-cover.jpg",
  "location": "Main Church Auditorium",
  "start_time": "2025-11-01T09:00:00Z",
  "end_time": "2025-11-01T11:00:00Z",
  "user_id": 5,
  "musics": [
    {
      "id": 1,
      "title": "Amazing Grace",
      "cover": "...",
      "audio": "..."
    }
  ],
  "created_at": "2025-10-25T10:30:00Z",
  "updated_at": "2025-10-25T10:30:00Z"
}
```

---

### Get All Events
Retrieve all events with pagination.

**Endpoint:** `GET /api/events`  
**Authentication:** Required (Bearer Token)

**Query Parameters:**
- `page` (optional, default: 1)
- `page_size` (optional, default: 20)

**Success Response (200):** Returns paginated event list

---

### Get User's Events
Retrieve events created by the current user.

**Endpoint:** `GET /api/events/user`  
**Authentication:** Required (Bearer Token)

**Query Parameters:**
- `page` (optional, default: 1)
- `page_size` (optional, default: 20)

---

### Get Event by ID
Retrieve a specific event.

**Endpoint:** `GET /api/events/:id`  
**Authentication:** Required (Bearer Token)

**Success Response (200):** Returns event details with music list

---

### Update Event
Update an existing event.

**Endpoint:** `PUT /api/events/:id`  
**Authentication:** Required (Bearer Token, must be owner)

**Request Body:** Same as Create Event

**Success Response (200):** Returns updated event

---

### Delete Event
Delete an event.

**Endpoint:** `DELETE /api/events/:id`  
**Authentication:** Required (Bearer Token, must be owner)

**Success Response (200):**
```json
{
  "message": "event deleted successfully"
}
```

---

## Bookings

### Create Booking
Register for an event.

**Endpoint:** `POST /api/bookings`  
**Authentication:** Required (Bearer Token)

**Request Body:**
```json
{
  "event_id": 1,
  "notes": "Looking forward to attending!"
}
```

**Success Response (201):**
```json
{
  "id": 1,
  "event_id": 1,
  "user_id": 5,
  "status": "pending",
  "notes": "Looking forward to attending!",
  "created_at": "2025-10-25T10:30:00Z",
  "updated_at": "2025-10-25T10:30:00Z"
}
```

---

### Get All Bookings
Retrieve all bookings (admin).

**Endpoint:** `GET /api/bookings`  
**Authentication:** Required (Bearer Token)

**Query Parameters:**
- `page` (optional, default: 1)
- `page_size` (optional, default: 20)

---

### Get User's Bookings
Retrieve current user's bookings.

**Endpoint:** `GET /api/bookings/user`  
**Authentication:** Required (Bearer Token)

**Query Parameters:**
- `page` (optional, default: 1)
- `page_size` (optional, default: 20)

---

### Get Bookings by Event
Retrieve all bookings for a specific event.

**Endpoint:** `GET /api/bookings/event/:event_id`  
**Authentication:** Required (Bearer Token)

**Query Parameters:**
- `page` (optional, default: 1)
- `page_size` (optional, default: 20)

---

### Get Booking by ID
Retrieve a specific booking.

**Endpoint:** `GET /api/bookings/:id`  
**Authentication:** Required (Bearer Token)

---

### Update Booking
Update booking status or notes.

**Endpoint:** `PUT /api/bookings/:id`  
**Authentication:** Required (Bearer Token)

**Request Body:**
```json
{
  "status": "confirmed",
  "notes": "Confirmed attendance"
}
```

**Valid Status Values:** `pending`, `confirmed`, `cancelled`

**Success Response (200):** Returns updated booking

---

### Delete Booking
Cancel/delete a booking.

**Endpoint:** `DELETE /api/bookings/:id`  
**Authentication:** Required (Bearer Token, must be owner)

**Success Response (200):**
```json
{
  "message": "booking deleted successfully"
}
```

---

## Favorites

### Get User's Favorites
Retrieve current user's favorite music.

**Endpoint:** `GET /api/favorites`  
**Authentication:** Required (Bearer Token)

**Query Parameters:**
- `page` (optional, default: 1)
- `page_size` (optional, default: 20)

**Success Response (200):**
```json
{
  "data": [
    {
      "id": 1,
      "title": "Amazing Grace",
      "cover": "...",
      "audio": "...",
      "user_id": 3,
      "created_at": "2025-10-25T10:30:00Z",
      "updated_at": "2025-10-25T10:30:00Z"
    }
  ],
  "pagination": { ... }
}
```

---

### Add to Favorites
Add a music track to favorites.

**Endpoint:** `POST /api/favorites/music/:id`  
**Authentication:** Required (Bearer Token)

**Success Response (200):**
```json
{
  "message": "music added to favorites"
}
```

---

### Remove from Favorites
Remove a music track from favorites.

**Endpoint:** `DELETE /api/favorites/music/:id`  
**Authentication:** Required (Bearer Token)

**Success Response (200):**
```json
{
  "message": "music removed from favorites"
}
```

---

### Check if Favorite
Check if a music track is in user's favorites.

**Endpoint:** `GET /api/favorites/music/:id/check`  
**Authentication:** Required (Bearer Token)

**Success Response (200):**
```json
{
  "is_favorite": true
}
```

---

### Get Favorite Count
Get the number of users who favorited a music track.

**Endpoint:** `GET /api/favorites/music/:id/count`  
**Authentication:** Required (Bearer Token)

**Success Response (200):**
```json
{
  "count": 42
}
```

---

## Bands

### Create Band
Create a new band/music collection.

**Endpoint:** `POST /api/bands`  
**Authentication:** Required (Bearer Token)

**Request Body:**
```json
{
  "name": "Worship Team A",
  "description": "Main worship team for Sunday services",
  "cover": "https://example.com/band-cover.jpg",
  "is_public": true
}
```

**Success Response (201):**
```json
{
  "id": 1,
  "name": "Worship Team A",
  "description": "Main worship team for Sunday services",
  "cover": "https://example.com/band-cover.jpg",
  "is_public": true,
  "user_id": 5,
  "created_at": "2025-10-25T10:30:00Z",
  "updated_at": "2025-10-25T10:30:00Z"
}
```

---

### Get All Bands
Retrieve all bands with pagination.

**Endpoint:** `GET /api/bands`  
**Authentication:** Required (Bearer Token)

**Query Parameters:**
- `page` (optional, default: 1)
- `page_size` (optional, default: 20)

---

### Get User's Bands
Retrieve bands created by the current user.

**Endpoint:** `GET /api/bands/user`  
**Authentication:** Required (Bearer Token)

**Query Parameters:**
- `page` (optional, default: 1)
- `page_size` (optional, default: 20)

---

### Get Public Bands
Retrieve all public bands.

**Endpoint:** `GET /api/bands/public`  
**Authentication:** Required (Bearer Token)

**Query Parameters:**
- `page` (optional, default: 1)
- `page_size` (optional, default: 20)

---

### Get Band by ID
Retrieve a specific band.

**Endpoint:** `GET /api/bands/:id`  
**Authentication:** Required (Bearer Token)

---

### Update Band
Update band details.

**Endpoint:** `PUT /api/bands/:id`  
**Authentication:** Required (Bearer Token, must be owner)

**Request Body:** Same as Create Band

---

### Delete Band
Delete a band.

**Endpoint:** `DELETE /api/bands/:id`  
**Authentication:** Required (Bearer Token, must be owner)

**Success Response (200):**
```json
{
  "message": "band deleted successfully"
}
```

---

### Get Band Music
Retrieve all music in a band.

**Endpoint:** `GET /api/bands/:id/musics`  
**Authentication:** Required (Bearer Token)

**Success Response (200):**
```json
{
  "data": [
    {
      "id": 1,
      "title": "Amazing Grace",
      "cover": "...",
      "audio": "..."
    }
  ]
}
```

---

### Add Music to Band
Add multiple music tracks to a band.

**Endpoint:** `POST /api/bands/:id/musics`  
**Authentication:** Required (Bearer Token, must be owner)

**Request Body:**
```json
{
  "music_ids": [1, 2, 3, 4]
}
```

**Success Response (200):**
```json
{
  "message": "music added to band successfully"
}
```

---

### Remove Music from Band
Remove a music track from a band.

**Endpoint:** `DELETE /api/bands/:id/musics/:music_id`  
**Authentication:** Required (Bearer Token, must be owner)

**Success Response (200):**
```json
{
  "message": "music removed from band"
}
```

---

### Reorder Band Music
Change the display order of music in a band.

**Endpoint:** `PUT /api/bands/:id/musics/reorder`  
**Authentication:** Required (Bearer Token, must be owner)

**Request Body:**
```json
{
  "music_orders": [
    { "music_id": 3, "display_order": 1 },
    { "music_id": 1, "display_order": 2 },
    { "music_id": 2, "display_order": 3 }
  ]
}
```

**Success Response (200):**
```json
{
  "message": "music order updated successfully"
}
```

---

### Get Band Members
Retrieve all members of a band.

**Endpoint:** `GET /api/bands/:id/members`  
**Authentication:** Required (Bearer Token)

**Success Response (200):**
```json
{
  "data": [
    {
      "id": 5,
      "username": "johndoe",
      "fullname": "John Doe",
      "email": "john@example.com",
      "role": "user",
      "created_at": "2025-10-01T10:00:00Z"
    }
  ]
}
```

---

## Settings

### Get User Settings
Retrieve current user's settings.

**Endpoint:** `GET /api/settings`  
**Authentication:** Required (Bearer Token)

**Success Response (200):**
```json
{
  "id": 1,
  "user_id": 5,
  "language": "en",
  "theme": "dark",
  "notify_on_booking": true,
  "notify_on_music": false,
  "notify_on_event": true,
  "enable_push_notifications": true,
  "created_at": "2025-10-25T10:30:00Z",
  "updated_at": "2025-10-25T12:00:00Z"
}
```

---

### Update Settings
Update user preferences.

**Endpoint:** `PUT /api/settings`  
**Authentication:** Required (Bearer Token)

**Request Body:**
```json
{
  "language": "kh",
  "theme": "dark",
  "notify_on_booking": true,
  "notify_on_music": true,
  "notify_on_event": true,
  "enable_push_notifications": false
}
```

**Valid Values:**
- `language`: `en`, `kh`, `kr`, `cn`
- `theme`: `light`, `dark`, `auto`
- Notification flags: `true` or `false`

**Success Response (200):** Returns updated settings

**Error Response (400):**
```json
{
  "error": "invalid settings data"
}
```

---

### Reset Settings to Defaults
Reset all settings to default values.

**Endpoint:** `POST /api/settings/reset`  
**Authentication:** Required (Bearer Token)

**Success Response (200):**
```json
{
  "id": 1,
  "user_id": 5,
  "language": "en",
  "theme": "light",
  "notify_on_booking": true,
  "notify_on_music": false,
  "notify_on_event": true,
  "enable_push_notifications": true,
  "created_at": "2025-10-25T10:30:00Z",
  "updated_at": "2025-10-25T12:30:00Z"
}
```

**Default Values:**
- Language: `en`
- Theme: `light`
- Notify on Booking: `true`
- Notify on Music: `false`
- Notify on Event: `true`
- Push Notifications: `true`

---

## Health Check

### Ping
Check if the API is running.

**Endpoint:** `GET /ping`  
**Authentication:** None (Public)

**Success Response (200):**
```json
{
  "message": "api work..."
}
```

---

## Common Error Responses

### 401 Unauthorized
```json
{
  "error": "user not authenticated"
}
```

### 403 Forbidden
```json
{
  "error": "permission denied"
}
```

### 404 Not Found
```json
{
  "error": "resource not found"
}
```

### 500 Internal Server Error
```json
{
  "error": "internal server error"
}
```

---

## Authentication Flow

1. **Register or Login** to get a JWT token
2. **Include the token** in all subsequent requests:
   ```
   Authorization: Bearer <your_jwt_token>
   ```
3. The token is automatically validated by the middleware
4. **Logout** to invalidate the token

---

## Notes

- All authenticated endpoints require a valid JWT token in the `Authorization` header
- Timestamps are in ISO 8601 format (UTC)
- Pagination defaults: `page=1`, `page_size=20`, `max=100`
- Settings are automatically created when a user registers (via database trigger)
- Users can only modify/delete their own resources (music, events, bands)
- All POST/PUT requests should have `Content-Type: application/json` header

---

## Environment Variables

```bash
PORT=8080                    # Server port
POSTGRES_URL=postgres://...  # Database connection string
SECRET=your_secret_key       # JWT secret key
GOOGLE_CLIENT_ID=...         # Google OAuth client ID (optional)
```

---

**Last Updated:** October 25, 2025  
**API Version:** 1.0

