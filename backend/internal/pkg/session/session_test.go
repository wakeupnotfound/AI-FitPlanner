package session

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestRedis(t *testing.T) (*redis.Client, *miniredis.Miniredis) {
	mr, err := miniredis.Run()
	require.NoError(t, err)

	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	return client, mr
}

func TestCreateSession(t *testing.T) {
	client, mr := setupTestRedis(t)
	defer mr.Close()
	defer client.Close()

	manager := NewSessionManager(client)
	ctx := context.Background()

	userID := int64(123)
	sessionID := "test-session-id"
	username := "testuser"
	ttl := time.Hour
	ipAddress := "192.168.1.1"
	userAgent := "Mozilla/5.0"

	err := manager.CreateSession(ctx, userID, sessionID, username, ttl, ipAddress, userAgent)
	assert.NoError(t, err)

	// Verify session was created
	session, err := manager.GetSession(ctx, sessionID)
	assert.NoError(t, err)
	assert.NotNil(t, session)
	assert.Equal(t, sessionID, session.SessionID)
	assert.Equal(t, userID, session.UserID)
	assert.Equal(t, username, session.Username)
	assert.Equal(t, ipAddress, session.IPAddress)
	assert.Equal(t, userAgent, session.UserAgent)
}

func TestGetSession_NotFound(t *testing.T) {
	client, mr := setupTestRedis(t)
	defer mr.Close()
	defer client.Close()

	manager := NewSessionManager(client)
	ctx := context.Background()

	session, err := manager.GetSession(ctx, "non-existent-session")
	assert.NoError(t, err)
	assert.Nil(t, session)
}

func TestDeleteSession(t *testing.T) {
	client, mr := setupTestRedis(t)
	defer mr.Close()
	defer client.Close()

	manager := NewSessionManager(client)
	ctx := context.Background()

	userID := int64(123)
	sessionID := "test-session-id"
	username := "testuser"
	ttl := time.Hour

	// Create session
	err := manager.CreateSession(ctx, userID, sessionID, username, ttl, "", "")
	require.NoError(t, err)

	// Delete session
	err = manager.DeleteSession(ctx, sessionID)
	assert.NoError(t, err)

	// Verify session was deleted
	session, err := manager.GetSession(ctx, sessionID)
	assert.NoError(t, err)
	assert.Nil(t, session)
}

func TestDeleteSession_NotFound(t *testing.T) {
	client, mr := setupTestRedis(t)
	defer mr.Close()
	defer client.Close()

	manager := NewSessionManager(client)
	ctx := context.Background()

	// Deleting non-existent session should not error
	err := manager.DeleteSession(ctx, "non-existent-session")
	assert.NoError(t, err)
}

func TestDeleteAllUserSessions(t *testing.T) {
	client, mr := setupTestRedis(t)
	defer mr.Close()
	defer client.Close()

	manager := NewSessionManager(client)
	ctx := context.Background()

	userID := int64(123)
	ttl := time.Hour

	// Create multiple sessions for the same user
	session1 := "session-1"
	session2 := "session-2"
	session3 := "session-3"

	err := manager.CreateSession(ctx, userID, session1, "testuser", ttl, "", "")
	require.NoError(t, err)
	err = manager.CreateSession(ctx, userID, session2, "testuser", ttl, "", "")
	require.NoError(t, err)
	err = manager.CreateSession(ctx, userID, session3, "testuser", ttl, "", "")
	require.NoError(t, err)

	// Delete all user sessions
	err = manager.DeleteAllUserSessions(ctx, userID)
	assert.NoError(t, err)

	// Verify all sessions were deleted
	s1, err := manager.GetSession(ctx, session1)
	assert.NoError(t, err)
	assert.Nil(t, s1)

	s2, err := manager.GetSession(ctx, session2)
	assert.NoError(t, err)
	assert.Nil(t, s2)

	s3, err := manager.GetSession(ctx, session3)
	assert.NoError(t, err)
	assert.Nil(t, s3)
}

func TestDeleteAllUserSessions_NoSessions(t *testing.T) {
	client, mr := setupTestRedis(t)
	defer mr.Close()
	defer client.Close()

	manager := NewSessionManager(client)
	ctx := context.Background()

	// Deleting sessions for user with no sessions should not error
	err := manager.DeleteAllUserSessions(ctx, int64(999))
	assert.NoError(t, err)
}

func TestSessionExpiration(t *testing.T) {
	client, mr := setupTestRedis(t)
	defer mr.Close()
	defer client.Close()

	manager := NewSessionManager(client)
	ctx := context.Background()

	userID := int64(123)
	sessionID := "test-session-id"
	username := "testuser"
	ttl := 2 * time.Second

	// Create session with short TTL
	err := manager.CreateSession(ctx, userID, sessionID, username, ttl, "", "")
	require.NoError(t, err)

	// Session should exist immediately
	session, err := manager.GetSession(ctx, sessionID)
	assert.NoError(t, err)
	assert.NotNil(t, session)

	// Fast-forward time in miniredis
	mr.FastForward(3 * time.Second)

	// Session should be expired
	session, err = manager.GetSession(ctx, sessionID)
	assert.NoError(t, err)
	assert.Nil(t, session)
}
