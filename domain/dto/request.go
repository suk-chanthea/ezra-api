package dto

import "time"

// RegisterRequest represents registration input
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=100"`
	Fullname string `json:"fullname" binding:"required,min=1,max=100"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// LoginRequest represents login input
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// CreateEventRequest represents event creation input
type CreateEventRequest struct {
	Title     string    `json:"title" binding:"required,min=1,max=200"`
	Content   string    `json:"content"`
	Cover     string    `json:"cover" binding:"required,min=1,max=200"`
	Location  string    `json:"location" binding:"required,min=1,max=200"`
	StartTime time.Time `json:"start_time" binding:"required"`
	EndTime   time.Time `json:"end_time" binding:"required,gtfield=StartTime"`
	MusicIDs  []uint    `json:"music_ids"`  // Add this field
}