package persistence

import (
	"context"
	"time"

	"github.com/suk-chanthea/ezra/domain/entity"
	"github.com/suk-chanthea/ezra/domain/repository"
	"gorm.io/gorm"
)

// BandModel is the GORM model for database
type BandModel struct {
	ID          uint      `gorm:"primaryKey"`
	Name        string    `gorm:"size:255;not null"`
	Description string    `gorm:"type:text"`
	Cover       string    `gorm:"size:255"`
	IsPublic    bool      `gorm:"default:false"`
	UserID      uint      `gorm:"not null"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}

func (BandModel) TableName() string {
	return "bands"
}

// BandMusicModel represents the junction table
type BandMusicModel struct {
	ID           uint      `gorm:"primaryKey"`
	BandID       uint      `gorm:"not null"`
	MusicID      uint      `gorm:"not null"`
	DisplayOrder int       `gorm:"default:0"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
}

func (BandMusicModel) TableName() string {
	return "band_musics"
}

type bandRepositoryImpl struct {
	db *gorm.DB
}

func NewBandRepository(db *gorm.DB) repository.BandRepository {
	return &bandRepositoryImpl{db: db}
}

// Basic CRUD operations

func (r *bandRepositoryImpl) Save(ctx context.Context, band *entity.Band) error {
	model := r.entityToModel(band)
	if err := r.db.WithContext(ctx).Create(&model).Error; err != nil {
		return err
	}
	band.ID = model.ID
	band.CreatedAt = model.CreatedAt
	band.UpdatedAt = model.UpdatedAt
	return nil
}

func (r *bandRepositoryImpl) FindAll(ctx context.Context) ([]*entity.Band, error) {
	var models []BandModel
	if err := r.db.WithContext(ctx).Find(&models).Error; err != nil {
		return nil, err
	}
	return r.modelsToEntities(models), nil
}

