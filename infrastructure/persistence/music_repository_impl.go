package persistence

import (
	"time"

	"github.com/suk-chanthea/ezra/domain/entity"
	"github.com/suk-chanthea/ezra/domain/repository"
	"gorm.io/gorm"
)

// MusicModel is the GORM model for database
type MusicModel struct {
	ID        uint      `gorm:"primaryKey"`
	Title     string    `gorm:"size:255;not null"`
	Cover     string    `gorm:"size:255;not null"`
	Audio     string    `gorm:"size:255"`
	UserID    uint      `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func (MusicModel) TableName() string {
	return "musics"
}

type musicRepositoryImpl struct {
	db *gorm.DB
}

func NewMusicRepository(db *gorm.DB) repository.MusicRepository {
	return &musicRepositoryImpl{db: db}
}

func (r *musicRepositoryImpl) Save(music *entity.Music) error {
	model := r.entityToModel(music)
	if err := r.db.Create(&model).Error; err != nil {
		return err
	}
	music.ID = model.ID
	music.CreatedAt = model.CreatedAt
	music.UpdatedAt = model.UpdatedAt
	return nil
}

func (r *musicRepositoryImpl) FindAll() ([]*entity.Music, error) {
	var models []MusicModel
	if err := r.db.Find(&models).Error; err != nil {
		return nil, err
	}
	return r.modelsToEntities(models), nil
}

func (r *musicRepositoryImpl) FindByID(id uint) (*entity.Music, error) {
	var model MusicModel
	if err := r.db.First(&model, id).Error; err != nil {
		return nil, err
	}
	return r.modelToEntity(&model), nil
}

func (r *musicRepositoryImpl) FindByIDs(ids []uint) ([]*entity.Music, error) {
	var models []MusicModel
	if err := r.db.Where("id IN ?", ids).Find(&models).Error; err != nil {
		return nil, err
	}
	return r.modelsToEntities(models), nil
}

func (r *musicRepositoryImpl) FindByUserID(userID uint) ([]*entity.Music, error) {
	var models []MusicModel
	if err := r.db.Where("user_id = ?", userID).Find(&models).Error; err != nil {
		return nil, err
	}
	return r.modelsToEntities(models), nil
}

func (r *musicRepositoryImpl) Update(music *entity.Music) error {
	model := r.entityToModel(music)
	return r.db.Save(&model).Error
}

func (r *musicRepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&MusicModel{}, id).Error
}

func (r *musicRepositoryImpl) entityToModel(music *entity.Music) *MusicModel {
	return &MusicModel{
		ID:        music.ID,
		Title:     music.Title,
		Cover:     music.Cover,
		Audio:     music.Audio,
		UserID:    music.UserID,
		CreatedAt: music.CreatedAt,
		UpdatedAt: music.UpdatedAt,
	}
}

func (r *musicRepositoryImpl) modelToEntity(model *MusicModel) *entity.Music {
	return &entity.Music{
		ID:        model.ID,
		Title:     model.Title,
		Cover:     model.Cover,
		Audio:     model.Audio,
		UserID:    model.UserID,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}
}

func (r *musicRepositoryImpl) modelsToEntities(models []MusicModel) []*entity.Music {
	entities := make([]*entity.Music, len(models))
	for i, model := range models {
		entities[i] = r.modelToEntity(&model)
	}
	return entities
}