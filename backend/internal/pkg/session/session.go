package session

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/ai-fitness-planner/backend/internal/model"
	"github.com/redis/go-redis/v9"
)

// SessionManager interface defines methods for session management
type SessionManager interface {
	CreateSession(ctx context.Context, userID int64, sessionID string, username string, ttl time.Duration, ipAddress string, userAgent string) error
	GetSession(ctx context.Context, sessionID string) (*model.Session, error)
	DeleteSession(ctx context.Context, sessionID string) error
	DeleteAllUserSessions(ctx context.Context, userID int64) error
}

// RedisSessionManager implements SessionManager using Redis
type RedisSessionManager struct {
	client *redis.Client
}

// NewSessionManager creates a new session manager with Redis client
func NewSessionManager(client *redis.Client) SessionManager {
	return &RedisSessionManager{
		client: client,
	}
}

// CreateSession creates a new session in Redis
func (m *RedisSessionManager) CreateSession(ctx context.Context, userID int64, sessionID string, username string, ttl time.Duration, ipAddress string, userAgent string) error {
	session := &model.Session{
		SessionID: sessionID,
		UserID:    userID,
		Username:  username,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(ttl),
		IPAddress: ipAddress,
		UserAgent: userAgent,
	}

	// Serialize session to JSON
	sessionData, err := json.Marshal(session)
	if err != nil {
		return fmt.Errorf("failed to marshal session: %w", err)
	}

	// Store session in Redis with TTL
	sessionKey := fmt.Sprintf("session:%s", sessionID)
	if err := m.client.Set(ctx, sessionKey, sessionData, ttl).Err(); err != nil {
		return fmt.Errorf("failed to store session in Redis: %w", err)
	}

	// Add session to user's session set for tracking all user sessions
	userSessionsKey := fmt.Sprintf("user_sessions:%d", userID)
	if err := m.client.SAdd(ctx, userSessionsKey, sessionID).Err(); err != nil {
		return fmt.Errorf("failed to add session to user sessions set: %w", err)
	}

	// Set TTL on user sessions set (slightly longer than session TTL)
	if err := m.client.Expire(ctx, userSessionsKey, ttl+time.Hour).Err(); err != nil {
		return fmt.Errorf("failed to set TTL on user sessions set: %w", err)
	}

	return nil
}

// GetSession retrieves a session from Redis
func (m *RedisSessionManager) GetSession(ctx context.Context, sessionID string) (*model.Session, error) {
	sessionKey := fmt.Sprintf("session:%s", sessionID)

	sessionData, err := m.client.Get(ctx, sessionKey).Result()
	if err == redis.Nil {
		return nil, nil // Session not found
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get session from Redis: %w", err)
	}

	var session model.Session
	if err := json.Unmarshal([]byte(sessionData), &session); err != nil {
		return nil, fmt.Errorf("failed to unmarshal session: %w", err)
	}

	return &session, nil
}

// DeleteSession deletes a session from Redis
func (m *RedisSessionManager) DeleteSession(ctx context.Context, sessionID string) error {
	sessionKey := fmt.Sprintf("session:%s", sessionID)

	// Get session to find user ID before deleting
	session, err := m.GetSession(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("failed to get session before deletion: %w", err)
	}

	// If session doesn't exist, nothing to delete
	if session == nil {
		return nil
	}

	// Delete session from Redis
	if err := m.client.Del(ctx, sessionKey).Err(); err != nil {
		return fmt.Errorf("failed to delete session from Redis: %w", err)
	}

	// Remove session from user's session set
	userSessionsKey := fmt.Sprintf("user_sessions:%d", session.UserID)
	if err := m.client.SRem(ctx, userSessionsKey, sessionID).Err(); err != nil {
		return fmt.Errorf("failed to remove session from user sessions set: %w", err)
	}

	return nil
}

// DeleteAllUserSessions deletes all sessions for a specific user
func (m *RedisSessionManager) DeleteAllUserSessions(ctx context.Context, userID int64) error {
	userSessionsKey := fmt.Sprintf("user_sessions:%d", userID)

	// Get all session IDs for the user
	sessionIDs, err := m.client.SMembers(ctx, userSessionsKey).Result()
	if err != nil {
		return fmt.Errorf("failed to get user sessions: %w", err)
	}

	// Delete each session
	for _, sessionID := range sessionIDs {
		sessionKey := fmt.Sprintf("session:%s", sessionID)
		if err := m.client.Del(ctx, sessionKey).Err(); err != nil {
			// Log error but continue deleting other sessions
			fmt.Printf("Warning: failed to delete session %s: %v\n", sessionID, err)
		}
	}

	// Delete the user sessions set
	if err := m.client.Del(ctx, userSessionsKey).Err(); err != nil {
		return fmt.Errorf("failed to delete user sessions set: %w", err)
	}

	return nil
}
