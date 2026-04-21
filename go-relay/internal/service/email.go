package service

import (
	"fmt"
	"net/mail"
	"net/smtp"
	"strings"
	"time"

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

// SendVerificationEmail sends an email verification email, or logs the URL if SMTP is not configured.
func SendVerificationEmail(cfg *config.Config, to, verifyToken string) {
	verifyURL := fmt.Sprintf("%s/auth/verify?token=%s", cfg.FrontendURL, verifyToken)

	if cfg.SMTPHost == "" {
		log.Info().Str("email", to).Str("verifyURL", verifyURL).Msg("Email verification (SMTP not configured)")
		return
	}

	subject := "Verify Your Email Address"
	body := fmt.Sprintf(`<!DOCTYPE html>
<html>
<body style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto;">
  <h2>Verify Your Email</h2>
  <p>Thanks for registering! Please verify your email address by clicking the link below:</p>
  <p><a href="%s" style="background: #4CAF50; color: white; padding: 10px 20px; text-decoration: none; border-radius: 4px;">Verify Email</a></p>
  <p>This link expires in 24 hours.</p>
  <p>If you didn't create an account, you can safely ignore this email.</p>
</body>
</html>`, verifyURL)

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
		var lastErr error
		for attempt := 0; attempt < 3; attempt++ {
			if attempt > 0 {
				time.Sleep(time.Duration(attempt*2) * time.Second)
			}
			if err := smtp.SendMail(addr, auth, extractEmailAddress(cfg.SMTPFrom), []string{to}, []byte(msg)); err != nil {
				lastErr = err
				log.Warn().Err(err).Int("attempt", attempt+1).Str("to", to).Msg("Failed to send verification email, retrying...")
				continue
			}
			log.Info().Str("to", to).Msg("Verification email sent")
			return
		}
		log.Error().Err(lastErr).Str("to", to).Msg("Failed to send verification email after 3 attempts")
	}()
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
		var lastErr error
		for attempt := 0; attempt < 3; attempt++ {
			if attempt > 0 {
				time.Sleep(time.Duration(attempt*2) * time.Second) // 2s, 4s backoff
			}
			if err := smtp.SendMail(addr, auth, extractEmailAddress(cfg.SMTPFrom), []string{to}, []byte(msg)); err != nil {
				lastErr = err
				log.Warn().Err(err).Int("attempt", attempt+1).Str("to", to).Msg("Failed to send password reset email, retrying...")
				continue
			}
			log.Info().Str("to", to).Msg("Password reset email sent")
			return
		}
		log.Error().Err(lastErr).Str("to", to).Msg("Failed to send password reset email after 3 attempts")
	}()
}
