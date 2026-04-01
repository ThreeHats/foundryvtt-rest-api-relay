package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/config"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/database"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/handler/helpers"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/model"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
	"github.com/stripe/stripe-go/v78"
	checkoutSession "github.com/stripe/stripe-go/v78/checkout/session"
	stripeCustomer "github.com/stripe/stripe-go/v78/customer"
	"github.com/stripe/stripe-go/v78/webhook"
)

// statusPriority maps subscription statuses to a priority value.
// Higher priority statuses should not be overwritten by lower ones.
var statusPriority = map[string]int{
	"incomplete":         0,
	"incomplete_expired": 0,
	"past_due":           1,
	"canceled":           2,
	"active":             3,
}

// StripeRouter creates Stripe subscription management routes.
// Stripe integration is disabled in local (memory/sqlite) mode.
func StripeRouter(db *database.DB, cfg *config.Config) chi.Router {
	r := chi.NewRouter()

	isDisabled := cfg.StripeSecretKey == ""

	// Set Stripe API key
	if !isDisabled {
		stripe.Key = cfg.StripeSecretKey
	}

	// GET /api/subscriptions/status
	r.Get("/status", func(w http.ResponseWriter, r *http.Request) {
		if isDisabled {
			helpers.WriteJSON(w, http.StatusOK, map[string]interface{}{
				"subscriptionStatus": "free",
				"message":            "Stripe is disabled in local mode",
			})
			return
		}

		reqCtx := helpers.GetRequestContext(r)
		if reqCtx == nil {
			helpers.WriteError(w, http.StatusUnauthorized, "Authentication required")
			return
		}

		helpers.WriteJSON(w, http.StatusOK, map[string]interface{}{
			"subscriptionStatus": reqCtx.SubscriptionStatus,
		})
	})

	// POST /api/subscriptions/create-checkout-session
	r.Post("/create-checkout-session", func(w http.ResponseWriter, r *http.Request) {
		if isDisabled {
			helpers.WriteError(w, http.StatusServiceUnavailable, "Stripe is not configured")
			return
		}

		reqCtx := helpers.GetRequestContext(r)
		if reqCtx == nil {
			helpers.WriteError(w, http.StatusUnauthorized, "Authentication required")
			return
		}

		user, ok := reqCtx.User.(*model.User)
		if !ok || user == nil {
			helpers.WriteError(w, http.StatusUnauthorized, "Invalid user")
			return
		}

		log.Info().Int64("userId", user.ID).Msg("Creating checkout session")

		ctx := r.Context()

		// Get or create Stripe customer
		customerId := user.StripeCustomerID.String
		if !user.StripeCustomerID.Valid || customerId == "" {
			customerParams := &stripe.CustomerParams{
				Email: stripe.String(user.Email),
				Params: stripe.Params{
					Metadata: map[string]string{
						"userId": fmt.Sprintf("%d", user.ID),
					},
				},
			}
			customer, err := stripeCustomer.New(customerParams)
			if err != nil {
				log.Error().Err(err).Int64("userId", user.ID).Msg("Failed to create Stripe customer")
				helpers.WriteError(w, http.StatusInternalServerError, "Failed to create checkout session")
				return
			}

			customerId = customer.ID
			user.StripeCustomerID = sql.NullString{String: customerId, Valid: true}
			if err := db.UserStore().Update(ctx, user); err != nil {
				log.Error().Err(err).Int64("userId", user.ID).Msg("Failed to save Stripe customer ID")
				helpers.WriteError(w, http.StatusInternalServerError, "Failed to create checkout session")
				return
			}
		}

		// Create checkout session
		sessionParams := &stripe.CheckoutSessionParams{
			Customer:           stripe.String(customerId),
			PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
			LineItems: []*stripe.CheckoutSessionLineItemParams{
				{
					Price:    stripe.String(cfg.StripePriceID),
					Quantity: stripe.Int64(1),
				},
			},
			Mode:       stripe.String(string(stripe.CheckoutSessionModeSubscription)),
			SuccessURL: stripe.String(cfg.FrontendURL + "/subscription-success?session_id={CHECKOUT_SESSION_ID}"),
			CancelURL:  stripe.String(cfg.FrontendURL + "/subscription-cancel"),
			Params: stripe.Params{
				Metadata: map[string]string{
					"userId": fmt.Sprintf("%d", user.ID),
				},
			},
		}

		session, err := checkoutSession.New(sessionParams)
		if err != nil {
			log.Error().Err(err).Int64("userId", user.ID).Msg("Failed to create checkout session")
			helpers.WriteError(w, http.StatusInternalServerError, "Failed to create checkout session")
			return
		}

		helpers.WriteJSON(w, http.StatusOK, map[string]string{"url": session.URL})
	})

	// POST /api/subscriptions/create-portal-session
	r.Post("/create-portal-session", func(w http.ResponseWriter, r *http.Request) {
		if isDisabled {
			helpers.WriteError(w, http.StatusServiceUnavailable, "Stripe is not configured")
			return
		}

		reqCtx := helpers.GetRequestContext(r)
		if reqCtx == nil {
			helpers.WriteError(w, http.StatusUnauthorized, "Authentication required")
			return
		}

		user, ok := reqCtx.User.(*model.User)
		if !ok || user == nil {
			helpers.WriteError(w, http.StatusUnauthorized, "Invalid user")
			return
		}

		log.Info().Int64("userId", user.ID).Msg("Redirecting user to customer portal")

		if cfg.StripePortalURL != "" {
			helpers.WriteJSON(w, http.StatusOK, map[string]string{"url": cfg.StripePortalURL})
			return
		}
		helpers.WriteError(w, http.StatusServiceUnavailable, "Stripe portal URL not configured")
	})

	return r
}

