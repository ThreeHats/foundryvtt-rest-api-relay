// Package alerts dispatches admin-level alerts to configured destinations.
//
// PII safety: Alert messages MUST NOT include user emails or full IPs in the
// body — use opaque identifiers (user IDs, client IDs) instead.
package alerts

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/mail"
	"net/smtp"
	"strings"
	"sync"
	"time"

	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/database"
	"github.com/rs/zerolog/log"
)

// Alert types — Security
const (
	TypeFailedAuthSpike          = "failed_auth_spike"
	TypeFlaggedConnection        = "flagged_connection"
	TypeRateLimitBurst           = "rate_limit_burst"
	TypeAdminLogin               = "admin_login"
	TypeDuplicateConnectionSpike = "duplicate_connection_spike"
	TypeMetadataMismatchSpike    = "metadata_mismatch_spike"
	TypeExecuteJsBurst           = "execute_js_burst"
	TypePasswordResetFlood       = "password_reset_flood"
	TypeInvalidTokenSpike        = "invalid_token_spike"
	TypeNewIPForToken            = "new_ip_for_token"
	TypeAccountPasswordChange    = "account_password_change"

	TypeWorldIDCrossAccount      = "world_id_cross_account"
	TypeRegistrationSpike        = "registration_spike"

	// Alert types — Operations
	TypeClientDisconnectSpike = "client_disconnect_spike"
	TypeSystemUnhealthy       = "system_unhealthy"
	TypeHeadlessSessionFlood  = "headless_session_flood"
	TypeStripePaymentFailed   = "stripe_payment_failed"

	// Alert types — Analytics
	TypeNewUserRegistration      = "new_user_registration"
	TypeNewSubscription          = "new_subscription"
	TypeSubscriptionCancelled    = "subscription_cancelled"
	TypeUserMonthlyLimitApproach      = "user_monthly_limit_approaching"
	TypeScopedKeyMonthlyLimitApproach = "scoped_key_monthly_limit_approaching"
)

// Event represents a triggered alert.
type Event struct {
	Type      string                 `json:"type"`
	Severity  string                 `json:"severity"` // "info" | "warning" | "critical"
	Message   string                 `json:"message"`
	Details   map[string]interface{} `json:"details,omitempty"`
	Timestamp time.Time              `json:"timestamp"`
}

// --- package-level globals ---

var globalDB *database.DB

type smtpCfg struct {
	Host        string
	Port        int
	User        string
	Pass        string
	From        string
	FrontendURL string
}

var globalSMTP *smtpCfg

// Init stores the DB and SMTP config for use by Fire/Track/Test without
// threading them through every callsite.
func Init(db *database.DB, host string, port int, user, pass, from, frontendURL string) {
	globalDB = db
	globalSMTP = &smtpCfg{host, port, user, pass, from, frontendURL}
}

// --- recent events buffer ---

var recentEvents = struct {
	mu     sync.Mutex
	events []Event
	max    int
}{max: 200}

// RecentEvents returns a copy of the last N alert events.
func RecentEvents() []Event {
	recentEvents.mu.Lock()
	defer recentEvents.mu.Unlock()
	out := make([]Event, len(recentEvents.events))
	copy(out, recentEvents.events)
	return out
}

// --- spike tracker ---

type trackerEntry struct {
	mu        sync.Mutex
	count     int
	windowEnd time.Time
	lastFired time.Time
}

var trackerMap sync.Map // key string → *trackerEntry

// Track increments a sliding-window counter for key.
// Returns true when count >= threshold AND cooldown since last fire has passed.
func Track(key string, threshold int, window, cooldown time.Duration) bool {
	v, _ := trackerMap.LoadOrStore(key, &trackerEntry{})
	e := v.(*trackerEntry)

	e.mu.Lock()
	defer e.mu.Unlock()

	now := time.Now()
	if now.After(e.windowEnd) {
		e.count = 0
		e.windowEnd = now.Add(window)
	}
	e.count++

	if e.count >= threshold && now.After(e.lastFired.Add(cooldown)) {
		e.lastFired = now
		return true
	}
	return false
}

// --- Fire ---

// Fire dispatches an alert event through all configured subscriptions.
// Best-effort; failures are logged but do not block the caller.
func Fire(event Event) {
	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now()
	}

	// Track in recent events buffer
	recentEvents.mu.Lock()
	recentEvents.events = append(recentEvents.events, event)
	if len(recentEvents.events) > recentEvents.max {
		recentEvents.events = recentEvents.events[len(recentEvents.events)-recentEvents.max:]
	}
	recentEvents.mu.Unlock()

	if globalDB == nil {
		return
	}

	// Look up subscribers and dispatch in goroutines
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		cfg, err := globalDB.AlertConfigStore().Get(ctx)
		if err != nil {
			log.Warn().Err(err).Msg("alerts: failed to load config")
			return
		}
		subs, err := globalDB.AlertSubscriptionStore().FindByAlertType(ctx, event.Type)
		if err != nil {
			log.Warn().Err(err).Str("type", event.Type).Msg("alerts: failed to load subscriptions")
			return
		}
		for _, sub := range subs {
			dest := ""
			switch sub.Channel {
			case "discord":
				dest = cfg.DiscordWebhookURL
			case "email":
				dest = cfg.EmailDestination
			}
			if dest == "" {
				continue
			}
			ch := sub.Channel
			go dispatch(dest, ch, event)
		}
	}()
}

