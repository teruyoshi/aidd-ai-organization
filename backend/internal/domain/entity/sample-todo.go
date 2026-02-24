package entity

import (
	"time"

	"gorm.io/gorm"
)

// Todo represents a todo item entity
type Todo struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	UserID      uint           `gorm:"not null;index" json:"user_id"`
	Title       string         `gorm:"size:255;not null" json:"title" validate:"required,min=1,max=255"`
	Description *string        `gorm:"type:text" json:"description,omitempty"`
	Done        bool           `gorm:"default:false" json:"done"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName specifies the table name for GORM
func (Todo) TableName() string {
	return "todos"
}

// ValidationError represents a validation error
type ValidationError struct {
	Message string
}

func (e ValidationError) Error() string {
	return e.Message
}

func NewValidationError(message string) *ValidationError {
	return &ValidationError{Message: message}
}
