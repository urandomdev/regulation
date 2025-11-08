package advisor

import (
	"regulation/internal/protocol"

	advisorsvc "regulation/internal/advisor"

	"github.com/gofiber/fiber/v3"
)

// BudgetPlan is a lightweight test endpoint that lets clients verify
// GPT-5-mini structured outputs without going through the authenticated flow.
func (h *Handler) BudgetPlan(ctx fiber.Ctx, req *advisorsvc.BudgetPlanRequest) (*advisorsvc.BudgetPlanResponse, error) {
	if h.service == nil {
		return nil, protocol.ErrorResponse{
			Code:    protocol.InternalError,
			Message: "budget advisor is currently unavailable",
		}
	}

	resp, err := h.service.GeneratePlan(ctx.Context(), req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
