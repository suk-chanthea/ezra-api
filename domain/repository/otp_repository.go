package repository

import "github.com/suk-chanthea/ezra/domain/entity"

// OTPRepository defines the interface for OTP data operations
type OTPRepository interface {
	Save(otp *entity.OTP) error
	FindByEmailAndPurpose(email string, purpose entity.OTPPurpose) (*entity.OTP, error)
	FindByEmailCodeAndPurpose(email, code string, purpose entity.OTPPurpose) (*entity.OTP, error)
	Update(otp *entity.OTP) error
	DeleteByEmail(email string) error
	DeleteByEmailAndPurpose(email string, purpose entity.OTPPurpose) error
	DeleteExpired() error
}

