package notification

import (
	"github.com/gofiber/fiber/v3"

	"regulation/internal/ent"
	"regulation/internal/ent/pushsubscription"
	"regulation/internal/protocol"
	"regulation/server/services/request_context"
)

// SubscribeRequest represents a web push subscription request
type SubscribeRequest struct {
	Endpoint string `cbor:"endpoint" json:"endpoint"`
	P256DH   string `cbor:"p256dh" json:"p256dh"`
	Auth     string `cbor:"auth" json:"auth"`
}

// Subscribe registers a new web push subscription for the authenticated user
// @Route POST /notification/subscribe
func (h *Handler) Subscribe(ctx fiber.Ctx, req *SubscribeRequest) error {
	user := request_context.User(ctx)

	// Validate request
	if req.Endpoint == "" || req.P256DH == "" || req.Auth == "" {
		return protocol.ErrorResponse{
			Code:    protocol.InvalidRequest,
			Message: "endpoint, p256dh, and auth are required",
		}
	}

	// Check if subscription already exists
	existing, err := h.db.PushSubscription.Query().
		Where(
			pushsubscription.EndpointEQ(req.Endpoint),
		).
		Only(ctx)

	if err == nil {
		// Subscription exists - update to active if needed
		if !existing.Active || existing.UserID != user.ID {
			_, err = existing.Update().
				SetActive(true).
				Save(ctx)
			if err != nil {
				return protocol.ErrorResponse{
					Code:    protocol.InternalError,
					Message: "failed to update subscription",
				}
			}
		}
		return nil
	}

	// Create new subscription
	_, err = h.db.PushSubscription.Create().
		SetUserID(user.ID).
		SetEndpoint(req.Endpoint).
		SetP256dh(req.P256DH).
		SetAuth(req.Auth).
		Save(ctx)

	if err != nil {
		if ent.IsConstraintError(err) {
			return protocol.ErrorResponse{
				Code:    protocol.InvalidRequest,
				Message: "subscription already exists",
			}
		}
		return protocol.ErrorResponse{
			Code:    protocol.InternalError,
			Message: "failed to create subscription",
		}
	}

	return nil
}

// UnsubscribeRequest represents a web push unsubscribe request
type UnsubscribeRequest struct {
	Endpoint string `cbor:"endpoint" json:"endpoint"`
}

// Unsubscribe removes a web push subscription for the authenticated user
// @Route DELETE /notification/subscribe
func (h *Handler) Unsubscribe(ctx fiber.Ctx, req *UnsubscribeRequest) error {
	user := request_context.User(ctx)

	if req.Endpoint == "" {
		return protocol.ErrorResponse{
			Code:    protocol.InvalidRequest,
			Message: "endpoint is required",
		}
	}

	// Find and delete the subscription
	deleted, err := h.db.PushSubscription.Delete().
		Where(
			pushsubscription.UserIDEQ(user.ID),
			pushsubscription.EndpointEQ(req.Endpoint),
		).
		Exec(ctx)

	if err != nil {
		return protocol.ErrorResponse{
			Code:    protocol.InternalError,
			Message: "failed to delete subscription",
		}
	}

	if deleted == 0 {
		return protocol.ErrorResponse{
			Code:    protocol.InvalidRequest,
			Message: "subscription not found",
		}
	}

	return nil
}
