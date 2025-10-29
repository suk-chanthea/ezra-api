package entity

import "time"

// Notification represents a user notification
type Notification struct {
	ID            uint
	UserID        *uint  // NULL for broadcast/team
	SenderID      *uint  // Who sent the notification
	BandID        *uint  // For team notifications
	RecipientType string // user, band, all
	Title         string
	Message       string
	Type          string // info, success, warning, error, booking, music, event
	RelatedType   string // music, event, booking, band
	RelatedID     *uint
	IsRead        bool
	ReadAt        *time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// NewNotification creates a new notification for a specific user
func NewNotification(userID uint, title, message, notifType string) *Notification {
	return &Notification{
		UserID:        &userID,
		RecipientType: "user",
		Title:         title,
		Message:       message,
		Type:          notifType,
		IsRead:        false,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
}

// NewBandNotification creates a new notification for a band/team
func NewBandNotification(bandID uint, title, message, notifType string) *Notification {
	return &Notification{
		BandID:        &bandID,
		RecipientType: "band",
		Title:         title,
		Message:       message,
		Type:          notifType,
		IsRead:        false,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
}

// NewBroadcastNotification creates a new notification for all users
func NewBroadcastNotification(title, message, notifType string) *Notification {
	return &Notification{
		RecipientType: "all",
		Title:         title,
		Message:       message,
		Type:          notifType,
		IsRead:        false,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
}

// IsValid validates notification entity
func (n *Notification) IsValid() bool {
	validTypes := map[string]bool{
		"info":    true,
		"success": true,
		"warning": true,
		"error":   true,
		"booking": true,
		"music":   true,
		"event":   true,
	}
	
	validRecipientTypes := map[string]bool{
		"user": true,
		"band": true,
		"all":  true,
	}
	
	// Basic validation
	if n.Title == "" || n.Message == "" || !validTypes[n.Type] || !validRecipientTypes[n.RecipientType] {
		return false
	}
	
	// Validate based on recipient type
	switch n.RecipientType {
	case "user":
		return n.UserID != nil && *n.UserID > 0 && n.BandID == nil
	case "band":
		return n.BandID != nil && *n.BandID > 0 && n.UserID == nil
	case "all":
		return n.UserID == nil && n.BandID == nil
	default:
		return false
	}
}

// MarkAsRead marks the notification as read
func (n *Notification) MarkAsRead() {
	n.IsRead = true
	now := time.Now()
	n.ReadAt = &now
}

