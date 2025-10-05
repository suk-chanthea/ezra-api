package domain

import "time"

type Event struct {
    ID        uint      `json:"id"`
    Title     string    `json:"title"`
    Content   string    `json:"content"`
    Location  string    `json:"location"`
    StartTime time.Time `json:"start_time"`
    EndTime   time.Time `json:"end_time"`
    UserID    uint      `json:"user_id"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
