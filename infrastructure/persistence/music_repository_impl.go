package persistence

import (
	"time"

	"github.com/suk-chanthea/ezra/domain/entity"
	"github.com/suk-chanthea/ezra/domain/repository"
	"gorm.io/gorm"
)

// MusicModel is the GORM model for database
type MusicModel struct {
	ID          uint      `gorm:"primaryKey"`
	Title       string    `gorm:"size:255;not null"`
	Artist      string    `gorm:"size:255"`
	Album       string    `gorm:"size:255"`
	Genre       string    `gorm:"size:100"`
	Duration    int       `gorm:"type:integer"` // in seconds
	BPM         int       `gorm:"type:integer"` // beats per minute
	Key         string    `gorm:"size:10"`      // musical key
	Cover       string    `gorm:"size:255"`
	Lyrics      string    `gorm:"type:text"`
	Description string    `gorm:"type:text"`
	UserID      uint      `gorm:"not null"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
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

func (r *musicRepositoryImpl) FindAllPaginated(offset, limit int) ([]*entity.Music, int64, error) {
	var models []MusicModel
	var total int64
	
	// Get total count
	if err := r.db.Model(&MusicModel{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	// Get paginated results
	if err := r.db.Offset(offset).Limit(limit).Find(&models).Error; err != nil {
		return nil, 0, err
	}
	
	return r.modelsToEntities(models), total, nil
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
		ID:          music.ID,
		Title:       music.Title,
		Artist:      music.Artist,
		Album:       music.Album,
		Genre:       music.Genre,
		Duration:    music.Duration,
		BPM:         music.BPM,
		Key:         music.Key,
		Cover:       music.Cover,
		Lyrics:      music.Lyrics,
		Description: music.Description,
		UserID:      music.UserID,
		CreatedAt:   music.CreatedAt,
		UpdatedAt:   music.UpdatedAt,
	}
}

func (r *musicRepositoryImpl) modelToEntity(model *MusicModel) *entity.Music {
	return &entity.Music{
		ID:          model.ID,
		Title:       model.Title,
		Artist:      model.Artist,
		Album:       model.Album,
		Genre:       model.Genre,
		Duration:    model.Duration,
		BPM:         model.BPM,
		Key:         model.Key,
		Cover:       model.Cover,
		Lyrics:      model.Lyrics,
		Description: model.Description,
		UserID:      model.UserID,
		CreatedAt:   model.CreatedAt,
		UpdatedAt:   model.UpdatedAt,
	}
}

func (r *musicRepositoryImpl) modelsToEntities(models []MusicModel) []*entity.Music {
	entities := make([]*entity.Music, len(models))
	for i, model := range models {
		entities[i] = r.modelToEntity(&model)
	}
	return entities
}