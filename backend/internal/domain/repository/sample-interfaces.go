package repository

import (
	"context"
	"sample-todo-backend/internal/domain/entity"
)

// UserRepository defines the interface for user data operations
type UserRepository interface {
	// Basic CRUD operations
	Create(ctx context.Context, user *entity.User) error
	GetByID(ctx context.Context, id uint) (*entity.User, error)
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
	Update(ctx context.Context, user *entity.User) error
	Delete(ctx context.Context, id uint) error

	// Advanced operations
	List(ctx context.Context, filter *entity.UserFilter) ([]*entity.User, error)
	Count(ctx context.Context, filter *entity.UserFilter) (int64, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)

	// Profile operations
	CreateProfile(ctx context.Context, profile *entity.Profile) error
	UpdateProfile(ctx context.Context, profile *entity.Profile) error
	GetProfile(ctx context.Context, userID uint) (*entity.Profile, error)

	// Role operations
	AssignRole(ctx context.Context, userRole *entity.UserRole) error
	RevokeRole(ctx context.Context, userID uint, role entity.Role) error
	GetRoles(ctx context.Context, userID uint) ([]*entity.UserRole, error)

	// Session operations
	CreateSession(ctx context.Context, session *entity.UserSession) error
	GetSession(ctx context.Context, token string) (*entity.UserSession, error)
	UpdateSession(ctx context.Context, session *entity.UserSession) error
	DeleteSession(ctx context.Context, token string) error
	DeleteExpiredSessions(ctx context.Context) error
	GetUserSessions(ctx context.Context, userID uint) ([]*entity.UserSession, error)
}

// TodoRepository defines the interface for todo data operations
type TodoRepository interface {
	// Basic CRUD operations
	Create(ctx context.Context, todo *entity.Todo) error
	GetByID(ctx context.Context, id uint) (*entity.Todo, error)
	GetByIDWithRelations(ctx context.Context, id uint) (*entity.Todo, error)
	Update(ctx context.Context, todo *entity.Todo) error
	Delete(ctx context.Context, id uint) error

	// Advanced operations
	List(ctx context.Context, filter *entity.TodoFilter) ([]*entity.Todo, error)
	Count(ctx context.Context, filter *entity.TodoFilter) (int64, error)

	// Status operations
	UpdateStatus(ctx context.Context, id uint, status entity.TodoStatus) error
	UpdatePosition(ctx context.Context, id uint, position int) error
	MarkCompleted(ctx context.Context, id uint) error
	MarkPending(ctx context.Context, id uint) error

	// Bulk operations
	BulkUpdateStatus(ctx context.Context, ids []uint, status entity.TodoStatus) error
	BulkDelete(ctx context.Context, ids []uint) error

	// Analytics
	GetCompletionStats(ctx context.Context, userID uint) (*TodoStats, error)
	GetOverdueTodos(ctx context.Context, userID uint) ([]*entity.Todo, error)
	GetTodosByPriority(ctx context.Context, userID uint, priority entity.Priority) ([]*entity.Todo, error)

	// Tag operations
	AddTag(ctx context.Context, todoID uint, tagID uint) error
	RemoveTag(ctx context.Context, todoID uint, tagID uint) error
	GetTodoTags(ctx context.Context, todoID uint) ([]*entity.Tag, error)
}

// CategoryRepository defines the interface for category data operations
type CategoryRepository interface {
	// Basic CRUD operations
	Create(ctx context.Context, category *entity.Category) error
	GetByID(ctx context.Context, id uint) (*entity.Category, error)
	Update(ctx context.Context, category *entity.Category) error
	Delete(ctx context.Context, id uint) error

	// Advanced operations
	List(ctx context.Context, filter *entity.CategoryFilter) ([]*entity.Category, error)
	Count(ctx context.Context, filter *entity.CategoryFilter) (int64, error)
	GetByUserID(ctx context.Context, userID uint) ([]*entity.Category, error)

	// Default category operations
	GetDefaultCategory(ctx context.Context, userID uint) (*entity.Category, error)
	SetAsDefault(ctx context.Context, categoryID uint, userID uint) error

	// Position operations
	UpdatePosition(ctx context.Context, id uint, position int) error
	ReorderCategories(ctx context.Context, userID uint, categoryIDs []uint) error

	// Statistics
	GetCategoryWithTodoCount(ctx context.Context, id uint) (*entity.Category, error)
	GetCategoriesWithTodoCount(ctx context.Context, userID uint) ([]*entity.Category, error)
}

