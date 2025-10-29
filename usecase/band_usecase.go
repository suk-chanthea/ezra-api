package usecase

import (
	"context"
	"errors"

	"github.com/suk-chanthea/ezra/domain/dto"
	"github.com/suk-chanthea/ezra/domain/entity"
	"github.com/suk-chanthea/ezra/domain/repository"
)

type BandUseCase interface {
	CreateBand(ctx context.Context, req *dto.CreateBandRequest, userID uint) error
	GetAllBands(ctx context.Context) ([]*dto.BandResponse, error)
	GetAllBandsPaginated(ctx context.Context, page, pageSize int) ([]*dto.BandResponse, *dto.PaginationMetadata, error)
	GetBandByID(ctx context.Context, id uint) (*dto.BandResponse, error)
	GetBandsByUserID(ctx context.Context, userID uint) ([]*dto.BandResponse, error)
	GetPublicBands(ctx context.Context) ([]*dto.BandResponse, error)
	GetPublicBandsPaginated(ctx context.Context, page, pageSize int) ([]*dto.BandResponse, *dto.PaginationMetadata, error)
	UpdateBand(ctx context.Context, id uint, req *dto.UpdateBandRequest, userID uint) error
	DeleteBand(ctx context.Context, id uint, userID uint) error
	
	// Music management
	AddMusicsToBand(ctx context.Context, bandID uint, musicIDs []uint, userID uint) error
	RemoveMusicFromBand(ctx context.Context, bandID, musicID, userID uint) error
	GetBandMusics(ctx context.Context, bandID uint) ([]*dto.MusicResponse, error)
	ReorderBandMusics(ctx context.Context, bandID uint, musicOrders []dto.MusicOrder, userID uint) error
	
	// Member management
	GetBandMembers(ctx context.Context, bandID uint) ([]*dto.UserResponse, error)
}

type bandUseCase struct {
	bandRepo  repository.BandRepository
	musicRepo repository.MusicRepository
}

func NewBandUseCase(bandRepo repository.BandRepository, musicRepo repository.MusicRepository) BandUseCase {
	return &bandUseCase{
		bandRepo:  bandRepo,
		musicRepo: musicRepo,
	}
}

func (uc *bandUseCase) CreateBand(ctx context.Context, req *dto.CreateBandRequest, userID uint) error {
	band := entity.NewBand(req.Name, req.Description, req.Cover, req.IsPublic, userID)
	
	if !band.IsValid() {
		return errors.New("invalid band data")
	}
	
	return uc.bandRepo.Save(ctx, band)
}

func (uc *bandUseCase) GetAllBands(ctx context.Context) ([]*dto.BandResponse, error) {
	bands, err := uc.bandRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	
	return uc.entitiesToResponses(ctx, bands), nil
}

func (uc *bandUseCase) GetAllBandsPaginated(ctx context.Context, page, pageSize int) ([]*dto.BandResponse, *dto.PaginationMetadata, error) {
	offset := (page - 1) * pageSize
	bands, total, err := uc.bandRepo.FindAllPaginated(ctx, offset, pageSize)
	if err != nil {
		return nil, nil, err
	}
	
	pagination := dto.NewPaginationMetadata(page, pageSize, total)
	return uc.entitiesToResponses(ctx, bands), pagination, nil
}

func (uc *bandUseCase) GetBandByID(ctx context.Context, id uint) (*dto.BandResponse, error) {
	band, err := uc.bandRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	
	return uc.entityToResponse(ctx, band, true), nil
}

func (uc *bandUseCase) GetBandsByUserID(ctx context.Context, userID uint) ([]*dto.BandResponse, error) {
	bands, err := uc.bandRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	
	return uc.entitiesToResponses(ctx, bands), nil
}

func (uc *bandUseCase) GetPublicBands(ctx context.Context) ([]*dto.BandResponse, error) {
	bands, err := uc.bandRepo.FindPublicBands(ctx)
	if err != nil {
		return nil, err
	}
	
	return uc.entitiesToResponses(ctx, bands), nil
}

func (uc *bandUseCase) GetPublicBandsPaginated(ctx context.Context, page, pageSize int) ([]*dto.BandResponse, *dto.PaginationMetadata, error) {
	offset := (page - 1) * pageSize
	bands, total, err := uc.bandRepo.FindPublicBandsPaginated(ctx, offset, pageSize)
	if err != nil {
		return nil, nil, err
	}
	
	pagination := dto.NewPaginationMetadata(page, pageSize, total)
	return uc.entitiesToResponses(ctx, bands), pagination, nil
}

func (uc *bandUseCase) UpdateBand(ctx context.Context, id uint, req *dto.UpdateBandRequest, userID uint) error {
	// Check if band exists and belongs to user
	existingBand, err := uc.bandRepo.FindByID(ctx, id)
	if err != nil {
		return errors.New("band not found")
	}

	if existingBand.UserID != userID {
		return errors.New("unauthorized to update this band")
	}

	// Update fields
	existingBand.Name = req.Name
	existingBand.Description = req.Description
	existingBand.Cover = req.Cover
	existingBand.IsPublic = req.IsPublic

	// Validate
	if !existingBand.IsValid() {
		return errors.New("invalid band data")
	}

	return uc.bandRepo.Update(ctx, existingBand)
}

