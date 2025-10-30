package repository

import "github.com/suk-chanthea/ezra/domain/entity"

type SupporterRepository interface {
	Create(supporter *entity.Supporter) error
	FindByID(id uint) (*entity.Supporter, error)
	FindByEmail(email string) (*entity.Supporter, error)
	FindAll(limit, offset int) ([]*entity.Supporter, error)
	FindByType(supporterType entity.SupporterType, limit, offset int) ([]*entity.Supporter, error)
	FindByUser(userID uint, limit, offset int) ([]*entity.Supporter, error)
	Update(supporter *entity.Supporter) error
	Delete(id uint) error
	Count() (int64, error)
	CountByType(supporterType entity.SupporterType) (int64, error)
}