// --- Test ---

// TestChannel sends a synthetic test event to the given channel/destination.
func TestChannel(channel, destination string) error {
	if destination == "" {
		return fmt.Errorf("no destination configured for channel %q", channel)
	}
	testEvent := Event{
		Type:      "test_alert",
		Severity:  "info",
		Message:   fmt.Sprintf("Test alert via %s", channel),
		Timestamp: time.Now(),
	}
	return dispatch(destination, channel, testEvent)
}

// --- dispatch ---

func dispatch(destination, channel string, event Event) error {
	switch channel {
	case "discord":
		if err := sendDiscord(destination, event); err != nil {
			log.Warn().Err(err).Str("channel", channel).Msg("alerts: discord delivery failed")
			return err
		}
	case "email":
		if globalSMTP != nil && globalSMTP.Host != "" {
			subject := fmt.Sprintf("[%s] Relay Alert: %s", strings.ToUpper(event.Severity), event.Type)
			body := buildAlertEmailBody(event)
			if err := sendAlertEmail(destination, subject, body); err != nil {
				log.Warn().Err(err).Str("channel", channel).Msg("alerts: email delivery failed")
				return err
			}
		} else {
			log.Info().Str("channel", channel).Str("type", event.Type).Msg("alerts: email delivery skipped (SMTP not configured)")
		}
	default:
		log.Warn().Str("channel", channel).Msg("alerts: unknown channel")
	}
	return nil
}

// --- Discord ---

func sendDiscord(webhookURL string, event Event) error {
	embed := map[string]interface{}{
		"title":       fmt.Sprintf("[%s] %s", event.Severity, event.Type),
		"description": event.Message,
		"timestamp":   event.Timestamp.UTC().Format(time.RFC3339),
		"color":       severityColor(event.Severity),
	}
	if len(event.Details) > 0 {
		fields := make([]map[string]interface{}, 0, len(event.Details))
		for k, v := range event.Details {
			fields = append(fields, map[string]interface{}{
				"name":   k,
				"value":  fmt.Sprintf("%v", v),
				"inline": true,
			})
		}
		embed["fields"] = fields
	}
	payload := map[string]interface{}{
		"embeds": []interface{}{embed},
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", webhookURL, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return fmt.Errorf("discord returned status %d", resp.StatusCode)
	}
	return nil
}

// --- Email ---

func buildAlertEmailBody(event Event) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(`<!DOCTYPE html>
<html>
<body style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto;">
  <h2 style="color: %s;">[%s] %s</h2>
  <p>%s</p>
  <p style="color: #888; font-size: 0.85em;">%s</p>`,
		severityHexColor(event.Severity),
		strings.ToUpper(event.Severity),
		event.Type,
		event.Message,
		event.Timestamp.UTC().Format(time.RFC1123),
	))
	if len(event.Details) > 0 {
		sb.WriteString("<table style=\"border-collapse:collapse;width:100%;margin-top:1em;\">")
		for k, v := range event.Details {
			sb.WriteString(fmt.Sprintf(
				"<tr><td style=\"padding:4px 8px;border:1px solid #ddd;font-weight:bold;\">%s</td><td style=\"padding:4px 8px;border:1px solid #ddd;\">%v</td></tr>",
				k, v,
			))
		}
		sb.WriteString("</table>")
	}
	sb.WriteString("</body></html>")
	return sb.String()
}

func sendAlertEmail(to, subject, body string) error {
	if globalSMTP == nil {
		return fmt.Errorf("SMTP not configured")
	}
	msg := strings.Join([]string{
		fmt.Sprintf("From: %s", globalSMTP.From),
		fmt.Sprintf("To: %s", to),
		fmt.Sprintf("Subject: %s", subject),
		"MIME-Version: 1.0",
		"Content-Type: text/html; charset=UTF-8",
		"",
		body,
	}, "\r\n")

	addr := fmt.Sprintf("%s:%d", globalSMTP.Host, globalSMTP.Port)
	var auth smtp.Auth
	if globalSMTP.User != "" {
		auth = smtp.PlainAuth("", globalSMTP.User, globalSMTP.Pass, globalSMTP.Host)
	}
	fromAddr := extractEmailAddress(globalSMTP.From)
	return smtp.SendMail(addr, auth, fromAddr, []string{to}, []byte(msg))
}

func extractEmailAddress(from string) string {
	addr, err := mail.ParseAddress(from)
	if err != nil {
		return from
	}
	return addr.Address
}

func severityColor(severity string) int {
	switch severity {
	case "critical":
		return 0xE53935
	case "warning":
		return 0xFB8C00
	default:
		return 0x1E88E5
	}
}

func severityHexColor(severity string) string {
	switch severity {
	case "critical":
		return "#e53935"
	case "warning":
		return "#fb8c00"
	default:
		return "#1e88e5"
	}
}
