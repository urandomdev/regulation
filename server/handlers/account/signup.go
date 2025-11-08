package account

import (
	"fmt"

	"github.com/DeltaLaboratory/contrib/hooks"
	"github.com/DeltaLaboratory/password"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/gofiber/fiber/v3"
)

// Signup creates a new user account
// @Route POST /account/signup
func (h *Handler) Signup(ctx fiber.Ctx, req *SignupRequest) (ret error) {
	tx, err := h.db.Tx(ctx)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer hooks.Rollback(tx, &ret)

	hash, err := password.Generate([]byte(req.Password))
	if err != nil {
		return fmt.Errorf("failed to hash hash: %w", err)
	}

	if err := tx.User.Create().
		SetEmail(req.Email).
		SetPassword(hash).
		SetNickname(req.Nickname).
		Exec(ctx); err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

type SignupRequest struct {
	Email    string `cbor:"email"`
	Password string `cbor:"password"`
	Nickname string `cbor:"nickname"`
}

func (r *SignupRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.Email, validation.Required, is.Email),
		validation.Field(&r.Password, validation.Required, validation.RuneLength(8, 128)),
		validation.Field(&r.Nickname, validation.Required, validation.RuneLength(3, 30)),
	)
}
