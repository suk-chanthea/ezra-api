package dto

import "time"

// AuthResponse represents authentication response
type AuthResponse struct {
	Message string `json:"message,omitempty"`
	Token   string `json:"token"`
}

// MusicResponse represents music output
type MusicResponse struct {
	ID        uint      `json:"id"`
	Title     string    `json:"title"`
	Cover     string    `json:"cover"`
	Audio     string    `json:"audio"`
	UserID    uint      `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
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