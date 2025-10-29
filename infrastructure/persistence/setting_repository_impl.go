package persistence

import (
	"errors"
	"time"

	"github.com/suk-chanthea/ezra/domain/entity"
	"github.com/suk-chanthea/ezra/domain/repository"
	"gorm.io/gorm"
)

// SettingModel is the GORM model for database
type SettingModel struct {
	ID                     uint      `gorm:"primaryKey;autoIncrement"`
	UserID                 uint      `gorm:"not null;uniqueIndex"`
	Language               string    `gorm:"size:10;default:en"`
	Theme                  string    `gorm:"size:20;default:light"`
	NotifyOnBooking        bool      `gorm:"default:true"`
	NotifyOnMusic          bool      `gorm:"default:false"`
	NotifyOnEvent          bool      `gorm:"default:true"`
	EnablePushNotifications bool      `gorm:"default:true"`
	CreatedAt              time.Time `gorm:"autoCreateTime"`
	UpdatedAt              time.Time `gorm:"autoUpdateTime"`
}

func (SettingModel) TableName() string {
	return "settings"
}

type settingRepositoryImpl struct {
	db *gorm.DB
}

// NewSettingRepository creates a new setting repository instance
func NewSettingRepository(db *gorm.DB) repository.SettingRepository {
	return &settingRepositoryImpl{db: db}
}

func (r *settingRepositoryImpl) Save(setting *entity.Setting) error {
	model := r.entityToModel(setting)
	if err := r.db.Create(&model).Error; err != nil {
		return err
	}
	setting.ID = model.ID
	setting.CreatedAt = model.CreatedAt
	setting.UpdatedAt = model.UpdatedAt
	return nil
}

func (r *settingRepositoryImpl) FindByUserID(userID uint) (*entity.Setting, error) {
	var model SettingModel
	if err := r.db.Where("user_id = ?", userID).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("settings not found")
		}
		return nil, err
	}
	return r.modelToEntity(&model), nil
}

func (r *settingRepositoryImpl) Update(setting *entity.Setting) error {
	model := r.entityToModel(setting)
	// Use Select to explicitly update all fields including zero values (like false)
	if err := r.db.Model(&SettingModel{}).
		Where("user_id = ?", setting.UserID).
		Select("language", "theme", "notify_on_booking", "notify_on_music", "notify_on_event", "enable_push_notifications").
		Updates(model).Error; err != nil {
		return err
	}
	// Retrieve the updated record to get the new timestamp
	var updated SettingModel
	if err := r.db.Where("user_id = ?", setting.UserID).First(&updated).Error; err != nil {
		return err
	}
	setting.UpdatedAt = updated.UpdatedAt
	return nil
}

func (r *settingRepositoryImpl) Delete(userID uint) error {
	result := r.db.Where("user_id = ?", userID).Delete(&SettingModel{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("settings not found")
	}
	return nil
}

// Helper methods to convert between entity and model
func (r *settingRepositoryImpl) entityToModel(setting *entity.Setting) *SettingModel {
	return &SettingModel{
		ID:                     setting.ID,
		UserID:                 setting.UserID,
		Language:               setting.Language,
		Theme:                  setting.Theme,
		NotifyOnBooking:        setting.NotifyOnBooking,
		NotifyOnMusic:          setting.NotifyOnMusic,
		NotifyOnEvent:          setting.NotifyOnEvent,
		EnablePushNotifications: setting.EnablePushNotifications,
		CreatedAt:              setting.CreatedAt,
		UpdatedAt:              setting.UpdatedAt,
	}
}

func (r *settingRepositoryImpl) modelToEntity(model *SettingModel) *entity.Setting {
	return &entity.Setting{
		ID:                     model.ID,
		UserID:                 model.UserID,
		Language:               model.Language,
		Theme:                  model.Theme,
		NotifyOnBooking:        model.NotifyOnBooking,
		NotifyOnMusic:          model.NotifyOnMusic,
		NotifyOnEvent:          model.NotifyOnEvent,
		EnablePushNotifications: model.EnablePushNotifications,
		CreatedAt:              model.CreatedAt,
		UpdatedAt:              model.UpdatedAt,
	}
}
