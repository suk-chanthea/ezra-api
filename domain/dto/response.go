package dto

import "time"

// AuthResponse represents authentication response
type AuthResponse struct {
	Message string `json:"message,omitempty"`
	Token   string `json:"token"`
}

// OTPResponse represents OTP operation response
type OTPResponse struct {
	Message   string    `json:"message"`
	Email     string    `json:"email"`
	ExpiresAt time.Time `json:"expires_at,omitempty"`
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
	ID            uint            `json:"id"`
	Username      string          `json:"username"`
	Fullname      string          `json:"fullname"`
	Profile       string          `json:"profile,omitempty"`
	Email         string          `json:"email"`
	EmailVerified bool            `json:"email_verified"`
	Phone         string          `json:"phone,omitempty"`
	Role          string          `json:"role"`
	Birthday      *time.Time      `json:"birthday,omitempty"`
	ChurchID      *uint           `json:"church_id,omitempty"`
	Church        *ChurchResponse `json:"church,omitempty"`
	ChurchStatus  string          `json:"church_status,omitempty"` // pending, approved, rejected
	BandID        *uint           `json:"band_id,omitempty"`
	Bio           string          `json:"bio,omitempty"`
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
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

// DonationResponse represents donation output
type DonationResponse struct {
	ID            uint           `json:"id"`
	Type          string         `json:"type"`
	DonorType     string         `json:"donor_type"`
	UserID        *uint          `json:"user_id,omitempty"`
	SupporterID   *uint          `json:"supporter_id,omitempty"`
	CompanyName   string         `json:"company_name,omitempty"`
	CompanyEmail  string         `json:"company_email,omitempty"`
	CompanyPhone  string         `json:"company_phone,omitempty"`
	Amount        float64        `json:"amount"`
	Currency      string         `json:"currency"`
	Message       string         `json:"message,omitempty"`
	Status        string         `json:"status"`
	TransactionID string         `json:"transaction_id,omitempty"`
	PaymentMethod string         `json:"payment_method,omitempty"`
	EventID       *uint          `json:"event_id,omitempty"`
	User          *UserResponse  `json:"user,omitempty"`
	Supporter     *SupporterResponse `json:"supporter,omitempty"`
	Event         *EventResponse `json:"event,omitempty"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	
	// Payment info (included if initiate_payment=true)
	PaymentInfo   *InitiatePaymentResponse `json:"payment_info,omitempty"`
}

// DonationStatsResponse represents donation statistics
type DonationStatsResponse struct {
	TotalAmount        float64 `json:"total_amount"`
	TotalDonations     int64   `json:"total_donations"`
	TotalSponsors      int64   `json:"total_sponsors"`
	DonateAmount       float64 `json:"donate_amount"`
	SponsorAmount      float64 `json:"sponsor_amount"`
	UserDonations      int64   `json:"user_donations"`
	CompanyDonations   int64   `json:"company_donations"`
}

// InitiatePaymentResponse represents payment initiation response
type InitiatePaymentResponse struct {
	DonationID      uint   `json:"donation_id"`
	TransactionID   string `json:"transaction_id"`
	PaymentURL      string `json:"payment_url,omitempty"`      // For card payments
	QRCode          string `json:"qr_code,omitempty"`          // For QR payments (base64)
	Amount          string `json:"amount"`
	Currency        string `json:"currency"`
	PaymentMethod   string `json:"payment_method"` // "qr" or "card"
	ExpiresAt       *time.Time `json:"expires_at,omitempty"`   // QR expiration time (3 minutes)
	ExpiresInSeconds int    `json:"expires_in_seconds,omitempty"` // Seconds until expiration
	Message         string `json:"message"`
}

// SupporterResponse represents supporter output
type SupporterResponse struct {
	ID           uint          `json:"id"`
	Name         string        `json:"name"`
	Email        string        `json:"email"`
	Phone        string        `json:"phone,omitempty"`
	Type         string        `json:"type"`
	Website      string        `json:"website,omitempty"`
	Address      string        `json:"address,omitempty"`
	Logo         string        `json:"logo,omitempty"`
	Description  string        `json:"description,omitempty"`
	UserID       *uint         `json:"user_id,omitempty"`
	User         *UserResponse `json:"user,omitempty"`
	TotalDonations int         `json:"total_donations,omitempty"` // Count of donations from this supporter
	TotalAmount    float64     `json:"total_amount,omitempty"`    // Total donation amount
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`
}

// ChurchResponse represents church output
type ChurchResponse struct {
	ID              uint         `json:"id"`
	Fullname        string       `json:"fullname"`
	Address         string       `json:"address,omitempty"`
	Phone           string       `json:"phone,omitempty"`
	Email           string       `json:"email,omitempty"`
	Website         string       `json:"website,omitempty"`
	PastorName      string       `json:"pastor_name,omitempty"`
	Description     string       `json:"description,omitempty"`
	Logo            string       `json:"logo,omitempty"`
	EstablishedDate *time.Time   `json:"established_date,omitempty"`
	Denomination    string       `json:"denomination,omitempty"`
	OwnerID         *uint        `json:"owner_id,omitempty"`
	Owner           *UserResponse `json:"owner,omitempty"`
	MemberCount     int          `json:"member_count,omitempty"` // Computed: count of approved members
	PendingCount    int          `json:"pending_count,omitempty"` // Computed: count of pending requests
	CreatedAt       time.Time    `json:"created_at"`
	UpdatedAt       time.Time    `json:"updated_at"`
}