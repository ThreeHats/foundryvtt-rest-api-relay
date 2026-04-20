package handler

import (
	"net/http"

	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/alerts"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/database"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/handler/helpers"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/model"
	"github.com/go-chi/chi/v5"
)

// AdminAlertsRouter exposes the admin alert config and recent events.
func AdminAlertsRouter(db *database.DB) chi.Router {
	r := chi.NewRouter()

	validAlertTypes := map[string]bool{
		// Security
		alerts.TypeFailedAuthSpike:          true,
		alerts.TypeFlaggedConnection:        true,
		alerts.TypeWorldIDCrossAccount:      true,
		alerts.TypeRegistrationSpike:        true,
		alerts.TypeRateLimitBurst:           true,
		alerts.TypeAdminLogin:               true,
		alerts.TypeDuplicateConnectionSpike: true,
		alerts.TypeMetadataMismatchSpike:    true,
		alerts.TypeExecuteJsBurst:           true,
		alerts.TypePasswordResetFlood:       true,
		alerts.TypeInvalidTokenSpike:        true,
		alerts.TypeNewIPForToken:            true,
		alerts.TypeAccountPasswordChange:    true,
		// Operations
		alerts.TypeClientDisconnectSpike: true,
		alerts.TypeSystemUnhealthy:       true,
		alerts.TypeHeadlessSessionFlood:  true,
		alerts.TypeStripePaymentFailed:   true,
		// Analytics
		alerts.TypeNewUserRegistration:      true,
		alerts.TypeNewSubscription:          true,
		alerts.TypeSubscriptionCancelled:    true,
		alerts.TypeUserMonthlyLimitApproach:      true,
		alerts.TypeScopedKeyMonthlyLimitApproach: true,
	}

	// GET /admin/api/alerts/config
	r.Get("/config", func(w http.ResponseWriter, req *http.Request) {
		cfg, err := db.AlertConfigStore().Get(req.Context())
		if err != nil {
			helpers.WriteError(w, http.StatusInternalServerError, err.Error())
			return
		}
		subs, err := db.AlertSubscriptionStore().FindAll(req.Context())
		if err != nil {
			helpers.WriteError(w, http.StatusInternalServerError, err.Error())
			return
		}
		type subEntry struct {
			AlertType string `json:"alertType"`
			Channel   string `json:"channel"`
		}
		entries := make([]subEntry, 0, len(subs))
		for _, s := range subs {
			entries = append(entries, subEntry{AlertType: s.AlertType, Channel: s.Channel})
		}
		helpers.WriteJSON(w, http.StatusOK, map[string]interface{}{
			"discordWebhookUrl": cfg.DiscordWebhookURL,
			"emailDestination":  cfg.EmailDestination,
			"subscriptions":     entries,
		})
	})

	// PUT /admin/api/alerts/config
	r.Put("/config", func(w http.ResponseWriter, req *http.Request) {
		body, err := parseBody(req)
		if err != nil {
			helpers.WriteError(w, http.StatusBadRequest, err.Error())
			return
		}

		cfg := &model.AlertConfig{
			DiscordWebhookURL: bodyStr(body, "discordWebhookUrl"),
			EmailDestination:  bodyStr(body, "emailDestination"),
		}
		if err := db.AlertConfigStore().Save(req.Context(), cfg); err != nil {
			helpers.WriteError(w, http.StatusInternalServerError, err.Error())
			return
		}

		// Parse and validate subscriptions array
		rawSubs, _ := body["subscriptions"].([]interface{})
		newSubs := make([]*model.AlertSubscription, 0, len(rawSubs))
		for _, item := range rawSubs {
			m, ok := item.(map[string]interface{})
			if !ok {
				continue
			}
			at, _ := m["alertType"].(string)
			ch, _ := m["channel"].(string)
			if at == "" || ch == "" {
				helpers.WriteError(w, http.StatusBadRequest, "each subscription must have alertType and channel")
				return
			}
			if !validAlertTypes[at] {
				helpers.WriteError(w, http.StatusBadRequest, "invalid alertType: "+at)
				return
			}
			if ch != "discord" && ch != "email" {
				helpers.WriteError(w, http.StatusBadRequest, "channel must be 'discord' or 'email'")
				return
			}
			newSubs = append(newSubs, &model.AlertSubscription{AlertType: at, Channel: ch, Enabled: true})
		}

		if err := db.AlertSubscriptionStore().BulkReplace(req.Context(), newSubs); err != nil {
			helpers.WriteError(w, http.StatusInternalServerError, err.Error())
			return
		}

		auditAdmin(req, db, "alerts.config_updated", "alert", "", "")

		type subEntry struct {
			AlertType string `json:"alertType"`
			Channel   string `json:"channel"`
		}
		entries := make([]subEntry, 0, len(newSubs))
		for _, s := range newSubs {
			entries = append(entries, subEntry{AlertType: s.AlertType, Channel: s.Channel})
		}
		helpers.WriteJSON(w, http.StatusOK, map[string]interface{}{
			"discordWebhookUrl": cfg.DiscordWebhookURL,
			"emailDestination":  cfg.EmailDestination,
			"subscriptions":     entries,
		})
	})

	// POST /admin/api/alerts/test
	r.Post("/test", func(w http.ResponseWriter, req *http.Request) {
		body, err := parseBody(req)
		if err != nil {
			helpers.WriteError(w, http.StatusBadRequest, err.Error())
			return
		}
		channel := bodyStr(body, "channel")
		if channel != "discord" && channel != "email" {
			helpers.WriteError(w, http.StatusBadRequest, "channel must be 'discord' or 'email'")
			return
		}
		cfg, err := db.AlertConfigStore().Get(req.Context())
		if err != nil {
			helpers.WriteError(w, http.StatusInternalServerError, err.Error())
			return
		}
		dest := ""
		if channel == "discord" {
			dest = cfg.DiscordWebhookURL
		} else {
			dest = cfg.EmailDestination
		}
		if err := alerts.TestChannel(channel, dest); err != nil {
			helpers.WriteError(w, http.StatusInternalServerError, err.Error())
			return
		}
		helpers.WriteJSON(w, http.StatusOK, map[string]string{"message": "Test alert sent"})
	})

	// GET /admin/api/alerts/recent
	r.Get("/recent", func(w http.ResponseWriter, req *http.Request) {
		helpers.WriteJSON(w, http.StatusOK, map[string]interface{}{
			"events": alerts.RecentEvents(),
		})
	})

	return r
}
