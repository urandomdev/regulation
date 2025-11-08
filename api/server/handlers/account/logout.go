package account

import (
	"fmt"

	"regulation/server/middleware"
	"regulation/server/services/request_context"

	"github.com/gofiber/fiber/v3"
)

// Logout destroys the current session
// @Route POST /account/logout
func (h *Handler) Logout(ctx fiber.Ctx) error {
	sess := request_context.Session(ctx)

	// Delete session from Redis
	if err := h.sessionManager.Delete(ctx, sess.ID); err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}

	// Clear session cookie
	middleware.ClearSessionCookie(ctx)

	return nil
}
