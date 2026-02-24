package repository

import (
	"context"
	"fmt"
	"sample-todo-backend/internal/domain/entity"
	"sample-todo-backend/internal/domain/repository"
	"time"

	"gorm.io/gorm"
)

// todoRepository implements repository.TodoRepository
type todoRepository struct {
	db *gorm.DB
}

// NewTodoRepository creates a new todo repository
func NewTodoRepository(db *gorm.DB) repository.TodoRepository {
	return &todoRepository{db: db}
}

// Create creates a new todo
func (r *todoRepository) Create(ctx context.Context, todo *entity.Todo) error {
	if err := r.db.WithContext(ctx).Create(todo).Error; err != nil {
		return fmt.Errorf("failed to create todo: %w", err)
	}
	return nil
}

// GetByID retrieves a todo by ID
func (r *todoRepository) GetByID(ctx context.Context, id uint) (*entity.Todo, error) {
	var todo entity.Todo
	err := r.db.WithContext(ctx).First(&todo, id).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("todo not found")
		}
		return nil, fmt.Errorf("failed to get todo: %w", err)
	}

	return &todo, nil
}

// GetByIDWithRelations retrieves a todo by ID with all relations
func (r *todoRepository) GetByIDWithRelations(ctx context.Context, id uint) (*entity.Todo, error) {
	var todo entity.Todo
	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Category").
		Preload("Tags").
		Preload("Attachments").
		Preload("Comments.User").
		First(&todo, id).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("todo not found")
		}
		return nil, fmt.Errorf("failed to get todo: %w", err)
	}

	return &todo, nil
}

// Update updates a todo
func (r *todoRepository) Update(ctx context.Context, todo *entity.Todo) error {
	if err := r.db.WithContext(ctx).Save(todo).Error; err != nil {
		return fmt.Errorf("failed to update todo: %w", err)
	}
	return nil
}

// Delete soft deletes a todo
func (r *todoRepository) Delete(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&entity.Todo{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete todo: %w", err)
	}
	return nil
}

// List retrieves todos with filtering
func (r *todoRepository) List(ctx context.Context, filter *entity.TodoFilter) ([]*entity.Todo, error) {
	query := r.db.WithContext(ctx).Model(&entity.Todo{})

	// Apply filters
	query = query.Where("user_id = ?", filter.UserID)

	if filter.CategoryID != nil {
		query = query.Where("category_id = ?", *filter.CategoryID)
	}

	if filter.Status != nil {
		query = query.Where("status = ?", *filter.Status)
	}

	if filter.Priority != nil {
		query = query.Where("priority = ?", *filter.Priority)
	}

	if filter.DueBefore != nil {
		query = query.Where("due_date <= ?", *filter.DueBefore)
	}

	if filter.DueAfter != nil {
		query = query.Where("due_date >= ?", *filter.DueAfter)
	}

	if filter.Search != nil && *filter.Search != "" {
		searchTerm := "%" + *filter.Search + "%"
		query = query.Where("title ILIKE ? OR description ILIKE ?", searchTerm, searchTerm)
	}

	if filter.Completed != nil {
		if *filter.Completed {
			query = query.Where("status = ?", entity.TodoStatusCompleted)
		} else {
			query = query.Where("status != ?", entity.TodoStatusCompleted)
		}
	}

	// Tag filtering
	if len(filter.Tags) > 0 {
		query = query.
			Joins("JOIN todo_tags ON todo_tags.todo_id = todos.id").
			Joins("JOIN tags ON tags.id = todo_tags.tag_id").
			Where("tags.name IN ?", filter.Tags).
			Group("todos.id").
			Having("COUNT(DISTINCT tags.id) = ?", len(filter.Tags))
	}

	// Apply sorting
	if filter.SortBy != "" {
		order := filter.SortBy
		if filter.SortOrder == "desc" {
			order += " DESC"
		}
		query = query.Order(order)
	} else {
		// Default sorting: incomplete todos first, then by priority, then by due date
		query = query.Order("CASE WHEN status = 'completed' THEN 1 ELSE 0 END, priority DESC, due_date ASC NULLS LAST, position ASC")
	}

	// Apply pagination
	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
	}
	if filter.Offset > 0 {
		query = query.Offset(filter.Offset)
	}

	var todos []*entity.Todo
	if err := query.Preload("Category").Preload("Tags").Find(&todos).Error; err != nil {
		return nil, fmt.Errorf("failed to list todos: %w", err)
	}

	return todos, nil
}

