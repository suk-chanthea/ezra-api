package repository

import (
	"context"

	"github.com/suk-chanthea/ezra/domain/entity"
)

// FavoriteRepository defines the interface for favorite data access
type FavoriteRepository interface {
	Create(ctx context.Context, favorite *entity.Favorite) error
	Delete(ctx context.Context, userID, musicID uint) error
	GetByUserID(ctx context.Context, userID uint) ([]*entity.Music, error)
	GetByUserIDPaginated(ctx context.Context, userID uint, offset, limit int) ([]*entity.Music, int64, error)
	IsFavorite(ctx context.Context, userID, musicID uint) (bool, error)
	GetFavoriteCount(ctx context.Context, musicID uint) (int64, error)
}

