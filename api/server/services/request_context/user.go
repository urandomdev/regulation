//revive:disable:var-naming // package name mirrors directory naming used across server
package request_context

import (
	"context"

	"regulation/internal/ent"
	"regulation/internal/session"

	"github.com/gofiber/fiber/v3"
)

// User retrieves the authenticated user from the context
// Panics if user is not found
func User(ctx context.Context) *ent.User {
	user, ok := ctx.Value(keyUser).(*ent.User)
	if !ok {
		panic("value for keyUser not found in context")
	}
	return user
}

// UserOptional retrieves the authenticated user from the context
// Returns nil if user is not found
func UserOptional(ctx context.Context) *ent.User {
	user, ok := ctx.Value(keyUser).(*ent.User)
	if !ok {
		return nil
	}
	return user
}

// SetUser stores the authenticated user in the context
func SetUser(ctx fiber.Ctx, user *ent.User) {
	ctx.RequestCtx().SetUserValue(keyUser, user)
}

// Session retrieves the session from the context
// Panics if session is not found
func Session(ctx context.Context) *session.Session {
	sess, ok := ctx.Value(keySession).(*session.Session)
	if !ok {
		panic("value for keySession not found in context")
	}
	return sess
}

// SessionOptional retrieves the session from the context
// Returns nil if session is not found
func SessionOptional(ctx context.Context) *session.Session {
	sess, ok := ctx.Value(keySession).(*session.Session)
	if !ok {
		return nil
	}
	return sess
}

// SetSession stores the session in the context
func SetSession(ctx fiber.Ctx, sess *session.Session) {
	ctx.RequestCtx().SetUserValue(keySession, sess)
}
