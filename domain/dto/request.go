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

// UpdateEventRequest represents event update input
type UpdateEventRequest struct {
	Title     string    `json:"title" binding:"required,min=1,max=200"`
	Content   string    `json:"content"`
	Cover     string    `json:"cover" binding:"required,min=1,max=200"`
	Location  string    `json:"location" binding:"required,min=1,max=200"`
	StartTime time.Time `json:"start_time" binding:"required"`
	EndTime   time.Time `json:"end_time" binding:"required,gtfield=StartTime"`
	MusicIDs  []uint    `json:"music_ids"`  // Add this field
}

// CreateBookingRequest represents booking creation input
type CreateBookingRequest struct {
	EventID uint   `json:"event_id" binding:"required"`
	Notes   string `json:"notes"`
}

// UpdateBookingRequest represents booking update input
type UpdateBookingRequest struct {
	Status string `json:"status" binding:"required,oneof=pending confirmed cancelled"`
	Notes  string `json:"notes"`
}

// GoogleLoginRequest represents Google OAuth login input
type GoogleLoginRequest struct {
	GoogleID       string `json:"google_id" binding:"required"`
	Email          string `json:"email" binding:"required,email"`
	Fullname       string `json:"fullname" binding:"required"`
	ProfilePicture string `json:"profile_picture"`
}

// PaginationRequest represents pagination parameters
type PaginationRequest struct {
	Page     int `form:"page" binding:"omitempty,min=1"`
	PageSize int `form:"page_size" binding:"omitempty,min=1,max=100"`
}

// GetPage returns the page number, defaults to 1
func (p *PaginationRequest) GetPage() int {
	if p.Page < 1 {
		return 1
	}
	return p.Page
}

// GetPageSize returns the page size, defaults to 20
func (p *PaginationRequest) GetPageSize() int {
	if p.PageSize < 1 {
		return 20
	}
	if p.PageSize > 100 {
		return 100
	}
	return p.PageSize
}

// GetOffset calculates the offset for database queries
func (p *PaginationRequest) GetOffset() int {
	return (p.GetPage() - 1) * p.GetPageSize()
}