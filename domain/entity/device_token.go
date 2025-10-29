package entity

import "time"

type DeviceToken struct {
	ID        uint
	UserID    uint
	Token     string
	Platform  string // ios, android, web
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewDeviceToken(userID uint, token, platform string) *DeviceToken {
	return &DeviceToken{
		UserID:   userID,
		Token:    token,
		Platform: platform,
		IsActive: true,
	}
}

func (dt *DeviceToken) IsValid() bool {
	if dt.UserID == 0 || dt.Token == "" {
		return false
	}

	// Validate platform
	validPlatforms := map[string]bool{
		"ios":     true,
		"android": true,
		"web":     true,
	}

	return validPlatforms[dt.Platform]
}

func (dt *DeviceToken) Deactivate() {
	dt.IsActive = false
}

