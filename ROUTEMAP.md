# Ezra API - Route Documentation

Complete API route reference for the Ezra music management system.

## Base URL
```
Production: http://your-domain.com
Development: http://localhost:8080
```

## Table of Contents  
- [Device Tokens (FCM Push Notifications)](#device-tokens-fcm-push-notifications)
- [Authentication](#authentication)
- [User Management](#user-management)
- [Music](#music)
- [Events](#events)
- [Bookings](#bookings)
- [Favorites](#favorites)
- [Bands](#bands)
- [Settings](#settings)
- [Notifications](#notifications)
- [Health Check](#health-check)

---

## Authentication

### Register
Create a new user account with email/password.

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

**Validation Rules:**
- `username`: required, min 1 char, max 100 chars
- `fullname`: required, min 1 char, max 100 chars
- `email`: required, valid email format, max 100 chars
- `password`: required, min 6 chars

**Success Response (201):**
```json
{
  "message": "user registered successfully",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Error Responses:**
- **400 Bad Request:** Invalid input or validation error
- **400 Bad Request:** Username or email already exists

```json
{
  "error": "username already exists"
}
```

---

### Login
Authenticate with username/email and password to get JWT token.

**Endpoint:** `POST /login`  
**Authentication:** None (Public)

**Request Body:**
```json
{
  "username": "johndoe",
  "password": "password123"
}
```

Or use email:
```json
{
  "username": "john@example.com",
  "password": "password123"
}
```

**Success Response (200):**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Error Responses:**
- **401 Unauthorized:** Invalid credentials
- **400 Bad Request:** Invalid request format

```json
{
  "error": "invalid credentials"
}
```

---

### Google Login
Authenticate using Google OAuth ID token.

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

**Error Responses:**
- **401 Unauthorized:** Invalid Google ID token
- **400 Bad Request:** Invalid token format or missing claims

```json
{
  "error": "invalid Google ID token"
}
```

**Note:** This endpoint validates the Google ID token, extracts user information (sub, email, name, picture), and either creates a new user or logs in an existing user linked to that Google account.

---

## User Management

### Logout
Invalidate the current user's JWT token.

**Endpoint:** `POST /api/logout`  
**Authentication:** Required (Bearer Token)

**Headers:**
```
Authorization: Bearer <your_jwt_token>
```

**Success Response (200):**
```json
{
  "message": "logged out successfully"
}
```

**Error Responses:**
- **401 Unauthorized:** User not authenticated
- **400 Bad Request:** Logout failed

---

### Delete User
Delete the current user's account and all associated data.

**Endpoint:** `DELETE /api/user`  
**Authentication:** Required (Bearer Token)

**Headers:**
```
Authorization: Bearer <your_jwt_token>
```

**Success Response (200):**
```json
{
  "message": "user deleted successfully"
}
```

**Error Responses:**
- **401 Unauthorized:** User not authenticated
- **400 Bad Request:** Deletion failed

**Note:** This is a permanent action. All user data including music, events, bookings, favorites, bands, and settings will be deleted due to CASCADE constraints.

---

## Music

### Create Music
Add a new music track to the system.

**Endpoint:** `POST /api/musics`  
**Authentication:** Required (Bearer Token)

**Headers:**
```
Authorization: Bearer <your_jwt_token>
Content-Type: application/json
```

**Request Body:**
```json
{
  "title": "Amazing Grace",
  "artist": "John Newton",
  "album": "Hymns of Grace",
  "genre": "Hymn",
  "duration": 240,
  "bpm": 80,
  "key": "G",
  "cover": "https://example.com/covers/amazing-grace.jpg",
  "lyrics": "Amazing grace, how sweet the sound...",
  "description": "A classic hymn about redemption"
}
```

**Validation Rules:**
- `title`: required, min 1 char, max 255 chars
- `artist`: optional, max 255 chars
- `album`: optional, max 255 chars
- `genre`: optional, max 100 chars
- `duration`: optional, integer (in seconds)
- `bpm`: optional, integer (beats per minute)
- `key`: optional, max 10 chars (musical key like "C", "Am", "G", etc.)
- `cover`: optional, max 255 chars
- `lyrics`: optional, text
- `description`: optional, text

**Success Response (201):**
```json
{
  "message": "music created successfully"
}
```

**Error Responses:**
- **400 Bad Request:** Invalid input or validation error
- **401 Unauthorized:** User not authenticated

---

### Get All Music
Retrieve all music tracks with optional pagination.

**Endpoint:** `GET /api/musics`  
**Authentication:** Required (Bearer Token)

**Query Parameters:**
- `page` (optional, default: 1) - Page number
- `page_size` (optional, default: 20, max: 100) - Items per page

**Examples:**
- Get all: `GET /api/musics`
- With pagination: `GET /api/musics?page=1&page_size=20`

**Success Response (200) - Paginated:**
```json
{
  "data": [
    {
      "id": 1,
      "title": "Amazing Grace",
      "artist": "John Newton",
      "album": "Hymns of Grace",
      "genre": "Hymn",
      "duration": 240,
      "bpm": 80,
      "key": "G",
      "cover": "https://example.com/covers/amazing-grace.jpg",
      "lyrics": "Amazing grace, how sweet the sound...",
      "description": "A classic hymn about redemption",
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

**Success Response (200) - All Results:**
```json
[
  {
    "id": 1,
    "title": "Amazing Grace",
    "artist": "John Newton",
    "album": "Hymns of Grace",
    "genre": "Hymn",
    "duration": 240,
    "bpm": 80,
    "key": "G",
    "cover": "https://example.com/covers/amazing-grace.jpg",
    "lyrics": "Amazing grace, how sweet the sound...",
    "description": "A classic hymn about redemption",
    "user_id": 5,
    "created_at": "2025-10-25T10:30:00Z",
    "updated_at": "2025-10-25T10:30:00Z"
  }
]
```

---

### Get User's Music
Retrieve all music tracks created by the authenticated user.

**Endpoint:** `GET /api/musics/user`  
**Authentication:** Required (Bearer Token)

**Success Response (200):**
```json
[
  {
    "id": 1,
    "title": "Amazing Grace",
    "artist": "John Newton",
    "album": "Hymns of Grace",
    "genre": "Hymn",
    "duration": 240,
    "bpm": 80,
    "key": "G",
    "cover": "https://example.com/covers/amazing-grace.jpg",
    "lyrics": "Amazing grace, how sweet the sound...",
    "description": "A classic hymn about redemption",
    "user_id": 5,
    "created_at": "2025-10-25T10:30:00Z",
    "updated_at": "2025-10-25T10:30:00Z"
  }
]
```

**Error Responses:**
- **401 Unauthorized:** User not authenticated
- **500 Internal Server Error:** Server error

---

### Get Music by ID
Retrieve a specific music track by its ID.

**Endpoint:** `GET /api/musics/:id`  
**Authentication:** Required (Bearer Token)

**Path Parameters:**
- `id` (required) - Music ID (integer)

**Success Response (200):**
```json
{
  "id": 1,
  "title": "Amazing Grace",
  "artist": "John Newton",
  "album": "Hymns of Grace",
  "genre": "Hymn",
  "duration": 240,
  "bpm": 80,
  "key": "G",
  "cover": "https://example.com/covers/amazing-grace.jpg",
  "lyrics": "Amazing grace, how sweet the sound...",
  "description": "A classic hymn about redemption",
  "user_id": 5,
  "created_at": "2025-10-25T10:30:00Z",
  "updated_at": "2025-10-25T10:30:00Z"
}
```

**Error Responses:**
- **400 Bad Request:** Invalid ID format
- **404 Not Found:** Music not found

```json
{
  "error": "music not found"
}
```

---

### Update Music
Update an existing music track. Only the owner can update.

**Endpoint:** `PUT /api/musics/:id`  
**Authentication:** Required (Bearer Token, must be owner)

**Path Parameters:**
- `id` (required) - Music ID (integer)

**Request Body:**
```json
{
  "title": "Amazing Grace (Updated)",
  "artist": "John Newton",
  "album": "Hymns of Grace - Remastered",
  "genre": "Hymn",
  "duration": 245,
  "bpm": 82,
  "key": "G",
  "cover": "https://example.com/covers/amazing-grace-v2.jpg",
  "lyrics": "Amazing grace, how sweet the sound...",
  "description": "A classic hymn about redemption - updated version"
}
```

**Success Response (200):**
```json
{
  "message": "music updated successfully"
}
```

**Error Responses:**
- **400 Bad Request:** Invalid input or user doesn't own this music
- **401 Unauthorized:** User not authenticated
- **404 Not Found:** Music not found

---

### Delete Music
Delete a music track. Only the owner can delete.

**Endpoint:** `DELETE /api/musics/:id`  
**Authentication:** Required (Bearer Token, must be owner)

**Path Parameters:**
- `id` (required) - Music ID (integer)

**Success Response (200):**
```json
{
  "message": "music deleted successfully"
}
```

**Error Responses:**
- **400 Bad Request:** Invalid ID or user doesn't own this music
- **401 Unauthorized:** User not authenticated
- **404 Not Found:** Music not found

---

## Events

### Create Event
Create a new event with optional music tracks.

**Endpoint:** `POST /api/events`  
**Authentication:** Required (Bearer Token)

**Note:** Creating an event automatically sends a broadcast notification to all users informing them about the new event. The event creator will not receive the notification themselves.

**Request Body:**
```json
{
  "title": "Sunday Worship Service",
  "content": "Weekly Sunday morning worship service",
  "cover": "https://example.com/events/sunday-worship.jpg",
  "location": "Main Church Auditorium",
  "start_time": "2025-11-01T09:00:00Z",
  "end_time": "2025-11-01T11:00:00Z",
  "music_ids": [1, 2, 3, 4]
}
```

**Validation Rules:**
- `title`: required, min 1 char, max 255 chars
- `content`: optional, text
- `cover`: optional, max 255 chars
- `location`: required
- `start_time`: required, ISO 8601 format
- `end_time`: required, ISO 8601 format
- `music_ids`: optional, array of music IDs

**Success Response (201):**
```json
{
  "id": 1,
  "title": "Sunday Worship Service",
  "content": "Weekly Sunday morning worship service",
  "cover": "https://example.com/events/sunday-worship.jpg",
  "location": "Main Church Auditorium",
  "start_time": "2025-11-01T09:00:00Z",
  "end_time": "2025-11-01T11:00:00Z",
  "user_id": 5,
  "musics": [
    {
      "id": 1,
      "title": "Amazing Grace",
      "artist": "John Newton",
      "genre": "Hymn",
      "cover": "..."
    }
  ],
  "created_at": "2025-10-29T10:30:00Z",
  "updated_at": "2025-10-29T10:30:00Z"
}
```

---

### Get All Events
Retrieve all events with optional pagination.

**Endpoint:** `GET /api/events`  
**Authentication:** Required (Bearer Token)

**Query Parameters:**
- `page` (optional, default: 1)
- `page_size` (optional, default: 20, max: 100)

**Success Response (200):** Similar to Get All Music, returns paginated event list with associated music tracks.

---

### Get User's Events
Retrieve events created by the authenticated user.

**Endpoint:** `GET /api/events/user`  
**Authentication:** Required (Bearer Token)

**Query Parameters:**
- `page` (optional)
- `page_size` (optional)

**Success Response (200):** Returns list of user's events with music.

---

### Get Event by ID
Retrieve a specific event with its associated music tracks.

**Endpoint:** `GET /api/events/:id`  
**Authentication:** Required (Bearer Token)

**Path Parameters:**
- `id` (required) - Event ID (integer)

**Success Response (200):** Returns event details including all associated music tracks.

---

### Update Event
Update an existing event. Only the owner can update.

**Endpoint:** `PUT /api/events/:id`  
**Authentication:** Required (Bearer Token, must be owner)

**Path Parameters:**
- `id` (required) - Event ID (integer)

**Request Body:** Same format as Create Event

**Success Response (200):** Returns updated event with music.

---

### Delete Event
Delete an event. Only the owner can delete.

**Endpoint:** `DELETE /api/events/:id`  
**Authentication:** Required (Bearer Token, must be owner)

**Path Parameters:**
- `id` (required) - Event ID (integer)

**Success Response (200):**
```json
{
  "message": "event deleted successfully"
}
```

---

## Bookings

### Create Booking
Register/book attendance for an event.

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
  "created_at": "2025-10-29T10:30:00Z",
  "updated_at": "2025-10-29T10:30:00Z"
}
```

**Error Responses:**
- **400 Bad Request:** Invalid input or duplicate booking
- **401 Unauthorized:** User not authenticated

---

### Get All Bookings
Retrieve all bookings in the system (typically admin use).

**Endpoint:** `GET /api/bookings`  
**Authentication:** Required (Bearer Token)

**Query Parameters:**
- `page` (optional, default: 1)
- `page_size` (optional, default: 20)

---

### Get User's Bookings
Retrieve all bookings made by the authenticated user.

**Endpoint:** `GET /api/bookings/user`  
**Authentication:** Required (Bearer Token)

**Query Parameters:**
- `page` (optional)
- `page_size` (optional)

**Success Response (200):** Returns paginated list of user's bookings.

---

### Get Bookings by Event
Retrieve all bookings for a specific event.

**Endpoint:** `GET /api/bookings/event/:event_id`  
**Authentication:** Required (Bearer Token)

**Path Parameters:**
- `event_id` (required) - Event ID (integer)

**Query Parameters:**
- `page` (optional)
- `page_size` (optional)

---

### Get Booking by ID
Retrieve a specific booking.

**Endpoint:** `GET /api/bookings/:id`  
**Authentication:** Required (Bearer Token)

**Path Parameters:**
- `id` (required) - Booking ID (integer)

---

### Update Booking
Update booking status or notes.

**Endpoint:** `PUT /api/bookings/:id`  
**Authentication:** Required (Bearer Token)

**Path Parameters:**
- `id` (required) - Booking ID (integer)

**Request Body:**
```json
{
  "status": "confirmed",
  "notes": "Attendance confirmed"
}
```

**Valid Status Values:**
- `pending` - Awaiting confirmation
- `confirmed` - Attendance confirmed
- `cancelled` - Booking cancelled

**Success Response (200):** Returns updated booking

---

### Delete Booking
Cancel/delete a booking. Only the booking owner can delete.

**Endpoint:** `DELETE /api/bookings/:id`  
**Authentication:** Required (Bearer Token, must be owner)

**Path Parameters:**
- `id` (required) - Booking ID (integer)

**Success Response (200):**
```json
{
  "message": "booking deleted successfully"
}
```

---

## Favorites

### Get User's Favorites
Retrieve all music tracks favorited by the authenticated user.

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
      "artist": "John Newton",
      "album": "Hymns of Grace",
      "genre": "Hymn",
      "duration": 240,
      "bpm": 80,
      "key": "G",
      "cover": "https://example.com/covers/amazing-grace.jpg",
      "lyrics": "Amazing grace, how sweet the sound...",
      "description": "A classic hymn about redemption",
      "user_id": 3,
      "created_at": "2025-10-25T10:30:00Z",
      "updated_at": "2025-10-25T10:30:00Z"
    }
  ],
  "pagination": {
    "current_page": 1,
    "page_size": 20,
    "total_pages": 2,
    "total_records": 35,
    "has_next_page": true,
    "has_prev_page": false
  }
}
```

---

### Add to Favorites
Add a music track to the user's favorites list.

**Endpoint:** `POST /api/favorites/music/:id`  
**Authentication:** Required (Bearer Token)

**Path Parameters:**
- `id` (required) - Music ID (integer)

**Success Response (200):**
```json
{
  "message": "music added to favorites"
}
```

**Error Responses:**
- **400 Bad Request:** Music already in favorites or music doesn't exist
- **401 Unauthorized:** User not authenticated

---

### Remove from Favorites
Remove a music track from the user's favorites list.

**Endpoint:** `DELETE /api/favorites/music/:id`  
**Authentication:** Required (Bearer Token)

**Path Parameters:**
- `id` (required) - Music ID (integer)

**Success Response (200):**
```json
{
  "message": "music removed from favorites"
}
```

**Error Responses:**
- **400 Bad Request:** Music not in favorites
- **401 Unauthorized:** User not authenticated

---

### Check if Favorite
Check whether a specific music track is in the user's favorites.

**Endpoint:** `GET /api/favorites/music/:id/check`  
**Authentication:** Required (Bearer Token)

**Path Parameters:**
- `id` (required) - Music ID (integer)

**Success Response (200):**
```json
{
  "is_favorite": true
}
```

Or:
```json
{
  "is_favorite": false
}
```

---

### Get Favorite Count
Get the total number of users who have favorited a specific music track.

**Endpoint:** `GET /api/favorites/music/:id/count`  
**Authentication:** Required (Bearer Token)

**Path Parameters:**
- `id` (required) - Music ID (integer)

**Success Response (200):**
```json
{
  "count": 42
}
```

**Note:** This is useful for displaying popularity metrics.

---

## Bands

### Create Band
Create a new band/music collection/library.

**Endpoint:** `POST /api/bands`  
**Authentication:** Required (Bearer Token)

**Request Body:**
```json
{
  "name": "Worship Team A",
  "description": "Main worship team for Sunday services",
  "cover": "https://example.com/bands/worship-team-a.jpg",
  "is_public": true
}
```

**Validation Rules:**
- `name`: required, min 1 char, max 255 chars
- `description`: optional, text
- `cover`: optional, max 255 chars
- `is_public`: optional, boolean (default: false)

**Success Response (201):**
```json
{
  "id": 1,
  "name": "Worship Team A",
  "description": "Main worship team for Sunday services",
  "cover": "https://example.com/bands/worship-team-a.jpg",
  "is_public": true,
  "user_id": 5,
  "created_at": "2025-10-29T10:30:00Z",
  "updated_at": "2025-10-29T10:30:00Z"
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

**Success Response (200):** Returns paginated list of all bands.

---

### Get User's Bands
Retrieve bands created by the authenticated user.

**Endpoint:** `GET /api/bands/user`  
**Authentication:** Required (Bearer Token)

**Query Parameters:**
- `page` (optional)
- `page_size` (optional)

---

### Get Public Bands
Retrieve all public bands (bands with `is_public = true`).

**Endpoint:** `GET /api/bands/public`  
**Authentication:** Required (Bearer Token)

**Query Parameters:**
- `page` (optional)
- `page_size` (optional)

---

### Get Band by ID
Retrieve a specific band's details.

**Endpoint:** `GET /api/bands/:id`  
**Authentication:** Required (Bearer Token)

**Path Parameters:**
- `id` (required) - Band ID (integer)

**Success Response (200):** Returns band details.

---

### Update Band
Update band information. Only the owner can update.

**Endpoint:** `PUT /api/bands/:id`  
**Authentication:** Required (Bearer Token, must be owner)

**Path Parameters:**
- `id` (required) - Band ID (integer)

**Request Body:** Same format as Create Band

**Success Response (200):** Returns updated band.

---

### Delete Band
Delete a band. Only the owner can delete.

**Endpoint:** `DELETE /api/bands/:id`  
**Authentication:** Required (Bearer Token, must be owner)

**Path Parameters:**
- `id` (required) - Band ID (integer)

**Success Response (200):**
```json
{
  "message": "band deleted successfully"
}
```

---

### Get Band Music
Retrieve all music tracks in a specific band.

**Endpoint:** `GET /api/bands/:id/musics`  
**Authentication:** Required (Bearer Token)

**Path Parameters:**
- `id` (required) - Band ID (integer)

**Success Response (200):**
```json
{
  "data": [
    {
      "id": 1,
      "title": "Amazing Grace",
      "artist": "John Newton",
      "album": "Hymns of Grace",
      "genre": "Hymn",
      "duration": 240,
      "bpm": 80,
      "key": "G",
      "cover": "https://example.com/covers/amazing-grace.jpg",
      "lyrics": "Amazing grace, how sweet the sound...",
      "description": "A classic hymn about redemption",
      "user_id": 5,
      "created_at": "2025-10-25T10:30:00Z",
      "updated_at": "2025-10-25T10:30:00Z"
    }
  ]
}
```

---

### Add Music to Band
Add one or more music tracks to a band. Only the owner can add.

**Endpoint:** `POST /api/bands/:id/musics`  
**Authentication:** Required (Bearer Token, must be owner)

**Path Parameters:**
- `id` (required) - Band ID (integer)

**Request Body:**
```json
{
  "music_ids": [1, 2, 3, 4, 5]
}
```

**Success Response (200):**
```json
{
  "message": "music added to band successfully"
}
```

**Error Responses:**
- **400 Bad Request:** Invalid input, music doesn't exist, or already in band
- **401 Unauthorized:** User not authenticated or not band owner

---

### Remove Music from Band
Remove a music track from a band. Only the owner can remove.

**Endpoint:** `DELETE /api/bands/:id/musics/:music_id`  
**Authentication:** Required (Bearer Token, must be owner)

**Path Parameters:**
- `id` (required) - Band ID (integer)
- `music_id` (required) - Music ID (integer)

**Success Response (200):**
```json
{
  "message": "music removed from band"
}
```

---

### Reorder Band Music
Change the display order of music tracks in a band. Only the owner can reorder.

**Endpoint:** `PUT /api/bands/:id/musics/reorder`  
**Authentication:** Required (Bearer Token, must be owner)

**Path Parameters:**
- `id` (required) - Band ID (integer)

**Request Body:**
```json
{
  "music_orders": [
    { "music_id": 3, "display_order": 1 },
    { "music_id": 1, "display_order": 2 },
    { "music_id": 5, "display_order": 3 },
    { "music_id": 2, "display_order": 4 }
  ]
}
```

**Success Response (200):**
```json
{
  "message": "music order updated successfully"
}
```

**Note:** This allows you to customize the order in which music appears in the band's playlist.

---

### Get Band Members
Retrieve all users who are members of a specific band.

**Endpoint:** `GET /api/bands/:id/members`  
**Authentication:** Required (Bearer Token)

**Path Parameters:**
- `id` (required) - Band ID (integer)

**Success Response (200):**
```json
{
  "data": [
    {
      "id": 5,
      "username": "johndoe",
      "fullname": "John Doe",
      "email": "john@example.com",
      "profile": "https://example.com/profiles/johndoe.jpg",
      "role": "user",
      "created_at": "2025-10-01T10:00:00Z"
    }
  ]
}
```

**Note:** Members are users who have their `band_id` set to this band.

---

## Device Tokens (FCM Push Notifications)

Firebase Cloud Messaging (FCM) integration for real-time push notifications to mobile apps and web browsers.

### Features
- Push notifications to iOS, Android, and Web platforms
- Automatic token management and cleanup
- Background notification delivery
- Works seamlessly with the notification system

### Setup Requirements
1. Create a Firebase project at [Firebase Console](https://console.firebase.google.com/)
2. Download the service account key (JSON file)
3. Set environment variable: `FIREBASE_CREDENTIALS_PATH=/path/to/firebase-adminsdk.json`
4. If not set, the API will run with a dummy FCM service (notifications still work, but no push)

---

### Register Device Token
Register a device token to receive push notifications.

**Endpoint:** `POST /api/device-tokens/register`  
**Authentication:** Required (Bearer Token)

**Request Body:**
```json
{
  "token": "fcm_device_token_string_here",
  "platform": "ios"
}
```

**Parameters:**
- `token` (required) - FCM device token from Firebase SDK
- `platform` (required) - Platform type: `ios`, `android`, or `web`

**Success Response (200):**
```json
{
  "message": "Device token registered successfully",
  "token_id": 42
}
```

**Notes:**
- Tokens are automatically upserted (updated if they already exist)
- One user can have multiple device tokens (multiple devices)
- Inactive tokens are automatically cleaned up when FCM reports them as invalid

---

### Unregister Device Token
Remove a device token (e.g., on logout).

**Endpoint:** `POST /api/device-tokens/unregister`  
**Authentication:** Required (Bearer Token)

**Request Body:**
```json
{
  "token": "fcm_device_token_string_here"
}
```

**Success Response (200):**
```json
{
  "message": "Device token unregistered successfully"
}
```

---

### Clear All Device Tokens
Remove all device tokens for the current user.

**Endpoint:** `DELETE /api/device-tokens/clear`  
**Authentication:** Required (Bearer Token)

**Success Response (200):**
```json
{
  "message": "All device tokens deleted successfully"
}
```

**Use Cases:**
- User logs out from all devices
- Account settings: "Sign out everywhere"
- Privacy/security feature

---

### Mobile Integration Example

#### React Native / Expo
```javascript
import messaging from '@react-native-firebase/messaging';
import axios from 'axios';

// Request permission and get token
async function registerForPushNotifications() {
  const authStatus = await messaging().requestPermission();
  
  if (authStatus === messaging.AuthorizationStatus.AUTHORIZED) {
    const token = await messaging().getToken();
    
    // Send token to your API
    await axios.post('https://your-api.com/api/device-tokens/register', {
      token: token,
      platform: Platform.OS === 'ios' ? 'ios' : 'android'
    }, {
      headers: { 'Authorization': `Bearer ${yourJWTToken}` }
    });
  }
}

// Handle foreground notifications
messaging().onMessage(async remoteMessage => {
  console.log('Notification received:', remoteMessage);
  // Show in-app notification
});

// Handle notification tap
messaging().onNotificationOpenedApp(remoteMessage => {
  // Navigate to relevant screen
  if (remoteMessage.data?.related_type === 'event') {
    navigation.navigate('EventDetail', { 
      id: remoteMessage.data.related_id 
    });
  }
});
```

#### Flutter
```dart
import 'package:firebase_messaging/firebase_messaging.dart';
import 'package:http/http.dart' as http;

Future<void> registerForPushNotifications() async {
  FirebaseMessaging messaging = FirebaseMessaging.instance;
  
  NotificationSettings settings = await messaging.requestPermission(
    alert: true,
    badge: true,
    sound: true,
  );

  if (settings.authorizationStatus == AuthorizationStatus.authorized) {
    String? token = await messaging.getToken();
    
    // Send to your API
    await http.post(
      Uri.parse('https://your-api.com/api/device-tokens/register'),
      headers: {
        'Authorization': 'Bearer $yourJWTToken',
        'Content-Type': 'application/json',
      },
      body: jsonEncode({
        'token': token,
        'platform': Platform.isIOS ? 'ios' : 'android',
      }),
    );
  }
}

// Handle foreground messages
FirebaseMessaging.onMessage.listen((RemoteMessage message) {
  // Show notification
});

// Handle notification tap
FirebaseMessaging.onMessageOpenedApp.listen((RemoteMessage message) {
  // Navigate to screen
});
```

#### Web (PWA)
```javascript
import { getMessaging, getToken, onMessage } from 'firebase/messaging';

async function registerForPushNotifications() {
  const messaging = getMessaging();
  
  try {
    const token = await getToken(messaging, { 
      vapidKey: 'YOUR_VAPID_KEY' 
    });
    
    // Send token to your API
    await fetch('https://your-api.com/api/device-tokens/register', {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${yourJWTToken}`,
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        token: token,
        platform: 'web'
      })
    });
  } catch (err) {
    console.error('Error getting FCM token:', err);
  }
}

// Handle foreground messages
onMessage(messaging, (payload) => {
  console.log('Message received:', payload);
  // Show notification
});
```

---

## Settings

### Get User Settings
Retrieve the authenticated user's settings/preferences.

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
  "updated_at": "2025-10-29T12:00:00Z"
}
```

**Note:** Settings are automatically created when a user registers via a database trigger.

---

### Update Settings
Update the authenticated user's preferences.

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
- `language`: `en` (English), `kh` (Khmer), `kr` (Korean), `cn` (Chinese)
- `theme`: `light`, `dark`, `auto`
- Notification flags: `true` or `false` (boolean)

**Success Response (200):** Returns updated settings object.

**Error Responses:**
- **400 Bad Request:** Invalid settings data or theme not in allowed values
- **401 Unauthorized:** User not authenticated

---

### Reset Settings to Defaults
Reset all user settings to their default values.

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
  "updated_at": "2025-10-29T12:30:00Z"
}
```

**Default Values:**
- Language: `en`
- Theme: `light`
- Notify on Booking: `true`
- Notify on Music: `false`
- Notify on Event: `true`
- Enable Push Notifications: `true`

---

## Notifications

**Recipient Types:**
- `user` - Send to a specific user
- `band` - Send to all members of a band/team
- `all` - Broadcast to all users in the system

**Features:**
- Users automatically receive notifications sent to them, their band, and all broadcasts
- Read tracking per notification
- Automatic filtering based on user's band membership

---

### Create Notification (To Specific User)
Create a new notification for a specific user.

**Endpoint:** `POST /api/notifications`  
**Authentication:** Required (Bearer Token)

**Request Body:**
```json
{
  "user_id": 10,
  "title": "New Booking Confirmation",
  "message": "Your booking for Sunday Service has been confirmed",
  "type": "booking",
  "related_type": "booking",
  "related_id": 15
}
```

**Validation Rules:**
- `user_id`: required for user notifications
- `title`: required, min 1 char, max 255 chars
- `message`: required, min 1 char
- `type`: required, one of: `info`, `success`, `warning`, `error`, `booking`, `music`, `event`
- `related_type`: optional, one of: `music`, `event`, `booking`, `band`
- `related_id`: optional, integer (ID of related resource)

**Success Response (201):**
```json
{
  "id": 1,
  "user_id": 10,
  "sender_id": 5,
  "recipient_type": "user",
  "title": "New Booking Confirmation",
  "message": "Your booking for Sunday Service has been confirmed",
  "type": "booking",
  "related_type": "booking",
  "related_id": 15,
  "is_read": false,
  "created_at": "2025-10-29T10:30:00Z",
  "updated_at": "2025-10-29T10:30:00Z"
}
```

**Error Responses:**
- **400 Bad Request:** Invalid input or validation error
- **401 Unauthorized:** User not authenticated

---

### Create Band Notification (To Team)
Create a notification for all members of a band/team.

**Endpoint:** `POST /api/notifications/band/:band_id`  
**Authentication:** Required (Bearer Token)

**Path Parameters:**
- `band_id` (required) - Band ID (integer)

**Request Body:**
```json
{
  "title": "New Song Added",
  "message": "Amazing Grace has been added to the band repertoire",
  "type": "music",
  "related_type": "music",
  "related_id": 42
}
```

**Success Response (201):**
```json
{
  "id": 2,
  "band_id": 3,
  "sender_id": 5,
  "recipient_type": "band",
  "title": "New Song Added",
  "message": "Amazing Grace has been added to the band repertoire",
  "type": "music",
  "related_type": "music",
  "related_id": 42,
  "is_read": false,
  "created_at": "2025-10-29T10:35:00Z",
  "updated_at": "2025-10-29T10:35:00Z"
}
```

**Note:** All users who are members of this band (have `band_id` set) will receive this notification.

---

### Create Broadcast Notification (To All Users)
Create a broadcast notification that all users in the system will receive.

**Endpoint:** `POST /api/notifications/broadcast`  
**Authentication:** Required (Bearer Token)

**Request Body:**
```json
{
  "title": "System Maintenance",
  "message": "The system will undergo maintenance on Sunday at 2 AM",
  "type": "warning"
}
```

**Success Response (201):**
```json
{
  "id": 3,
  "sender_id": 1,
  "recipient_type": "all",
  "title": "System Maintenance",
  "message": "The system will undergo maintenance on Sunday at 2 AM",
  "type": "warning",
  "is_read": false,
  "created_at": "2025-10-29T10:40:00Z",
  "updated_at": "2025-10-29T10:40:00Z"
}
```

**Note:** Every user in the system will see this notification. Use sparingly for important announcements.

---

### Get All Notifications
Retrieve all notifications for the authenticated user with pagination.

**Endpoint:** `GET /api/notifications`  
**Authentication:** Required (Bearer Token)

**Query Parameters:**
- `page` (optional, default: 1)
- `page_size` (optional, default: 20, max: 100)

**Example:** `GET /api/notifications?page=1&page_size=20`

**Success Response (200):**
```json
{
  "data": [
    {
      "id": 5,
      "recipient_type": "all",
      "sender_id": 1,
      "title": "System Maintenance",
      "message": "Scheduled maintenance tonight",
      "type": "warning",
      "is_read": false,
      "created_at": "2025-10-29T11:30:00Z",
      "updated_at": "2025-10-29T11:30:00Z"
    },
    {
      "id": 4,
      "band_id": 3,
      "recipient_type": "band",
      "sender_id": 2,
      "title": "Band Practice",
      "message": "Practice rescheduled to Saturday",
      "type": "info",
      "is_read": false,
      "created_at": "2025-10-29T11:15:00Z",
      "updated_at": "2025-10-29T11:15:00Z"
    },
    {
      "id": 3,
      "user_id": 5,
      "recipient_type": "user",
      "sender_id": 2,
      "title": "New Music Added",
      "message": "Amazing Grace has been added to your band",
      "type": "music",
      "related_type": "music",
      "related_id": 12,
      "is_read": false,
      "created_at": "2025-10-29T11:00:00Z",
      "updated_at": "2025-10-29T11:00:00Z"
    },
    {
      "id": 2,
      "user_id": 5,
      "recipient_type": "user",
      "title": "Event Reminder",
      "message": "Sunday Service starts in 1 hour",
      "type": "event",
      "related_type": "event",
      "related_id": 8,
      "is_read": true,
      "read_at": "2025-10-29T10:45:00Z",
      "created_at": "2025-10-29T10:30:00Z",
      "updated_at": "2025-10-29T10:45:00Z"
    }
  ],
  "pagination": {
    "current_page": 1,
    "page_size": 20,
    "total_pages": 1,
    "total_records": 2,
    "has_next_page": false,
    "has_prev_page": false
  }
}
```

---

### Get Unread Notifications
Retrieve all unread notifications for the authenticated user.

**Endpoint:** `GET /api/notifications/unread`  
**Authentication:** Required (Bearer Token)

**Success Response (200):**
```json
[
  {
    "id": 3,
    "user_id": 5,
    "title": "New Music Added",
    "message": "Amazing Grace has been added to your band",
    "type": "music",
    "related_type": "music",
    "related_id": 12,
    "is_read": false,
    "created_at": "2025-10-29T11:00:00Z",
    "updated_at": "2025-10-29T11:00:00Z"
  }
]
```

---

### Get Unread Count
Get the count of unread notifications for the authenticated user.

**Endpoint:** `GET /api/notifications/unread/count`  
**Authentication:** Required (Bearer Token)

**Success Response (200):**
```json
{
  "count": 5
}
```

**Note:** This is useful for displaying notification badges in the UI.

---

### Get Notification by ID
Retrieve a specific notification.

**Endpoint:** `GET /api/notifications/:id`  
**Authentication:** Required (Bearer Token, must be owner)

**Path Parameters:**
- `id` (required) - Notification ID (integer)

**Success Response (200):**
```json
{
  "id": 3,
  "user_id": 5,
  "title": "New Music Added",
  "message": "Amazing Grace has been added to your band",
  "type": "music",
  "related_type": "music",
  "related_id": 12,
  "is_read": false,
  "created_at": "2025-10-29T11:00:00Z",
  "updated_at": "2025-10-29T11:00:00Z"
}
```

**Error Responses:**
- **400 Bad Request:** Invalid ID format
- **404 Not Found:** Notification not found
- **401 Unauthorized:** User doesn't own this notification

---

### Mark Notification as Read
Mark a specific notification as read.

**Endpoint:** `PUT /api/notifications/:id/read`  
**Authentication:** Required (Bearer Token, must be owner)

**Path Parameters:**
- `id` (required) - Notification ID (integer)

**Success Response (200):**
```json
{
  "message": "notification marked as read"
}
```

**Error Responses:**
- **400 Bad Request:** Invalid ID or operation failed
- **401 Unauthorized:** User doesn't own this notification
- **404 Not Found:** Notification not found

---

### Mark All Notifications as Read
Mark all notifications as read for the authenticated user.

**Endpoint:** `PUT /api/notifications/read-all`  
**Authentication:** Required (Bearer Token)

**Success Response (200):**
```json
{
  "message": "all notifications marked as read"
}
```

**Note:** This operation affects all unread notifications for the user.

---

### Delete Notification
Delete a specific notification.

**Endpoint:** `DELETE /api/notifications/:id`  
**Authentication:** Required (Bearer Token, must be owner)

**Path Parameters:**
- `id` (required) - Notification ID (integer)

**Success Response (200):**
```json
{
  "message": "notification deleted successfully"
}
```

**Error Responses:**
- **400 Bad Request:** Invalid ID or operation failed
- **401 Unauthorized:** User doesn't own this notification
- **404 Not Found:** Notification not found

---

### Delete All Notifications
Delete all notifications for the authenticated user.

**Endpoint:** `DELETE /api/notifications`  
**Authentication:** Required (Bearer Token)

**Success Response (200):**
```json
{
  "message": "all notifications deleted successfully"
}
```

**Note:** This permanently deletes all notifications for the user. This action cannot be undone.

---

## Health Check

### Ping
Simple endpoint to check if the API is running and responsive.

**Endpoint:** `GET /ping`  
**Authentication:** None (Public)

**Success Response (200):**
```json
{
  "message": "api work..."
}
```

**Usage:** Use this for health checks, monitoring, load balancer health checks, etc.

---

## Common HTTP Status Codes

### Success Codes
- **200 OK** - Request succeeded
- **201 Created** - Resource created successfully

### Client Error Codes
- **400 Bad Request** - Invalid input, validation error, or business logic error
- **401 Unauthorized** - Missing or invalid authentication token
- **403 Forbidden** - User doesn't have permission (e.g., not the owner)
- **404 Not Found** - Resource not found

### Server Error Codes
- **500 Internal Server Error** - Unexpected server error

---

## Error Response Format

All error responses follow this format:

```json
{
  "error": "error message here"
}
```

Or with validation errors:

```json
{
  "errors": "Username is required"
}
```

---

## Authentication Flow

### Standard Flow (Email/Password)
1. **Register** → `POST /register` → Receive JWT token
2. **Login** → `POST /login` → Receive JWT token
3. **Use Token** → Include in `Authorization: Bearer <token>` header for all protected routes
4. **Logout** → `POST /api/logout` → Token invalidated

### OAuth Flow (Google)
1. **Get Google ID Token** → From Google OAuth flow on client side
2. **Authenticate** → `POST /auth/google` with ID token → Receive JWT token
3. **Use Token** → Same as standard flow

### Token Usage
Include the JWT token in the `Authorization` header for all protected routes:

```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

---

## Pagination

### Request
Most list endpoints support pagination via query parameters:

```
GET /api/musics?page=2&page_size=25
```

- `page`: Page number (default: 1)
- `page_size`: Items per page (default: 20, max: 100)

### Response
Paginated responses include metadata:

```json
{
  "data": [ /* array of items */ ],
  "pagination": {
    "current_page": 2,
    "page_size": 25,
    "total_pages": 10,
    "total_records": 245,
    "has_next_page": true,
    "has_prev_page": true
  }
}
```

---

## Database Schema Notes

### Tables
- **roles** - System roles with permissions (JSON)
- **users** - User accounts with authentication data
- **tokens** - Multi-device token management
- **settings** - User preferences (auto-created on registration)
- **musics** - Core music metadata
- **music_audio** - Multiple audio files per music (Original, Instrumental, etc.)
- **music_sheets** - Sheet music files (Lead Sheet, Chord Chart, etc.)
- **events** - Event scheduling
- **event_musics** - Many-to-many: Events ↔ Music
- **bookings** - Event registrations
- **bands** - Music collections/libraries
- **band_musics** - Many-to-many: Bands ↔ Music
- **favorites** - User favorites (Many-to-many: Users ↔ Music)
- **notifications** - User notifications with read tracking

### Automatic Features
- **Timestamps** - All tables have `created_at` and `updated_at` with auto-update triggers
- **Cascade Deletes** - Deleting a user/music/event automatically deletes related records
- **Settings Creation** - New users automatically get default settings via trigger
- **Timezone Support** - All timestamps use `TIMESTAMPTZ` (UTC, converted to Asia/Phnom_Penh)

---

## Environment Variables

```bash
# Server Configuration
PORT=8080

# Database
POSTGRES_URL=postgres://user:password@host:5432/dbname?sslmode=disable

# Security
SECRET=your_jwt_secret_key_here

# OAuth (Optional)
GOOGLE_CLIENT_ID=your_google_client_id_here

# Gin Mode
GIN_MODE=release  # or 'debug' for development
```

---

## Notes

- All authenticated endpoints require a valid JWT token in the `Authorization` header
- All timestamps are in ISO 8601 format with timezone (UTC)
- The API uses UTC internally but converts to `Asia/Phnom_Penh` timezone
- Users can only modify/delete their own resources (music, events, bookings, bands)
- All POST/PUT requests should include `Content-Type: application/json` header
- Music can belong to multiple bands and events (many-to-many relationships)
- Settings are automatically created when a user registers via database trigger
- Audio files are stored in the `music_audio` table (supports multiple files per music: Original, Instrumental, Acapella, Live, Acoustic, Remix, Cover)
- Sheet music files are stored in the `music_sheets` table (supports multiple sheets per music with different types and languages)
- Music metadata includes: artist, album, genre, duration, BPM, key, cover, lyrics, and description
- Notifications support multiple types (info, success, warning, error, booking, music, event) and can be linked to related resources
- Notifications track read status and read timestamps for user experience
- Notifications can be sent to specific users, bands/teams, or broadcast to all users
- Users automatically see notifications sent directly to them, to their band, and all broadcast notifications
- Sender tracking to know who created each notification
- Users do not receive notifications they created themselves (self-notifications are filtered out)
- Creating an event automatically sends a broadcast notification to all users (except the event creator)
- Firebase Cloud Messaging (FCM) is integrated for push notifications to mobile and web
- Device tokens are stored in the `device_tokens` table and support iOS, Android, and Web platforms
- Push notifications are sent automatically when notifications are created (works in background)
- Invalid device tokens are automatically cleaned up when FCM reports them

---

**Last Updated:** October 29, 2025  
**API Version:** 1.0  
**Database:** PostgreSQL 16  
**Framework:** Gin (Go)
