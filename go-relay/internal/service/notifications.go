package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/smtp"
	"time"

	"github.com/rs/zerolog/log"
)

// ClientConnectionInfo holds metadata about a client connection for notifications.
type ClientConnectionInfo struct {
	ClientID       string
	WorldID        string
	WorldTitle     string
	SystemID       string
	FoundryVersion string
	IPAddress      string
	CustomName     string
}

// NotificationConfig holds SMTP and notification configuration.
type NotificationConfig struct {
	SMTPHost    string
	SMTPPort    int
	SMTPUser    string
	SMTPPass    string
	SMTPFrom    string
	FrontendURL string
}

// SendConnectionNotification sends a notification about a new client connection.
// Sends to Discord webhook and/or email depending on settings.
// This should be called asynchronously (in a goroutine).
func SendConnectionNotification(discordWebhookURL, notifyEmail string, smtpCfg *NotificationConfig, info ClientConnectionInfo) {
	if discordWebhookURL != "" {
		if err := sendDiscordNotification(discordWebhookURL, info, "connect"); err != nil {
			log.Warn().Err(err).Str("clientId", info.ClientID).Msg("Failed to send Discord connection notification")
		}
	}

	if notifyEmail != "" && smtpCfg != nil && smtpCfg.SMTPHost != "" {
		if err := sendEmailNotification(notifyEmail, smtpCfg, info, "connect"); err != nil {
			log.Warn().Err(err).Str("clientId", info.ClientID).Msg("Failed to send email connection notification")
		}
	}
}

// SendDisconnectNotification sends a notification about a client disconnection.
// Mirror of SendConnectionNotification but with a different event type and styling.
// Should be called asynchronously.
func SendDisconnectNotification(discordWebhookURL, notifyEmail string, smtpCfg *NotificationConfig, info ClientConnectionInfo) {
	if discordWebhookURL != "" {
		if err := sendDiscordNotification(discordWebhookURL, info, "disconnect"); err != nil {
			log.Warn().Err(err).Str("clientId", info.ClientID).Msg("Failed to send Discord disconnect notification")
		}
	}

	if notifyEmail != "" && smtpCfg != nil && smtpCfg.SMTPHost != "" {
		if err := sendEmailNotification(notifyEmail, smtpCfg, info, "disconnect"); err != nil {
			log.Warn().Err(err).Str("clientId", info.ClientID).Msg("Failed to send email disconnect notification")
		}
	}
}

// SendTestNotification sends a test notification to verify settings.
func SendTestNotification(discordWebhookURL, notifyEmail string, smtpCfg *NotificationConfig) error {
	testInfo := ClientConnectionInfo{
		ClientID:   "test-connection",
		WorldTitle: "Test Notification",
		IPAddress:  "127.0.0.1",
	}

	if discordWebhookURL != "" {
		if err := sendDiscordNotification(discordWebhookURL, testInfo, "test"); err != nil {
			return fmt.Errorf("discord notification failed: %w", err)
		}
	}

	if notifyEmail != "" && smtpCfg != nil && smtpCfg.SMTPHost != "" {
		if err := sendEmailNotification(notifyEmail, smtpCfg, testInfo, "test"); err != nil {
			return fmt.Errorf("email notification failed: %w", err)
		}
	}

	return nil
}

