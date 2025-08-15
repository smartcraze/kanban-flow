package models

import "gorm.io/gorm"

// Board represents a kanban board
type Board struct {
	gorm.Model
	Title       string `json:"title" gorm:"not null"`
	Description string `json:"description"`
	OwnerID     uint   `json:"owner_id" gorm:"not null"`
	Lists       []List `json:"lists,omitempty" gorm:"foreignKey:BoardID"`
}
