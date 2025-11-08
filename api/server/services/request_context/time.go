//revive:disable:var-naming // package name mirrors existing directory structure
package request_context

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v3"
)

func Time(ctx context.Context) time.Time {
	value, ok := ctx.Value(keyTime).(time.Time)
	if !ok {
		return time.Time{}
	}

	return value
}

func SetTime(ctx fiber.Ctx, value time.Time) {
	ctx.RequestCtx().SetUserValue(keyTime, value)
}
