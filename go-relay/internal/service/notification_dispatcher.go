package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/model"
	"github.com/rs/zerolog/log"
)

// NotificationEvent identifies a notification event type. Each event corresponds
// to a per-event toggle on NotificationSettings (account-level) and/or
// ApiKeyNotificationSettings (per-key).
type NotificationEvent string

const (
	EventConnect                    NotificationEvent = "connect"
	EventDisconnect                 NotificationEvent = "disconnect"
	EventMetadataMismatch           NotificationEvent = "metadata-mismatch"
	EventSettingsChange             NotificationEvent = "settings-change"
	EventExecuteJs                  NotificationEvent = "execute-js"
	EventMacroExecute               NotificationEvent = "macro-execute"
	EventRateLimit                  NotificationEvent = "rate-limit"
	EventError                      NotificationEvent = "error"
	// EventRemoteRequest is dispatched (via the batcher, not directly) when
	// cross-world remote-request activity has been batched over the 5-minute
	// window. Always-on (same as EventDuplicateConnectionRejected).
	EventRemoteRequest NotificationEvent = "remote-request"
	// EventDuplicateConnectionRejected fires when a WS client tries to connect
	// with a clientId that is already held by another live connection. This is
	// the silent-alarm: an attacker holding a stolen connection token can't
	// take over while the legitimate client is online, but the rejection event
	// alerts the account owner that someone tried.
	EventDuplicateConnectionRejected NotificationEvent = "duplicate-connection-rejected"
	// EventNewClientConnect fires the first time a clientId is seen connecting.
	// Distinct from EventConnect (which fires on every connect/reconnect).
	EventNewClientConnect NotificationEvent = "new-client-connect"
)

// NotificationContext describes the metadata for a notification event.
//
// All fields are optional except Event. The dispatcher uses whatever fields
// are present to build the notification body.
type NotificationContext struct {
	Event       NotificationEvent
	UserID      int64  // Account that owns the affected client/key (always required)
	ApiKeyID    int64  // Optional: scoped key associated with this event (0 if N/A)
	ClientID    string // Foundry client ID, if applicable
	WorldID     string
	WorldTitle  string
	SystemID    string
	IPAddress   string
	Description string // Free-form description (e.g., the script for execute-js)
	Reason      string // For metadata mismatches: details of the mismatch
	Severity    string // "info" | "warning" | "alert"
}

// NotificationStores bundles the database stores the dispatcher needs.
type NotificationStores struct {
	NotificationSettings            NotificationSettingsLookup
	ApiKeyNotificationSettings      ApiKeyNotificationSettingsLookup
	KnownClientNotificationSettings KnownClientNotificationSettingsLookup
	KnownClients                    KnownClientLookup
}

// NotificationSettingsLookup is the subset of NotificationSettingsStore the
// dispatcher needs. Defined here as an interface to avoid an import cycle
// (the model package can't import this package).
type NotificationSettingsLookup interface {
	FindByUser(ctx context.Context, userID int64) (*model.NotificationSettings, error)
}

// ApiKeyNotificationSettingsLookup is the subset of ApiKeyNotificationSettingsStore needed.
type ApiKeyNotificationSettingsLookup interface {
	FindByApiKey(ctx context.Context, apiKeyID int64) (*model.ApiKeyNotificationSettings, error)
}

// KnownClientNotificationSettingsLookup is the subset of KnownClientNotificationSettingsStore needed.
type KnownClientNotificationSettingsLookup interface {
	FindByKnownClient(ctx context.Context, knownClientID int64) (*model.KnownClientNotificationSettings, error)
}

// KnownClientLookup is the subset of KnownClientStore needed by the dispatcher.
type KnownClientLookup interface {
	FindByClientID(ctx context.Context, userID int64, clientID string) (*model.KnownClient, error)
}

// Dispatcher routes notification events to the appropriate destinations.
//
// For each event:
//  1. Look up the user's account-level NotificationSettings — if the event's
//     toggle is enabled, send to the account-level Discord webhook + email.
//  2. If the event has an associated ApiKeyID, look up that key's
//     ApiKeyNotificationSettings — if the event's toggle is enabled, send to
//     the per-key Discord webhook + email.
//
// Both destinations may fire for the same event. The dispatcher is fire-and-
// forget — failures are logged but do not bubble up to callers.
//
// A per-user debounce map suppresses duplicate events of the same type for
// the same client within the user-configured window (NotificationDebounceWindowSecs).
type Dispatcher struct {
	stores      NotificationStores
	smtpCfg     *NotificationConfig
	frontendURL string
	// debounce tracks the last time a notification was sent per (userID, event, clientID).
	// Key format: "{userID}:{event}:{clientID}"
	debounce sync.Map
}

