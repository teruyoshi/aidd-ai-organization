package service

import "time"

// Validator interface for request validation
type Validator interface {
	Validate(interface{}) error
}

// User service DTOs

// RegisterUserRequest represents user registration request
type RegisterUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=100"`
	Name     string `json:"name" validate:"required,min=2,max=100"`
}

// AuthenticateUserRequest represents user authentication request
type AuthenticateUserRequest struct {
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required"`
	IPAddress string `json:"-"` // Set by middleware
	UserAgent string `json:"-"` // Set by middleware
}

// UpdateProfileRequest represents profile update request
type UpdateProfileRequest struct {
	Name        string     `json:"name,omitempty" validate:"omitempty,min=2,max=100"`
	Avatar      *string    `json:"avatar,omitempty" validate:"omitempty,url"`
	Bio         *string    `json:"bio,omitempty" validate:"omitempty,max=500"`
	Location    *string    `json:"location,omitempty" validate:"omitempty,max=100"`
	Website     *string    `json:"website,omitempty" validate:"omitempty,url"`
	DateOfBirth *time.Time `json:"date_of_birth,omitempty"`
	Timezone    string     `json:"timezone,omitempty" validate:"omitempty,timezone"`
	Language    string     `json:"language,omitempty" validate:"omitempty,len=2"`
}

// ChangePasswordRequest represents password change request
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=8,max=100"`
}

// UserResponse represents user response
type UserResponse struct {
	ID        uint             `json:"id"`
	Email     string           `json:"email"`
	Name      string           `json:"name"`
	Avatar    *string          `json:"avatar,omitempty"`
	Status    string           `json:"status"`
	Profile   *ProfileResponse `json:"profile,omitempty"`
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt time.Time        `json:"updated_at"`
}

