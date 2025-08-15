package models

import "time"

// BoardMember represents user access to boards
type BoardMember struct {
	BoardID  uint      `json:"board_id" gorm:"primaryKey"`
	UserID   uint      `json:"user_id" gorm:"primaryKey"`
	Role     string    `json:"role" gorm:"not null"` // "owner", "editor", "viewer"
	JoinedAt time.Time `json:"joined_at"`
}
