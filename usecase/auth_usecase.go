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
	ValidateToken(token string) (*jwt.Token, error)
}

type authUseCase struct {
	userRepo  repository.UserRepository
	secretKey []byte
}

func NewAuthUseCase(repo repository.UserRepository, secret string) AuthUseCase {
	return &authUseCase{
		userRepo:  repo,
		secretKey: []byte(secret),
	}
}

func (uc *authUseCase) Register(req *dto.RegisterRequest) (*dto.AuthResponse, error) {
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

func (uc *authUseCase) ValidateToken(tokenStr string) (*jwt.Token, error) {
	return jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return uc.secretKey, nil
	})
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