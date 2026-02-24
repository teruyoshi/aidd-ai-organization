package entity

import (
	"time"

	"gorm.io/gorm"
)

// Todo represents a todo item entity
type Todo struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	UserID      uint           `gorm:"not null;index:idx_todos_user_id" json:"user_id"`
	CategoryID  *uint          `gorm:"index:idx_todos_category_id" json:"category_id,omitempty"`
	Title       string         `gorm:"size:255;not null" json:"title" validate:"required,min=1,max=255"`
	Description *string        `gorm:"type:text" json:"description,omitempty"`
	Status      TodoStatus     `gorm:"type:enum('pending','in_progress','completed','cancelled');default:'pending';index:idx_todos_status" json:"status"`
	Priority    Priority       `gorm:"type:enum('low','medium','high','urgent');default:'medium';index:idx_todos_priority" json:"priority"`
	DueDate     *time.Time     `gorm:"index:idx_todos_due_date" json:"due_date,omitempty"`
	CompletedAt *time.Time     `json:"completed_at,omitempty"`
	Position    int            `gorm:"default:0" json:"position"` // ソート順
	Tags        []Tag          `gorm:"many2many:todo_tags;" json:"tags,omitempty"`
	Attachments []Attachment   `gorm:"foreignKey:TodoID;constraint:OnDelete:CASCADE" json:"attachments,omitempty"`
	Comments    []Comment      `gorm:"foreignKey:TodoID;constraint:OnDelete:CASCADE" json:"comments,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	User     User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user,omitempty"`
	Category *Category `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
}

// TodoStatus represents the status of a todo item
type TodoStatus string

const (
	TodoStatusPending    TodoStatus = "pending"
	TodoStatusInProgress TodoStatus = "in_progress"
	TodoStatusCompleted  TodoStatus = "completed"
	TodoStatusCancelled  TodoStatus = "cancelled"
)

// Priority represents the priority level of a todo item
type Priority string

const (
	PriorityLow    Priority = "low"
	PriorityMedium Priority = "medium"
	PriorityHigh   Priority = "high"
	PriorityUrgent Priority = "urgent"
)

// Tag represents a tag that can be attached to todos
type Tag struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	UserID      uint           `gorm:"not null;index" json:"user_id"`
	Name        string         `gorm:"size:100;not null" json:"name" validate:"required,min=1,max=100"`
	Color       string         `gorm:"size:7;default:'#3B82F6'" json:"color"` // HEX color code
	Description *string        `gorm:"type:text" json:"description,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	User  User   `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Todos []Todo `gorm:"many2many:todo_tags;"`
}

// Attachment represents a file attachment for a todo
type Attachment struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	TodoID       uint      `gorm:"not null;index" json:"todo_id"`
	FileName     string    `gorm:"size:255;not null" json:"file_name"`
	OriginalName string    `gorm:"size:255;not null" json:"original_name"`
	MimeType     string    `gorm:"size:100;not null" json:"mime_type"`
	FileSize     int64     `gorm:"not null" json:"file_size"`
	FilePath     string    `gorm:"size:500;not null" json:"-"` // ファイルパスはJSONに含めない
	UploadedBy   uint      `gorm:"not null" json:"uploaded_by"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	// Relations
	Todo       Todo `gorm:"foreignKey:TodoID;constraint:OnDelete:CASCADE"`
	UploadedByUser User `gorm:"foreignKey:UploadedBy"`
}

// Comment represents a comment on a todo item
type Comment struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	TodoID    uint      `gorm:"not null;index" json:"todo_id"`
	UserID    uint      `gorm:"not null;index" json:"user_id"`
	Content   string    `gorm:"type:text;not null" json:"content" validate:"required,min=1"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relations
	Todo Todo `gorm:"foreignKey:TodoID;constraint:OnDelete:CASCADE"`
	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

// TodoFilter represents filtering options for todos
type TodoFilter struct {
	UserID     uint         `json:"user_id"`
	CategoryID *uint        `json:"category_id"`
	Status     *TodoStatus  `json:"status"`
	Priority   *Priority    `json:"priority"`
	Tags       []string     `json:"tags"`
	DueBefore  *time.Time   `json:"due_before"`
	DueAfter   *time.Time   `json:"due_after"`
	Search     *string      `json:"search"`
	Completed  *bool        `json:"completed"`
	Limit      int          `json:"limit"`
	Offset     int          `json:"offset"`
	SortBy     string       `json:"sort_by"`
	SortOrder  string       `json:"sort_order"`
}

// TodoMethods contains business logic methods for Todo entity

// IsCompleted checks if the todo is completed
func (t *Todo) IsCompleted() bool {
	return t.Status == TodoStatusCompleted
}

// IsOverdue checks if the todo is overdue
func (t *Todo) IsOverdue() bool {
	if t.DueDate == nil || t.IsCompleted() {
		return false
	}
	return t.DueDate.Before(time.Now())
}

// MarkCompleted marks the todo as completed
func (t *Todo) MarkCompleted() {
	t.Status = TodoStatusCompleted
	now := time.Now()
	t.CompletedAt = &now
}

// MarkPending marks the todo as pending
func (t *Todo) MarkPending() {
	t.Status = TodoStatusPending
	t.CompletedAt = nil
}

// GetPriorityValue returns numeric value for priority (for sorting)
func (t *Todo) GetPriorityValue() int {
	switch t.Priority {
	case PriorityLow:
		return 1
	case PriorityMedium:
		return 2
	case PriorityHigh:
		return 3
	case PriorityUrgent:
		return 4
	default:
		return 2
	}
}

// Validate performs business logic validation
func (t *Todo) Validate() error {
	if t.Title == "" {
		return NewValidationError("title is required")
	}

	if len(t.Title) > 255 {
		return NewValidationError("title must be less than 255 characters")
	}

	if t.UserID == 0 {
		return NewValidationError("user_id is required")
	}

	return nil
}

// CanEdit checks if the todo can be edited by the given user
func (t *Todo) CanEdit(userID uint) bool {
	return t.UserID == userID
}

// CanView checks if the todo can be viewed by the given user
func (t *Todo) CanView(userID uint) bool {
	// 現在の実装では作成者のみ表示可能
	// 将来的には共有機能を追加予定
	return t.UserID == userID
}

// TableName specifies the table name for GORM
func (Todo) TableName() string {
	return "todos"
}

func (Tag) TableName() string {
	return "tags"
}

func (Attachment) TableName() string {
	return "todo_attachments"
}

func (Comment) TableName() string {
	return "todo_comments"
}

// GetStatusColor returns the color associated with the todo status
func (t *Todo) GetStatusColor() string {
	switch t.Status {
	case TodoStatusPending:
		return "#6B7280" // gray
	case TodoStatusInProgress:
		return "#3B82F6" // blue
	case TodoStatusCompleted:
		return "#10B981" // green
	case TodoStatusCancelled:
		return "#EF4444" // red
	default:
		return "#6B7280"
	}
}

// GetPriorityColor returns the color associated with the todo priority
func (t *Todo) GetPriorityColor() string {
	switch t.Priority {
	case PriorityLow:
		return "#10B981" // green
	case PriorityMedium:
		return "#F59E0B" // yellow
	case PriorityHigh:
		return "#EF4444" // red
	case PriorityUrgent:
		return "#DC2626" // dark red
	default:
		return "#F59E0B"
	}
}