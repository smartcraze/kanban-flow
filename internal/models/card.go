package models

import (
	"time"

	"gorm.io/gorm"
)

// Card represents a task in a kanban list
type Card struct {
	gorm.Model
	Title       string     `json:"title" gorm:"not null"`
	Description string     `json:"description"`
	Position    int        `json:"position" gorm:"not null"`
	ListID      uint       `json:"list_id" gorm:"not null"`
	AssigneeID  *uint      `json:"assignee_id"`
	DueDate     *time.Time `json:"due_date"`
}
