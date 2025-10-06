package domain

import "time"

type User struct {
    ID        uint     `gorm:"primaryKey;autoIncrement" json:"id"`
    Username  string    `gorm:"size:100;not null" json:"username"`
    Fullname  string    `gorm:"size:100;not null" json:"fullname"`
    Profile   string    `gorm:"size:255" json:"profile"`
    Email     string    `gorm:"size:100;unique;not null" json:"email"`
    Password  string    `gorm:"size:255;not null" json:"-"`
    Role      string    `gorm:"size:20;default:user" json:"role"`
    Token     string    `gorm:"size:255" json:"token"`
    CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
    UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
