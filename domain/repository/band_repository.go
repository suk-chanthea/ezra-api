package repository

import (
	"context"

	"github.com/suk-chanthea/ezra/domain/entity"
)

// BandRepository defines the interface for band data operations
type BandRepository interface {
	// Basic CRUD
	Save(ctx context.Context, band *entity.Band) error
	FindAll(ctx context.Context) ([]*entity.Band, error)
	FindAllPaginated(ctx context.Context, offset, limit int) ([]*entity.Band, int64, error)
	FindByID(ctx context.Context, id uint) (*entity.Band, error)
	FindByUserID(ctx context.Context, userID uint) ([]*entity.Band, error)
	FindPublicBands(ctx context.Context) ([]*entity.Band, error)
	FindPublicBandsPaginated(ctx context.Context, offset, limit int) ([]*entity.Band, int64, error)
	Update(ctx context.Context, band *entity.Band) error
	Delete(ctx context.Context, id uint) error
	
	// Music management
	AddMusicsToBand(ctx context.Context, bandID uint, musicIDs []uint) error
	RemoveMusicFromBand(ctx context.Context, bandID, musicID uint) error
	GetBandMusics(ctx context.Context, bandID uint) ([]*entity.Music, error)
	ReorderBandMusics(ctx context.Context, bandID uint, musicOrders map[uint]int) error
	
	// Member management
	GetBandMemberCount(ctx context.Context, bandID uint) (int64, error)
	GetBandMembers(ctx context.Context, bandID uint) ([]*entity.User, error)
}

