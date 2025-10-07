package entity

import "time"

// User represents the core business entity
type User struct {
	ID        uint
	Username  string
	Fullname  string
	Profile   string
	Email     string
	Password  string
	Role      string
	Token     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewUser creates a new user entity
func NewUser(username, fullname, email, password string) *User {
	return &User{
		Username:  username,
		Fullname:  fullname,
		Email:     email,
		Password:  password,
		Role:      "user",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// IsValid validates user entity
func (u *User) IsValid() bool {
	return u.Username != "" && u.Email != "" && u.Password != ""
}