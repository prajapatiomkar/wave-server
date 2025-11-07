package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Username  string         `gorm:"unique;not null;size:50" json:"username"`
	Email     string         `gorm:"unique;not null;size:100" json:"email"`
	Password  string         `gorm:"not null" json:"-"`
	FullName  string         `gorm:"size:100" json:"full_name"`
	Avatar    string         `json:"avatar"`
	IsOnline  bool           `gorm:"default:false" json:"is_online"`
	LastSeen  *time.Time     `json:"last_seen"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type UserResponse struct {
	ID       uint       `json:"id"`
	Username string     `json:"username"`
	Email    string     `json:"email"`
	FullName string     `json:"full_name"`
	Avatar   string     `json:"avatar"`
	IsOnline bool       `json:"is_online"`
	LastSeen *time.Time `json:"last_seen"`
}
