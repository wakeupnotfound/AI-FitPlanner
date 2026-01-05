package service

import (
	"context"
	"time"

	"github.com/ai-fitness-planner/backend/internal/errors"
	"github.com/ai-fitness-planner/backend/internal/model"
	"github.com/ai-fitness-planner/backend/internal/pkg/jwt"
	"github.com/ai-fitness-planner/backend/internal/pkg/session"
	"github.com/ai-fitness-planner/backend/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

// RegisterRequest represents the registration request data
type RegisterRequest struct {
	Username string  `json:"username" validate:"required,min=3,max=50"`
	Email    string  `json:"email" validate:"required,email,max=100"`
	Phone    *string `json:"phone" validate:"omitempty,max=20"`
	Password string  `json:"password" validate:"required,min=8,max=100"`
}

// LoginRequest represents the login request data
type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// AuthResponse represents the authentication response with tokens
type AuthResponse struct {
	AccessToken  string      `json:"access_token"`
	RefreshToken string      `json:"refresh_token"`
	User         *model.User `json:"user"`
}

// TokenResponse represents the token refresh response
type TokenResponse struct {
	AccessToken string `json:"access_token"`
}

// AuthService interface defines methods for authentication operations
type AuthService interface {
	Register(ctx context.Context, req *RegisterRequest) (*AuthResponse, error)
	Login(ctx context.Context, req *LoginRequest, ipAddress, userAgent string) (*AuthResponse, error)
	Logout(ctx context.Context, sessionID string) error
	RefreshToken(ctx context.Context, refreshToken string) (*TokenResponse, error)
	ValidateSession(ctx context.Context, sessionID string) (*model.Session, error)
}

// authService implements the AuthService interface
type authService struct {
	userRepo       repository.UserRepository
	jwtManager     jwt.JWTManager
	sessionManager session.SessionManager
}

// NewAuthService creates a new instance of AuthService
func NewAuthService(
	userRepo repository.UserRepository,
	jwtManager jwt.JWTManager,
	sessionManager session.SessionManager,
) AuthService {
	return &authService{
		userRepo:       userRepo,
		jwtManager:     jwtManager,
		sessionManager: sessionManager,
	}
}

// Register creates a new user account with encrypted password
// Validates: Requirements 1.1, 1.2
func (s *authService) Register(ctx context.Context, req *RegisterRequest) (*AuthResponse, error) {
	// Check if username already exists
	existingUser, err := s.userRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrDatabase, "failed to check username")
	}
	if existingUser != nil {
		return nil, errors.ErrUsernameExists
	}

	// Check if email already exists
	existingEmail, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrDatabase, "failed to check email")
	}
	if existingEmail != nil {
		return nil, errors.ErrEmailExists
	}

	// Hash password using bcrypt
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrInternalServer, "failed to hash password")
	}

	// Create user
	user := &model.User{
		Username:     req.Username,
		Email:        req.Email,
		Phone:        req.Phone,
		PasswordHash: string(passwordHash),
		Status:       1, // Active
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, errors.Wrap(err, errors.ErrDatabase, "failed to create user")
	}

	// Generate tokens
	accessToken, err := s.jwtManager.GenerateAccessToken(user.ID, user.Username)
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrInternalServer, "failed to generate access token")
	}

	refreshToken, err := s.jwtManager.GenerateRefreshToken(user.ID, user.Username)
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrInternalServer, "failed to generate refresh token")
	}

	// Extract session ID from access token
	claims, err := s.jwtManager.ValidateToken(accessToken)
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrInternalServer, "failed to validate generated token")
	}

	// Create session in Redis
	if err := s.sessionManager.CreateSession(
		ctx,
		user.ID,
		claims.SessionID,
		user.Username,
		time.Hour*24*7, // 7 days TTL
		"",             // IP address will be set by handler
		"",             // User agent will be set by handler
	); err != nil {
		return nil, errors.Wrap(err, errors.ErrCache, "failed to create session")
	}

	// Remove password hash from response
	user.PasswordHash = ""

	return &AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
	}, nil
}

