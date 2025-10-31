package usecase

import (
	"errors"
	"time"

	"github.com/suk-chanthea/ezra/domain/dto"
	"github.com/suk-chanthea/ezra/domain/entity"
	"github.com/suk-chanthea/ezra/domain/repository"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthUseCase interface {
	Register(req *dto.RegisterRequest) (*dto.AuthResponse, error)
	Login(req *dto.LoginRequest) (*dto.AuthResponse, error)
	Logout(userID uint) error
	DeleteUser(userID uint) error
	GoogleLogin(googleID, email, fullname, profilePicture string) (*dto.AuthResponse, error)
	ResetPassword(req *dto.ResetPasswordRequest) (*dto.SuccessResponse, error)
	ValidateToken(token string) (*jwt.Token, error)
	VerifyTokenInDatabase(userID uint, token string) error
}

type authUseCase struct {
	userRepo       repository.UserRepository
	otpRepo        repository.OTPRepository
	secretKey      []byte
	googleClientID string
}

func NewAuthUseCase(repo repository.UserRepository, otpRepo repository.OTPRepository, secret, googleClientID string) AuthUseCase {
	return &authUseCase{
		userRepo:       repo,
		otpRepo:        otpRepo,
		secretKey:      []byte(secret),
		googleClientID: googleClientID,
	}
}

func (uc *authUseCase) Register(req *dto.RegisterRequest) (*dto.AuthResponse, error) {
	// Verify OTP first (email_verification purpose)
	otp, err := uc.otpRepo.FindByEmailCodeAndPurpose(req.Email, req.OTPCode, entity.OTPPurpose("email_verification"))
	if err != nil {
		return nil, errors.New("invalid OTP code")
	}

	// Check if OTP is verified and not expired
	if !otp.Verified {
		return nil, errors.New("OTP not verified. Please verify OTP first via /otp/verify")
	}

	if otp.IsExpired() {
		return nil, errors.New("OTP has expired. Please request a new one")
	}

	// Check if user already exists
	existing, _ := uc.userRepo.FindByUsername(req.Username)
	if existing != nil {
		return nil, errors.New("username already exists")
	}

	existing, _ = uc.userRepo.FindByEmail(req.Email)
	if existing != nil {
		return nil, errors.New("email already exists")
	}

	// Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create entity
	user := entity.NewUser(req.Username, req.Fullname, req.Email, string(hash))
	
	// Mark email as verified since OTP was verified
	user.EmailVerified = true

	// Save to database
	if err := uc.userRepo.Save(user); err != nil {
		return nil, err
	}

	// Generate JWT
	token, err := uc.generateToken(user)
	if err != nil {
		return nil, err
	}

	// Update user with token
	if err := uc.userRepo.UpdateToken(user.ID, token); err != nil {
		return nil, err
	}

	// Delete the used OTP (only for this purpose)
	go uc.otpRepo.DeleteByEmailAndPurpose(req.Email, entity.OTPPurpose("email_verification"))

	return &dto.AuthResponse{
		Message: "user registered successfully",
		Token:   token,
	}, nil
}

func (uc *authUseCase) Login(req *dto.LoginRequest) (*dto.AuthResponse, error) {
	// Find user
	user, err := uc.userRepo.FindByUsername(req.Username)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Optional 2FA: If OTP code is provided, verify it
	if req.OTPCode != "" {
		otp, err := uc.otpRepo.FindByEmailCodeAndPurpose(user.Email, req.OTPCode, entity.OTPPurpose("login"))
		if err != nil {
			return nil, errors.New("invalid 2FA OTP code")
		}

		// Check if OTP is verified and not expired
		if !otp.Verified {
			return nil, errors.New("OTP not verified. Please verify OTP first via /otp/verify")
		}

		if otp.IsExpired() {
			return nil, errors.New("OTP has expired. Please request a new one")
		}

		// Delete the used OTP (only for this purpose)
		go uc.otpRepo.DeleteByEmailAndPurpose(user.Email, entity.OTPPurpose("login"))
	}

	// Generate JWT
	token, err := uc.generateToken(user)
	if err != nil {
		return nil, err
	}

	// Update token in database
	if err := uc.userRepo.UpdateToken(user.ID, token); err != nil {
		return nil, err
	}

	return &dto.AuthResponse{
		Token: token,
	}, nil
}

func (uc *authUseCase) Logout(userID uint) error {
	// Clear the user's token
	return uc.userRepo.UpdateToken(userID, "")
}

func (uc *authUseCase) DeleteUser(userID uint) error {
	// Check if user exists
	user, err := uc.userRepo.FindByID(userID)
	if err != nil {
		return errors.New("user not found")
	}

	// Delete the user
	if err := uc.userRepo.Delete(user.ID); err != nil {
		return errors.New("failed to delete user")
	}

	return nil
}

func (uc *authUseCase) GoogleLogin(googleID, email, fullname, profilePicture string) (*dto.AuthResponse, error) {
	// Check if user already exists with this Google ID
	user, err := uc.userRepo.FindByProviderID("google", googleID)

	if err != nil {
		// User doesn't exist, check if email is already registered
		existingUser, _ := uc.userRepo.FindByEmail(email)
		if existingUser != nil {
			// Email exists but with different provider
			if existingUser.Provider != "google" {
				return nil, errors.New("email already registered with different provider")
			}
		}

		// Create new user with Google OAuth
		user = entity.NewOAuthUser(email, fullname, "google", googleID)
		user.Profile = profilePicture

		// Save to database
		if err := uc.userRepo.Save(user); err != nil {
			return nil, err
		}
	} else {
		// User exists - update their information from Google
		user.Fullname = fullname
		user.Profile = profilePicture
		user.Email = email // In case Google email changed
		user.UpdatedAt = time.Now()

		// Update user data in database
		if err := uc.userRepo.Update(user); err != nil {
			return nil, err
		}
	}

	// Generate JWT
	token, err := uc.generateToken(user)
	if err != nil {
		return nil, err
	}

	// Update token in database
	if err := uc.userRepo.UpdateToken(user.ID, token); err != nil {
		return nil, err
	}

	return &dto.AuthResponse{
		Message: "google login successful",
		Token:   token,
	}, nil
}

func (uc *authUseCase) ResetPassword(req *dto.ResetPasswordRequest) (*dto.SuccessResponse, error) {
	// First, verify the OTP is valid and verified
	otp, err := uc.otpRepo.FindByEmailCodeAndPurpose(req.Email, req.OTPCode, entity.OTPPurpose("password_reset"))
	if err != nil {
		return nil, errors.New("invalid OTP code")
	}

	// Check if OTP is verified and not expired
	if !otp.Verified {
		return nil, errors.New("OTP not verified. Please verify OTP first")
	}

	if otp.IsExpired() {
		return nil, errors.New("OTP has expired. Please request a new one")
	}

	// Find user by email
	user, err := uc.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Hash new password
	hash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	// Update password
	user.Password = string(hash)
	user.UpdatedAt = time.Now()

	if err := uc.userRepo.Update(user); err != nil {
		return nil, errors.New("failed to update password")
	}

	// Invalidate all existing tokens (logout all sessions)
	uc.userRepo.UpdateToken(user.ID, "")

	// Delete the used OTP (only for this purpose)
	go uc.otpRepo.DeleteByEmailAndPurpose(req.Email, entity.OTPPurpose("password_reset"))

	return &dto.SuccessResponse{
		Message: "password reset successfully",
	}, nil
}

func (uc *authUseCase) ValidateToken(tokenStr string) (*jwt.Token, error) {
	return jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return uc.secretKey, nil
	})
}

func (uc *authUseCase) VerifyTokenInDatabase(userID uint, token string) error {
	// Find user by ID
	user, err := uc.userRepo.FindByID(userID)
	if err != nil {
		return errors.New("user not found")
	}

	// Check if token matches the one in database
	if user.Token == "" {
		return errors.New("user has been logged out")
	}

	if user.Token != token {
		return errors.New("token has been invalidated")
	}

	return nil
}

func (uc *authUseCase) generateToken(user *entity.User) (string, error) {
	claims := jwt.MapClaims{
		"sub":      user.ID,
		"username": user.Username,
		"role":     user.Role,
		"exp":      time.Now().AddDate(0, 3, 0).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(uc.secretKey)
}
