package middleware

import (
	"regulation/internal/request_context"

	"github.com/gofiber/fiber/v3"
	"github.com/rs/zerolog"

	"github.com/DeltaLaboratory/contrib/u22"
)

func NewLoggerMiddleware(logger zerolog.Logger) fiber.Handler {
	return func(ctx fiber.Ctx) error {
		request_context.SetLogger(ctx, logger.With().
			Time("time", request_context.Time(ctx)).
			Str("path", ctx.Path()).
			Str("request_id", u22.Encode(request_context.RequestID(ctx))).Logger())
		return ctx.Next()
	}
}
