package account

import (
	"fmt"

	"regulation/internal/request_context"
	"regulation/server/middleware"

	"github.com/gofiber/fiber/v3"
)

// Logout destroys the current session
// @Route POST /account/logout
func (h *Handler) Logout(ctx fiber.Ctx) error {
	sess, ok := request_context.GetSession(ctx)
	if !ok {
		// Session not found in context, but cookie might still be set
		// Clear the cookie anyway
		middleware.ClearSessionCookie(ctx)
		return nil
	}

	// Delete session from Redis
	if err := h.sessionManager.Delete(ctx, sess.ID); err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}

	// Clear session cookie
	middleware.ClearSessionCookie(ctx)

	return nil
}
