package middleware

import (
	"strconv"
	"time"

	"regulation/server/services/request_context"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"

	"github.com/DeltaLaboratory/contrib/u22"
)

func NewRequestInfoMiddleware() fiber.Handler {
	return func(ctx fiber.Ctx) error {
		requestID := uuid.Must(uuid.NewV7())
		requestTime := time.Now()

		request_context.SetTime(ctx, requestTime)
		request_context.SetRequestID(ctx, requestID)

		ctx.Set("Request-ID", u22.Encode(requestID))
		err := ctx.Next()
		ctx.Set("Server-Timing", "server;dur="+strconv.FormatInt(time.Since(requestTime).Milliseconds(), 10))
		return err
	}
}
