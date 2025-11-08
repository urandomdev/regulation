package middleware

import (
	"fmt"

	"regulation/internal/ent"
	"regulation/internal/protocol"
	"regulation/internal/session"
	"regulation/server/services/request_context"

	"github.com/gofiber/fiber/v3"
)

const sessionCookieName = "sessionid"

// Auth middleware validates the session cookie and loads the user
type Auth struct {
	db             *ent.Client
	sessionManager *session.Manager
}

// NewAuth creates a new authentication middleware
func NewAuth(db *ent.Client, sessionManager *session.Manager) *Auth {
	return &Auth{
		db:             db,
		sessionManager: sessionManager,
	}
}

// Handle validates the session and loads the authenticated user
func (a *Auth) Handle(ctx fiber.Ctx) error {
	sessionID := ctx.Cookies(sessionCookieName)
	if sessionID == "" {
		return protocol.ErrorResponse{
			Code:    protocol.UnauthorizedError,
			Message: "authentication required",
		}
	}

	sess, err := a.sessionManager.Get(ctx, sessionID)
	if err != nil {
		return protocol.ErrorResponse{
			Code:    protocol.UnauthorizedError,
			Message: "invalid or expired session",
		}
	}

	user, err := a.db.User.Get(ctx, sess.UserID)
	if err != nil {
		if ent.IsNotFound(err) {
			// User was deleted, invalidate session
			_ = a.sessionManager.Delete(ctx, sessionID)
			return protocol.ErrorResponse{
				Code:    protocol.UnauthorizedError,
				Message: "user not found",
			}
		}
		return fmt.Errorf("failed to load user: %w", err)
	}

	request_context.SetUser(ctx, user)
	request_context.SetSession(ctx, sess)

	return ctx.Next()
}

// SetSessionCookie sets the session cookie on the response
func SetSessionCookie(ctx fiber.Ctx, sessionID string) {
	ctx.Cookie(&fiber.Cookie{
		Name:     sessionCookieName,
		Value:    sessionID,
		Path:     "/",
		HTTPOnly: true,
		Secure:   true, // Set to true in production with HTTPS
		SameSite: fiber.CookieSameSiteStrictMode,
		MaxAge:   86400, // 24 hours in seconds
	})
}

// ClearSessionCookie removes the session cookie
func ClearSessionCookie(ctx fiber.Ctx) {
	ctx.Cookie(&fiber.Cookie{
		Name:     sessionCookieName,
		Value:    "",
		Path:     "/",
		HTTPOnly: true,
		Secure:   true,
		SameSite: fiber.CookieSameSiteStrictMode,
		MaxAge:   -1, // Expire immediately
	})
}
