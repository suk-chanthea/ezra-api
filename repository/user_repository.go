package repository

import (
    "github.com/suk-chanthea/ezra/domain"
    "gorm.io/gorm"
)

type UserRepository interface {
    GetByUsername(username string) (*domain.User, error)
    Create(user *domain.User) error
	UpdateToken(id uint, token string) error
}

type userRepository struct {
    db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
    return &userRepository{db}
}

func (r *userRepository) UpdateToken(id uint, token string) error {
    return r.db.Model(&domain.User{}).
        Where("id = ?", id).
        Updates(map[string]interface{}{
            "token":      token,
            "updated_at": gorm.Expr("NOW()"),
        }).Error
}


func (r *userRepository) GetByUsername(username string) (*domain.User, error) {
    var user domain.User
    if err := r.db.Where("username = ?", username).First(&user).Error; err != nil {
        return nil, err
    }
    return &user, nil
}

func (r *userRepository) Create(user *domain.User) error {
    return r.db.Create(user).Error
}
