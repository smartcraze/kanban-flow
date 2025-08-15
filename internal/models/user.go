package models

type User struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	Email    string `json:"email" gorm:"unique;not null"`
	Password string `json:"-" gorm:"not null "` // Password should not be exposed
	Name     string `json:"name"`
}
