package dto

import "time"

// AuthResponse represents authentication response
type AuthResponse struct {
	Message string `json:"message,omitempty"`
	Token   string `json:"token"`
}

// MusicResponse represents music output
type MusicResponse struct {
	ID          uint      `json:"id"`
	Title       string    `json:"title"`
	Artist      string    `json:"artist,omitempty"`
	Album       string    `json:"album,omitempty"`
	Genre       string    `json:"genre,omitempty"`
	Duration    int       `json:"duration,omitempty"`    // in seconds
	BPM         int       `json:"bpm,omitempty"`         // beats per minute
	Key         string    `json:"key,omitempty"`         // musical key
	Cover       string    `json:"cover,omitempty"`
	Lyrics      string    `json:"lyrics,omitempty"`
	Description string    `json:"description,omitempty"`
	UserID      uint      `json:"user_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// EventResponse represents event output
type EventResponse struct {
	ID        uint             `json:"id"`
	Title     string           `json:"title"`
	Content   string           `json:"content"`
	Cover     string           `json:"cover"`
	Location  string           `json:"location"`
	StartTime time.Time        `json:"start_time"`
	EndTime   time.Time        `json:"end_time"`
	UserID    uint             `json:"user_id"`
	Musics    []*MusicResponse `json:"musics,omitempty"`  // Add this field
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt time.Time        `json:"updated_at"`
}

// ErrorResponse represents error output
type ErrorResponse struct {
	Error   string `json:"error,omitempty"`
	Errors  string `json:"errors,omitempty"`
	Message string `json:"message,omitempty"`
}

// SuccessResponse represents success output
type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// UserResponse represents user output (without sensitive data)
type UserResponse struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Fullname  string    `json:"fullname"`
	Profile   string    `json:"profile,omitempty"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// BookingResponse represents booking output
type BookingResponse struct {
	ID        uint           `json:"id"`
	EventID   uint           `json:"event_id"`
	UserID    uint           `json:"user_id"`
	Status    string         `json:"status"`
	Notes     string         `json:"notes,omitempty"`
	Event     *EventResponse `json:"event,omitempty"`
	User      *UserResponse  `json:"user,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

// PaginationMetadata represents pagination information
type PaginationMetadata struct {
	CurrentPage  int   `json:"current_page"`
	PageSize     int   `json:"page_size"`
	TotalPages   int   `json:"total_pages"`
	TotalRecords int64 `json:"total_records"`
	HasNextPage  bool  `json:"has_next_page"`
	HasPrevPage  bool  `json:"has_prev_page"`
}

// PaginatedResponse represents a paginated response wrapper
type PaginatedResponse struct {
	Data       interface{}         `json:"data"`
	Pagination *PaginationMetadata `json:"pagination"`
}

// NewPaginationMetadata creates pagination metadata
func NewPaginationMetadata(page, pageSize int, totalRecords int64) *PaginationMetadata {
	totalPages := int((totalRecords + int64(pageSize) - 1) / int64(pageSize))
	if totalPages < 1 {
		totalPages = 1
	}
	
	return &PaginationMetadata{
		CurrentPage:  page,
		PageSize:     pageSize,
		TotalPages:   totalPages,
		TotalRecords: totalRecords,
		HasNextPage:  page < totalPages,
		HasPrevPage:  page > 1,
	}
}

// BandResponse represents band output
type BandResponse struct {
	ID           uint             `json:"id"`
	Name         string           `json:"name"`
	Description  string           `json:"description"`
	Cover        string           `json:"cover"`
	IsPublic     bool             `json:"is_public"`
	UserID       uint             `json:"user_id"`
	MemberCount  int64            `json:"member_count,omitempty"`
	MusicCount   int              `json:"music_count,omitempty"`
	Musics       []*MusicResponse `json:"musics,omitempty"`
	Members      []*UserResponse  `json:"members,omitempty"`
	CreatedAt    time.Time        `json:"created_at"`
	UpdatedAt    time.Time        `json:"updated_at"`
}

// SettingResponse represents user settings output
type SettingResponse struct {
	ID                     uint      `json:"id"`
	UserID                 uint      `json:"user_id"`
	Language               string    `json:"language"`
	Theme                  string    `json:"theme"`
	NotifyOnBooking        bool      `json:"notify_on_booking"`
	NotifyOnMusic          bool      `json:"notify_on_music"`
	NotifyOnEvent          bool      `json:"notify_on_event"`
	EnablePushNotifications bool      `json:"enable_push_notifications"`
	CreatedAt              time.Time `json:"created_at"`
	UpdatedAt              time.Time `json:"updated_at"`
}

// NotificationResponse represents notification output
type NotificationResponse struct {
	ID            uint       `json:"id"`
	UserID        *uint      `json:"user_id,omitempty"`
	SenderID      *uint      `json:"sender_id,omitempty"`
	BandID        *uint      `json:"band_id,omitempty"`
	RecipientType string     `json:"recipient_type"` // user, band, all
	Title         string     `json:"title"`
	Message       string     `json:"message"`
	Type          string     `json:"type"`
	RelatedType   string     `json:"related_type,omitempty"`
	RelatedID     *uint      `json:"related_id,omitempty"`
	IsRead        bool       `json:"is_read"`
	ReadAt        *time.Time `json:"read_at,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}