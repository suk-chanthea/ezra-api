package domain

import "time"

type Event struct {
    ID          uint        `gorm:"primaryKey" json:"id"`
    Title       string      `json:"title" binding:"required" validate:"required,min=1,max=200"`
    Content     string      `json:"content"`
    Cover       string      `json:"cover" binding:"required" validate:"required,min=1,max=200"`
    Location    string      `json:"location" binding:"required" validate:"required,min=1,max=200"`
    StartTime   time.Time   `json:"start_time" binding:"required" validate:"required"`
    EndTime     time.Time   `json:"end_time" binding:"required" validate:"required,gtfield=StartTime"`
    UserID      uint        `json:"user_id"`
    CreatedAt   time.Time   `gorm:"autoCreateTime" json:"created_at"` // GORM sets automatically
    UpdatedAt   time.Time   `gorm:"autoUpdateTime" json:"updated_at"`
}
