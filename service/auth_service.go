package service

import (
    "errors"
    "time"

    "github.com/suk-chanthea/ezra/domain"
    "github.com/suk-chanthea/ezra/repository"

    "github.com/golang-jwt/jwt/v5"
    "golang.org/x/crypto/bcrypt"
)

type AuthService interface {
    Register(username, fullname, email, password string) (string,error)
    Login(username, password string) (string, error)
    ValidateToken(token string) (*jwt.Token, error)
}

type authService struct {
    repo repository.UserRepository
    key  []byte
}

func NewAuthService(repo repository.UserRepository, secret string) AuthService {
    return &authService{repo: repo, key: []byte(secret)}
}

func (s *authService) Register(username, fullname, email, password string) (string, error) {
    // hash password
    hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

    user := &domain.User{
        Username: username,
        Fullname: fullname,
        Email:    email,
        Password: string(hash),
        Role:     "user",
    }

    // save user in DB first
    if err := s.repo.Create(user); err != nil {
        return "", err
    }

    // generate JWT token after registration
    claims := jwt.MapClaims{
        "sub":      user.ID,
        "username": user.Username,
        "role":     user.Role,
        "exp":      time.Now().Add(time.Hour * 1).Unix(),
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    signedToken, err := token.SignedString(s.key)
    if err != nil {
        return "", err
    }

    // save token to user
    user.Token = signedToken
    // update DB with token
    if err := s.repo.UpdateToken(user.ID, signedToken); err != nil {
        return "", err
    }

    return signedToken, nil
}


func (s *authService) Login(username, password string) (string, error) {
    user, err := s.repo.GetByUsername(username)
    if err != nil {
        return "", errors.New("user not found")
    }

    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
        return "", errors.New("invalid credentials")
    }

    claims := jwt.MapClaims{
        "sub":      user.ID,
        "username": user.Username,
        "role":     user.Role,
        "exp":      time.Now().Add(time.Hour * 1).Unix(),
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(s.key)
}

func (s *authService) ValidateToken(tokenStr string) (*jwt.Token, error) {
    return jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
        return s.key, nil
    })
}
