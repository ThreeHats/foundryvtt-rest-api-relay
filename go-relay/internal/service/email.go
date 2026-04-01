package service

import (
	"fmt"
	"net/mail"
	"net/smtp"
	"strings"

	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/config"
	"github.com/rs/zerolog/log"
)

// extractEmailAddress parses a "Name <addr>" string and returns just the addr.
// If parsing fails, returns the input as-is.
func extractEmailAddress(from string) string {
	addr, err := mail.ParseAddress(from)
	if err != nil {
		return from
	}
	return addr.Address
}

// SendPasswordResetEmail sends a password reset email, or logs the URL if SMTP is not configured.
func SendPasswordResetEmail(cfg *config.Config, to, resetToken string) {
	if cfg.SMTPHost == "" {
		// Fallback: log the reset URL
		resetURL := fmt.Sprintf("%s?reset-token=%s", cfg.FrontendURL, resetToken)
		log.Info().Str("email", to).Str("resetURL", resetURL).Msg("Password reset (SMTP not configured)")
		return
	}

	resetURL := fmt.Sprintf("%s?reset-token=%s", cfg.FrontendURL, resetToken)

	subject := "Password Reset Request"
	body := fmt.Sprintf(`<!DOCTYPE html>
<html>
<body style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto;">
  <h2>Password Reset</h2>
  <p>You requested a password reset. Click the link below to reset your password:</p>
  <p><a href="%s" style="background: #4CAF50; color: white; padding: 10px 20px; text-decoration: none; border-radius: 4px;">Reset Password</a></p>
  <p>This link expires in 1 hour.</p>
  <p>If you didn't request this, you can safely ignore this email.</p>
</body>
</html>`, resetURL)

	msg := strings.Join([]string{
		fmt.Sprintf("From: %s", cfg.SMTPFrom),
		fmt.Sprintf("To: %s", to),
		fmt.Sprintf("Subject: %s", subject),
		"MIME-Version: 1.0",
		"Content-Type: text/html; charset=UTF-8",
		"",
		body,
	}, "\r\n")

	addr := fmt.Sprintf("%s:%d", cfg.SMTPHost, cfg.SMTPPort)
	var auth smtp.Auth
	if cfg.SMTPUser != "" {
		auth = smtp.PlainAuth("", cfg.SMTPUser, cfg.SMTPPass, cfg.SMTPHost)
	}

	go func() {
		if err := smtp.SendMail(addr, auth, extractEmailAddress(cfg.SMTPFrom), []string{to}, []byte(msg)); err != nil {
			log.Error().Err(err).Str("to", to).Msg("Failed to send password reset email")
		} else {
			log.Info().Str("to", to).Msg("Password reset email sent")
		}
	}()
}
