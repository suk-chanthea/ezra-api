package usecase

import (
	"context"
	"errors"

	"github.com/suk-chanthea/ezra/domain/dto"
	"github.com/suk-chanthea/ezra/domain/entity"
	"github.com/suk-chanthea/ezra/domain/repository"
)

type FavoriteUseCase interface {
	AddFavorite(ctx context.Context, userID, musicID uint) error
	RemoveFavorite(ctx context.Context, userID, musicID uint) error
	GetUserFavorites(ctx context.Context, userID uint) ([]*dto.MusicResponse, error)
	GetUserFavoritesPaginated(ctx context.Context, userID uint, page, pageSize int) ([]*dto.MusicResponse, *dto.PaginationMetadata, error)
	IsFavorite(ctx context.Context, userID, musicID uint) (bool, error)
	GetFavoriteCount(ctx context.Context, musicID uint) (int64, error)
}

type favoriteUseCase struct {
	favoriteRepo repository.FavoriteRepository
	musicRepo    repository.MusicRepository
}

func NewFavoriteUseCase(favoriteRepo repository.FavoriteRepository, musicRepo repository.MusicRepository) FavoriteUseCase {
	return &favoriteUseCase{
		favoriteRepo: favoriteRepo,
		musicRepo:    musicRepo,
	}
}

func (uc *favoriteUseCase) AddFavorite(ctx context.Context, userID, musicID uint) error {
	// Check if music exists
	_, err := uc.musicRepo.FindByID(musicID)
	if err != nil {
		return errors.New("music not found")
	}

	// Check if already favorited
	isFav, err := uc.favoriteRepo.IsFavorite(ctx, userID, musicID)
	if err != nil {
		return err
	}
	if isFav {
		return errors.New("already favorited")
	}

	favorite := entity.NewFavorite(userID, musicID)
	if !favorite.IsValid() {
		return errors.New("invalid favorite data")
	}

	return uc.favoriteRepo.Create(ctx, favorite)
}

func (uc *favoriteUseCase) RemoveFavorite(ctx context.Context, userID, musicID uint) error {
	// Check if favorite exists
	isFav, err := uc.favoriteRepo.IsFavorite(ctx, userID, musicID)
	if err != nil {
		return err
	}
	if !isFav {
		return errors.New("favorite not found")
	}

	return uc.favoriteRepo.Delete(ctx, userID, musicID)
}

func (uc *favoriteUseCase) GetUserFavorites(ctx context.Context, userID uint) ([]*dto.MusicResponse, error) {
	musics, err := uc.favoriteRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return uc.entitiesToResponses(musics), nil
}

func (uc *favoriteUseCase) GetUserFavoritesPaginated(ctx context.Context, userID uint, page, pageSize int) ([]*dto.MusicResponse, *dto.PaginationMetadata, error) {
	offset := (page - 1) * pageSize
	musics, total, err := uc.favoriteRepo.GetByUserIDPaginated(ctx, userID, offset, pageSize)
	if err != nil {
		return nil, nil, err
	}
	
	pagination := dto.NewPaginationMetadata(page, pageSize, total)
	return uc.entitiesToResponses(musics), pagination, nil
}

func (uc *favoriteUseCase) IsFavorite(ctx context.Context, userID, musicID uint) (bool, error) {
	return uc.favoriteRepo.IsFavorite(ctx, userID, musicID)
}

func (uc *favoriteUseCase) GetFavoriteCount(ctx context.Context, musicID uint) (int64, error) {
	return uc.favoriteRepo.GetFavoriteCount(ctx, musicID)
}

func (uc *favoriteUseCase) entityToResponse(music *entity.Music) *dto.MusicResponse {
	return &dto.MusicResponse{
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

func (uc *favoriteUseCase) entitiesToResponses(musics []*entity.Music) []*dto.MusicResponse {
	responses := make([]*dto.MusicResponse, len(musics))
	for i, music := range musics {
		responses[i] = uc.entityToResponse(music)
	}
	return responses
}

