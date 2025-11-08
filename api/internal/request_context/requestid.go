//revive:disable:var-naming // package name mirrors directory naming used across server
package request_context

import (
	"context"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

func RequestID(ctx context.Context) uuid.UUID {
	id, ok := ctx.Value(keyRequestID).(uuid.UUID)
	if !ok {
		return uuid.Nil
	}

	return id
}

func SetRequestID(ctx fiber.Ctx, id uuid.UUID) {
	ctx.RequestCtx().SetUserValue(keyRequestID, id)
}
