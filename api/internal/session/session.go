package session

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/rueidis"
)

const (
	sessionKeyPrefix = "session:"
	sessionTTL       = 24 * time.Hour
)

// Manager handles session storage and retrieval
type Manager struct {
	cache rueidis.Client
}

// NewManager creates a new session manager
func NewManager(cache rueidis.Client) *Manager {
	return &Manager{
		cache: cache,
	}
}

// Session represents a user session
type Session struct {
	ID           string
	UserID       uuid.UUID
	CreatedAt    time.Time
	LastAccessed time.Time
}

// Create creates a new session for the given user ID
func (m *Manager) Create(ctx context.Context, userID uuid.UUID) (*Session, error) {
	sessionID, err := generateSessionID()
	if err != nil {
		return nil, fmt.Errorf("failed to generate session ID: %w", err)
	}

	session := &Session{
		ID:           sessionID,
		UserID:       userID,
		CreatedAt:    time.Now().UTC(),
		LastAccessed: time.Now().UTC(),
	}

	if err := m.save(ctx, session); err != nil {
		return nil, err
	}

	return session, nil
}

// Get retrieves a session by ID and refreshes its TTL
func (m *Manager) Get(ctx context.Context, sessionID string) (*Session, error) {
	key := sessionKey(sessionID)

	cmd := m.cache.B().Hgetall().Key(key).Build()
	result, err := m.cache.Do(ctx, cmd).AsStrMap()
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	if len(result) == 0 {
		return nil, errors.New("session not found")
	}

	userID, err := uuid.Parse(result["user_id"])
	if err != nil {
		return nil, fmt.Errorf("invalid user ID in session: %w", err)
	}

	createdAt, err := time.Parse(time.RFC3339, result["created_at"])
	if err != nil {
		return nil, fmt.Errorf("invalid created_at in session: %w", err)
	}

	session := &Session{
		ID:           sessionID,
		UserID:       userID,
		CreatedAt:    createdAt,
		LastAccessed: time.Now().UTC(),
	}

	// Refresh session TTL
	if err := m.save(ctx, session); err != nil {
		return nil, fmt.Errorf("failed to refresh session: %w", err)
	}

	return session, nil
}

// Delete removes a session from storage
func (m *Manager) Delete(ctx context.Context, sessionID string) error {
	key := sessionKey(sessionID)

	cmd := m.cache.B().Del().Key(key).Build()
	if err := m.cache.Do(ctx, cmd).Error(); err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}

	return nil
}

// save stores or updates a session in Redis
func (m *Manager) save(ctx context.Context, session *Session) error {
	key := sessionKey(session.ID)

	cmds := make(rueidis.Commands, 0, 2)
	cmds = append(cmds,
		m.cache.B().Hset().
			Key(key).
			FieldValue().
			FieldValue("user_id", session.UserID.String()).
			FieldValue("created_at", session.CreatedAt.Format(time.RFC3339)).
			FieldValue("last_accessed", session.LastAccessed.Format(time.RFC3339)).
			Build(),
		m.cache.B().Expire().Key(key).Seconds(int64(sessionTTL.Seconds())).Build(),
	)

	for _, resp := range m.cache.DoMulti(ctx, cmds...) {
		if err := resp.Error(); err != nil {
			return fmt.Errorf("failed to save session: %w", err)
		}
	}

	return nil
}

// sessionKey returns the Redis key for a session
func sessionKey(sessionID string) string {
	return sessionKeyPrefix + sessionID
}

// generateSessionID creates a cryptographically secure random session ID
func generateSessionID() (string, error) {
	bytes := make([]byte, 32) // 256 bits
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
