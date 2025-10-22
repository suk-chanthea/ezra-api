package entity

import "time"

// User represents the core business entity
type User struct {
	ID         uint
	Username   string
	Fullname   string
	Profile    string
	Email      string
	Password   string
	Role       string
	Token      string
	Provider   string // "local", "google", etc.
	ProviderID string // Google ID, Facebook ID, etc.
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// NewUser creates a new user entity for local registration
func NewUser(username, fullname, email, password string) *User {
	return &User{
		Username:  username,
		Fullname:  fullname,
		Email:     email,
		Password:  password,
		Role:      "user",
		Provider:  "local",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// NewOAuthUser creates a new user entity for OAuth providers
func NewOAuthUser(email, fullname, provider, providerID string) *User {
	return &User{
		Username:   email, // Use email as username for OAuth users
		Fullname:   fullname,
		Email:      email,
		Password:   "", // No password for OAuth users
		Role:       "user",
		Provider:   provider,
		ProviderID: providerID,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
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