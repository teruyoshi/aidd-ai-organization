package repository

import (
	"context"
	"fmt"
	"sample-todo-backend/internal/domain/entity"
	"sample-todo-backend/internal/domain/repository"
	"time"

	"gorm.io/gorm"
)

// userRepository implements repository.UserRepository
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return &userRepository{db: db}
}

// Create creates a new user
func (r *userRepository) Create(ctx context.Context, user *entity.User) error {
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

// GetByID retrieves a user by ID
func (r *userRepository) GetByID(ctx context.Context, id uint) (*entity.User, error) {
	var user entity.User
	err := r.db.WithContext(ctx).
		Preload("Profile").
		First(&user, id).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

// GetByEmail retrieves a user by email
func (r *userRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	var user entity.User
	err := r.db.WithContext(ctx).
		Preload("Profile").
		Where("email = ?", email).
		First(&user).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

// Update updates a user
func (r *userRepository) Update(ctx context.Context, user *entity.User) error {
	if err := r.db.WithContext(ctx).Save(user).Error; err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

// Delete soft deletes a user
func (r *userRepository) Delete(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&entity.User{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

// List retrieves users with filtering
func (r *userRepository) List(ctx context.Context, filter *repository.UserFilter) ([]*entity.User, error) {
	query := r.db.WithContext(ctx).Model(&entity.User{})

	// Apply filters
	if filter.Status != nil {
		query = query.Where("status = ?", *filter.Status)
	}

	if filter.Search != nil && *filter.Search != "" {
		searchTerm := "%" + *filter.Search + "%"
		query = query.Where("name ILIKE ? OR email ILIKE ?", searchTerm, searchTerm)
	}

	if filter.CreatedAfter != nil {
		query = query.Where("created_at >= ?", *filter.CreatedAfter)
	}

	if filter.CreatedBefore != nil {
		query = query.Where("created_at <= ?", *filter.CreatedBefore)
	}

	// Apply sorting
	if filter.SortBy != "" {
		order := filter.SortBy
		if filter.SortOrder == "desc" {
			order += " DESC"
		}
		query = query.Order(order)
	} else {
		query = query.Order("created_at DESC")
	}

	// Apply pagination
	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
	}
	if filter.Offset > 0 {
		query = query.Offset(filter.Offset)
	}

	var users []*entity.User
	if err := query.Preload("Profile").Find(&users).Error; err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	return users, nil
}

// Count counts users with filtering
func (r *userRepository) Count(ctx context.Context, filter *repository.UserFilter) (int64, error) {
	query := r.db.WithContext(ctx).Model(&entity.User{})

	// Apply filters (same as List method)
	if filter.Status != nil {
		query = query.Where("status = ?", *filter.Status)
	}

	if filter.Search != nil && *filter.Search != "" {
		searchTerm := "%" + *filter.Search + "%"
		query = query.Where("name ILIKE ? OR email ILIKE ?", searchTerm, searchTerm)
	}

	if filter.CreatedAfter != nil {
		query = query.Where("created_at >= ?", *filter.CreatedAfter)
	}

	if filter.CreatedBefore != nil {
		query = query.Where("created_at <= ?", *filter.CreatedBefore)
	}

	var count int64
	if err := query.Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count users: %w", err)
	}

	return count, nil
}

// ExistsByEmail checks if a user exists by email
func (r *userRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entity.User{}).
		Where("email = ?", email).
		Count(&count).Error

	if err != nil {
		return false, fmt.Errorf("failed to check user existence: %w", err)
	}

	return count > 0, nil
}

// Profile operations

// CreateProfile creates a user profile
func (r *userRepository) CreateProfile(ctx context.Context, profile *entity.Profile) error {
	if err := r.db.WithContext(ctx).Create(profile).Error; err != nil {
		return fmt.Errorf("failed to create profile: %w", err)
	}
	return nil
}

// UpdateProfile updates a user profile
func (r *userRepository) UpdateProfile(ctx context.Context, profile *entity.Profile) error {
	if err := r.db.WithContext(ctx).Save(profile).Error; err != nil {
		return fmt.Errorf("failed to update profile: %w", err)
	}
	return nil
}

// GetProfile retrieves a user profile
func (r *userRepository) GetProfile(ctx context.Context, userID uint) (*entity.Profile, error) {
	var profile entity.Profile
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		First(&profile).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("profile not found")
		}
		return nil, fmt.Errorf("failed to get profile: %w", err)
	}

	return &profile, nil
}

// Role operations

// AssignRole assigns a role to a user
func (r *userRepository) AssignRole(ctx context.Context, userRole *entity.UserRole) error {
	// Check if role already exists
	var existing entity.UserRole
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND role = ?", userRole.UserID, userRole.Role).
		First(&existing).Error

	if err == nil {
		return fmt.Errorf("role already assigned to user")
	}

	if err != gorm.ErrRecordNotFound {
		return fmt.Errorf("failed to check existing role: %w", err)
	}

	// Create new role assignment
	if err := r.db.WithContext(ctx).Create(userRole).Error; err != nil {
		return fmt.Errorf("failed to assign role: %w", err)
	}

	return nil
}

// RevokeRole revokes a role from a user
func (r *userRepository) RevokeRole(ctx context.Context, userID uint, role entity.Role) error {
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND role = ?", userID, role).
		Delete(&entity.UserRole{}).Error

	if err != nil {
		return fmt.Errorf("failed to revoke role: %w", err)
	}

	return nil
}

// GetRoles retrieves all roles for a user
func (r *userRepository) GetRoles(ctx context.Context, userID uint) ([]*entity.UserRole, error) {
	var roles []*entity.UserRole
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Find(&roles).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get user roles: %w", err)
	}

	return roles, nil
}

