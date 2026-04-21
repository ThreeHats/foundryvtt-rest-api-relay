package model

import (
	"slices"
	"strings"
)

// Scope constants define the canonical scope strings for API key permissions.
const (
	ScopeEntityRead      = "entity:read"
	ScopeEntityWrite     = "entity:write"
	ScopeRollRead        = "roll:read"
	ScopeRollExecute     = "roll:execute"
	ScopeChatRead        = "chat:read"
	ScopeChatWrite       = "chat:write"
	ScopeEncounterRead   = "encounter:read"
	ScopeEncounterManage = "encounter:manage"
	ScopeMacroList       = "macro:list"
	ScopeMacroExecute    = "macro:execute"
	ScopeMacroWrite      = "macro:write"
	ScopeSceneRead       = "scene:read"
	ScopeSceneWrite      = "scene:write"
	ScopeCanvasRead      = "canvas:read"
	ScopeCanvasWrite     = "canvas:write"
	ScopeEffectsRead     = "effects:read"
	ScopeEffectsWrite    = "effects:write"
	ScopeUserRead        = "user:read"
	ScopeUserWrite       = "user:write"
	ScopeFileRead        = "file:read"
	ScopeFileWrite       = "file:write"
	ScopePlaylistControl = "playlist:control"
	ScopeWorldInfo       = "world:info"
	ScopeClientsRead     = "clients:read"
	ScopeSheetRead       = "sheet:read"
	ScopeEventsSubscribe = "events:subscribe"
	ScopeSessionManage   = "session:manage"
	ScopeExecuteJS       = "execute-js"
	ScopeSearch          = "search"
	ScopeStructureRead   = "structure:read"
	ScopeStructureWrite  = "structure:write"
	ScopeDnd5e           = "dnd5e"
)

// DefaultScopes is the safe, read-only set granted to new API keys.
// Does NOT include execute-js, macro:execute, macro:write, or any write scopes.
var DefaultScopes = []string{
	ScopeEntityRead,
	ScopeRollRead,
	ScopeChatRead,
	ScopeEncounterRead,
	ScopeMacroList,
	ScopeSceneRead,
	ScopeCanvasRead,
	ScopeEffectsRead,
	ScopeUserRead,
	ScopeFileRead,
	ScopeWorldInfo,
	ScopeClientsRead,
	ScopeSheetRead,
	ScopeEventsSubscribe,
	ScopeSearch,
	ScopeStructureRead,
	ScopeDnd5e,
}

// AllScopes is the exhaustive set of every scope (for master keys and explicit opt-in).
var AllScopes = []string{
	ScopeEntityRead, ScopeEntityWrite,
	ScopeRollRead, ScopeRollExecute,
	ScopeChatRead, ScopeChatWrite,
	ScopeEncounterRead, ScopeEncounterManage,
	ScopeMacroList, ScopeMacroExecute, ScopeMacroWrite,
	ScopeSceneRead, ScopeSceneWrite,
	ScopeCanvasRead, ScopeCanvasWrite,
	ScopeEffectsRead, ScopeEffectsWrite,
	ScopeUserRead, ScopeUserWrite,
	ScopeFileRead, ScopeFileWrite,
	ScopePlaylistControl,
	ScopeWorldInfo,
	ScopeClientsRead,
	ScopeSheetRead,
	ScopeEventsSubscribe,
	ScopeSessionManage,
	ScopeExecuteJS,
	ScopeSearch,
	ScopeStructureRead, ScopeStructureWrite,
	ScopeDnd5e,
}