// WebhookRouter creates Stripe webhook handler.
func WebhookRouter(db *database.DB, cfg *config.Config) chi.Router {
	r := chi.NewRouter()

	isDisabled := cfg.StripeSecretKey == ""

	// Set Stripe API key
	if !isDisabled {
		stripe.Key = cfg.StripeSecretKey
	}

	r.Post("/stripe", func(w http.ResponseWriter, r *http.Request) {
		if isDisabled {
			helpers.WriteJSON(w, http.StatusOK, map[string]string{"received": "true"})
			return
		}

		// Read raw body for signature verification
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Error().Err(err).Msg("Failed to read webhook body")
			helpers.WriteError(w, http.StatusBadRequest, "Failed to read request body")
			return
		}

		// Verify webhook signature
		event, err := webhook.ConstructEvent(body, r.Header.Get("Stripe-Signature"), cfg.StripeWebhookSecret)
		if err != nil {
			log.Error().Err(err).Msg("Webhook signature verification failed")
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Webhook Error: %v", err)
			return
		}

		// Handle the event
		var success bool
		switch event.Type {
		case "customer.subscription.created", "customer.subscription.updated":
			success = handleSubscriptionUpdated(db, event)
		case "customer.subscription.deleted":
			success = handleSubscriptionDeleted(db, event)
		case "invoice.payment_succeeded":
			success = handlePaymentSucceeded(db, event)
		case "invoice.payment_failed":
			success = handlePaymentFailed(db, event)
		default:
			log.Info().Str("eventType", string(event.Type)).Msg("Unhandled event type")
			success = true
		}

		if success {
			w.WriteHeader(http.StatusOK)
		} else {
			// Return 500 so Stripe will retry
			helpers.WriteError(w, http.StatusInternalServerError, "Failed to process webhook")
		}
	})

	return r
}

// handleSubscriptionUpdated processes customer.subscription.created and customer.subscription.updated events.
func handleSubscriptionUpdated(db *database.DB, event stripe.Event) bool {
	var subscription struct {
		ID               string `json:"id"`
		Customer         string `json:"customer"`
		Status           string `json:"status"`
		CurrentPeriodEnd int64  `json:"current_period_end"`
	}
	if err := json.Unmarshal(event.Data.Raw, &subscription); err != nil {
		log.Error().Err(err).Msg("Failed to parse subscription from webhook event")
		return false
	}

	log.Info().
		Str("customerId", subscription.Customer).
		Str("status", subscription.Status).
		Msg("Processing subscription update")

	ctx := context.Background()
	user, err := db.UserStore().FindByStripeCustomerID(ctx, subscription.Customer)
	if err != nil || user == nil {
		log.Error().Str("customerId", subscription.Customer).Msg("User not found for customer")
		return false
	}

	// Map Stripe statuses: 'trialing' should grant full access
	effectiveStatus := subscription.Status
	if effectiveStatus == "trialing" {
		effectiveStatus = "active"
	}

	// Status priority check to prevent downgrade race conditions
	newPriority, ok := statusPriority[effectiveStatus]
	if !ok {
		newPriority = 1
	}
	currentStatus := user.GetSubscriptionStatus()
	currentPriority, ok := statusPriority[currentStatus]
	if !ok {
		currentPriority = -1
	}

	if newPriority < currentPriority {
		log.Info().
			Int64("userId", user.ID).
			Str("currentStatus", currentStatus).
			Str("newStatus", effectiveStatus).
			Str("stripeStatus", subscription.Status).
			Msg("Skipping downgrade, updating subscription ID and period end only")

		// Still update subscription ID and period end
		user.SubscriptionID = sql.NullString{String: subscription.ID, Valid: true}
		endsAt := time.Unix(subscription.CurrentPeriodEnd, 0)
		user.SubscriptionEndsAt = &model.SQLiteTime{Time: endsAt, Valid: true}
		if err := db.UserStore().Update(ctx, user); err != nil {
			log.Error().Err(err).Int64("userId", user.ID).Msg("Failed to update subscription metadata")
			return false
		}
		return true
	}

	log.Info().
		Int64("userId", user.ID).
		Str("customerId", subscription.Customer).
		Str("effectiveStatus", effectiveStatus).
		Str("stripeStatus", subscription.Status).
		Msg("Updating subscription status")

	user.SubscriptionStatus = sql.NullString{String: effectiveStatus, Valid: true}
	user.SubscriptionID = sql.NullString{String: subscription.ID, Valid: true}
	endsAt := time.Unix(subscription.CurrentPeriodEnd, 0)
	user.SubscriptionEndsAt = &model.SQLiteTime{Time: endsAt, Valid: true}

	if err := db.UserStore().Update(ctx, user); err != nil {
		log.Error().Err(err).Int64("userId", user.ID).Msg("Failed to update subscription")
		return false
	}

	log.Info().Int64("userId", user.ID).Str("status", effectiveStatus).Msg("Successfully updated subscription")
	return true
}

