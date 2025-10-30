package repository

import "github.com/suk-chanthea/ezra/domain/entity"

type ChurchRepository interface {
	Create(church *entity.Church) error
	FindByID(id uint) (*entity.Church, error)
	FindByName(name string) (*entity.Church, error)
	FindAll(limit, offset int) ([]*entity.Church, error)
	FindByDenomination(denomination string, limit, offset int) ([]*entity.Church, error)
	Update(church *entity.Church) error
	Delete(id uint) error
	Count() (int64, error)
	CountMembers(churchID uint, status string) (int64, error) // Count members by status
	FindMembers(churchID uint, status string, limit, offset int) ([]*entity.User, error) // Get members by status
}