// ProfileResponse represents profile response
type ProfileResponse struct {
	ID          uint       `json:"id"`
	Bio         *string    `json:"bio,omitempty"`
	Location    *string    `json:"location,omitempty"`
	Website     *string    `json:"website,omitempty"`
	DateOfBirth *time.Time `json:"date_of_birth,omitempty"`
	Timezone    string     `json:"timezone"`
	Language    string     `json:"language"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// AuthResponse represents authentication response
type AuthResponse struct {
	Token string        `json:"token"`
	User  *UserResponse `json:"user"`
}

// SessionResponse represents session response
type SessionResponse struct {
	ID           uint      `json:"id"`
	IPAddress    string    `json:"ip_address"`
	UserAgent    string    `json:"user_agent"`
	LastActivity time.Time `json:"last_activity"`
	ExpiresAt    time.Time `json:"expires_at"`
	IsActive     bool      `json:"is_active"`
}

// Todo service DTOs

// CreateTodoRequest represents todo creation request
type CreateTodoRequest struct {
	CategoryID  *uint      `json:"category_id,omitempty"`
	Title       string     `json:"title" validate:"required,min=1,max=255"`
	Description *string    `json:"description,omitempty" validate:"omitempty,max=1000"`
	Priority    string     `json:"priority" validate:"required,oneof=low medium high urgent"`
	DueDate     *time.Time `json:"due_date,omitempty"`
	Tags        []string   `json:"tags,omitempty"`
}

// UpdateTodoRequest represents todo update request
type UpdateTodoRequest struct {
	CategoryID  *uint      `json:"category_id,omitempty"`
	Title       string     `json:"title,omitempty" validate:"omitempty,min=1,max=255"`
	Description *string    `json:"description,omitempty" validate:"omitempty,max=1000"`
	Status      string     `json:"status,omitempty" validate:"omitempty,oneof=pending in_progress completed cancelled"`
	Priority    string     `json:"priority,omitempty" validate:"omitempty,oneof=low medium high urgent"`
	DueDate     *time.Time `json:"due_date,omitempty"`
	Position    *int       `json:"position,omitempty"`
	Tags        []string   `json:"tags,omitempty"`
}

// TodoListRequest represents todo list request
type TodoListRequest struct {
	CategoryID *uint      `json:"category_id,omitempty"`
	Status     *string    `json:"status,omitempty" validate:"omitempty,oneof=pending in_progress completed cancelled"`
	Priority   *string    `json:"priority,omitempty" validate:"omitempty,oneof=low medium high urgent"`
	Tags       []string   `json:"tags,omitempty"`
	DueBefore  *time.Time `json:"due_before,omitempty"`
	DueAfter   *time.Time `json:"due_after,omitempty"`
	Search     *string    `json:"search,omitempty"`
	Completed  *bool      `json:"completed,omitempty"`
	Limit      int        `json:"limit,omitempty" validate:"omitempty,min=1,max=100"`
	Offset     int        `json:"offset,omitempty" validate:"omitempty,min=0"`
	SortBy     string     `json:"sort_by,omitempty" validate:"omitempty,oneof=created_at updated_at due_date priority position title"`
	SortOrder  string     `json:"sort_order,omitempty" validate:"omitempty,oneof=asc desc"`
}

// TodoResponse represents todo response
type TodoResponse struct {
	ID          uint                  `json:"id"`
	UserID      uint                  `json:"user_id"`
	CategoryID  *uint                 `json:"category_id,omitempty"`
	Title       string                `json:"title"`
	Description *string               `json:"description,omitempty"`
	Status      string                `json:"status"`
	Priority    string                `json:"priority"`
	DueDate     *time.Time            `json:"due_date,omitempty"`
	CompletedAt *time.Time            `json:"completed_at,omitempty"`
	Position    int                   `json:"position"`
	Tags        []*TagResponse        `json:"tags,omitempty"`
	Attachments []*AttachmentResponse `json:"attachments,omitempty"`
	Category    *CategoryResponse     `json:"category,omitempty"`
	CreatedAt   time.Time             `json:"created_at"`
	UpdatedAt   time.Time             `json:"updated_at"`
}

// TodoListResponse represents paginated todo list response
type TodoListResponse struct {
	Todos      []*TodoResponse `json:"todos"`
	Total      int64           `json:"total"`
	Limit      int             `json:"limit"`
	Offset     int             `json:"offset"`
	HasMore    bool            `json:"has_more"`
}

// BulkUpdateTodosRequest represents bulk todo update request
type BulkUpdateTodosRequest struct {
	TodoIDs []uint `json:"todo_ids" validate:"required,min=1"`
	Status  string `json:"status" validate:"required,oneof=pending in_progress completed cancelled"`
}

// Category service DTOs

// CreateCategoryRequest represents category creation request
type CreateCategoryRequest struct {
	Name        string  `json:"name" validate:"required,min=1,max=100"`
	Description *string `json:"description,omitempty" validate:"omitempty,max=500"`
	Color       string  `json:"color,omitempty" validate:"omitempty,hexcolor"`
	Icon        *string `json:"icon,omitempty" validate:"omitempty,max=50"`
	Position    int     `json:"position,omitempty"`
}

// UpdateCategoryRequest represents category update request
type UpdateCategoryRequest struct {
	Name        string  `json:"name,omitempty" validate:"omitempty,min=1,max=100"`
	Description *string `json:"description,omitempty" validate:"omitempty,max=500"`
	Color       string  `json:"color,omitempty" validate:"omitempty,hexcolor"`
	Icon        *string `json:"icon,omitempty" validate:"omitempty,max=50"`
	Position    *int    `json:"position,omitempty"`
}

// CategoryResponse represents category response
type CategoryResponse struct {
	ID          uint    `json:"id"`
	UserID      uint    `json:"user_id"`
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
	Color       string  `json:"color"`
	Icon        *string `json:"icon,omitempty"`
	Position    int     `json:"position"`
	IsDefault   bool    `json:"is_default"`
	TodoCount   int     `json:"todo_count"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ReorderCategoriesRequest represents category reordering request
type ReorderCategoriesRequest struct {
	CategoryIDs []uint `json:"category_ids" validate:"required,min=1"`
}

// Tag service DTOs

// CreateTagRequest represents tag creation request
type CreateTagRequest struct {
	Name        string  `json:"name" validate:"required,min=1,max=100"`
	Color       string  `json:"color,omitempty" validate:"omitempty,hexcolor"`
	Description *string `json:"description,omitempty" validate:"omitempty,max=500"`
}

// UpdateTagRequest represents tag update request
type UpdateTagRequest struct {
	Name        string  `json:"name,omitempty" validate:"omitempty,min=1,max=100"`
	Color       string  `json:"color,omitempty" validate:"omitempty,hexcolor"`
	Description *string `json:"description,omitempty" validate:"omitempty,max=500"`
}

// TagResponse represents tag response
type TagResponse struct {
	ID          uint    `json:"id"`
	UserID      uint    `json:"user_id"`
	Name        string  `json:"name"`
	Color       string  `json:"color"`
	Description *string `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Attachment service DTOs

// AttachmentResponse represents attachment response
type AttachmentResponse struct {
	ID           uint      `json:"id"`
	TodoID       uint      `json:"todo_id"`
	FileName     string    `json:"file_name"`
	OriginalName string    `json:"original_name"`
	MimeType     string    `json:"mime_type"`
	FileSize     int64     `json:"file_size"`
	UploadedBy   uint      `json:"uploaded_by"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Comment service DTOs

// CreateCommentRequest represents comment creation request
type CreateCommentRequest struct {
	Content string `json:"content" validate:"required,min=1,max=1000"`
}

// UpdateCommentRequest represents comment update request
type UpdateCommentRequest struct {
	Content string `json:"content" validate:"required,min=1,max=1000"`
}

// CommentResponse represents comment response
type CommentResponse struct {
	ID        uint         `json:"id"`
	TodoID    uint         `json:"todo_id"`
	UserID    uint         `json:"user_id"`
	Content   string       `json:"content"`
	User      *UserResponse `json:"user,omitempty"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
}

// Statistics and analytics DTOs

// TodoStatsResponse represents todo statistics response
type TodoStatsResponse struct {
	Total          int64   `json:"total"`
	Completed      int64   `json:"completed"`
	InProgress     int64   `json:"in_progress"`
	Pending        int64   `json:"pending"`
	Cancelled      int64   `json:"cancelled"`
	Overdue        int64   `json:"overdue"`
	CompletionRate float64 `json:"completion_rate"`
}

// DashboardResponse represents dashboard data response
type DashboardResponse struct {
	Stats            *TodoStatsResponse `json:"stats"`
	OverdueTodos     []*TodoResponse    `json:"overdue_todos"`
	UpcomingTodos    []*TodoResponse    `json:"upcoming_todos"`
	RecentTodos      []*TodoResponse    `json:"recent_todos"`
	PopularTags      []*TagResponse     `json:"popular_tags"`
	CategoriesWithCount []*CategoryResponse `json:"categories_with_count"`
}

// Common response wrappers

// SuccessResponse represents a successful operation response
type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Code    string `json:"code,omitempty"`
	Details interface{} `json:"details,omitempty"`
}

// ValidationError represents validation error details
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Value   interface{} `json:"value,omitempty"`
}