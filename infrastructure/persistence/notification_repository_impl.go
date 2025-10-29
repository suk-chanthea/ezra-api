package persistence

import (
	"context"
	"time"

	"github.com/suk-chanthea/ezra/domain/entity"
	"github.com/suk-chanthea/ezra/domain/repository"
	"gorm.io/gorm"
)

// NotificationModel is the GORM model for database
type NotificationModel struct {
	ID            uint       `gorm:"primaryKey"`
	UserID        *uint      `gorm:"index"`
	SenderID      *uint      `gorm:"index"`
	BandID        *uint      `gorm:"index"`
	RecipientType string     `gorm:"size:20;not null;default:user;index"`
	Title         string     `gorm:"size:255;not null"`
	Message       string     `gorm:"type:text;not null"`
	Type          string     `gorm:"size:50;not null;default:info;index"`
	RelatedType   string     `gorm:"size:50"`
	RelatedID     *uint      `gorm:"type:integer"`
	IsRead        bool       `gorm:"default:false;index"`
	ReadAt        *time.Time `gorm:"type:timestamptz"`
	CreatedAt     time.Time  `gorm:"autoCreateTime;index:idx_notifications_created_at,sort:desc"`
	UpdatedAt     time.Time  `gorm:"autoUpdateTime"`
}

func (NotificationModel) TableName() string {
	return "notifications"
}

type notificationRepositoryImpl struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) repository.NotificationRepository {
	return &notificationRepositoryImpl{db: db}
}

func (r *notificationRepositoryImpl) Create(ctx context.Context, notification *entity.Notification) error {
	model := r.entityToModel(notification)
	if err := r.db.WithContext(ctx).Create(&model).Error; err != nil {
		return err
	}
	notification.ID = model.ID
	notification.CreatedAt = model.CreatedAt
	notification.UpdatedAt = model.UpdatedAt
	return nil
}

func (r *notificationRepositoryImpl) FindByID(ctx context.Context, id uint) (*entity.Notification, error) {
	var model NotificationModel
	if err := r.db.WithContext(ctx).First(&model, id).Error; err != nil {
		return nil, err
	}
	return r.modelToEntity(&model), nil
}

func (r *notificationRepositoryImpl) FindByUserID(ctx context.Context, userID uint) ([]*entity.Notification, error) {
	var models []NotificationModel
	
	// Get user's band_id
	var user struct {
		BandID *uint
	}
	r.db.WithContext(ctx).Table("users").Select("band_id").Where("id = ?", userID).Scan(&user)
	
	// Build query to get:
	// 1. Direct notifications (user_id = userID)
	// 2. Band notifications (band_id = user's band_id)
	// 3. Broadcast notifications (recipient_type = 'all')
	// Exclude notifications created by the user themselves (sender_id != userID)
	query := r.db.WithContext(ctx).Where("user_id = ?", userID)
	
	if user.BandID != nil {
		query = query.Or("band_id = ?", *user.BandID)
	}
	
	query = query.Or("recipient_type = ?", "all")
	
	// Exclude own notifications
	query = query.Where("(sender_id IS NULL OR sender_id != ?)", userID)
	
	if err := query.Order("created_at DESC").Find(&models).Error; err != nil {
		return nil, err
	}
	
	return r.modelsToEntities(models), nil
}

func (r *notificationRepositoryImpl) FindByUserIDPaginated(ctx context.Context, userID uint, offset, limit int) ([]*entity.Notification, int64, error) {
	var models []NotificationModel
	var total int64
	
	// Get user's band_id
	var user struct {
		BandID *uint
	}
	r.db.WithContext(ctx).Table("users").Select("band_id").Where("id = ?", userID).Scan(&user)
	
	// Build query
	query := r.db.WithContext(ctx).Model(&NotificationModel{}).Where("user_id = ?", userID)
	
	if user.BandID != nil {
		query = query.Or("band_id = ?", *user.BandID)
	}
	
	query = query.Or("recipient_type = ?", "all")
	
	// Exclude own notifications
	query = query.Where("(sender_id IS NULL OR sender_id != ?)", userID)

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	query = r.db.WithContext(ctx).Where("user_id = ?", userID)
	
	if user.BandID != nil {
		query = query.Or("band_id = ?", *user.BandID)
	}
	
	query = query.Or("recipient_type = ?", "all")
	
	// Exclude own notifications
	query = query.Where("(sender_id IS NULL OR sender_id != ?)", userID)
	
	if err := query.Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&models).Error; err != nil {
		return nil, 0, err
	}

	return r.modelsToEntities(models), total, nil
}

func (r *notificationRepositoryImpl) FindUnreadByUserID(ctx context.Context, userID uint) ([]*entity.Notification, error) {
	var models []NotificationModel
	
	// Get user's band_id
	var user struct {
		BandID *uint
	}
	r.db.WithContext(ctx).Table("users").Select("band_id").Where("id = ?", userID).Scan(&user)
	
	query := r.db.WithContext(ctx).Where("user_id = ? AND is_read = ?", userID, false)
	
	if user.BandID != nil {
		query = query.Or("band_id = ? AND is_read = ?", *user.BandID, false)
	}
	
	query = query.Or("recipient_type = ? AND is_read = ?", "all", false)
	
	// Exclude own notifications
	query = query.Where("(sender_id IS NULL OR sender_id != ?)", userID)
	
	if err := query.Order("created_at DESC").Find(&models).Error; err != nil {
		return nil, err
	}
	
	return r.modelsToEntities(models), nil
}

