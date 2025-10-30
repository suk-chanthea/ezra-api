package entity

import "time"

// OTPPurpose represents the purpose of the OTP
type OTPPurpose string

const (
	OTPPurposeEmailVerification OTPPurpose = "email_verification"
	OTPPurposePasswordReset     OTPPurpose = "password_reset"
	OTPPurposeLogin             OTPPurpose = "login"
)

// OTP represents an OTP code for verification
type OTP struct {
	ID        uint
	Email     string
	Code      string
	Purpose   OTPPurpose
	ExpiresAt time.Time
	Verified  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewOTP creates a new OTP entity
func NewOTP(email, code string, purpose OTPPurpose, expiryMinutes int) *OTP {
	return &OTP{
		Email:     email,
		Code:      code,
		Purpose:   purpose,
		ExpiresAt: time.Now().Add(time.Duration(expiryMinutes) * time.Minute),
		Verified:  false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// IsExpired checks if the OTP has expired
func (o *OTP) IsExpired() bool {
	return time.Now().After(o.ExpiresAt)
}

// IsValid checks if the OTP is valid for verification
func (o *OTP) IsValid() bool {
	return !o.IsExpired() && !o.Verified
}

// MarkAsVerified marks the OTP as verified
func (o *OTP) MarkAsVerified() {
	o.Verified = true
	o.UpdatedAt = time.Now()
}