// Count counts todos with filtering
func (r *todoRepository) Count(ctx context.Context, filter *entity.TodoFilter) (int64, error) {
	query := r.db.WithContext(ctx).Model(&entity.Todo{})

	// Apply same filters as List method
	query = query.Where("user_id = ?", filter.UserID)

	if filter.CategoryID != nil {
		query = query.Where("category_id = ?", *filter.CategoryID)
	}

	if filter.Status != nil {
		query = query.Where("status = ?", *filter.Status)
	}

	if filter.Priority != nil {
		query = query.Where("priority = ?", *filter.Priority)
	}

	if filter.DueBefore != nil {
		query = query.Where("due_date <= ?", *filter.DueBefore)
	}

	if filter.DueAfter != nil {
		query = query.Where("due_date >= ?", *filter.DueAfter)
	}

	if filter.Search != nil && *filter.Search != "" {
		searchTerm := "%" + *filter.Search + "%"
		query = query.Where("title ILIKE ? OR description ILIKE ?", searchTerm, searchTerm)
	}

	if filter.Completed != nil {
		if *filter.Completed {
			query = query.Where("status = ?", entity.TodoStatusCompleted)
		} else {
			query = query.Where("status != ?", entity.TodoStatusCompleted)
		}
	}

	if len(filter.Tags) > 0 {
		query = query.
			Joins("JOIN todo_tags ON todo_tags.todo_id = todos.id").
			Joins("JOIN tags ON tags.id = todo_tags.tag_id").
			Where("tags.name IN ?", filter.Tags).
			Group("todos.id").
			Having("COUNT(DISTINCT tags.id) = ?", len(filter.Tags))
	}

	var count int64
	if err := query.Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count todos: %w", err)
	}

	return count, nil
}

// Status operations

// UpdateStatus updates todo status
func (r *todoRepository) UpdateStatus(ctx context.Context, id uint, status entity.TodoStatus) error {
	updates := map[string]interface{}{
		"status": status,
	}

	// Update completion timestamp
	if status == entity.TodoStatusCompleted {
		now := time.Now()
		updates["completed_at"] = &now
	} else {
		updates["completed_at"] = nil
	}

	err := r.db.WithContext(ctx).
		Model(&entity.Todo{}).
		Where("id = ?", id).
		Updates(updates).Error

	if err != nil {
		return fmt.Errorf("failed to update todo status: %w", err)
	}

	return nil
}

// UpdatePosition updates todo position
func (r *todoRepository) UpdatePosition(ctx context.Context, id uint, position int) error {
	err := r.db.WithContext(ctx).
		Model(&entity.Todo{}).
		Where("id = ?", id).
		Update("position", position).Error

	if err != nil {
		return fmt.Errorf("failed to update todo position: %w", err)
	}

	return nil
}

// MarkCompleted marks todo as completed
func (r *todoRepository) MarkCompleted(ctx context.Context, id uint) error {
	return r.UpdateStatus(ctx, id, entity.TodoStatusCompleted)
}

// MarkPending marks todo as pending
func (r *todoRepository) MarkPending(ctx context.Context, id uint) error {
	return r.UpdateStatus(ctx, id, entity.TodoStatusPending)
}

// Bulk operations

// BulkUpdateStatus updates status for multiple todos
func (r *todoRepository) BulkUpdateStatus(ctx context.Context, ids []uint, status entity.TodoStatus) error {
	if len(ids) == 0 {
		return nil
	}

	updates := map[string]interface{}{
		"status": status,
	}

	if status == entity.TodoStatusCompleted {
		now := time.Now()
		updates["completed_at"] = &now
	} else {
		updates["completed_at"] = nil
	}

	err := r.db.WithContext(ctx).
		Model(&entity.Todo{}).
		Where("id IN ?", ids).
		Updates(updates).Error

	if err != nil {
		return fmt.Errorf("failed to bulk update todo status: %w", err)
	}

	return nil
}

// BulkDelete soft deletes multiple todos
func (r *todoRepository) BulkDelete(ctx context.Context, ids []uint) error {
	if len(ids) == 0 {
		return nil
	}

	err := r.db.WithContext(ctx).
		Where("id IN ?", ids).
		Delete(&entity.Todo{}).Error

	if err != nil {
		return fmt.Errorf("failed to bulk delete todos: %w", err)
	}

	return nil
}

// Analytics