// Session operations

// CreateSession creates a user session
func (r *userRepository) CreateSession(ctx context.Context, session *entity.UserSession) error {
	if err := r.db.WithContext(ctx).Create(session).Error; err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}
	return nil
}

// GetSession retrieves a session by token
func (r *userRepository) GetSession(ctx context.Context, token string) (*entity.UserSession, error) {
	var session entity.UserSession
	err := r.db.WithContext(ctx).
		Preload("User").
		Where("session_token = ? AND is_active = ? AND expires_at > ?", token, true, time.Now()).
		First(&session).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("session not found or expired")
		}
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	return &session, nil
}

// UpdateSession updates a session
func (r *userRepository) UpdateSession(ctx context.Context, session *entity.UserSession) error {
	if err := r.db.WithContext(ctx).Save(session).Error; err != nil {
		return fmt.Errorf("failed to update session: %w", err)
	}
	return nil
}

// DeleteSession deletes a session by token
func (r *userRepository) DeleteSession(ctx context.Context, token string) error {
	err := r.db.WithContext(ctx).
		Where("session_token = ?", token).
		Delete(&entity.UserSession{}).Error

	if err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}

	return nil
}

// DeleteExpiredSessions deletes all expired sessions
func (r *userRepository) DeleteExpiredSessions(ctx context.Context) error {
	err := r.db.WithContext(ctx).
		Where("expires_at <= ? OR is_active = ?", time.Now(), false).
		Delete(&entity.UserSession{}).Error

	if err != nil {
		return fmt.Errorf("failed to delete expired sessions: %w", err)
	}

	return nil
}

// GetUserSessions retrieves all active sessions for a user
func (r *userRepository) GetUserSessions(ctx context.Context, userID uint) ([]*entity.UserSession, error) {
	var sessions []*entity.UserSession
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND is_active = ? AND expires_at > ?", userID, true, time.Now()).
		Order("last_activity DESC").
		Find(&sessions).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get user sessions: %w", err)
	}

	return sessions, nil
}