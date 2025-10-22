package repository

import "github.com/suk-chanthea/ezra/domain/entity"

// UserRepository defines the interface for user data operations
type UserRepository interface {
	Save(user *entity.User) error
	FindByUsername(username string) (*entity.User, error)
	FindByEmail(email string) (*entity.User, error)
	FindByID(id uint) (*entity.User, error)
	FindByProviderID(provider, providerID string) (*entity.User, error)
	Update(user *entity.User) error
	UpdateToken(id uint, token string) error
	Delete(id uint) error
}