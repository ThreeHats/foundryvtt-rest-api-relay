package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/database"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/ws"
	"github.com/rs/zerolog/log"
)

// APIKeyValidationResult holds the result of API key validation.
type APIKeyValidationResult struct {
	Valid          bool
	UserID         int64  // relay DB user ID; 0 if lookup failed
	MasterAPIKey   string
	ScopedClientID string
	ScopedUserID   string
}

// ValidateAPIKeyDetailed validates an API key against the database.
// Supports both master keys and scoped API keys.
func ValidateAPIKeyDetailed(ctx context.Context, db *database.DB, apiKey string) *APIKeyValidationResult {
	// 1. Try master key lookup
	user, err := db.UserStore().FindByAPIKey(ctx, apiKey)
	if err != nil {
		log.Error().Err(err).Msg("Error validating API key")
		return &APIKeyValidationResult{Valid: false}
	}
	if user != nil {
		return &APIKeyValidationResult{Valid: true, UserID: user.ID, MasterAPIKey: user.APIKeyHash.String}
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

	// MasterAPIKey is now the SHA-256 hash of the user's master key — used as
	// an opaque per-account identifier in ClientManager / Redis. The plaintext
	// master key only exists in the registration/regeneration response and is
	// never stored in the DB or accessible via lookup.
	return &APIKeyValidationResult{
		Valid:          true,
		UserID:         parentUser.ID,
		MasterAPIKey:   parentUser.APIKeyHash.String,
		ScopedClientID: scopedKey.ScopedClientID.String,
		ScopedUserID:   scopedKey.ScopedUserID.String,
	}
}

// ValidateAPIKey performs a simple valid/invalid check.
func ValidateAPIKey(ctx context.Context, db *database.DB, apiKey string) bool {
	return ValidateAPIKeyDetailed(ctx, db, apiKey).Valid
}

// TrackWSAPIUsage tracks per-message API usage for WebSocket connections.
// apiKey is the per-account identifier (the SHA-256 hash of the master API key),
// as returned by ValidateAPIKeyDetailed or MakeWSValidateConnectionToken.
func TrackWSAPIUsage(ctx context.Context, db *database.DB, apiKey string) (bool, string) {
	user, err := db.UserStore().FindByAPIKeyHash(ctx, apiKey)
	if err != nil || user == nil {
		return false, "Invalid API key"
	}

	// Increment counters (best effort)
	if err := db.UserStore().IncrementRequests(ctx, user.ID); err != nil {
		log.Warn().Err(err).Msg("Failed to track WS API usage")
	}

	return true, ""
}

// MakeWSValidateConnectionToken creates a validation function for connection tokens.
// It hashes the raw token, looks it up in the database, and returns the user's
// master API key, the token's IP allowlist, the token's human-readable name,
// and the token's database ID.
//
// The token ID is returned so the relay can record it on the Client struct,
// enabling targeted disconnects when the token is revoked. The token name is
// stored on the client so it can be included in connection log entries.
func MakeWSValidateConnectionToken(db *database.DB) func(token string) (string, string, string, int64, error) {
	return func(token string) (string, string, string, int64, error) {
		// Hash the token
		hash := sha256.Sum256([]byte(token))
		tokenHash := hex.EncodeToString(hash[:])

		ctx := context.Background()
		ct, err := db.ConnectionTokenStore().FindByTokenHash(ctx, tokenHash)
		if err != nil {
			log.Warn().Err(err).Str("hash", tokenHash[:8]+"…").Msg("Connection token lookup DB error")
			return "", "", "", 0, fmt.Errorf("invalid connection token")
		}
		if ct == nil {
			log.Warn().Str("hash", tokenHash[:8]+"…").Msg("Connection token not found in DB")
			return "", "", "", 0, fmt.Errorf("invalid connection token")
		}

		// Look up user to get master API key
		user, err := db.UserStore().FindByID(ctx, ct.UserID)
		if err != nil || user == nil {
			return "", "", "", 0, fmt.Errorf("invalid user for connection token")
		}

		// Update last used asynchronously
		go db.ConnectionTokenStore().UpdateLastUsed(context.Background(), ct.ID)

		allowedIPs := ""
		if ct.AllowedIPs.Valid {
			allowedIPs = ct.AllowedIPs.String
		}

		// Return the user's apiKeyHash as the per-account identifier under
		// which the WS client gets registered in ClientManager. We can't
		// return the plaintext master key because it isn't stored — the hash
		// is the only stable per-account identifier we have.
		return user.APIKeyHash.String, allowedIPs, ct.Name, ct.ID, nil
	}
}

// MakeWSValidateAPIKey creates a validation function for the WebSocket relay.
func MakeWSValidateAPIKey(db *database.DB) func(token string) (*ws.APIKeyValidation, error) {
	return func(token string) (*ws.APIKeyValidation, error) {
		result := ValidateAPIKeyDetailed(context.Background(), db, token)
		return &ws.APIKeyValidation{
			Valid:          result.Valid,
			UserID:         result.UserID,
			MasterAPIKey:   result.MasterAPIKey,
			ScopedClientID: result.ScopedClientID,
			ScopedUserID:   result.ScopedUserID,
		}, nil
	}
}