// TagRepository defines the interface for tag data operations
type TagRepository interface {
	// Basic CRUD operations
	Create(ctx context.Context, tag *entity.Tag) error
	GetByID(ctx context.Context, id uint) (*entity.Tag, error)
	Update(ctx context.Context, tag *entity.Tag) error
	Delete(ctx context.Context, id uint) error

	// Advanced operations
	List(ctx context.Context, userID uint, search string) ([]*entity.Tag, error)
	GetByUserID(ctx context.Context, userID uint) ([]*entity.Tag, error)
	GetByName(ctx context.Context, userID uint, name string) (*entity.Tag, error)

	// Usage statistics
	GetPopularTags(ctx context.Context, userID uint, limit int) ([]*entity.Tag, error)
	GetUnusedTags(ctx context.Context, userID uint) ([]*entity.Tag, error)
}

// AttachmentRepository defines the interface for attachment data operations
type AttachmentRepository interface {
	// Basic CRUD operations
	Create(ctx context.Context, attachment *entity.Attachment) error
	GetByID(ctx context.Context, id uint) (*entity.Attachment, error)
	Delete(ctx context.Context, id uint) error

	// Todo-specific operations
	GetByTodoID(ctx context.Context, todoID uint) ([]*entity.Attachment, error)
	DeleteByTodoID(ctx context.Context, todoID uint) error

	// File operations
	GetByFileName(ctx context.Context, fileName string) (*entity.Attachment, error)
	GetTotalFileSize(ctx context.Context, userID uint) (int64, error)

	// Cleanup operations
	GetOrphanedAttachments(ctx context.Context) ([]*entity.Attachment, error)
	CleanupOrphanedAttachments(ctx context.Context) error
}

// CommentRepository defines the interface for comment data operations
type CommentRepository interface {
	// Basic CRUD operations
	Create(ctx context.Context, comment *entity.Comment) error
	GetByID(ctx context.Context, id uint) (*entity.Comment, error)
	Update(ctx context.Context, comment *entity.Comment) error
	Delete(ctx context.Context, id uint) error

	// Todo-specific operations
	GetByTodoID(ctx context.Context, todoID uint) ([]*entity.Comment, error)
	DeleteByTodoID(ctx context.Context, todoID uint) error

	// Advanced operations
	Count(ctx context.Context, todoID uint) (int64, error)
	GetRecent(ctx context.Context, userID uint, limit int) ([]*entity.Comment, error)
}

// Repository aggregates
type Repositories struct {
	User       UserRepository
	Todo       TodoRepository
	Category   CategoryRepository
	Tag        TagRepository
	Attachment AttachmentRepository
	Comment    CommentRepository
}

// Data transfer objects and helper types

// TodoStats represents todo completion statistics
type TodoStats struct {
	Total       int64 `json:"total"`
	Completed   int64 `json:"completed"`
	InProgress  int64 `json:"in_progress"`
	Pending     int64 `json:"pending"`
	Cancelled   int64 `json:"cancelled"`
	Overdue     int64 `json:"overdue"`
	CompletionRate float64 `json:"completion_rate"`
}

// UserFilter extends entity.User with additional filtering options
type UserFilter struct {
	Status    *entity.UserStatus `json:"status"`
	Search    *string           `json:"search"`
	CreatedAfter  *time.Time   `json:"created_after"`
	CreatedBefore *time.Time   `json:"created_before"`
	Limit     int              `json:"limit"`
	Offset    int              `json:"offset"`
	SortBy    string           `json:"sort_by"`
	SortOrder string           `json:"sort_order"`
}

// Transaction defines the interface for database transactions
type Transaction interface {
	Commit() error
	Rollback() error
}

// TransactionManager defines the interface for transaction management
type TransactionManager interface {
	WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error
	BeginTransaction(ctx context.Context) (Transaction, error)
}