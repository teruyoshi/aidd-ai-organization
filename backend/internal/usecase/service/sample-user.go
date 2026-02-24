package service

import (
	"context"
	"fmt"
	"sample-todo-backend/internal/domain/entity"
	"sample-todo-backend/internal/domain/repository"
	"sample-todo-backend/pkg/auth"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// UserService handles user-related business logic
type UserService struct {
	userRepo   repository.UserRepository
	validator  Validator
	jwtService *auth.JWTService
}

// NewUserService creates a new user service
func NewUserService(userRepo repository.UserRepository, validator Validator, jwtService *auth.JWTService) *UserService {
	return &UserService{
		userRepo:   userRepo,
		validator:  validator,
		jwtService: jwtService,
	}
}

// User management operations

// RegisterUser registers a new user
func (s *UserService) RegisterUser(ctx context.Context, req *RegisterUserRequest) (*UserResponse, error) {
	// Validate request
	if err := s.validator.Validate(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Check if email already exists
	exists, err := s.userRepo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to check email existence: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("email already registered")
	}

	// Hash password
	hashedPassword, err := s.hashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user entity
	user := &entity.User{
		Email:    req.Email,
		Password: hashedPassword,
		Name:     req.Name,
		Status:   entity.UserStatusActive,
	}

	// Validate business rules
	if err := user.Validate(); err != nil {
		return nil, fmt.Errorf("user validation failed: %w", err)
	}

	// Save user
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Create default profile
	profile := &entity.Profile{
		UserID:   user.ID,
		Timezone: "Asia/Tokyo",
		Language: "ja",
	}
	if err := s.userRepo.CreateProfile(ctx, profile); err != nil {
		// Log warning but don't fail registration
		// In real app, you might want to use proper logging
	}

	return s.toUserResponse(user, profile), nil
}

// AuthenticateUser authenticates a user and returns JWT token
func (s *UserService) AuthenticateUser(ctx context.Context, req *AuthenticateUserRequest) (*AuthResponse, error) {
	// Validate request
	if err := s.validator.Validate(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Get user by email
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Check if user is active
	if !user.IsActive() {
		return nil, fmt.Errorf("user account is not active")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Generate JWT token
	token, err := s.jwtService.GenerateToken(user.ID, user.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// Create session
	session := &entity.UserSession{
		UserID:       user.ID,
		SessionToken: token,
		IPAddress:    req.IPAddress,
		UserAgent:    req.UserAgent,
		ExpiresAt:    time.Now().Add(24 * time.Hour), // 24 hours
		LastActivity: time.Now(),
		IsActive:     true,
	}

	if err := s.userRepo.CreateSession(ctx, session); err != nil {
		// Log warning but don't fail authentication
	}

	return &AuthResponse{
		Token: token,
		User:  s.toUserResponse(user, user.Profile),
	}, nil
}

// GetUserProfile retrieves user profile
func (s *UserService) GetUserProfile(ctx context.Context, userID uint) (*UserResponse, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	profile, err := s.userRepo.GetProfile(ctx, userID)
	if err != nil {
		// Profile might not exist, create empty one
		profile = &entity.Profile{UserID: userID}
	}

	return s.toUserResponse(user, profile), nil
}

// UpdateUserProfile updates user profile
func (s *UserService) UpdateUserProfile(ctx context.Context, userID uint, req *UpdateProfileRequest) (*UserResponse, error) {
	// Validate request
	if err := s.validator.Validate(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Get existing user
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Update user fields if provided
	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Avatar != nil {
		user.Avatar = req.Avatar
	}

	// Validate and save user
	if err := user.Validate(); err != nil {
		return nil, fmt.Errorf("user validation failed: %w", err)
	}
	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	// Get or create profile
	profile, err := s.userRepo.GetProfile(ctx, userID)
	if err != nil {
		// Create new profile
		profile = &entity.Profile{
			UserID:   userID,
			Timezone: "Asia/Tokyo",
			Language: "ja",
		}
	}

	// Update profile fields if provided
	if req.Bio != nil {
		profile.Bio = req.Bio
	}
	if req.Location != nil {
		profile.Location = req.Location
	}
	if req.Website != nil {
		profile.Website = req.Website
	}
	if req.DateOfBirth != nil {
		profile.DateOfBirth = req.DateOfBirth
	}
	if req.Timezone != "" {
		profile.Timezone = req.Timezone
	}
	if req.Language != "" {
		profile.Language = req.Language
	}

	// Save profile
	if profile.ID == 0 {
		if err := s.userRepo.CreateProfile(ctx, profile); err != nil {
			return nil, fmt.Errorf("failed to create profile: %w", err)
		}
	} else {
		if err := s.userRepo.UpdateProfile(ctx, profile); err != nil {
			return nil, fmt.Errorf("failed to update profile: %w", err)
		}
	}

	return s.toUserResponse(user, profile), nil
}

// ChangePassword changes user password
func (s *UserService) ChangePassword(ctx context.Context, userID uint, req *ChangePasswordRequest) error {
	// Validate request
	if err := s.validator.Validate(req); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Get user
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Verify current password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.CurrentPassword)); err != nil {
		return fmt.Errorf("current password is incorrect")
	}

	// Hash new password
	hashedPassword, err := s.hashPassword(req.NewPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Update password
	user.Password = hashedPassword
	if err := s.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	// Invalidate all sessions except current one (optional)
	// This would require session token in the request to exclude current session

	return nil
}

// Session management

// ValidateSession validates a session token
func (s *UserService) ValidateSession(ctx context.Context, token string) (*UserResponse, error) {
	session, err := s.userRepo.GetSession(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("invalid session: %w", err)
	}

	// Update last activity
	session.LastActivity = time.Now()
	if err := s.userRepo.UpdateSession(ctx, session); err != nil {
		// Log warning but don't fail validation
	}

	profile, _ := s.userRepo.GetProfile(ctx, session.UserID)
	if profile == nil {
		profile = &entity.Profile{UserID: session.UserID}
	}

	return s.toUserResponse(&session.User, profile), nil
}

// LogoutUser logs out a user by invalidating session
func (s *UserService) LogoutUser(ctx context.Context, token string) error {
	return s.userRepo.DeleteSession(ctx, token)
}

// GetUserSessions retrieves all active sessions for a user
func (s *UserService) GetUserSessions(ctx context.Context, userID uint) ([]*SessionResponse, error) {
	sessions, err := s.userRepo.GetUserSessions(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get sessions: %w", err)
	}

	responses := make([]*SessionResponse, len(sessions))
	for i, session := range sessions {
		responses[i] = &SessionResponse{
			ID:           session.ID,
			IPAddress:    session.IPAddress,
			UserAgent:    session.UserAgent,
			LastActivity: session.LastActivity,
			ExpiresAt:    session.ExpiresAt,
			IsActive:     session.IsActive,
		}
	}

	return responses, nil
}

// RevokeSession revokes a specific session
func (s *UserService) RevokeSession(ctx context.Context, userID uint, sessionID uint) error {
	// First verify the session belongs to the user
	sessions, err := s.userRepo.GetUserSessions(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get sessions: %w", err)
	}

	var targetSession *entity.UserSession
	for _, session := range sessions {
		if session.ID == sessionID {
			targetSession = session
			break
		}
	}

	if targetSession == nil {
		return fmt.Errorf("session not found")
	}

	return s.userRepo.DeleteSession(ctx, targetSession.SessionToken)
}

// Utility methods

// hashPassword hashes a password using bcrypt
func (s *UserService) hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// toUserResponse converts entity.User to UserResponse
func (s *UserService) toUserResponse(user *entity.User, profile *entity.Profile) *UserResponse {
	resp := &UserResponse{
		ID        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		Avatar:    user.Avatar,
		Status:    string(user.Status),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	if profile != nil {
		resp.Profile = &ProfileResponse{
			ID:          profile.ID,
			Bio:         profile.Bio,
			Location:    profile.Location,
			Website:     profile.Website,
			DateOfBirth: profile.DateOfBirth,
			Timezone:    profile.Timezone,
			Language:    profile.Language,
			CreatedAt:   profile.CreatedAt,
			UpdatedAt:   profile.UpdatedAt,
		}
	}

	return resp
}

// Cleanup operations

// CleanupExpiredSessions removes expired sessions
func (s *UserService) CleanupExpiredSessions(ctx context.Context) error {
	return s.userRepo.DeleteExpiredSessions(ctx)
}