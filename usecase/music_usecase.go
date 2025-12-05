package usecase

import (
	"errors"

	"github.com/suk-chanthea/ezra/domain/dto"
	"github.com/suk-chanthea/ezra/domain/entity"
	"github.com/suk-chanthea/ezra/domain/repository"
)

type MusicUseCase interface {
	CreateMusic(title, cover string, userID uint) error
	GetAllMusics() ([]*dto.MusicResponse, error)
	GetAllMusicsPaginated(page, pageSize int) ([]*dto.MusicResponse, *dto.PaginationMetadata, error)
	GetMusicByID(id uint) (*dto.MusicResponse, error)
	GetMusicsByUserID(userID uint) ([]*dto.MusicResponse, error)
	UpdateMusic(id uint, title, cover string, userID uint) error
	DeleteMusic(id uint, userID uint) error
}

type musicUseCase struct {
	musicRepo repository.MusicRepository
}

func NewMusicUseCase(repo repository.MusicRepository) MusicUseCase {
	return &musicUseCase{
		musicRepo: repo,
	}
}

func (uc *musicUseCase) CreateMusic(title, cover string, userID uint) error {
	music := entity.NewMusic(title, cover, userID)
	
	if !music.IsValid() {
		return errors.New("invalid music data")
	}
	
	return uc.musicRepo.Save(music)
}

func (uc *musicUseCase) GetAllMusics() ([]*dto.MusicResponse, error) {
	musics, err := uc.musicRepo.FindAll()
	if err != nil {
		return nil, err
	}
	
	return uc.entitiesToResponses(musics), nil
}

func (uc *musicUseCase) GetAllMusicsPaginated(page, pageSize int) ([]*dto.MusicResponse, *dto.PaginationMetadata, error) {
	offset := (page - 1) * pageSize
	musics, total, err := uc.musicRepo.FindAllPaginated(offset, pageSize)
	if err != nil {
		return nil, nil, err
	}
	
	pagination := dto.NewPaginationMetadata(page, pageSize, total)
	return uc.entitiesToResponses(musics), pagination, nil
}

func (uc *musicUseCase) GetMusicByID(id uint) (*dto.MusicResponse, error) {
	music, err := uc.musicRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	
	return uc.entityToResponse(music), nil
}

func (uc *musicUseCase) GetMusicsByUserID(userID uint) ([]*dto.MusicResponse, error) {
	musics, err := uc.musicRepo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}
	
	return uc.entitiesToResponses(musics), nil
}

func (uc *musicUseCase) UpdateMusic(id uint, title, cover string, userID uint) error {
	music, err := uc.musicRepo.FindByID(id)
	if err != nil {
		return err
	}
	
	// Check ownership
	if music.UserID != userID {
		return errors.New("unauthorized")
	}
	
	music.Title = title
	music.Cover = cover
	
	return uc.musicRepo.Update(music)
}

func (uc *musicUseCase) DeleteMusic(id uint, userID uint) error {
	music, err := uc.musicRepo.FindByID(id)
	if err != nil {
		return err
	}
	
	// Check ownership
	if music.UserID != userID {
		return errors.New("unauthorized")
	}
	
	return uc.musicRepo.Delete(id)
}

func (uc *musicUseCase) entityToResponse(music *entity.Music) *dto.MusicResponse {
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
		CreatedAt:   dto.NewLocalTime(music.CreatedAt),
		UpdatedAt:   dto.NewLocalTime(music.UpdatedAt),
	}
}

func (uc *musicUseCase) entitiesToResponses(musics []*entity.Music) []*dto.MusicResponse {
	responses := make([]*dto.MusicResponse, len(musics))
	for i, music := range musics {
		responses[i] = uc.entityToResponse(music)
	}
	return responses
}