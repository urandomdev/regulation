package account

import (
	"regulation/internal/request_context"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

// Me returns the currently authenticated user's information
// @Route GET /account/me
func (h *Handler) Me(ctx fiber.Ctx) (*UserResponse, error) {
	user := request_context.MustGetUser(ctx)

	return &UserResponse{
		ID:       user.ID,
		Email:    user.Email,
		Nickname: user.Nickname,
	}, nil
}

type UserResponse struct {
	ID       uuid.UUID `cbor:"id"`
	Email    string    `cbor:"email"`
	Nickname string    `cbor:"nickname"`
}