// GetCompletionStats retrieves completion statistics for a user
func (r *todoRepository) GetCompletionStats(ctx context.Context, userID uint) (*repository.TodoStats, error) {
	stats := &repository.TodoStats{}

	// Get total count
	if err := r.db.WithContext(ctx).
		Model(&entity.Todo{}).
		Where("user_id = ?", userID).
		Count(&stats.Total).Error; err != nil {
		return nil, fmt.Errorf("failed to count total todos: %w", err)
	}

	// Get counts by status
	statusCounts := []struct {
		Status entity.TodoStatus
		Count  int64
	}{}

	err := r.db.WithContext(ctx).
		Model(&entity.Todo{}).
		Select("status, COUNT(*) as count").
		Where("user_id = ?", userID).
		Group("status").
		Scan(&statusCounts).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get status counts: %w", err)
	}

	// Map status counts
	for _, sc := range statusCounts {
		switch sc.Status {
		case entity.TodoStatusCompleted:
			stats.Completed = sc.Count
		case entity.TodoStatusInProgress:
			stats.InProgress = sc.Count
		case entity.TodoStatusPending:
			stats.Pending = sc.Count
		case entity.TodoStatusCancelled:
			stats.Cancelled = sc.Count
		}
	}

	// Get overdue count
	now := time.Now()
	if err := r.db.WithContext(ctx).
		Model(&entity.Todo{}).
		Where("user_id = ? AND due_date < ? AND status != ?", userID, now, entity.TodoStatusCompleted).
		Count(&stats.Overdue).Error; err != nil {
		return nil, fmt.Errorf("failed to count overdue todos: %w", err)
	}

	// Calculate completion rate
	if stats.Total > 0 {
		stats.CompletionRate = float64(stats.Completed) / float64(stats.Total) * 100
	}

	return stats, nil
}

// GetOverdueTodos retrieves overdue todos for a user
func (r *todoRepository) GetOverdueTodos(ctx context.Context, userID uint) ([]*entity.Todo, error) {
	var todos []*entity.Todo
	now := time.Now()

	err := r.db.WithContext(ctx).
		Where("user_id = ? AND due_date < ? AND status != ?", userID, now, entity.TodoStatusCompleted).
		Order("due_date ASC").
		Preload("Category").
		Find(&todos).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get overdue todos: %w", err)
	}

	return todos, nil
}

// GetTodosByPriority retrieves todos by priority for a user
func (r *todoRepository) GetTodosByPriority(ctx context.Context, userID uint, priority entity.Priority) ([]*entity.Todo, error) {
	var todos []*entity.Todo

	err := r.db.WithContext(ctx).
		Where("user_id = ? AND priority = ? AND status != ?", userID, priority, entity.TodoStatusCompleted).
		Order("due_date ASC NULLS LAST, created_at DESC").
		Preload("Category").
		Preload("Tags").
		Find(&todos).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get todos by priority: %w", err)
	}

	return todos, nil
}

// Tag operations

// AddTag adds a tag to a todo
func (r *todoRepository) AddTag(ctx context.Context, todoID uint, tagID uint) error {
	// Check if association already exists
	var count int64
	if err := r.db.WithContext(ctx).
		Table("todo_tags").
		Where("todo_id = ? AND tag_id = ?", todoID, tagID).
		Count(&count).Error; err != nil {
		return fmt.Errorf("failed to check tag association: %w", err)
	}

	if count > 0 {
		return fmt.Errorf("tag already associated with todo")
	}

	// Create association
	if err := r.db.WithContext(ctx).
		Exec("INSERT INTO todo_tags (todo_id, tag_id) VALUES (?, ?)", todoID, tagID).Error; err != nil {
		return fmt.Errorf("failed to add tag: %w", err)
	}

	return nil
}

// RemoveTag removes a tag from a todo
func (r *todoRepository) RemoveTag(ctx context.Context, todoID uint, tagID uint) error {
	err := r.db.WithContext(ctx).
		Exec("DELETE FROM todo_tags WHERE todo_id = ? AND tag_id = ?", todoID, tagID).Error

	if err != nil {
		return fmt.Errorf("failed to remove tag: %w", err)
	}

	return nil
}

// GetTodoTags retrieves all tags for a todo
func (r *todoRepository) GetTodoTags(ctx context.Context, todoID uint) ([]*entity.Tag, error) {
	var tags []*entity.Tag

	err := r.db.WithContext(ctx).
		Joins("JOIN todo_tags ON todo_tags.tag_id = tags.id").
		Where("todo_tags.todo_id = ?", todoID).
		Find(&tags).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get todo tags: %w", err)
	}

	return tags, nil
}