func (r *notificationRepositoryImpl) GetUnreadCount(ctx context.Context, userID uint) (int64, error) {
	var count int64
	
	// Get user's band_id
	var user struct {
		BandID *uint
	}
	r.db.WithContext(ctx).Table("users").Select("band_id").Where("id = ?", userID).Scan(&user)
	
	query := r.db.WithContext(ctx).Model(&NotificationModel{}).Where("user_id = ? AND is_read = ?", userID, false)
	
	if user.BandID != nil {
		query = query.Or("band_id = ? AND is_read = ?", *user.BandID, false)
	}
	
	query = query.Or("recipient_type = ? AND is_read = ?", "all", false)
	
	// Exclude own notifications
	query = query.Where("(sender_id IS NULL OR sender_id != ?)", userID)
	
	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}
	
	return count, nil
}

func (r *notificationRepositoryImpl) FindByBandID(ctx context.Context, bandID uint) ([]*entity.Notification, error) {
	var models []NotificationModel
	if err := r.db.WithContext(ctx).
		Where("band_id = ?", bandID).
		Order("created_at DESC").
		Find(&models).Error; err != nil {
		return nil, err
	}
	return r.modelsToEntities(models), nil
}

func (r *notificationRepositoryImpl) FindBroadcastNotifications(ctx context.Context) ([]*entity.Notification, error) {
	var models []NotificationModel
	if err := r.db.WithContext(ctx).
		Where("recipient_type = ?", "all").
		Order("created_at DESC").
		Find(&models).Error; err != nil {
		return nil, err
	}
	return r.modelsToEntities(models), nil
}

func (r *notificationRepositoryImpl) Update(ctx context.Context, notification *entity.Notification) error {
	model := r.entityToModel(notification)
	return r.db.WithContext(ctx).Save(&model).Error
}

func (r *notificationRepositoryImpl) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&NotificationModel{}, id).Error
}

func (r *notificationRepositoryImpl) DeleteAllByUserID(ctx context.Context, userID uint) error {
	return r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Delete(&NotificationModel{}).Error
}

func (r *notificationRepositoryImpl) MarkAsRead(ctx context.Context, id, userID uint) error {
	now := time.Now()
	
	// For broadcast/band notifications, we need to track per-user read status
	// For now, we'll just mark it as read for the specific notification
	// In future, consider a separate table for user-notification read status
	
	return r.db.WithContext(ctx).
		Model(&NotificationModel{}).
		Where("id = ?", id).
		Select("is_read", "read_at").
		Updates(map[string]interface{}{
			"is_read": true,
			"read_at": now,
		}).Error
}

func (r *notificationRepositoryImpl) MarkAllAsRead(ctx context.Context, userID uint) error {
	now := time.Now()
	
	// Get user's band_id
	var user struct {
		BandID *uint
	}
	r.db.WithContext(ctx).Table("users").Select("band_id").Where("id = ?", userID).Scan(&user)
	
	query := r.db.WithContext(ctx).Model(&NotificationModel{}).
		Where("user_id = ? AND is_read = ?", userID, false)
	
	if user.BandID != nil {
		query = query.Or("band_id = ? AND is_read = ?", *user.BandID, false)
	}
	
	query = query.Or("recipient_type = ? AND is_read = ?", "all", false)
	
	return query.Select("is_read", "read_at").
		Updates(map[string]interface{}{
			"is_read": true,
			"read_at": now,
		}).Error
}

func (r *notificationRepositoryImpl) entityToModel(notification *entity.Notification) *NotificationModel {
	return &NotificationModel{
		ID:            notification.ID,
		UserID:        notification.UserID,
		SenderID:      notification.SenderID,
		BandID:        notification.BandID,
		RecipientType: notification.RecipientType,
		Title:         notification.Title,
		Message:       notification.Message,
		Type:          notification.Type,
		RelatedType:   notification.RelatedType,
		RelatedID:     notification.RelatedID,
		IsRead:        notification.IsRead,
		ReadAt:        notification.ReadAt,
		CreatedAt:     notification.CreatedAt,
		UpdatedAt:     notification.UpdatedAt,
	}
}

func (r *notificationRepositoryImpl) modelToEntity(model *NotificationModel) *entity.Notification {
	return &entity.Notification{
		ID:            model.ID,
		UserID:        model.UserID,
		SenderID:      model.SenderID,
		BandID:        model.BandID,
		RecipientType: model.RecipientType,
		Title:         model.Title,
		Message:       model.Message,
		Type:          model.Type,
		RelatedType:   model.RelatedType,
		RelatedID:     model.RelatedID,
		IsRead:        model.IsRead,
		ReadAt:        model.ReadAt,
		CreatedAt:     model.CreatedAt,
		UpdatedAt:     model.UpdatedAt,
	}
}

func (r *notificationRepositoryImpl) modelsToEntities(models []NotificationModel) []*entity.Notification {
	entities := make([]*entity.Notification, len(models))
	for i, model := range models {
		entities[i] = r.modelToEntity(&model)
	}
	return entities
}

