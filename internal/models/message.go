package models

import (
	"time"

	"gorm.io/gorm"
)

type Message struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	RoomID    string         `gorm:"index;not null" json:"room_id"`
	UserID    uint           `gorm:"index;not null" json:"user_id"`
	Content   string         `gorm:"type:text;not null" json:"content"`
	Type      string         `gorm:"default:'text'" json:"type"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	User User `gorm:"foreignKey:UserID" json:"user"`
}

type MessageResponse struct {
	ID        uint         `json:"id"`
	RoomID    string       `json:"room_id"`
	UserID    uint         `json:"user_id"`
	Content   string       `json:"content"`
	Type      string       `json:"type"`
	CreatedAt time.Time    `json:"created_at"`
	User      UserResponse `json:"user"`
}
