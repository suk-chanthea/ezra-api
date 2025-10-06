package domain

import "time"

type Event struct {
    ID          uint        `gorm:"primaryKey" json:"-"`
    Title       string      `json:"title"`
    Content     string      `json:"content"`
    Cover       string      `json:"cover"`
    Location    string      `json:"location"`
    StartTime   time.Time   `json:"start_time"`
    EndTime     time.Time   `json:"end_time"`
    UserID      uint        `json:"-"`
    CreatedAt   time.Time   `gorm:"autoCreateTime" json:"-"` // GORM sets automatically
    UpdatedAt   time.Time   `gorm:"autoUpdateTime" json:"-"`
}
