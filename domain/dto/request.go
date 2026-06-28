package dto

import "time"

// RegisterRequest represents registration input
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=100"`
	Fullname string `json:"fullname" binding:"required,min=1,max=100"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	// OTPCode  string `json:"otp_code" binding:"required,min=6,max=6"` // Must verify email via OTP before registration
}

// LoginRequest represents login input
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	// OTPCode  string `json:"otp_code,omitempty" binding:"omitempty,min=6,max=6"` // Optional 2FA OTP code
}

// GoogleLoginRequest represents Google OAuth login input
type GoogleLoginRequest struct {
	IDToken string `json:"id_token" binding:"required"`
}

// SendOTPRequest represents OTP generation request
type SendOTPRequest struct {
	Email   string `json:"email" binding:"required,email"`
	Purpose string `json:"purpose" binding:"required,oneof=email_verification password_reset login"`
}

// VerifyOTPRequest represents OTP verification request
type VerifyOTPRequest struct {
	Email   string `json:"email" binding:"required,email"`
	Code    string `json:"code" binding:"required,min=6,max=6"`
	Purpose string `json:"purpose" binding:"required,oneof=email_verification password_reset login"`
}

// ResetPasswordRequest represents password reset request (after OTP verification)
type ResetPasswordRequest struct {
	Email       string `json:"email" binding:"required,email"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
	OTPCode     string `json:"otp_code" binding:"required,min=6,max=6"` // Must provide OTP for verification
}

