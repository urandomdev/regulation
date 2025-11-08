package config

import (
	"errors"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

// WebPush holds Web Push notification configuration
type WebPush struct {
	// VAPID (Voluntary Application Server Identification) keys
	VAPIDPublicKey  string `json:"vapid_public_key"`
	VAPIDPrivateKey string `json:"vapid_private_key"`

	// Subscriber email (format: mailto:admin@example.com)
	Subscriber string `json:"subscriber"`

	// TTL (Time To Live) in seconds - how long to keep the notification
	// Default: 86400 (24 hours)
	TTL int `json:"ttl"`
}

func (wp *WebPush) Validate() error {
	if wp == nil {
		return errors.New("webpush config is nil")
	}

	return validation.ValidateStruct(wp,
		validation.Field(&wp.VAPIDPublicKey, validation.Required, validation.Length(87, 87)),
		validation.Field(&wp.VAPIDPrivateKey, validation.Required, validation.Length(43, 43)),
		validation.Field(&wp.Subscriber, validation.Required, is.Email),
		validation.Field(&wp.TTL, validation.Min(0), validation.Max(2419200)), // max 28 days
	)
}

// GetTTL returns the TTL with a sensible default
func (wp *WebPush) GetTTL() int {
	if wp == nil || wp.TTL == 0 {
		return 86400 // 24 hours default
	}
	return wp.TTL
}

// HasWebPush returns true if Web Push is configured
func (wp *WebPush) HasWebPush() bool {
	return wp != nil &&
		wp.VAPIDPublicKey != "" &&
		wp.VAPIDPrivateKey != "" &&
		wp.Subscriber != ""
}
