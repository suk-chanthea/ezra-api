package usecase

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/suk-chanthea/ezra/domain/dto"
	"github.com/suk-chanthea/ezra/domain/entity"
	"github.com/suk-chanthea/ezra/domain/repository"
	"github.com/suk-chanthea/ezra/infrastructure/firebase"
)

type NotificationUseCase interface {
	CreateNotification(ctx context.Context, senderID uint, req *dto.CreateNotificationRequest) (*dto.NotificationResponse, error)
	CreateBandNotification(ctx context.Context, senderID, bandID uint, req *dto.CreateNotificationRequest) (*dto.NotificationResponse, error)
	CreateBroadcastNotification(ctx context.Context, senderID uint, req *dto.CreateNotificationRequest) (*dto.NotificationResponse, error)
	GetNotifications(ctx context.Context, userID uint, page, pageSize int) ([]*dto.NotificationResponse, *dto.PaginationMetadata, error)
	GetUnreadNotifications(ctx context.Context, userID uint) ([]*dto.NotificationResponse, error)
	GetUnreadCount(ctx context.Context, userID uint) (int64, error)
	GetNotificationByID(ctx context.Context, userID, notifID uint) (*dto.NotificationResponse, error)
	MarkAsRead(ctx context.Context, userID, notifID uint) error
	MarkAllAsRead(ctx context.Context, userID uint) error
	DeleteNotification(ctx context.Context, userID, notifID uint) error
	DeleteAllNotifications(ctx context.Context, userID uint) error
}

type notificationUseCase struct {
	notificationRepo repository.NotificationRepository
	fcmService       firebase.FCMService
}

func NewNotificationUseCase(repo repository.NotificationRepository, fcmService firebase.FCMService) NotificationUseCase {
	return &notificationUseCase{
		notificationRepo: repo,
		fcmService:       fcmService,
	}
}

func (uc *notificationUseCase) CreateNotification(ctx context.Context, senderID uint, req *dto.CreateNotificationRequest) (*dto.NotificationResponse, error) {
	if req.UserID == nil {
		return nil, errors.New("user_id is required for user notifications")
	}
	
	notification := entity.NewNotification(*req.UserID, req.Title, req.Message, req.Type)
	notification.SenderID = &senderID
	
	if req.RelatedType != "" {
		notification.RelatedType = req.RelatedType
	}
	if req.RelatedID != nil {
		notification.RelatedID = req.RelatedID
	}

	if !notification.IsValid() {
		return nil, errors.New("invalid notification data")
	}

	if err := uc.notificationRepo.Create(ctx, notification); err != nil {
		return nil, err
	}

	// Send FCM push notification in background
	go uc.sendFCMToUser(*req.UserID, notification)

	return uc.entityToResponse(notification), nil
}

func (uc *notificationUseCase) CreateBandNotification(ctx context.Context, senderID, bandID uint, req *dto.CreateNotificationRequest) (*dto.NotificationResponse, error) {
	notification := entity.NewBandNotification(bandID, req.Title, req.Message, req.Type)
	notification.SenderID = &senderID
	
	if req.RelatedType != "" {
		notification.RelatedType = req.RelatedType
	}
	if req.RelatedID != nil {
		notification.RelatedID = req.RelatedID
	}

	if !notification.IsValid() {
		return nil, errors.New("invalid notification data")
	}

	if err := uc.notificationRepo.Create(ctx, notification); err != nil {
		return nil, err
	}

	// Send FCM push notification to band members in background
	go uc.sendFCMToBand(bandID, notification)

	return uc.entityToResponse(notification), nil
}

func (uc *notificationUseCase) CreateBroadcastNotification(ctx context.Context, senderID uint, req *dto.CreateNotificationRequest) (*dto.NotificationResponse, error) {
	notification := entity.NewBroadcastNotification(req.Title, req.Message, req.Type)
	notification.SenderID = &senderID
	
	if req.RelatedType != "" {
		notification.RelatedType = req.RelatedType
	}
	if req.RelatedID != nil {
		notification.RelatedID = req.RelatedID
	}

	if !notification.IsValid() {
		return nil, errors.New("invalid notification data")
	}

	if err := uc.notificationRepo.Create(ctx, notification); err != nil {
		return nil, err
	}

	// Send FCM push notification to all users (except sender) in background
	go uc.sendFCMBroadcast(senderID, notification)

	return uc.entityToResponse(notification), nil
}

func (uc *notificationUseCase) GetNotifications(ctx context.Context, userID uint, page, pageSize int) ([]*dto.NotificationResponse, *dto.PaginationMetadata, error) {
	offset := (page - 1) * pageSize
	notifications, total, err := uc.notificationRepo.FindByUserIDPaginated(ctx, userID, offset, pageSize)
	if err != nil {
		return nil, nil, err
	}

	pagination := dto.NewPaginationMetadata(page, pageSize, total)
	return uc.entitiesToResponses(notifications), pagination, nil
}

