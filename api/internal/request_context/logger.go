//revive:disable:var-naming // package name mirrors directory naming used across server
package request_context

import (
	"context"

	"github.com/gofiber/fiber/v3"
	"github.com/rs/zerolog"
)

func Logger(ctx context.Context) zerolog.Logger {
	logger, ok := ctx.Value(keyLogger).(zerolog.Logger)
	if !ok {
		return zerolog.Nop()
	}

	return logger
}

func SetLogger(ctx fiber.Ctx, logger zerolog.Logger) {
	ctx.RequestCtx().SetUserValue(keyLogger, logger)
}
