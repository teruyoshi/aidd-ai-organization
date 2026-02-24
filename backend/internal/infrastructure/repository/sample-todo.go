package repository

import (
	"context"
	"fmt"
	"sample-todo-backend/internal/domain/entity"
	"sample-todo-backend/internal/domain/repository"

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

// ListByUserID retrieves all todos for a user
func (r *todoRepository) ListByUserID(ctx context.Context, userID uint) ([]*entity.Todo, error) {
	var todos []*entity.Todo
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&todos).Error
	if err != nil {
		return nil, fmt.Errorf("failed to list todos: %w", err)
	}
	return todos, nil
}
