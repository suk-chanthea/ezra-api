package entity

import "time"

// ChurchMembershipStatus represents the status of a user's church membership
type ChurchMembershipStatus string

const (
	ChurchStatusPending  ChurchMembershipStatus = "pending"
	ChurchStatusApproved ChurchMembershipStatus = "approved"
	ChurchStatusRejected ChurchMembershipStatus = "rejected"
)

// User represents the core business entity
type User struct {
	ID            uint
	Username      string
	Name      string
	Profile       string
	Email         string
	EmailVerified bool                   // Whether email has been verified via OTP
	Phone         string                 // User's phone number
	Password      string
	Role          string
	Token         string
	Provider      string                 // "local", "google", etc.
	ProviderID    string                 // Google ID, Facebook ID, etc.
	Birthday      *time.Time             // User's date of birth
	ChurchID      *uint                  // Optional link to church
	Church        *Church                // Related church
	ChurchStatus  ChurchMembershipStatus // pending, approved, rejected
	BandID        *uint                  // Optional link to band
	Bio           string                 // User biography/description
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// NewUser creates a new user entity for local registration
func NewUser(username, name, email, password string) *User {
	return &User{
		Username:      username,
		Name:      name,
		Email:         email,
		EmailVerified: false, // Must verify email via OTP
		Password:      password,
		Role:          "user",
		Provider:      "local",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
}

// NewOAuthUser creates a new user entity for OAuth providers
func NewOAuthUser(email, name, provider, providerID string) *User {
	return &User{
		Username:      email, // Use email as username for OAuth users
		Name:      name,
		Email:         email,
		EmailVerified: true, // OAuth providers verify email automatically
		Password:      "",   // No password for OAuth users
		Role:          "user",
		Provider:      provider,
		ProviderID:    providerID,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
}

// IsValid validates user entity
func (u *User) IsValid() bool {
	if u.Provider == "local" {
		return u.Username != "" && u.Email != "" && u.Password != ""
	}
	// For OAuth users, password is not required
	return u.Username != "" && u.Email != "" && u.ProviderID != ""
}