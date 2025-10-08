package repository

import "github.com/suk-chanthea/ezra/domain/entity"

// MusicRepository defines the interface for music data operations
type MusicRepository interface {
	Save(music *entity.Music) error
	FindAll() ([]*entity.Music, error)
	FindByID(id uint) (*entity.Music, error)
	FindByIDs(ids []uint) ([]*entity.Music, error)
	FindByUserID(userID uint) ([]*entity.Music, error)
	Update(music *entity.Music) error
	Delete(id uint) error
}