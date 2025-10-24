package persistence

import (
	"context"
	"time"

	"github.com/suk-chanthea/ezra/domain/entity"
	"github.com/suk-chanthea/ezra/domain/repository"
	"gorm.io/gorm"
)

// FavoriteModel is the GORM model for database
type FavoriteModel struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"not null"`
	MusicID   uint      `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

func (FavoriteModel) TableName() string {
	return "favorites"
}

type favoriteRepositoryImpl struct {
	db *gorm.DB
}

func NewFavoriteRepository(db *gorm.DB) repository.FavoriteRepository {
	return &favoriteRepositoryImpl{db: db}
}

func (r *favoriteRepositoryImpl) Create(ctx context.Context, favorite *entity.Favorite) error {
	model := r.entityToModel(favorite)
	if err := r.db.WithContext(ctx).Create(&model).Error; err != nil {
		return err
	}
	favorite.ID = model.ID
	favorite.CreatedAt = model.CreatedAt
	return nil
}

func (r *favoriteRepositoryImpl) Delete(ctx context.Context, userID, musicID uint) error {
	return r.db.WithContext(ctx).
		Where("user_id = ? AND music_id = ?", userID, musicID).
		Delete(&FavoriteModel{}).Error
}

func (r *favoriteRepositoryImpl) GetByUserID(ctx context.Context, userID uint) ([]*entity.Music, error) {
	var musics []MusicModel
	err := r.db.WithContext(ctx).
		Table("musics").
		Joins("JOIN favorites ON musics.id = favorites.music_id").
		Where("favorites.user_id = ?", userID).
		Order("favorites.created_at DESC").
		Find(&musics).Error
	
	if err != nil {
		return nil, err
	}

	// Convert to entities
	entities := make([]*entity.Music, len(musics))
	for i, model := range musics {
		entities[i] = &entity.Music{
			ID:        model.ID,
			Title:     model.Title,
			Cover:     model.Cover,
			Audio:     model.Audio,
			UserID:    model.UserID,
			CreatedAt: model.CreatedAt,
			UpdatedAt: model.UpdatedAt,
		}
	}
	return entities, nil
}

func (r *favoriteRepositoryImpl) GetByUserIDPaginated(ctx context.Context, userID uint, offset, limit int) ([]*entity.Music, int64, error) {
	var musics []MusicModel
	var total int64
	
	// Get total count
	err := r.db.WithContext(ctx).
		Table("favorites").
		Where("user_id = ?", userID).
		Count(&total).Error
	
	if err != nil {
		return nil, 0, err
	}
	
	// Get paginated results
	err = r.db.WithContext(ctx).
		Table("musics").
		Joins("JOIN favorites ON musics.id = favorites.music_id").
		Where("favorites.user_id = ?", userID).
		Order("favorites.created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&musics).Error
	
	if err != nil {
		return nil, 0, err
	}

	// Convert to entities
	entities := make([]*entity.Music, len(musics))
	for i, model := range musics {
		entities[i] = &entity.Music{
			ID:        model.ID,
			Title:     model.Title,
			Cover:     model.Cover,
			Audio:     model.Audio,
			UserID:    model.UserID,
			CreatedAt: model.CreatedAt,
			UpdatedAt: model.UpdatedAt,
		}
	}
	return entities, total, nil
}

func (r *favoriteRepositoryImpl) IsFavorite(ctx context.Context, userID, musicID uint) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&FavoriteModel{}).
		Where("user_id = ? AND music_id = ?", userID, musicID).
		Count(&count).Error
	
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *favoriteRepositoryImpl) GetFavoriteCount(ctx context.Context, musicID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&FavoriteModel{}).
		Where("music_id = ?", musicID).
		Count(&count).Error
	
	return count, err
}

func (r *favoriteRepositoryImpl) entityToModel(favorite *entity.Favorite) *FavoriteModel {
	return &FavoriteModel{
		ID:        favorite.ID,
		UserID:    favorite.UserID,
		MusicID:   favorite.MusicID,
		CreatedAt: favorite.CreatedAt,
	}
}

