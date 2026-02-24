package entity

import (
	"time"

	"gorm.io/gorm"
)

// User represents a user entity
type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Email     string         `gorm:"uniqueIndex;size:255;not null" json:"email" validate:"required,email"`
	Password  string         `gorm:"size:255;not null" json:"-"` // JSON出力時は除外
	Name      string         `gorm:"size:100;not null" json:"name" validate:"required,min=2,max=100"`
	Todos     []Todo         `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"todos,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName specifies the table name for GORM
func (User) TableName() string {
	return "users"
}
