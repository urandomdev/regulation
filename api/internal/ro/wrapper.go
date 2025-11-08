package ro

import (
	"errors"
	"fmt"

	"regulation/internal/config"
	"regulation/internal/protocol"
	"regulation/server/services/request_context"

	"github.com/DeltaLaboratory/contrib/u22"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/utils/v2"
)

type Validatable interface {
	Validate() error
}

type ValidatableCtx interface {
	Validate(ctx fiber.Ctx) error
}

func Error(ctx fiber.Ctx, err error) error {
	logger := request_context.Logger(ctx)

	ctx.Status(fiber.StatusTeapot)

	var errorResponse protocol.ErrorResponse
	if errors.As(err, &errorResponse) {
		logger.Debug().Msg(errorResponse.Error())

		return ctx.CBOR(errorResponse)
	}

	logger.Error().Err(err).Msg("internal server error")

	//goland:noinspection GoBoolExpressions
	if config.Version == "local" {
		return ctx.CBOR(protocol.ErrorResponse{
			Code:    protocol.InternalError,
			Message: err.Error(),
			Meta: map[string]any{
				"request_id": u22.Encode(request_context.RequestID(ctx)),
			},
		})
	}

	return ctx.CBOR(protocol.ErrorResponse{
		Code:    protocol.InternalError,
		Message: "internal server error",
		Meta: map[string]any{
			"request_id": u22.Encode(request_context.RequestID(ctx)),
		},
	})
}

// ParseBody parse user request and bind to T
// Codes: CodeInvalidRequest
func ParseBody[T any](ctx fiber.Ctx, request *T) error {
	contentType := utils.ToLower(ctx.Get("Content-Type"))
	contentType = utils.ParseVendorSpecificContentType(contentType)

	switch contentType {
	case "application/cbor":
		if err := ctx.Bind().CBOR(request); err != nil {
			return protocol.ErrorResponse{
				Code:    protocol.InvalidRequest,
				Message: fmt.Sprintf("failed to parse request: %v", err),
			}
		}
	default:
		return protocol.ErrorResponse{
			Code:    protocol.InvalidRequest,
			Message: fmt.Sprintf("unsupported content type: %q", contentType),
		}
	}

	if validatable, ok := any(request).(Validatable); ok {
		if err := validatable.Validate(); err != nil {
			if errors.As(err, &protocol.ErrorResponse{}) {
				return err
			}
			return protocol.ErrorResponse{
				Code:    protocol.InvalidRequest,
				Message: fmt.Sprintf("validation failed: %v", err),
			}
		}
	}

	if validatable, ok := any(request).(ValidatableCtx); ok {
		if err := validatable.Validate(ctx); err != nil {
			if errors.As(err, &protocol.ErrorResponse{}) {
				return err
			}
			return protocol.ErrorResponse{
				Code:    protocol.InvalidRequest,
				Message: fmt.Sprintf("validation failed: %v", err),
			}
		}
	}

	return nil
}

// Handler types with different signatures
type Handler[R, T any] func(fiber.Ctx, *R) (*T, error)
type Handler2[R any] func(fiber.Ctx, *R) error
type Handler3[T any] func(fiber.Ctx) (*T, error)
type Handler4 func(fiber.Ctx) error

// WrapHandlerBypass wraps a Handler function (with no request and response handling except for error)
func WrapHandlerBypass(f Handler4) func(fiber.Ctx) error {
	return func(ctx fiber.Ctx) error {
		err := f(ctx)
		if err != nil {
			return Error(ctx, err)
		}

		return nil
	}
}

// WrapHandler wraps a Handler function
func WrapHandler[R, T any](f Handler[R, T]) func(fiber.Ctx) error {
	return func(ctx fiber.Ctx) error {
		var req R

		err := ParseBody(ctx, &req)
		if err != nil {
			return Error(ctx, err)
		}

		res, err := f(ctx, &req)
		if err != nil {
			return Error(ctx, err)
		}

		if res == nil {
			return ctx.Status(fiber.StatusNoContent).Send(nil)
		}

		return ctx.CBOR(res)
	}
}

// WrapHandler2 wraps a Handler2 function (with no response)
func WrapHandler2[R any](f Handler2[R]) func(fiber.Ctx) error {
	return func(ctx fiber.Ctx) error {
		var req R

		err := ParseBody(ctx, &req)
		if err != nil {
			return Error(ctx, err)
		}

		err = f(ctx, &req)
		if err != nil {
			return Error(ctx, err)
		}

		return ctx.Status(fiber.StatusNoContent).Send(nil)
	}
}

// WrapHandler3 wraps a Handler3 function (no request)
func WrapHandler3[T any](f Handler3[T]) func(fiber.Ctx) error {
	return func(ctx fiber.Ctx) error {
		res, err := f(ctx)
		if err != nil {
			return Error(ctx, err)
		}

		if res == nil {
			return ctx.Status(fiber.StatusNoContent).Send(nil)
		}

		return ctx.CBOR(res)
	}
}

// WrapHandler4 wraps a Handler4 function (no request, no response)
func WrapHandler4(f Handler4) func(fiber.Ctx) error {
	return func(ctx fiber.Ctx) error {
		err := f(ctx)
		if err != nil {
			return Error(ctx, err)
		}
		return ctx.Status(fiber.StatusNoContent).Send(nil)
	}
}