// UpdateProfileRequest represents user profile update input
type UpdateProfileRequest struct {
	Username string `json:"username" binding:"required,min=3,max=100"`
	Fullname string `json:"fullname" binding:"required,min=1,max=100"`
	Profile  string `json:"profile"`
	Phone    string `json:"phone" binding:"omitempty,max=50"`
	Birthday string `json:"birthday"` // Format: YYYY-MM-DD
	Bio      string `json:"bio"`
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

// CreateBandRequest represents band creation input
type CreateBandRequest struct {
	Name        string `json:"name" binding:"required,min=1,max=255"`
	Description string `json:"description"`
	Cover       string `json:"cover"`
	IsPublic    bool   `json:"is_public"`
}

// UpdateBandRequest represents band update input
type UpdateBandRequest struct {
	Name        string `json:"name" binding:"required,min=1,max=255"`
	Description string `json:"description"`
	Cover       string `json:"cover"`
	IsPublic    bool   `json:"is_public"`
}

// AddMusicsRequest represents adding music to band input
type AddMusicsRequest struct {
	MusicIDs []uint `json:"music_ids" binding:"required,min=1"`
}

// ReorderMusicsRequest represents reordering music in band
type ReorderMusicsRequest struct {
	MusicOrders []MusicOrder `json:"music_orders" binding:"required,min=1"`
}

// MusicOrder represents a music ID and its display order
type MusicOrder struct {
	MusicID      uint `json:"music_id" binding:"required"`
	DisplayOrder int  `json:"display_order"`
}

// UpdateSettingRequest represents settings update input
type UpdateSettingRequest struct {
	Language               string `json:"language" binding:"required,oneof=en kh kr cn"`
	Theme                  string `json:"theme" binding:"required,oneof=light dark auto"`
	NotifyOnBooking        bool   `json:"notify_on_booking"`
	NotifyOnMusic          bool   `json:"notify_on_music"`
	NotifyOnEvent          bool   `json:"notify_on_event"`
	EnablePushNotifications bool   `json:"enable_push_notifications"`
}

// CreateNotificationRequest represents notification creation input
type CreateNotificationRequest struct {
	UserID      *uint  `json:"user_id"`      // For single user (required if recipient_type=user)
	BandID      *uint  `json:"band_id"`      // For band/team (required if recipient_type=band)
	Title       string `json:"title" binding:"required,min=1,max=255"`
	Message     string `json:"message" binding:"required,min=1"`
	Type        string `json:"type" binding:"required,oneof=info success warning error booking music event"`
	RelatedType string `json:"related_type" binding:"omitempty,oneof=music event booking band"`
	RelatedID   *uint  `json:"related_id"`
}

// CreateDonationRequest represents donation creation input
type CreateDonationRequest struct {
	Type           string  `json:"type" binding:"required,oneof=donate sponsor"`
	DonorType      string  `json:"donor_type" binding:"required,oneof=user company organization church"`
	CompanyName    string  `json:"company_name"`  // Required if donor_type=company, organization, or church
	CompanyEmail   string  `json:"company_email"` // Required if donor_type=company, organization, or church
	CompanyPhone   string  `json:"company_phone"`
	Amount         float64 `json:"amount" binding:"required,gt=0"`
	Currency       string  `json:"currency" binding:"required,oneof=USD KHR"`
	Message        string  `json:"message"`
	EventID        *uint   `json:"event_id"` // Optional: If provided, donation is for joining this event. If null, donation is for the app
	InitiatePayment bool   `json:"initiate_payment"` // If true, returns payment info (QR for donate, Visa for sponsor) immediately
}

// UpdateDonationStatusRequest represents donation status update input
type UpdateDonationStatusRequest struct {
	Status        string `json:"status" binding:"required,oneof=pending completed failed refunded"`
	TransactionID string `json:"transaction_id"`
	PaymentMethod string `json:"payment_method"`
}

// DonationFilterRequest represents donation filter parameters
type DonationFilterRequest struct {
	Type      string `form:"type" binding:"omitempty,oneof=donate sponsor"`
	DonorType string `form:"donor_type" binding:"omitempty,oneof=user company organization church"`
	Status    string `form:"status" binding:"omitempty,oneof=pending completed failed refunded"`
	EventID   *uint  `form:"event_id"`
	PaginationRequest
}

// CreateSupporterRequest represents supporter creation input
type CreateSupporterRequest struct {
	Name        string `json:"name" binding:"required,min=1,max=255"`
	Email       string `json:"email" binding:"required,email,max=255"`
	Phone       string `json:"phone" binding:"omitempty,max=50"`
	Type        string `json:"type" binding:"required,oneof=company organization church"`
	Website     string `json:"website" binding:"omitempty,max=255"`
	Address     string `json:"address"`
	Logo        string `json:"logo" binding:"omitempty,max=255"`
	Description string `json:"description"`
	SupporterID *uint  `json:"supporter_id"` // Optional: For linking donation to existing supporter
}

// UpdateSupporterRequest represents supporter update input
type UpdateSupporterRequest struct {
	Name        string `json:"name" binding:"required,min=1,max=255"`
	Email       string `json:"email" binding:"required,email,max=255"`
	Phone       string `json:"phone" binding:"omitempty,max=50"`
	Type        string `json:"type" binding:"required,oneof=company organization church"`
	Website     string `json:"website" binding:"omitempty,max=255"`
	Address     string `json:"address"`
	Logo        string `json:"logo" binding:"omitempty,max=255"`
	Description string `json:"description"`
}

// CreateChurchRequest represents church creation input
type CreateChurchRequest struct {
	Fullname        string `json:"fullname" binding:"required,min=1,max=255"`
	Address         string `json:"address"`
	Phone           string `json:"phone" binding:"omitempty,max=50"`
	Email           string `json:"email" binding:"omitempty,email,max=255"`
	Website         string `json:"website" binding:"omitempty,max=255"`
	PastorName      string `json:"pastor_name" binding:"omitempty,max=255"`
	Description     string `json:"description"`
	Logo            string `json:"logo" binding:"omitempty,max=255"`
	EstablishedDate string `json:"established_date"` // Format: YYYY-MM-DD
	Denomination    string `json:"denomination" binding:"omitempty,max=100"`
}

// JoinChurchRequest represents user request to join a church
type JoinChurchRequest struct {
	ChurchID uint `json:"church_id" binding:"required"`
}

// ApproveChurchMemberRequest represents approval/rejection of church membership
type ApproveChurchMemberRequest struct {
	UserID uint   `json:"user_id" binding:"required"`
	Status string `json:"status" binding:"required,oneof=approved rejected"`
}

// UpdateChurchRequest represents church update input
type UpdateChurchRequest struct {
	Fullname        string `json:"fullname" binding:"required,min=1,max=255"`
	Address         string `json:"address"`
	Phone           string `json:"phone" binding:"omitempty,max=50"`
	Email           string `json:"email" binding:"omitempty,email,max=255"`
	Website         string `json:"website" binding:"omitempty,max=255"`
	PastorName      string `json:"pastor_name" binding:"omitempty,max=255"`
	Description     string `json:"description"`
	Logo            string `json:"logo" binding:"omitempty,max=255"`
	EstablishedDate string `json:"established_date"` // Format: YYYY-MM-DD
	Denomination    string `json:"denomination" binding:"omitempty,max=100"`
}