// ActionToScopeRequired maps WS action message types to the scope a caller
// must hold to invoke them via remote-request. This is the canonical mapping
// the cross-world tunnel uses to gate which actions a connection token can
// perform on allowed-target clients.
//
// The keys here MUST match the message `type` strings the relay forwards to
// target Foundry modules (see internal/handler/*.go for the HTTP handlers
// that send the corresponding action messages over WS today).
//
// Actions not listed here are denied by default — adding a new action
// REQUIRES adding it to this map.
var ActionToScopeRequired = map[string]string{
	"get":              ScopeEntityRead,
	"create":           ScopeEntityWrite,
	"update":           ScopeEntityWrite,
	"delete":           ScopeEntityWrite,
	"give":             ScopeEntityWrite,
	"remove":           ScopeEntityWrite,
	"decrease":         ScopeEntityWrite,
	"increase":         ScopeEntityWrite,
	"kill":             ScopeEntityWrite,
	"search":           ScopeSearch,
	"rolls":                 ScopeRollRead,
	"last-roll":             ScopeRollRead,
	"roll":                  ScopeRollExecute,
	"chat-messages":         ScopeChatRead,
	"chat-send":             ScopeChatWrite,
	"encounters":            ScopeEncounterRead,
	"start-encounter":       ScopeEncounterManage,
	"next-turn":             ScopeEncounterManage,
	"next-round":            ScopeEncounterManage,
	"last-turn":             ScopeEncounterManage,
	"last-round":            ScopeEncounterManage,
	"end-encounter":         ScopeEncounterManage,
	"add-to-encounter":      ScopeEncounterManage,
	"remove-from-encounter": ScopeEncounterManage,
	"macros":                ScopeMacroList,
	"macro-execute":         ScopeMacroExecute,
	"create-macro":          ScopeMacroWrite,
	"update-macro":          ScopeMacroWrite,
	"delete-macro":          ScopeMacroWrite,
	"get-scene":             ScopeSceneRead,
	"create-scene":          ScopeSceneWrite,
	"update-scene":          ScopeSceneWrite,
	"delete-scene":          ScopeSceneWrite,
	"switch-scene":          ScopeSceneWrite,
	"get-users":             ScopeUserRead,
	"get-user":              ScopeUserRead,
	"create-user":           ScopeUserWrite,
	"update-user":           ScopeUserWrite,
	"delete-user":           ScopeUserWrite,
	"file-system":           ScopeFileRead,
	"download-file":         ScopeFileRead,
	"upload-file":           ScopeFileWrite,
	"create-folder":         ScopeStructureWrite,
	"delete-folder":         ScopeStructureWrite,
	"structure":             ScopeStructureRead,
	"get-folder":            ScopeStructureRead,
	"sheet-screenshot":      ScopeSheetRead,
	"scene-screenshot":      ScopeSceneRead,
	"scene-raw-image":       ScopeSceneRead,
	"get-playlists":         ScopePlaylistControl,
	"playlist-play":         ScopePlaylistControl,
	"playlist-stop":         ScopePlaylistControl,
	"playlist-next":         ScopePlaylistControl,
	"playlist-volume":       ScopePlaylistControl,
	"play-sound":            ScopePlaylistControl,
	"stop-sound":            ScopePlaylistControl,
	"world-info":            ScopeWorldInfo,
	"execute-js":            ScopeExecuteJS,
}

// ScopeForAction returns the scope required to invoke a given action via
// remote-request, plus a boolean indicating whether the action is known.
// Unknown actions MUST be denied — never allowed by default.
func ScopeForAction(action string) (string, bool) {
	scope, ok := ActionToScopeRequired[action]
	return scope, ok
}

// validScopes is a lookup set for validation.
var validScopes map[string]bool

func init() {
	validScopes = make(map[string]bool, len(AllScopes))
	for _, s := range AllScopes {
		validScopes[s] = true
	}
}

// IsValidScope returns true if the given scope string is a known scope.
func IsValidScope(scope string) bool {
	return validScopes[scope]
}

// HasScope returns true if the given scope is present in the list.
func HasScope(scopes []string, required string) bool {
	return slices.Contains(scopes, required)
}

// ParseScopes splits a comma-separated scope string into a slice.
// Returns nil for empty strings.
func ParseScopes(csv string) []string {
	csv = strings.TrimSpace(csv)
	if csv == "" {
		return nil
	}
	parts := strings.Split(csv, ",")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			result = append(result, p)
		}
	}
	if len(result) == 0 {
		return nil
	}
	return result
}

// ScopesString joins a scope slice into a comma-separated string.
func ScopesString(scopes []string) string {
	return strings.Join(scopes, ",")
}

// ValidateScopes checks that all scopes in the list are valid.
// Returns the first invalid scope found, or empty string if all valid.
func ValidateScopes(scopes []string) string {
	for _, s := range scopes {
		if !IsValidScope(s) {
			return s
		}
	}
	return ""
}

// DefaultScopesString returns the default scopes as a comma-separated string.
func DefaultScopesString() string {
	return ScopesString(DefaultScopes)
}