// handleSubscriptionDeleted processes customer.subscription.deleted events.
func handleSubscriptionDeleted(db *database.DB, event stripe.Event) bool {
	var subscription struct {
		Customer   string `json:"customer"`
		CanceledAt int64  `json:"canceled_at"`
	}
	if err := json.Unmarshal(event.Data.Raw, &subscription); err != nil {
		log.Error().Err(err).Msg("Failed to parse subscription deletion from webhook event")
		return false
	}

	ctx := context.Background()
	user, err := db.UserStore().FindByStripeCustomerID(ctx, subscription.Customer)
	if err != nil || user == nil {
		log.Error().Str("customerId", subscription.Customer).Msg("User not found for customer")
		return false
	}

	user.SubscriptionStatus = sql.NullString{String: "canceled", Valid: true}
	if subscription.CanceledAt > 0 {
		canceledAt := time.Unix(subscription.CanceledAt, 0)
		user.SubscriptionEndsAt = &model.SQLiteTime{Time: canceledAt, Valid: true}
	}

	if err := db.UserStore().Update(ctx, user); err != nil {
		log.Error().Err(err).Int64("userId", user.ID).Msg("Failed to update subscription deletion")
		return false
	}

	log.Info().Int64("userId", user.ID).Msg("Subscription canceled")
	return true
}

// handlePaymentSucceeded processes invoice.payment_succeeded events.
// Monthly reset is handled by cron, so this only logs.
func handlePaymentSucceeded(db *database.DB, event stripe.Event) bool {
	var invoice struct {
		Customer     string `json:"customer"`
		Subscription string `json:"subscription"`
	}
	if err := json.Unmarshal(event.Data.Raw, &invoice); err != nil {
		log.Error().Err(err).Msg("Failed to parse invoice from webhook event")
		return false
	}

	if invoice.Subscription != "" {
		ctx := context.Background()
		user, err := db.UserStore().FindByStripeCustomerID(ctx, invoice.Customer)
		if err != nil || user == nil {
			log.Error().Str("customerId", invoice.Customer).Msg("User not found for customer")
			return false
		}
		log.Info().Int64("userId", user.ID).Msg("Payment success recorded")
	}
	return true
}

// handlePaymentFailed processes invoice.payment_failed events.
func handlePaymentFailed(db *database.DB, event stripe.Event) bool {
	var invoice struct {
		Customer     string `json:"customer"`
		Subscription string `json:"subscription"`
	}
	if err := json.Unmarshal(event.Data.Raw, &invoice); err != nil {
		log.Error().Err(err).Msg("Failed to parse invoice from webhook event")
		return false
	}

	if invoice.Subscription != "" {
		ctx := context.Background()
		user, err := db.UserStore().FindByStripeCustomerID(ctx, invoice.Customer)
		if err != nil || user == nil {
			log.Error().Str("customerId", invoice.Customer).Msg("User not found for customer")
			return false
		}

		user.SubscriptionStatus = sql.NullString{String: "past_due", Valid: true}
		if err := db.UserStore().Update(ctx, user); err != nil {
			log.Error().Err(err).Int64("userId", user.ID).Msg("Failed to update subscription to past_due")
			return false
		}

		log.Info().Int64("userId", user.ID).Msg("Updated subscription status to past_due")
	}
	return true
}

