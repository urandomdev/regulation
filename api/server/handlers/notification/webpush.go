package notification

import (
	"github.com/gofiber/fiber/v3"

	"regulation/internal/protocol"
)

// GetVAPIDPublicKey returns the VAPID public key for web push notifications
// @Route GET /notification/vapid
func (h *Handler) GetVAPIDPublicKey(_ fiber.Ctx) (*VAPIDPublicKeyResponse, error) {
	if h.config.WebPush == nil || !h.config.WebPush.HasWebPush() {
		return nil, protocol.ErrorResponse{
			Code:    protocol.InvalidRequest,
			Message: "Web push notifications are not configured",
		}
	}

	return &VAPIDPublicKeyResponse{
		PublicKey: h.config.WebPush.VAPIDPublicKey,
	}, nil
}

// VAPIDPublicKeyResponse holds the VAPID public key for web push notifications
type VAPIDPublicKeyResponse struct {
	PublicKey string `cbor:"public_key" json:"public_key"`
}
