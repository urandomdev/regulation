package rule

import (
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"

	entrule "regulation/internal/ent/rule"
	"regulation/internal/protocol"
	"regulation/server/services/request_context"
)

// DeleteRule deletes a rule
// @Route DELETE /rules/:id
func (h *Handler) DeleteRule(ctx fiber.Ctx) error {
	session := request_context.Session(ctx)

	// Parse rule ID from path
	ruleID, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return protocol.ErrorResponse{
			Code:    protocol.InvalidParametersError,
			Message: "invalid rule ID format",
		}
	}

	// Delete the rule (verifying ownership via WHERE clause)
	deleted, err := h.db.Rule.
		Delete().
		Where(
			entrule.ID(ruleID),
			entrule.UserID(session.UserID),
		).
		Exec(ctx)

	if err != nil {
		return err
	}

	if deleted == 0 {
		return protocol.ErrorResponse{
			Code:    protocol.NotFoundError,
			Message: "rule not found",
		}
	}

	return nil
}