// NewDispatcher creates a notification dispatcher.
func NewDispatcher(stores NotificationStores, smtpCfg *NotificationConfig, frontendURL string) *Dispatcher {
	return &Dispatcher{
		stores:      stores,
		smtpCfg:     smtpCfg,
		frontendURL: frontendURL,
	}
}

// Dispatch routes the notification event to all matching destinations.
// This should typically be called from a goroutine — it makes blocking HTTP
// and SMTP calls.
func (d *Dispatcher) Dispatch(nc NotificationContext) {
	if nc.UserID == 0 {
		log.Warn().Str("event", string(nc.Event)).Msg("Dispatch called with no UserID")
		return
	}

	ctx := context.Background()

	// Check account-level debounce window before doing anything.
	// Key: "{userID}:{event}:{clientID}" — suppresses the same event type
	// for the same client within the configured window.
	if d.stores.NotificationSettings != nil {
		acct, err := d.stores.NotificationSettings.FindByUser(ctx, nc.UserID)
		if err == nil && acct != nil && acct.NotificationDebounceWindowSecs > 0 {
			debounceKey := fmt.Sprintf("%d:%s:%s", nc.UserID, nc.Event, nc.ClientID)
			window := time.Duration(acct.NotificationDebounceWindowSecs) * time.Second
			if lastRaw, loaded := d.debounce.Load(debounceKey); loaded {
				if last, ok := lastRaw.(time.Time); ok && time.Since(last) < window {
					log.Debug().
						Str("event", string(nc.Event)).
						Str("clientId", nc.ClientID).
						Int64("userId", nc.UserID).
						Msg("Notification suppressed by debounce window")
					return
				}
			}
			// Record send time before dispatching so concurrent calls are also suppressed.
			d.debounce.Store(debounceKey, time.Now())
		}

		// Account-level dispatch
		if acct != nil && d.accountEventEnabled(acct, nc.Event) {
			d.sendToDestination(
				acct.DiscordWebhookURL.String,
				acct.NotifyEmail.String,
				nc,
			)
		}
	}

	// Per-key dispatch (only if ApiKeyID is set)
	if nc.ApiKeyID != 0 && d.stores.ApiKeyNotificationSettings != nil {
		keyNS, err := d.stores.ApiKeyNotificationSettings.FindByApiKey(ctx, nc.ApiKeyID)
		if err == nil && keyNS != nil && d.keyEventEnabled(keyNS, nc.Event) {
			d.sendToDestination(
				keyNS.DiscordWebhookURL.String,
				keyNS.NotifyEmail.String,
				nc,
			)
		}
	}

	// World-level dispatch (connect, disconnect, execute-js, macro-execute for a specific world)
	if nc.ClientID != "" &&
		d.stores.KnownClients != nil &&
		d.stores.KnownClientNotificationSettings != nil {
		switch nc.Event {
		case EventConnect, EventDisconnect, EventExecuteJs, EventMacroExecute:
			kc, err := d.stores.KnownClients.FindByClientID(ctx, nc.UserID, nc.ClientID)
			if err == nil && kc != nil {
				worldNS, err := d.stores.KnownClientNotificationSettings.FindByKnownClient(ctx, kc.ID)
				if err == nil && worldNS != nil && d.worldEventEnabled(worldNS, nc.Event) {
					d.sendToDestination(
						worldNS.DiscordWebhookURL.String,
						worldNS.NotifyEmail.String,
						nc,
					)
				}
			}
		}
	}
}

// accountEventEnabled checks if an account-level setting allows the given event.
func (d *Dispatcher) accountEventEnabled(s *model.NotificationSettings, event NotificationEvent) bool {
	switch event {
	case EventNewClientConnect:
		return s.NotifyOnNewClientConnect
	case EventConnect:
		return s.NotifyOnConnect
	case EventDisconnect:
		return s.NotifyOnDisconnect
	case EventMetadataMismatch:
		return s.NotifyOnMetadataMismatch
	case EventDuplicateConnectionRejected, EventRemoteRequest:
		// Always-on for security-critical events — no per-event toggle.
		return true
	case EventSettingsChange:
		return s.NotifyOnSettingsChange
	case EventExecuteJs:
		return s.NotifyOnExecuteJs
	case EventMacroExecute:
		return s.NotifyOnMacroExecute
	}
	return false
}