func (r *bandRepositoryImpl) FindAllPaginated(ctx context.Context, offset, limit int) ([]*entity.Band, int64, error) {
	var models []BandModel
	var total int64
	
	// Get total count
	if err := r.db.WithContext(ctx).Model(&BandModel{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	// Get paginated results
	if err := r.db.WithContext(ctx).Offset(offset).Limit(limit).Find(&models).Error; err != nil {
		return nil, 0, err
	}
	
	return r.modelsToEntities(models), total, nil
}

func (r *bandRepositoryImpl) FindByID(ctx context.Context, id uint) (*entity.Band, error) {
	var model BandModel
	if err := r.db.WithContext(ctx).First(&model, id).Error; err != nil {
		return nil, err
	}
	return r.modelToEntity(&model), nil
}

func (r *bandRepositoryImpl) FindByUserID(ctx context.Context, userID uint) ([]*entity.Band, error) {
	var models []BandModel
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&models).Error; err != nil {
		return nil, err
	}
	return r.modelsToEntities(models), nil
}

func (r *bandRepositoryImpl) FindByUserIDPaginated(ctx context.Context, userID uint, offset, limit int) ([]*entity.Band, int64, error) {
	var models []BandModel
	var total int64
	
	// Get total count for this user
	if err := r.db.WithContext(ctx).Model(&BandModel{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	// Get paginated results
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Offset(offset).Limit(limit).Find(&models).Error; err != nil {
		return nil, 0, err
	}
	
	return r.modelsToEntities(models), total, nil
}

func (r *bandRepositoryImpl) FindPublicBands(ctx context.Context) ([]*entity.Band, error) {
	var models []BandModel
	if err := r.db.WithContext(ctx).Where("is_public = ?", true).Find(&models).Error; err != nil {
		return nil, err
	}
	return r.modelsToEntities(models), nil
}

func (r *bandRepositoryImpl) FindPublicBandsPaginated(ctx context.Context, offset, limit int) ([]*entity.Band, int64, error) {
	var models []BandModel
	var total int64
	
	// Get total count
	if err := r.db.WithContext(ctx).Model(&BandModel{}).Where("is_public = ?", true).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	// Get paginated results
	if err := r.db.WithContext(ctx).Where("is_public = ?", true).Offset(offset).Limit(limit).Find(&models).Error; err != nil {
		return nil, 0, err
	}
	
	return r.modelsToEntities(models), total, nil
}

func (r *bandRepositoryImpl) Update(ctx context.Context, band *entity.Band) error {
	model := r.entityToModel(band)
	return r.db.WithContext(ctx).Save(&model).Error
}

func (r *bandRepositoryImpl) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&BandModel{}, id).Error
}

// Music management

func (r *bandRepositoryImpl) AddMusicsToBand(ctx context.Context, bandID uint, musicIDs []uint) error {
	for _, musicID := range musicIDs {
		bandMusic := BandMusicModel{
			BandID:  bandID,
			MusicID: musicID,
		}
		// Use FirstOrCreate to avoid duplicates
		if err := r.db.WithContext(ctx).
			Where("band_id = ? AND music_id = ?", bandID, musicID).
			FirstOrCreate(&bandMusic).Error; err != nil {
			return err
		}
	}
	return nil
}

func (r *bandRepositoryImpl) RemoveMusicFromBand(ctx context.Context, bandID, musicID uint) error {
	return r.db.WithContext(ctx).
		Where("band_id = ? AND music_id = ?", bandID, musicID).
		Delete(&BandMusicModel{}).Error
}

func (r *bandRepositoryImpl) GetBandMusics(ctx context.Context, bandID uint) ([]*entity.Music, error) {
	var musics []MusicModel
	err := r.db.WithContext(ctx).
		Table("musics").
		Joins("JOIN band_musics ON musics.id = band_musics.music_id").
		Where("band_musics.band_id = ?", bandID).
		Order("band_musics.display_order ASC, band_musics.created_at ASC").
		Find(&musics).Error
	
	if err != nil {
		return nil, err
	}

	// Convert to entities
	entities := make([]*entity.Music, len(musics))
	for i, model := range musics {
		entities[i] = &entity.Music{
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
	return entities, nil
}

func (r *bandRepositoryImpl) ReorderBandMusics(ctx context.Context, bandID uint, musicOrders map[uint]int) error {
	for musicID, order := range musicOrders {
		if err := r.db.WithContext(ctx).
			Model(&BandMusicModel{}).
			Where("band_id = ? AND music_id = ?", bandID, musicID).
			Update("display_order", order).Error; err != nil {
			return err
		}
	}
	return nil
}

// Member management

func (r *bandRepositoryImpl) GetBandMemberCount(ctx context.Context, bandID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&UserModel{}).
		Where("band_id = ?", bandID).
		Count(&count).Error
	return count, err
}

func (r *bandRepositoryImpl) GetBandMembers(ctx context.Context, bandID uint) ([]*entity.User, error) {
	var users []UserModel
	err := r.db.WithContext(ctx).
		Where("band_id = ?", bandID).
		Find(&users).Error
	
	if err != nil {
		return nil, err
	}

	// Convert to entities
	entities := make([]*entity.User, len(users))
	for i, model := range users {
		entities[i] = &entity.User{
			ID:         model.ID,
			Username:   model.Username,
			Name:   model.Name,
			Profile:    model.Profile,
			Email:      model.Email,
			Role:       model.Role,
			Provider:   model.Provider,
			ProviderID: model.ProviderID,
			CreatedAt:  model.CreatedAt,
			UpdatedAt:  model.UpdatedAt,
		}
	}
	return entities, nil
}

// Helper methods

func (r *bandRepositoryImpl) entityToModel(band *entity.Band) *BandModel {
	return &BandModel{
		ID:          band.ID,
		Name:        band.Name,
		Description: band.Description,
		Cover:       band.Cover,
		IsPublic:    band.IsPublic,
		UserID:      band.UserID,
		CreatedAt:   band.CreatedAt,
		UpdatedAt:   band.UpdatedAt,
	}
}

func (r *bandRepositoryImpl) modelToEntity(model *BandModel) *entity.Band {
	return &entity.Band{
		ID:          model.ID,
		Name:        model.Name,
		Description: model.Description,
		Cover:       model.Cover,
		IsPublic:    model.IsPublic,
		UserID:      model.UserID,
		CreatedAt:   model.CreatedAt,
		UpdatedAt:   model.UpdatedAt,
	}
}

func (r *bandRepositoryImpl) modelsToEntities(models []BandModel) []*entity.Band {
	entities := make([]*entity.Band, len(models))
	for i, model := range models {
		entities[i] = r.modelToEntity(&model)
	}
	return entities
}