func sendDiscordNotification(webhookURL string, info ClientConnectionInfo, eventType string) error {
	title := "New Foundry Client Connected"
	color := 3066993 // Green
	switch eventType {
	case "disconnect":
		title = "Foundry Client Disconnected"
		color = 15105570 // Orange
	case "test":
		title = "Test Notification"
		color = 3447003 // Blue
	}

	description := fmt.Sprintf("**Client ID:** `%s`", info.ClientID)

	if info.WorldTitle != "" {
		description += fmt.Sprintf("\n**World:** %s", info.WorldTitle)
	}
	if info.SystemID != "" {
		description += fmt.Sprintf("\n**System:** %s", info.SystemID)
	}
	if info.FoundryVersion != "" {
		description += fmt.Sprintf("\n**Foundry Version:** %s", info.FoundryVersion)
	}
	if info.IPAddress != "" {
		description += fmt.Sprintf("\n**IP Address:** %s", info.IPAddress)
	}
	if info.CustomName != "" {
		description += fmt.Sprintf("\n**Custom Name:** %s", info.CustomName)
	}

	payload := map[string]interface{}{
		"embeds": []map[string]interface{}{
			{
				"title":       title,
				"description": description,
				"color":       color,
				"timestamp":   time.Now().UTC().Format(time.RFC3339),
				"footer": map[string]string{
					"text": "Foundry REST API Relay",
				},
			},
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal discord payload: %w", err)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Post(webhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("discord webhook request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("discord webhook returned status %d", resp.StatusCode)
	}

	return nil
}

func sendEmailNotification(to string, cfg *NotificationConfig, info ClientConnectionInfo, eventType string) error {
	subject := "Foundry REST API: New Client Connected"
	verb := "connected to"
	switch eventType {
	case "disconnect":
		subject = "Foundry REST API: Client Disconnected"
		verb = "disconnected from"
	case "test":
		subject = "Foundry REST API: Test Notification"
		verb = "test notification for"
	}

	worldInfo := info.ClientID
	if info.WorldTitle != "" {
		worldInfo = info.WorldTitle + " (" + info.ClientID + ")"
	}

	body := fmt.Sprintf(`A Foundry client has %s your relay.

Client: %s
IP Address: %s
Foundry Version: %s
System: %s
Time: %s

If this was not expected, log in to your dashboard to rotate the connection token.
%s
`, verb, worldInfo, info.IPAddress, info.FoundryVersion, info.SystemID,
		time.Now().UTC().Format(time.RFC3339), cfg.FrontendURL)

	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\n%s",
		cfg.SMTPFrom, to, subject, body)

	addr := fmt.Sprintf("%s:%d", cfg.SMTPHost, cfg.SMTPPort)
	auth := smtp.PlainAuth("", cfg.SMTPUser, cfg.SMTPPass, cfg.SMTPHost)

	if err := smtp.SendMail(addr, auth, cfg.SMTPFrom, []string{to}, []byte(msg)); err != nil {
		return fmt.Errorf("send email: %w", err)
	}

	return nil
}

// --- Generic helpers used by the unified Dispatcher ---

// sendDiscordEmbed sends a single Discord embed to the given webhook URL.
// Generic version used by the notification dispatcher.
func sendDiscordEmbed(webhookURL, title, description string, color int) error {
	payload := map[string]interface{}{
		"embeds": []map[string]interface{}{
			{
				"title":       title,
				"description": description,
				"color":       color,
				"timestamp":   time.Now().UTC().Format(time.RFC3339),
				"footer": map[string]string{
					"text": "Foundry REST API Relay",
				},
			},
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal discord payload: %w", err)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Post(webhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("discord webhook request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("discord webhook returned status %d", resp.StatusCode)
	}

	return nil
}

// sendEmail sends a plain-text email via the configured SMTP server.
// Generic version used by the notification dispatcher.
func sendEmail(to, subject, body string, cfg *NotificationConfig) error {
	if cfg == nil || cfg.SMTPHost == "" {
		return fmt.Errorf("SMTP not configured")
	}

	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\n%s",
		cfg.SMTPFrom, to, subject, body)

	addr := fmt.Sprintf("%s:%d", cfg.SMTPHost, cfg.SMTPPort)
	auth := smtp.PlainAuth("", cfg.SMTPUser, cfg.SMTPPass, cfg.SMTPHost)

	if err := smtp.SendMail(addr, auth, cfg.SMTPFrom, []string{to}, []byte(msg)); err != nil {
		return fmt.Errorf("send email: %w", err)
	}

	return nil
}
