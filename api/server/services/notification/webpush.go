package notification

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/SherClockHolmes/webpush-go"
	"github.com/fxamacker/cbor/v2"
	"github.com/google/uuid"

	"regulation/internal/config"
	"regulation/internal/ent"
	"regulation/internal/ent/pushsubscription"
)

// Service handles notification delivery
type Service struct {
	config *config.Config
	db     *ent.Client
}

// New creates a new notification service
func New(config *config.Config, db *ent.Client) *Service {
	return &Service{
		config: config,
		db:     db,
	}
}

// WebPushSubscription represents a web push subscription from a client
type WebPushSubscription struct {
	Endpoint string `json:"endpoint"`
	Auth     string `json:"auth"`
	P256DH   string `json:"p256dh"`
}

// Payload represents the notification data to send
type Payload struct {
	Title string `cbor:"title" json:"title"`
	Body  string `cbor:"body" json:"body"`
	URL   string `cbor:"url,omitempty" json:"url,omitempty"`
}

// SendWebPush sends a web push notification to a subscriber
func (s *Service) SendWebPush(ctx context.Context, subscription *WebPushSubscription, payload *Payload) error {
	if s.config.WebPush == nil || !s.config.WebPush.HasWebPush() {
		return errors.New("web push is not configured")
	}

	if subscription == nil {
		return errors.New("subscription is required")
	}

	if payload == nil {
		return errors.New("payload is required")
	}

	// Marshal payload to CBOR
	message, err := cbor.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal notification payload: %w", err)
	}

	// Create the subscription from credentials
	sub := &webpush.Subscription{
		Endpoint: subscription.Endpoint,
		Keys: webpush.Keys{
			Auth:   subscription.Auth,
			P256dh: subscription.P256DH,
		},
	}

	// Configure web push options
	options := &webpush.Options{
		TTL:             s.config.WebPush.GetTTL(),
		Urgency:         webpush.UrgencyNormal,
		VAPIDPublicKey:  s.config.WebPush.VAPIDPublicKey,
		VAPIDPrivateKey: s.config.WebPush.VAPIDPrivateKey,
		Subscriber:      s.config.WebPush.Subscriber,
	}

	// Send the notification
	resp, err := webpush.SendNotificationWithContext(ctx, message, sub, options)
	if err != nil {
		return fmt.Errorf("failed to send web push notification: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode == http.StatusGone {
		return errors.New("subscription expired")
	}

	if resp.StatusCode >= 400 {
		return fmt.Errorf("web push failed with status %d", resp.StatusCode)
	}

	return nil
}

// SendToUser sends a web push notification to all active subscriptions for a user
func (s *Service) SendToUser(ctx context.Context, userID uuid.UUID, payload *Payload) (int, []error) {
	if s.config.WebPush == nil || !s.config.WebPush.HasWebPush() {
		return 0, []error{errors.New("web push is not configured")}
	}

	if payload == nil {
		return 0, []error{errors.New("payload is required")}
	}

	// Get all active subscriptions for the user
	subscriptions, err := s.db.PushSubscription.Query().
		Where(
			pushsubscription.UserIDEQ(userID),
			pushsubscription.ActiveEQ(true),
		).
		All(ctx)

	if err != nil {
		return 0, []error{fmt.Errorf("failed to query subscriptions: %w", err)}
	}

	if len(subscriptions) == 0 {
		return 0, nil
	}

	// Marshal payload once
	message, err := cbor.Marshal(payload)
	if err != nil {
		return 0, []error{fmt.Errorf("failed to marshal notification payload: %w", err)}
	}

	// Send to all subscriptions
	var errs []error
	successCount := 0

	for _, sub := range subscriptions {
		err := s.sendToSubscription(ctx, sub, message)
		if err != nil {
			errs = append(errs, err)
		} else {
			successCount++
		}
	}

	return successCount, errs
}

// sendToSubscription sends a notification to a single subscription
func (s *Service) sendToSubscription(ctx context.Context, sub *ent.PushSubscription, message []byte) error {
	// Create the subscription
	subscription := &webpush.Subscription{
		Endpoint: sub.Endpoint,
		Keys: webpush.Keys{
			Auth:   sub.Auth,
			P256dh: sub.P256dh,
		},
	}

	// Configure web push options
	options := &webpush.Options{
		TTL:             s.config.WebPush.GetTTL(),
		Urgency:         webpush.UrgencyNormal,
		VAPIDPublicKey:  s.config.WebPush.VAPIDPublicKey,
		VAPIDPrivateKey: s.config.WebPush.VAPIDPrivateKey,
		Subscriber:      s.config.WebPush.Subscriber,
	}

	// Send the notification
	resp, err := webpush.SendNotificationWithContext(ctx, message, subscription, options)
	if err != nil {
		return fmt.Errorf("failed to send web push notification: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode == http.StatusGone {
		// Mark subscription as inactive
		_ = sub.Update().
			SetActive(false).
			Exec(ctx)
		return errors.New("subscription expired")
	}

	if resp.StatusCode >= 400 {
		return fmt.Errorf("web push failed with status %d", resp.StatusCode)
	}

	return nil
}
