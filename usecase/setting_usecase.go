package usecase

import (
	"errors"

	"github.com/suk-chanthea/ezra/domain/dto"
	"github.com/suk-chanthea/ezra/domain/repository"
)

type SettingUseCase interface {
	GetUserSettings(userID uint) (*dto.SettingResponse, error)
	UpdateSettings(userID uint, req *dto.UpdateSettingRequest) (*dto.SettingResponse, error)
	ResetToDefaults(userID uint) (*dto.SettingResponse, error)
}

type settingUseCase struct {
	settingRepo repository.SettingRepository
}

func NewSettingUseCase(repo repository.SettingRepository) SettingUseCase {
	return &settingUseCase{settingRepo: repo}
}

func (uc *settingUseCase) GetUserSettings(userID uint) (*dto.SettingResponse, error) {
	setting, err := uc.settingRepo.FindByUserID(userID)
	if err != nil {
		return nil, errors.New("settings not found")
	}

	return &dto.SettingResponse{
		ID:                      setting.ID,
		UserID:                  setting.UserID,
		Language:                setting.Language,
		Theme:                   setting.Theme,
		NotifyOnBooking:         setting.NotifyOnBooking,
		NotifyOnMusic:           setting.NotifyOnMusic,
		NotifyOnEvent:           setting.NotifyOnEvent,
		EnablePushNotifications: setting.EnablePushNotifications,
		CreatedAt:               dto.NewLocalTime(setting.CreatedAt),
		UpdatedAt:               dto.NewLocalTime(setting.UpdatedAt),
	}, nil
}

func (uc *settingUseCase) UpdateSettings(userID uint, req *dto.UpdateSettingRequest) (*dto.SettingResponse, error) {
	// Check if settings exist
	existing, err := uc.settingRepo.FindByUserID(userID)
	if err != nil {
		return nil, errors.New("settings not found")
	}

	// Update fields
	existing.Language = req.Language
	existing.Theme = req.Theme
	existing.NotifyOnBooking = req.NotifyOnBooking
	existing.NotifyOnMusic = req.NotifyOnMusic
	existing.NotifyOnEvent = req.NotifyOnEvent
	existing.EnablePushNotifications = req.EnablePushNotifications

	// Validate
	if !existing.IsValid() {
		return nil, errors.New("invalid settings data")
	}

	// Save changes
	if err := uc.settingRepo.Update(existing); err != nil {
		return nil, err
	}

	return &dto.SettingResponse{
		ID:                      existing.ID,
		UserID:                  existing.UserID,
		Language:                existing.Language,
		Theme:                   existing.Theme,
		NotifyOnBooking:         existing.NotifyOnBooking,
		NotifyOnMusic:           existing.NotifyOnMusic,
		NotifyOnEvent:           existing.NotifyOnEvent,
		EnablePushNotifications: existing.EnablePushNotifications,
		CreatedAt:               dto.NewLocalTime(existing.CreatedAt),
		UpdatedAt:               dto.NewLocalTime(existing.UpdatedAt),
	}, nil
}

func (uc *settingUseCase) ResetToDefaults(userID uint) (*dto.SettingResponse, error) {
	// Check if settings exist
	existing, err := uc.settingRepo.FindByUserID(userID)
	if err != nil {
		return nil, errors.New("settings not found")
	}

	// Reset to default values
	existing.Language = "en"
	existing.Theme = "light"
	existing.NotifyOnBooking = true
	existing.NotifyOnMusic = false
	existing.NotifyOnEvent = true
	existing.EnablePushNotifications = true

	// Save changes
	if err := uc.settingRepo.Update(existing); err != nil {
		return nil, err
	}

	return &dto.SettingResponse{
		ID:                      existing.ID,
		UserID:                  existing.UserID,
		Language:                existing.Language,
		Theme:                   existing.Theme,
		NotifyOnBooking:         existing.NotifyOnBooking,
		NotifyOnMusic:           existing.NotifyOnMusic,
		NotifyOnEvent:           existing.NotifyOnEvent,
		EnablePushNotifications: existing.EnablePushNotifications,
		CreatedAt:               dto.NewLocalTime(existing.CreatedAt),
		UpdatedAt:               dto.NewLocalTime(existing.UpdatedAt),
	}, nil
}
