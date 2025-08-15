package models

import "gorm.io/gorm"

// List represents a column in a kanban board
type List struct {
	gorm.Model
	Title    string `json:"title" gorm:"not null"`
	Position int    `json:"position" gorm:"not null"`
	BoardID  uint   `json:"board_id" gorm:"not null"`
	Cards    []Card `json:"cards,omitempty" gorm:"foreignKey:ListID"`
}
