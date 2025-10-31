package usecase

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"

	"github.com/suk-chanthea/ezra/domain/dto"
	"github.com/suk-chanthea/ezra/domain/entity"
	"github.com/suk-chanthea/ezra/domain/repository"
	"github.com/suk-chanthea/ezra/infrastructure/email"
)

type OTPUseCase interface {
	SendOTP(req *dto.SendOTPRequest) (*dto.OTPResponse, error)
	VerifyOTP(req *dto.VerifyOTPRequest) (*dto.SuccessResponse, error)
	ResendOTP(email string, purpose entity.OTPPurpose) (*dto.OTPResponse, error)
}

type otpUseCase struct {
	otpRepo      repository.OTPRepository
	userRepo     repository.UserRepository
	emailService email.EmailService
	otpExpiry    int // in minutes
}

func NewOTPUseCase(
	otpRepo repository.OTPRepository,
	userRepo repository.UserRepository,
	emailService email.EmailService,
	otpExpiry int,
) OTPUseCase {
	if otpExpiry <= 0 {
		otpExpiry = 10 // default 10 minutes
	}
	return &otpUseCase{
		otpRepo:      otpRepo,
		userRepo:     userRepo,
		emailService: emailService,
		otpExpiry:    otpExpiry,
	}
}

// SendOTP generates and sends an OTP to the user's email
func (uc *otpUseCase) SendOTP(req *dto.SendOTPRequest) (*dto.OTPResponse, error) {
	// Validate email format
	if !email.ValidateEmail(req.Email) {
		return nil, errors.New("invalid email format")
	}

	// For email verification, check if email already exists
	if req.Purpose == "email_verification" {
		existingUser, _ := uc.userRepo.FindByEmail(req.Email)
		if existingUser != nil {
			return nil, errors.New("email already registered, please login or reset password")
		}
	}

	// For login (2FA), ensure the email exists
	if req.Purpose == "login" {
		_, err := uc.userRepo.FindByEmail(req.Email)
		if err != nil {
			return nil, errors.New("email not found")
		}
	}

	// For password reset, check if email exists
	if req.Purpose == "password_reset" {
		_, err := uc.userRepo.FindByEmail(req.Email)
		if err != nil {
			return nil, errors.New("email not found")
		}
	}

	// Generate random 6-digit OTP
	code, err := uc.generateOTPCode()
	if err != nil {
		return nil, errors.New("failed to generate OTP")
	}

	// Delete any existing unverified OTPs for this email and purpose
	uc.otpRepo.DeleteByEmail(req.Email)

	// Create new OTP
	otp := entity.NewOTP(req.Email, code, entity.OTPPurpose(req.Purpose), uc.otpExpiry)

	// Save to database
	if err := uc.otpRepo.Save(otp); err != nil {
		return nil, errors.New("failed to save OTP")
	}

	// Send email with OTP
	if err := uc.emailService.SendOTP(req.Email, code, req.Purpose); err != nil {
		return nil, fmt.Errorf("failed to send email: %v", err)
	}

	return &dto.OTPResponse{
		Message:   "OTP sent successfully to your email",
		Email:     req.Email,
		ExpiresAt: otp.ExpiresAt,
	}, nil
}

// VerifyOTP verifies the OTP code
func (uc *otpUseCase) VerifyOTP(req *dto.VerifyOTPRequest) (*dto.SuccessResponse, error) {
	// Find OTP by email, code, and purpose
	otp, err := uc.otpRepo.FindByEmailCodeAndPurpose(req.Email, req.Code, entity.OTPPurpose(req.Purpose))
	if err != nil {
		return nil, errors.New("invalid or expired OTP")
	}

	// Check if OTP is already verified
	if otp.Verified {
		return nil, errors.New("OTP already used")
	}

	// Check if OTP is expired
	if otp.IsExpired() {
		return nil, errors.New("OTP has expired")
	}

	// Check if OTP is valid
	if !otp.IsValid() {
		return nil, errors.New("invalid OTP")
	}

	// Mark OTP as verified
	otp.MarkAsVerified()
	if err := uc.otpRepo.Update(otp); err != nil {
		return nil, errors.New("failed to verify OTP")
	}

	// Don't delete yet - OTP will be deleted after it's used in register/login/reset-password

	return &dto.SuccessResponse{
		Message: "OTP verified successfully",
		Data: map[string]interface{}{
			"email":   req.Email,
			"purpose": req.Purpose,
		},
	}, nil
}

// ResendOTP resends OTP to the user's email
func (uc *otpUseCase) ResendOTP(email string, purpose entity.OTPPurpose) (*dto.OTPResponse, error) {
	return uc.SendOTP(&dto.SendOTPRequest{
		Email:   email,
		Purpose: string(purpose),
	})
}

// generateOTPCode generates a random 6-digit OTP code
func (uc *otpUseCase) generateOTPCode() (string, error) {
	// Generate a random number between 100000 and 999999
	max := big.NewInt(900000)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}
	
	code := n.Int64() + 100000
	return fmt.Sprintf("%06d", code), nil
}

