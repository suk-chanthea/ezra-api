package repository

import "github.com/suk-chanthea/ezra/domain/entity"

// SettingRepository defines the interface for setting data operations
type SettingRepository interface {
	Save(setting *entity.Setting) error
	FindByUserID(userID uint) (*entity.Setting, error)
	Update(setting *entity.Setting) error
	Delete(userID uint) error
}

