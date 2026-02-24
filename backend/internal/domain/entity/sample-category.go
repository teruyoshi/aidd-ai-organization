package entity

import (
	"time"

	"gorm.io/gorm"
)

// Category represents a todo category entity
type Category struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	UserID      uint           `gorm:"not null;index:idx_categories_user_id" json:"user_id"`
	Name        string         `gorm:"size:100;not null" json:"name" validate:"required,min=1,max=100"`
	Description *string        `gorm:"type:text" json:"description,omitempty"`
	Color       string         `gorm:"size:7;default:'#3B82F6'" json:"color"` // HEX color code
	Icon        *string        `gorm:"size:50" json:"icon,omitempty"`        // Icon name or emoji
	Position    int            `gorm:"default:0" json:"position"`            // ソート順
	IsDefault   bool           `gorm:"default:false" json:"is_default"`      // デフォルトカテゴリかどうか
	TodoCount   int            `gorm:"-" json:"todo_count,omitempty"`        // 集計用（DBには保存しない）
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	User  User   `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user,omitempty"`
	Todos []Todo `gorm:"foreignKey:CategoryID" json:"todos,omitempty"`
}

// CategoryFilter represents filtering options for categories
type CategoryFilter struct {
	UserID    uint    `json:"user_id"`
	Search    *string `json:"search"`
	IsDefault *bool   `json:"is_default"`
	Limit     int     `json:"limit"`
	Offset    int     `json:"offset"`
	SortBy    string  `json:"sort_by"`
	SortOrder string  `json:"sort_order"`
}

// CategoryMethods contains business logic methods for Category entity

// Validate performs business logic validation
func (c *Category) Validate() error {
	if c.Name == "" {
		return NewValidationError("name is required")
	}

	if len(c.Name) > 100 {
		return NewValidationError("name must be less than 100 characters")
	}

	if c.UserID == 0 {
		return NewValidationError("user_id is required")
	}

	// Validate color format (HEX color)
	if c.Color != "" && len(c.Color) != 7 {
		return NewValidationError("color must be a valid HEX color code")
	}

	return nil
}

// CanEdit checks if the category can be edited by the given user
func (c *Category) CanEdit(userID uint) bool {
	return c.UserID == userID
}

// CanView checks if the category can be viewed by the given user
func (c *Category) CanView(userID uint) bool {
	// 現在の実装では作成者のみ表示可能
	// 将来的には共有機能を追加予定
	return c.UserID == userID
}

// CanDelete checks if the category can be deleted by the given user
func (c *Category) CanDelete(userID uint) bool {
	// デフォルトカテゴリは削除不可
	if c.IsDefault {
		return false
	}
	return c.UserID == userID
}

// TableName specifies the table name for GORM
func (Category) TableName() string {
	return "categories"
}

// GetColorValue returns the color value with default
func (c *Category) GetColorValue() string {
	if c.Color == "" {
		return "#3B82F6" // Default blue
	}
	return c.Color
}

// GetIconValue returns the icon value with default
func (c *Category) GetIconValue() string {
	if c.Icon == nil || *c.Icon == "" {
		return "📋" // Default clipboard emoji
	}
	return *c.Icon
}

// SetAsDefault marks this category as default and ensures only one default per user
func (c *Category) SetAsDefault() {
	c.IsDefault = true
}

// UnsetAsDefault removes default flag from this category
func (c *Category) UnsetAsDefault() {
	c.IsDefault = false
}

// GetDisplayName returns the category display name
func (c *Category) GetDisplayName() string {
	if c.Icon != nil && *c.Icon != "" {
		return *c.Icon + " " + c.Name
	}
	return c.Name
}

// HasTodos checks if category has any todos (based on TodoCount field)
func (c *Category) HasTodos() bool {
	return c.TodoCount > 0
}

// IsDeletable checks if category can be safely deleted
func (c *Category) IsDeletable() bool {
	// デフォルトカテゴリまたはTODOが存在する場合は削除不可
	return !c.IsDefault && !c.HasTodos()
}