package usecase

import (
	"errors"

	"github.com/suk-chanthea/ezra/domain/dto"
	"github.com/suk-chanthea/ezra/domain/entity"
	"github.com/suk-chanthea/ezra/domain/repository"
)

type MusicUseCase interface {
	CreateMusic(title, cover, audio string, userID uint) error
	GetAllMusics() ([]*dto.MusicResponse, error)
	GetMusicByID(id uint) (*dto.MusicResponse, error)
	GetMusicsByUserID(userID uint) ([]*dto.MusicResponse, error)
	UpdateMusic(id uint, title, cover, audio string, userID uint) error
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

func (uc *musicUseCase) CreateMusic(title, cover, audio string, userID uint) error {
	music := entity.NewMusic(title, cover, audio, userID)
	
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

func (uc *musicUseCase) UpdateMusic(id uint, title, cover, audio string, userID uint) error {
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
	music.Audio = audio
	
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
		ID:        music.ID,
		Title:     music.Title,
		Cover:     music.Cover,
		Audio:     music.Audio,
		UserID:    music.UserID,
		CreatedAt: music.CreatedAt,
		UpdatedAt: music.UpdatedAt,
	}
}

func (uc *musicUseCase) entitiesToResponses(musics []*entity.Music) []*dto.MusicResponse {
	responses := make([]*dto.MusicResponse, len(musics))
	for i, music := range musics {
		responses[i] = uc.entityToResponse(music)
	}
	return responses
}