// Login authenticates a user and returns tokens
// Validates: Requirements 1.2, 1.3
func (s *authService) Login(ctx context.Context, req *LoginRequest, ipAddress, userAgent string) (*AuthResponse, error) {
	// Get user by username
	user, err := s.userRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrDatabase, "failed to get user")
	}
	if user == nil {
		return nil, errors.New(errors.ErrInvalidCredentials, "invalid username or password")
	}

	// Check if user is disabled
	if user.Status != 1 {
		return nil, errors.ErrUserDisabled
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, errors.New(errors.ErrInvalidCredentials, "invalid username or password")
	}

	// Generate tokens
	accessToken, err := s.jwtManager.GenerateAccessToken(user.ID, user.Username)
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrInternalServer, "failed to generate access token")
	}

	refreshToken, err := s.jwtManager.GenerateRefreshToken(user.ID, user.Username)
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrInternalServer, "failed to generate refresh token")
	}

	// Extract session ID from access token
	claims, err := s.jwtManager.ValidateToken(accessToken)
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrInternalServer, "failed to validate generated token")
	}

	// Create session in Redis
	if err := s.sessionManager.CreateSession(
		ctx,
		user.ID,
		claims.SessionID,
		user.Username,
		time.Hour*24*7, // 7 days TTL
		ipAddress,
		userAgent,
	); err != nil {
		return nil, errors.Wrap(err, errors.ErrCache, "failed to create session")
	}

	// Remove password hash from response
	user.PasswordHash = ""

	return &AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
	}, nil
}

// Logout invalidates the current session
// Validates: Requirements 1.5
func (s *authService) Logout(ctx context.Context, sessionID string) error {
	if sessionID == "" {
		return errors.New(errors.ErrInvalidParam, "session ID is required")
	}

	// Delete session from Redis
	if err := s.sessionManager.DeleteSession(ctx, sessionID); err != nil {
		return errors.Wrap(err, errors.ErrCache, "failed to delete session")
	}

	return nil
}

// RefreshToken generates a new access token from a valid refresh token
// Validates: Requirements 1.4
func (s *authService) RefreshToken(ctx context.Context, refreshToken string) (*TokenResponse, error) {
	// Validate refresh token and extract claims
	claims, err := s.jwtManager.ValidateToken(refreshToken)
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrUnauthorized, "invalid refresh token")
	}

	// Verify it's a refresh token
	if claims.Type != "refresh" {
		return nil, errors.New(errors.ErrUnauthorized, "token is not a refresh token")
	}

	// Verify session still exists
	session, err := s.sessionManager.GetSession(ctx, claims.SessionID)
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrCache, "failed to get session")
	}
	if session == nil {
		return nil, errors.ErrSessionNotFound
	}

	// Generate new access token
	newAccessToken, err := s.jwtManager.GenerateAccessToken(claims.UserID, claims.Username)
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrInternalServer, "failed to generate access token")
	}

	return &TokenResponse{
		AccessToken: newAccessToken,
	}, nil
}

// ValidateSession checks if a session is valid
// Validates: Requirements 1.6
func (s *authService) ValidateSession(ctx context.Context, sessionID string) (*model.Session, error) {
	if sessionID == "" {
		return nil, errors.New(errors.ErrInvalidParam, "session ID is required")
	}

	session, err := s.sessionManager.GetSession(ctx, sessionID)
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrCache, "failed to get session")
	}

	if session == nil {
		return nil, errors.ErrSessionNotFound
	}

	// Check if session has expired
	if time.Now().After(session.ExpiresAt) {
		// Delete expired session
		_ = s.sessionManager.DeleteSession(ctx, sessionID)
		return nil, errors.New(errors.ErrTokenExpired, "session has expired")
	}

	return session, nil
}
