package entity

import (
	"time"

	"gorm.io/gorm"
)

// User represents a user entity
type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Email     string         `gorm:"uniqueIndex:idx_users_email;size:255;not null" json:"email" validate:"required,email"`
	Password  string         `gorm:"size:255;not null" json:"-"` // JSON出力時は除外
	Name      string         `gorm:"size:100;not null" json:"name" validate:"required,min=2,max=100"`
	Avatar    *string        `gorm:"size:500" json:"avatar,omitempty"`
	Status    UserStatus     `gorm:"type:enum('active','inactive','suspended');default:'active'" json:"status"`
	Profile   *Profile       `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"profile,omitempty"`
	Todos     []Todo         `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"todos,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// UserStatus represents the status of a user
type UserStatus string

const (
	UserStatusActive    UserStatus = "active"
	UserStatusInactive  UserStatus = "inactive"
	UserStatusSuspended UserStatus = "suspended"
)

// Profile represents additional user profile information
type Profile struct {
	ID          uint    `gorm:"primaryKey" json:"id"`
	UserID      uint    `gorm:"not null;index" json:"user_id"`
	Bio         *string `gorm:"type:text" json:"bio,omitempty"`
	Location    *string `gorm:"size:100" json:"location,omitempty"`
	Website     *string `gorm:"size:500" json:"website,omitempty"`
	DateOfBirth *time.Time `json:"date_of_birth,omitempty"`
	Timezone    string  `gorm:"size:50;default:'Asia/Tokyo'" json:"timezone"`
	Language    string  `gorm:"size:10;default:'ja'" json:"language"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// UserRole represents user roles for authorization
type UserRole struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	UserID      uint           `gorm:"not null;index" json:"user_id"`
	Role        Role           `gorm:"type:enum('admin','user','moderator');default:'user'" json:"role"`
	GrantedBy   *uint          `gorm:"index" json:"granted_by,omitempty"`
	GrantedAt   time.Time      `gorm:"not null" json:"granted_at"`
	ExpiresAt   *time.Time     `json:"expires_at,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	User      User  `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	GrantedByUser *User `gorm:"foreignKey:GrantedBy"`
}

type Role string

const (
	RoleAdmin     Role = "admin"
	RoleUser      Role = "user"
	RoleModerator Role = "moderator"
)

// UserSession represents user session information
type UserSession struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	UserID       uint      `gorm:"not null;index" json:"user_id"`
	SessionToken string    `gorm:"size:255;not null;uniqueIndex" json:"-"`
	IPAddress    string    `gorm:"size:45" json:"ip_address"`
	UserAgent    string    `gorm:"size:500" json:"user_agent"`
	ExpiresAt    time.Time `gorm:"not null;index" json:"expires_at"`
	LastActivity time.Time `gorm:"not null" json:"last_activity"`
	IsActive     bool      `gorm:"default:true" json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	// Relations
	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

// UserMethods contains business logic methods for User entity
type UserMethods struct{}

// HasRole checks if user has a specific role
func (u *User) HasRole(role Role) bool {
	// この実装では簡略化。実際はUserRoleテーブルを確認
	return role == RoleUser // デフォルトでは全てのユーザーはuser role
}

// IsActive checks if user is active
func (u *User) IsActive() bool {
	return u.Status == UserStatusActive
}

// GetDisplayName returns the display name (name or email)
func (u *User) GetDisplayName() string {
	if u.Name != "" {
		return u.Name
	}
	return u.Email
}

// GetRoles returns all roles assigned to the user
func (u *User) GetRoles() []Role {
	// 実際の実装ではデータベースから取得
	return []Role{RoleUser}
}

// Validate performs business logic validation
func (u *User) Validate() error {
	if u.Email == "" {
		return NewValidationError("email is required")
	}

	if u.Name == "" {
		return NewValidationError("name is required")
	}

	if len(u.Name) < 2 {
		return NewValidationError("name must be at least 2 characters")
	}

	return nil
}

// TableName specifies the table name for GORM
func (User) TableName() string {
	return "users"
}

func (Profile) TableName() string {
	return "user_profiles"
}

func (UserRole) TableName() string {
	return "user_roles"
}

func (UserSession) TableName() string {
	return "user_sessions"
}

// ValidationError represents a validation error
type ValidationError struct {
	Message string
	Field   string
}

func (e ValidationError) Error() string {
	return e.Message
}

func NewValidationError(message string) *ValidationError {
	return &ValidationError{Message: message}
}