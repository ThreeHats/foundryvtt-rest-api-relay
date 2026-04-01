package service

import (
	"context"

	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/database"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/ws"
	"github.com/rs/zerolog/log"
)

// APIKeyValidationResult holds the result of API key validation.
type APIKeyValidationResult struct {
	Valid          bool
	MasterAPIKey   string
	ScopedClientID string
	ScopedUserID   string
}

// ValidateAPIKeyDetailed validates an API key against the database.
// Supports both master keys and scoped API keys.
func ValidateAPIKeyDetailed(ctx context.Context, db *database.DB, apiKey string, isMemoryStore bool) *APIKeyValidationResult {
	if isMemoryStore {
		return &APIKeyValidationResult{Valid: true}
	}

	// 1. Try master key lookup
	user, err := db.UserStore().FindByAPIKey(ctx, apiKey)
	if err != nil {
		log.Error().Err(err).Msg("Error validating API key")
		return &APIKeyValidationResult{Valid: false}
	}
	if user != nil {
		return &APIKeyValidationResult{Valid: true, MasterAPIKey: apiKey}
	}

	// 2. Try scoped API key lookup
	scopedKey, err := db.ApiKeyStore().FindByKey(ctx, apiKey)
	if err != nil {
		log.Error().Err(err).Msg("Error looking up scoped key")
		return &APIKeyValidationResult{Valid: false}
	}
	if scopedKey == nil {
		return &APIKeyValidationResult{Valid: false}
	}

	if !scopedKey.Enabled {
		return &APIKeyValidationResult{Valid: false}
	}

	if scopedKey.IsExpired() {
		return &APIKeyValidationResult{Valid: false}
	}

	// Look up parent user for master key
	parentUser, err := db.UserStore().FindByID(ctx, scopedKey.UserID)
	if err != nil || parentUser == nil {
		return &APIKeyValidationResult{Valid: false}
	}

	return &APIKeyValidationResult{
		Valid:          true,
		MasterAPIKey:   parentUser.APIKey,
		ScopedClientID: scopedKey.ScopedClientID.String,
		ScopedUserID:   scopedKey.ScopedUserID.String,
	}
}

// ValidateAPIKey performs a simple valid/invalid check.
func ValidateAPIKey(ctx context.Context, db *database.DB, apiKey string, isMemoryStore bool) bool {
	return ValidateAPIKeyDetailed(ctx, db, apiKey, isMemoryStore).Valid
}

// TrackWSAPIUsage tracks per-message API usage for WebSocket connections.
func TrackWSAPIUsage(ctx context.Context, db *database.DB, apiKey string, dailyLimit int, isMemoryStore bool) (bool, string) {
	if isMemoryStore {
		return true, ""
	}

	user, err := db.UserStore().FindByAPIKey(ctx, apiKey)
	if err != nil || user == nil {
		return false, "Invalid API key"
	}

	// Daily rate limit check
	if user.RequestsToday >= dailyLimit {
		return false, "Daily API request limit reached"
	}

	// Increment counters (best effort)
	if err := db.UserStore().IncrementRequests(ctx, user.ID); err != nil {
		log.Warn().Err(err).Msg("Failed to track WS API usage")
	}

	return true, ""
}

// MakeWSValidateAPIKey creates a validation function for the WebSocket relay.
func MakeWSValidateAPIKey(db *database.DB, isMemoryStore bool) func(token string) (*ws.APIKeyValidation, error) {
	return func(token string) (*ws.APIKeyValidation, error) {
		result := ValidateAPIKeyDetailed(context.Background(), db, token, isMemoryStore)
		return &ws.APIKeyValidation{
			Valid:          result.Valid,
			MasterAPIKey:   result.MasterAPIKey,
			ScopedClientID: result.ScopedClientID,
			ScopedUserID:   result.ScopedUserID,
		}, nil
	}
}
