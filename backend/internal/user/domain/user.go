package domain

import (
	"gorm.io/gorm"
)

// User represents a user in the system
// swagger:model
type User struct {
	gorm.Model `swaggerignore:"true"`
	Name       string `json:"name" gorm:"not null"`
	Email      string `json:"email" gorm:"not null;unique"`
	Password   string `json:"password" gorm:"not null"`
	Role       string `json:"role" gorm:"not null;default:'user'"`
}
