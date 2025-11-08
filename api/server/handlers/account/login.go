package account

import (
	"fmt"

	"regulation/internal/ent"
	"regulation/internal/ent/user"
	"regulation/internal/protocol"
	"regulation/server/middleware"

	"github.com/DeltaLaboratory/password"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/gofiber/fiber/v3"
)

// Login authenticates a user and creates a session
// @Route POST /account/login
func (h *Handler) Login(ctx fiber.Ctx, req *LoginRequest) error {
	// Find user by email
	u, err := h.db.User.Query().
		Where(user.EmailEQ(req.Email)).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return protocol.ErrorResponse{
				Code:    protocol.InvalidCredentialsError,
				Message: "invalid email or password",
			}
		}
		return fmt.Errorf("failed to query user: %w", err)
	}

	// Verify password
	valid, err := password.Verify([]byte(req.Password), u.Password)
	if err != nil || !valid {
		return protocol.ErrorResponse{
			Code:    protocol.InvalidCredentialsError,
			Message: "invalid email or password",
		}
	}

	// Create session
	sess, err := h.sessionManager.Create(ctx, u.ID)
	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}

	// Set session cookie
	middleware.SetSessionCookie(ctx, sess.ID)

	return nil
}

type LoginRequest struct {
	Email    string `cbor:"email"`
	Password string `cbor:"password"`
}

func (r *LoginRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.Email, validation.Required, is.Email),
		validation.Field(&r.Password, validation.Required),
	)
}
