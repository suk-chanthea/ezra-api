package repository

import (
	"context"

	"github.com/suk-chanthea/ezra/domain/entity"
)

type NotificationRepository interface {
	// Create creates a new notification
	Create(ctx context.Context, notification *entity.Notification) error

	// FindByID finds a notification by ID
	FindByID(ctx context.Context, id uint) (*entity.Notification, error)

	// FindByUserID finds all notifications for a user (includes personal, band, and broadcast)
	FindByUserID(ctx context.Context, userID uint) ([]*entity.Notification, error)

	// FindByUserIDPaginated finds notifications for a user with pagination
	FindByUserIDPaginated(ctx context.Context, userID uint, offset, limit int) ([]*entity.Notification, int64, error)

	// FindUnreadByUserID finds all unread notifications for a user
	FindUnreadByUserID(ctx context.Context, userID uint) ([]*entity.Notification, error)

	// GetUnreadCount gets the count of unread notifications for a user
	GetUnreadCount(ctx context.Context, userID uint) (int64, error)

	// FindByBandID finds all notifications for a band
	FindByBandID(ctx context.Context, bandID uint) ([]*entity.Notification, error)

	// FindBroadcastNotifications finds all broadcast notifications
	FindBroadcastNotifications(ctx context.Context) ([]*entity.Notification, error)

	// Update updates a notification
	Update(ctx context.Context, notification *entity.Notification) error

	// Delete deletes a notification by ID
	Delete(ctx context.Context, id uint) error

	// DeleteAllByUserID deletes all notifications for a user
	DeleteAllByUserID(ctx context.Context, userID uint) error

	// MarkAsRead marks a notification as read
	MarkAsRead(ctx context.Context, id, userID uint) error

	// MarkAllAsRead marks all notifications as read for a user
	MarkAllAsRead(ctx context.Context, userID uint) error
}

