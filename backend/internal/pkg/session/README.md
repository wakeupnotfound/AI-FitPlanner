# Session Management Package

This package provides session management functionality using Redis as the backend store.

## Features

- Create sessions with TTL (Time To Live)
- Retrieve session information
- Delete individual sessions
- Delete all sessions for a specific user
- Automatic session expiration via Redis TTL

## Usage

```go
import (
    "context"
    "time"
    "github.com/ai-fitness-planner/backend/internal/pkg/session"
    "github.com/redis/go-redis/v9"
)

// Create a session manager
redisClient := redis.NewClient(&redis.Options{
    Addr: "localhost:6379",
})
manager := session.NewSessionManager(redisClient)

// Create a session
ctx := context.Background()
err := manager.CreateSession(
    ctx,
    userID,
    sessionID,
    username,
    time.Hour,      // TTL
    "192.168.1.1",  // IP address
    "Mozilla/5.0",  // User agent
)

// Get a session
session, err := manager.GetSession(ctx, sessionID)
if session != nil {
    fmt.Printf("User ID: %d\n", session.UserID)
}

// Delete a session
err = manager.DeleteSession(ctx, sessionID)

// Delete all sessions for a user
err = manager.DeleteAllUserSessions(ctx, userID)
```

## Implementation Details

### Redis Keys

- Session data: `session:{sessionID}` - Stores the full session object as JSON
- User sessions set: `user_sessions:{userID}` - Stores a set of session IDs for each user

### Session Model

The `Session` struct contains:
- `SessionID`: Unique identifier for the session
- `UserID`: ID of the user who owns the session
- `Username`: Username of the user
- `CreatedAt`: Timestamp when session was created
- `ExpiresAt`: Timestamp when session will expire
- `IPAddress`: IP address of the client (optional)
- `UserAgent`: User agent string of the client (optional)

### TTL Management

- Sessions automatically expire based on the TTL set during creation
- The user sessions set has a TTL slightly longer than the session TTL to ensure cleanup
- Redis handles automatic deletion of expired keys

## Testing

To run tests, you need to install test dependencies:

```bash
go get github.com/alicebob/miniredis/v2
go get github.com/stretchr/testify
go test -v ./internal/pkg/session/...
```

The tests use `miniredis` to provide an in-memory Redis implementation for testing without requiring a real Redis server.

## Requirements Validation

This implementation satisfies the following requirements:

- **Requirement 1.5**: Session invalidation on logout
- **Requirement 9.6**: Session expiration and cleanup

## Design Properties

This implementation supports the following correctness properties:

- **Property 3**: Session Invalidation - After logout, the session should not exist in Redis