func (uc *bandUseCase) DeleteBand(ctx context.Context, id uint, userID uint) error {
	// Check if band exists and belongs to user
	band, err := uc.bandRepo.FindByID(ctx, id)
	if err != nil {
		return errors.New("band not found")
	}

	if band.UserID != userID {
		return errors.New("unauthorized to delete this band")
	}

	return uc.bandRepo.Delete(ctx, id)
}

// Music management

func (uc *bandUseCase) AddMusicsToBand(ctx context.Context, bandID uint, musicIDs []uint, userID uint) error {
	// Check if band exists and belongs to user
	band, err := uc.bandRepo.FindByID(ctx, bandID)
	if err != nil {
		return errors.New("band not found")
	}

	if band.UserID != userID {
		return errors.New("unauthorized to modify this band")
	}

	// Validate all music IDs exist
	if len(musicIDs) > 0 {
		musics, err := uc.musicRepo.FindByIDs(musicIDs)
		if err != nil {
			return errors.New("failed to validate music IDs")
		}
		
		if len(musics) != len(musicIDs) {
			return errors.New("one or more music IDs do not exist")
		}
	}

	return uc.bandRepo.AddMusicsToBand(ctx, bandID, musicIDs)
}

func (uc *bandUseCase) RemoveMusicFromBand(ctx context.Context, bandID, musicID, userID uint) error {
	// Check if band exists and belongs to user
	band, err := uc.bandRepo.FindByID(ctx, bandID)
	if err != nil {
		return errors.New("band not found")
	}

	if band.UserID != userID {
		return errors.New("unauthorized to modify this band")
	}

	return uc.bandRepo.RemoveMusicFromBand(ctx, bandID, musicID)
}

func (uc *bandUseCase) GetBandMusics(ctx context.Context, bandID uint) ([]*dto.MusicResponse, error) {
	musics, err := uc.bandRepo.GetBandMusics(ctx, bandID)
	if err != nil {
		return nil, err
	}

	return uc.musicsToResponses(musics), nil
}

func (uc *bandUseCase) ReorderBandMusics(ctx context.Context, bandID uint, musicOrders []dto.MusicOrder, userID uint) error {
	// Check if band exists and belongs to user
	band, err := uc.bandRepo.FindByID(ctx, bandID)
	if err != nil {
		return errors.New("band not found")
	}

	if band.UserID != userID {
		return errors.New("unauthorized to modify this band")
	}

	// Convert to map
	orderMap := make(map[uint]int)
	for _, order := range musicOrders {
		orderMap[order.MusicID] = order.DisplayOrder
	}

	return uc.bandRepo.ReorderBandMusics(ctx, bandID, orderMap)
}

// Member management

func (uc *bandUseCase) GetBandMembers(ctx context.Context, bandID uint) ([]*dto.UserResponse, error) {
	members, err := uc.bandRepo.GetBandMembers(ctx, bandID)
	if err != nil {
		return nil, err
	}

	return uc.usersToResponses(members), nil
}

// Helper methods

func (uc *bandUseCase) entityToResponse(ctx context.Context, band *entity.Band, includeDetails bool) *dto.BandResponse {
	response := &dto.BandResponse{
		ID:          band.ID,
		Name:        band.Name,
		Description: band.Description,
		Cover:       band.Cover,
		IsPublic:    band.IsPublic,
		UserID:      band.UserID,
		CreatedAt:   band.CreatedAt,
		UpdatedAt:   band.UpdatedAt,
	}
	
	if includeDetails {
		// Get member count
		memberCount, _ := uc.bandRepo.GetBandMemberCount(ctx, band.ID)
		response.MemberCount = memberCount
		
		// Get musics
		musics, _ := uc.bandRepo.GetBandMusics(ctx, band.ID)
		response.Musics = uc.musicsToResponses(musics)
		response.MusicCount = len(musics)
	}
	
	return response
}

func (uc *bandUseCase) entitiesToResponses(ctx context.Context, bands []*entity.Band) []*dto.BandResponse {
	responses := make([]*dto.BandResponse, len(bands))
	for i, band := range bands {
		responses[i] = uc.entityToResponse(ctx, band, false)
		// Get member count for list view
		memberCount, _ := uc.bandRepo.GetBandMemberCount(ctx, band.ID)
		responses[i].MemberCount = memberCount
	}
	return responses
}

func (uc *bandUseCase) musicsToResponses(musics []*entity.Music) []*dto.MusicResponse {
	responses := make([]*dto.MusicResponse, len(musics))
	for i, music := range musics {
		responses[i] = &dto.MusicResponse{
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
	return responses
}

func (uc *bandUseCase) usersToResponses(users []*entity.User) []*dto.UserResponse {
	responses := make([]*dto.UserResponse, len(users))
	for i, user := range users {
		responses[i] = &dto.UserResponse{
			ID:        user.ID,
			Username:  user.Username,
			Fullname:  user.Fullname,
			Profile:   user.Profile,
			Email:     user.Email,
			Role:      user.Role,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		}
	}
	return responses
}