// keyEventEnabled checks if a per-key setting allows the given event.
func (d *Dispatcher) keyEventEnabled(s *model.ApiKeyNotificationSettings, event NotificationEvent) bool {
	switch event {
	case EventExecuteJs:
		return s.NotifyOnExecuteJs
	case EventMacroExecute:
		return s.NotifyOnMacroExecute
	case EventRateLimit:
		return s.NotifyOnRateLimit
	case EventError:
		return s.NotifyOnError
	}
	// Per-key settings don't track connect/disconnect/settings-change/metadata-mismatch
	// — those are inherently account-level concerns.
	return false
}

// worldEventEnabled checks if a per-world setting allows the given event.
func (d *Dispatcher) worldEventEnabled(s *model.KnownClientNotificationSettings, event NotificationEvent) bool {
	switch event {
	case EventConnect:
		return s.NotifyOnConnect
	case EventDisconnect:
		return s.NotifyOnDisconnect
	case EventExecuteJs:
		return s.NotifyOnExecuteJs
	case EventMacroExecute:
		return s.NotifyOnMacroExecute
	}
	return false
}

// sendToDestination sends a notification via Discord webhook and/or email.
func (d *Dispatcher) sendToDestination(discordWebhookURL, notifyEmail string, nc NotificationContext) {
	title, color := titleAndColor(nc.Event)
	description := buildDescription(nc)

	if discordWebhookURL != "" {
		if err := sendDiscordEmbed(discordWebhookURL, title, description, color); err != nil {
			log.Warn().Err(err).Str("event", string(nc.Event)).Msg("Discord notification failed")
		}
	}

	if notifyEmail != "" && d.smtpCfg != nil && d.smtpCfg.SMTPHost != "" {
		subject := fmt.Sprintf("Foundry REST API: %s", title)
		body := fmt.Sprintf("%s\n\n%s\n\n--\nFoundry REST API Relay\n%s",
			title, description, d.frontendURL)
		if err := sendEmail(notifyEmail, subject, body, d.smtpCfg); err != nil {
			log.Warn().Err(err).Str("event", string(nc.Event)).Msg("Email notification failed")
		}
	}
}

// titleAndColor returns the title and Discord embed color for an event type.
func titleAndColor(event NotificationEvent) (string, int) {
	switch event {
	case EventNewClientConnect:
		return "New Foundry World Connected", 3066993 // green
	case EventConnect:
		return "Foundry Client Connected", 3066993 // green
	case EventDisconnect:
		return "Foundry Client Disconnected", 15105570 // orange
	case EventMetadataMismatch:
		return "⚠ Suspicious Connection: Metadata Mismatch", 15158332 // red
	case EventDuplicateConnectionRejected:
		return "⚠ Rejected Duplicate Connection", 15158332 // red
	case EventRemoteRequest:
		return "Cross-World Activity Summary", 3447003 // blue
	case EventSettingsChange:
		return "⚠ Foundry Module Setting Changed", 15105570 // orange
	case EventExecuteJs:
		return "execute-js Called", 3447003 // blue
	case EventMacroExecute:
		return "macro-execute Called", 3447003 // blue
	case EventRateLimit:
		return "API Key Rate-Limited", 15105570 // orange
	case EventError:
		return "API Key Error", 15158332 // red
	}
	return "Foundry REST API Event", 3447003
}

// buildDescription builds a human-readable description for a notification.
func buildDescription(nc NotificationContext) string {
	desc := ""
	if nc.WorldTitle != "" {
		desc += fmt.Sprintf("**World:** %s\n", nc.WorldTitle)
	}
	if nc.ClientID != "" {
		desc += fmt.Sprintf("**Client:** `%s`\n", nc.ClientID)
	}
	if nc.SystemID != "" {
		desc += fmt.Sprintf("**System:** %s\n", nc.SystemID)
	}
	if nc.IPAddress != "" {
		desc += fmt.Sprintf("**IP:** %s\n", nc.IPAddress)
	}
	if nc.Reason != "" {
		desc += fmt.Sprintf("**Reason:** %s\n", nc.Reason)
	}
	if nc.Description != "" {
		desc += fmt.Sprintf("\n%s", nc.Description)
	}
	if desc == "" {
		desc = "(no additional details)"
	}
	return desc
}
