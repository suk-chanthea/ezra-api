package persistence

import (
	"context"
	"time"

	"github.com/suk-chanthea/ezra/domain/entity"
	"github.com/suk-chanthea/ezra/domain/repository"
	"gorm.io/gorm"
)

type DeviceTokenModel struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"not null;index"`
	Token     string    `gorm:"type:text;not null;index"`
	Platform  string    `gorm:"size:20;not null"`
	IsActive  bool      `gorm:"default:true;index"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func (DeviceTokenModel) TableName() string {
	return "device_tokens"
}

type deviceTokenRepositoryImpl struct {
	db *gorm.DB
}

func NewDeviceTokenRepository(db *gorm.DB) repository.DeviceTokenRepository {
	return &deviceTokenRepositoryImpl{db: db}
}

// Model to Entity conversion
func (r *deviceTokenRepositoryImpl) modelToEntity(model *DeviceTokenModel) *entity.DeviceToken {
	return &entity.DeviceToken{
		ID:        model.ID,
		UserID:    model.UserID,
		Token:     model.Token,
		Platform:  model.Platform,
		IsActive:  model.IsActive,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}
}

// Entity to Model conversion
func (r *deviceTokenRepositoryImpl) entityToModel(entity *entity.DeviceToken) *DeviceTokenModel {
	return &DeviceTokenModel{
		ID:        entity.ID,
		UserID:    entity.UserID,
		Token:     entity.Token,
		Platform:  entity.Platform,
		IsActive:  entity.IsActive,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
	}
}

func (r *deviceTokenRepositoryImpl) Save(ctx context.Context, token *entity.DeviceToken) error {
	model := r.entityToModel(token)
	
	// Use UPSERT: Insert or update if exists (based on unique constraint)
	result := r.db.WithContext(ctx).
		Where("user_id = ? AND token = ?", model.UserID, model.Token).
		Assign(DeviceTokenModel{
			Platform:  model.Platform,
			IsActive:  true,
			UpdatedAt: time.Now(),
		}).
		FirstOrCreate(model)
	
	if result.Error != nil {
		return result.Error
	}
	
	token.ID = model.ID
	token.CreatedAt = model.CreatedAt
	token.UpdatedAt = model.UpdatedAt
	
	return nil
}

func (r *deviceTokenRepositoryImpl) GetActiveTokensByUserID(ctx context.Context, userID uint) ([]string, error) {
	var tokens []string
	
	err := r.db.WithContext(ctx).
		Model(&DeviceTokenModel{}).
		Where("user_id = ? AND is_active = ?", userID, true).
		Pluck("token", &tokens).Error
	
	return tokens, err
}

func (r *deviceTokenRepositoryImpl) GetTokensByBandID(ctx context.Context, bandID uint) ([]string, error) {
	var tokens []string
	
	// Get tokens for all users in the band
	err := r.db.WithContext(ctx).
		Table("device_tokens").
		Select("device_tokens.token").
		Joins("JOIN users ON users.id = device_tokens.user_id").
		Where("users.band_id = ? AND device_tokens.is_active = ?", bandID, true).
		Pluck("token", &tokens).Error
	
	return tokens, err
}

func (r *deviceTokenRepositoryImpl) GetAllActiveTokens(ctx context.Context) ([]string, error) {
	var tokens []string
	
	err := r.db.WithContext(ctx).
		Model(&DeviceTokenModel{}).
		Where("is_active = ?", true).
		Pluck("token", &tokens).Error
	
	return tokens, err
}

func (r *deviceTokenRepositoryImpl) GetAllActiveTokensExcept(ctx context.Context, excludeUserID uint) ([]string, error) {
	var tokens []string
	
	err := r.db.WithContext(ctx).
		Model(&DeviceTokenModel{}).
		Where("is_active = ? AND user_id != ?", true, excludeUserID).
		Pluck("token", &tokens).Error
	
	return tokens, err
}

func (r *deviceTokenRepositoryImpl) DeleteToken(ctx context.Context, token string) error {
	return r.db.WithContext(ctx).
		Where("token = ?", token).
		Delete(&DeviceTokenModel{}).Error
}

func (r *deviceTokenRepositoryImpl) DeleteTokens(ctx context.Context, tokens []string) error {
	if len(tokens) == 0 {
		return nil
	}
	
	return r.db.WithContext(ctx).
		Where("token IN ?", tokens).
		Delete(&DeviceTokenModel{}).Error
}

func (r *deviceTokenRepositoryImpl) DeactivateToken(ctx context.Context, token string) error {
	return r.db.WithContext(ctx).
		Model(&DeviceTokenModel{}).
		Where("token = ?", token).
		Update("is_active", false).Error
}

func (r *deviceTokenRepositoryImpl) DeleteUserTokens(ctx context.Context, userID uint) error {
	return r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Delete(&DeviceTokenModel{}).Error
}