func (uc *notificationUseCase) GetUnreadNotifications(ctx context.Context, userID uint) ([]*dto.NotificationResponse, error) {
	notifications, err := uc.notificationRepo.FindUnreadByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return uc.entitiesToResponses(notifications), nil
}

func (uc *notificationUseCase) GetUnreadCount(ctx context.Context, userID uint) (int64, error) {
	return uc.notificationRepo.GetUnreadCount(ctx, userID)
}

func (uc *notificationUseCase) GetNotificationByID(ctx context.Context, userID, notifID uint) (*dto.NotificationResponse, error) {
	notification, err := uc.notificationRepo.FindByID(ctx, notifID)
	if err != nil {
		return nil, errors.New("notification not found")
	}

	// Check if user can access this notification
	// Allow if it's their personal notification, or a broadcast/band notification
	if notification.RecipientType == "user" && (notification.UserID == nil || *notification.UserID != userID) {
		return nil, errors.New("unauthorized")
	}

	return uc.entityToResponse(notification), nil
}

func (uc *notificationUseCase) MarkAsRead(ctx context.Context, userID, notifID uint) error {
	// Verify user can access this notification
	notification, err := uc.notificationRepo.FindByID(ctx, notifID)
	if err != nil {
		return errors.New("notification not found")
	}

	// Allow if it's their personal notification, or a broadcast/band notification
	if notification.UserID != nil && *notification.UserID != userID && notification.RecipientType == "user" {
		return errors.New("unauthorized")
	}

	return uc.notificationRepo.MarkAsRead(ctx, notifID, userID)
}

func (uc *notificationUseCase) MarkAllAsRead(ctx context.Context, userID uint) error {
	return uc.notificationRepo.MarkAllAsRead(ctx, userID)
}

func (uc *notificationUseCase) DeleteNotification(ctx context.Context, userID, notifID uint) error {
	// Verify user can access this notification
	notification, err := uc.notificationRepo.FindByID(ctx, notifID)
	if err != nil {
		return errors.New("notification not found")
	}

	// Only allow deleting personal notifications
	if notification.UserID == nil || *notification.UserID != userID {
		return errors.New("can only delete personal notifications")
	}

	return uc.notificationRepo.Delete(ctx, notifID)
}

func (uc *notificationUseCase) DeleteAllNotifications(ctx context.Context, userID uint) error {
	return uc.notificationRepo.DeleteAllByUserID(ctx, userID)
}

func (uc *notificationUseCase) entityToResponse(notification *entity.Notification) *dto.NotificationResponse {
	return &dto.NotificationResponse{
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
		ReadAt:        dto.NewLocalTimePtr(notification.ReadAt),
		CreatedAt:     dto.NewLocalTime(notification.CreatedAt),
		UpdatedAt:     dto.NewLocalTime(notification.UpdatedAt),
	}
}

func (uc *notificationUseCase) entitiesToResponses(notifications []*entity.Notification) []*dto.NotificationResponse {
	responses := make([]*dto.NotificationResponse, len(notifications))
	for i, notification := range notifications {
		responses[i] = uc.entityToResponse(notification)
	}
	return responses
}

// Helper functions for sending FCM notifications

func (uc *notificationUseCase) sendFCMToUser(userID uint, notification *entity.Notification) {
	data := uc.buildFCMData(notification)
	
	if err := uc.fcmService.SendToUser(context.Background(), userID, notification.Title, notification.Message, data); err != nil {
		log.Printf("❌ Failed to send FCM notification to user %d: %v", userID, err)
	}
}

func (uc *notificationUseCase) sendFCMToBand(bandID uint, notification *entity.Notification) {
	data := uc.buildFCMData(notification)
	
	if err := uc.fcmService.SendToBand(context.Background(), bandID, notification.Title, notification.Message, data); err != nil {
		log.Printf("❌ Failed to send FCM notification to band %d: %v", bandID, err)
	}
}

func (uc *notificationUseCase) sendFCMBroadcast(senderID uint, notification *entity.Notification) {
	data := uc.buildFCMData(notification)
	
	if err := uc.fcmService.SendToAllExcept(context.Background(), senderID, notification.Title, notification.Message, data); err != nil {
		log.Printf("❌ Failed to send FCM broadcast notification: %v", err)
	}
}

func (uc *notificationUseCase) buildFCMData(notification *entity.Notification) map[string]string {
	data := map[string]string{
		"notification_id": fmt.Sprintf("%d", notification.ID),
		"type":            notification.Type,
		"recipient_type":  notification.RecipientType,
	}
	
	if notification.RelatedType != "" {
		data["related_type"] = notification.RelatedType
	}
	if notification.RelatedID != nil {
		data["related_id"] = fmt.Sprintf("%d", *notification.RelatedID)
	}
	if notification.SenderID != nil {
		data["sender_id"] = fmt.Sprintf("%d", *notification.SenderID)
	}
	
	return data
}

