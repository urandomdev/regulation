//revive:disable:var-naming // package name mirrors directory naming used across server
package request_context

import (
	"regulation/internal/ent"
	"regulation/internal/session"

	"github.com/gofiber/fiber/v3"
)

const (
	userKey    = "user"
	sessionKey = "session"
)

// SetUser stores the authenticated user in the context
func SetUser(ctx fiber.Ctx, user *ent.User) {
	ctx.Locals(userKey, user)
}

// GetUser retrieves the authenticated user from the context
func GetUser(ctx fiber.Ctx) (*ent.User, bool) {
	user, ok := ctx.Locals(userKey).(*ent.User)
	return user, ok
}

// MustGetUser retrieves the authenticated user from the context, panics if not found
func MustGetUser(ctx fiber.Ctx) *ent.User {
	user, ok := GetUser(ctx)
	if !ok {
		panic("user not found in context")
	}
	return user
}

// SetSession stores the session in the context
func SetSession(ctx fiber.Ctx, sess *session.Session) {
	ctx.Locals(sessionKey, sess)
}

// GetSession retrieves the session from the context
func GetSession(ctx fiber.Ctx) (*session.Session, bool) {
	sess, ok := ctx.Locals(sessionKey).(*session.Session)
	return sess, ok
}
