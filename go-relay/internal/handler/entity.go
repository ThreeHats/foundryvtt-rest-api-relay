package handler

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/alerts"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/config"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/database"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/handler/helpers"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/model"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/service"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/ws"
)

// Common param definitions
func clientIDParam() helpers.ParamDef {
	return helpers.ParamDef{Name: "clientId", From: []helpers.ParamSource{helpers.FromQuery}, Type: helpers.TypeString, Description: "Client ID for the Foundry world"}
}
func userIDParam() helpers.ParamDef {
	return helpers.ParamDef{Name: "userId", From: []helpers.ParamSource{helpers.FromQuery, helpers.FromBody}, Type: helpers.TypeString, Description: "Foundry user ID or username to scope permissions (omit for GM-level access)"}
}
func uuidParam() helpers.ParamDef {
	return helpers.ParamDef{Name: "uuid", From: []helpers.ParamSource{helpers.FromQuery}, Type: helpers.TypeString, Description: "UUID of the entity to retrieve (optional if selected=true)"}
}
func selectedParam() helpers.ParamDef {
	return helpers.ParamDef{Name: "selected", From: []helpers.ParamSource{helpers.FromQuery, helpers.FromBody}, Type: helpers.TypeBoolean, Description: "Whether to get the selected entity"}
}

var bqParam = []helpers.ParamSource{helpers.FromBody, helpers.FromQuery}

// --- Entity route configs ---

// Get entity details
//
// This endpoint retrieves the details of a specific entity.
// @tag Entity
// @returns Entity details object containing requested information
func entityGet() helpers.APIRouteConfig {
	return helpers.APIRouteConfig{
		Type: "entity",
		OptionalParams: []helpers.ParamDef{clientIDParam(), uuidParam(), selectedParam(),
			{Name: "actor", From: []helpers.ParamSource{helpers.FromQuery}, Type: helpers.TypeBoolean, Description: "Return the actor of specified entity"}, userIDParam()},
	}
}

// Create a new entity
//
// This endpoint creates a new entity in the Foundry world.
// @tag Entity
// @returns Result of the entity creation operation
func entityCreate() helpers.APIRouteConfig {
	return helpers.APIRouteConfig{
		Type: "create",
		RequiredParams: []helpers.ParamDef{
			{Name: "entityType", From: []helpers.ParamSource{helpers.FromBody}, Type: helpers.TypeString, Required: true, Description: "Document type of entity to create (Scene, Actor, Item, JournalEntry, RollTable, Cards, Macro, Playlist, ext.)"},
			{Name: "data", From: []helpers.ParamSource{helpers.FromBody}, Type: helpers.TypeObject, Required: true, Description: "Data for the new entity"},
		},
		OptionalParams: []helpers.ParamDef{clientIDParam(), {Name: "folder", From: []helpers.ParamSource{helpers.FromBody}, Type: helpers.TypeString, Description: "Optional folder UUID to place the new entity in"}, {Name: "keepId", From: bqParam, Type: helpers.TypeBoolean, Description: "If true, preserve the _id from the provided data instead of generating a new one"}, {Name: "override", From: bqParam, Type: helpers.TypeBoolean, Description: "If true and keepId is set, replace any existing entity with the same ID"}, userIDParam()},
		ValidateParams: func(params helpers.Params, r *http.Request) (map[string]interface{}, bool) {
			if params.GetString("entityType") == "Macro" {
				reqCtx := helpers.GetRequestContext(r)
				if reqCtx != nil && reqCtx.ScopedKey != nil && !reqCtx.ScopedKey.HasScope(model.ScopeMacroWrite) {
					return map[string]interface{}{"error": "macro:write scope is required to create macros"}, true
				}
				if data, ok := params["data"].(map[string]interface{}); ok {
					if cmd, ok := data["command"].(string); ok && !helpers.ValidateScript(cmd) {
						return map[string]interface{}{"error": "Script contains forbidden patterns"}, true
					}
				}
			}
			return nil, false
		},
	}
}

// Update an entity
//
// This endpoint updates an existing entity in the Foundry world.
// @tag Entity
// @returns Result of the entity update operation
func entityUpdate() helpers.APIRouteConfig {
	return helpers.APIRouteConfig{
		Type:           "update",
		RequiredParams: []helpers.ParamDef{{Name: "data", From: []helpers.ParamSource{helpers.FromBody}, Type: helpers.TypeObject, Required: true, Description: "Object containing the fields to update"}},
		OptionalParams: []helpers.ParamDef{clientIDParam(), uuidParam(), selectedParam(), {Name: "actor", From: []helpers.ParamSource{helpers.FromQuery}, Type: helpers.TypeBoolean, Description: "Whether to update the actor of specified entity"}, userIDParam()},
	}
}

// Delete an entity
//
// This endpoint deletes an entity from the Foundry world.
// @tag Entity
// @returns Result of the entity deletion operation
func entityDelete() helpers.APIRouteConfig {
	return helpers.APIRouteConfig{Type: "delete", OptionalParams: []helpers.ParamDef{clientIDParam(), uuidParam(), selectedParam(), userIDParam()}}
}

// Give an item to another entity
//
// Transfers an item from one entity to another.
// @tag Entity
// @returns Result of the give operation
func entityGive() helpers.APIRouteConfig {
	return helpers.APIRouteConfig{
		Type: "give",
		OptionalParams: []helpers.ParamDef{
			clientIDParam(),
			{Name: "fromUuid", From: []helpers.ParamSource{helpers.FromBody}, Description: "UUID of the entity giving the item"},
			{Name: "toUuid", From: []helpers.ParamSource{helpers.FromBody}, Description: "UUID of the entity receiving the item"},
			selectedParam(),
			{Name: "itemUuid", From: []helpers.ParamSource{helpers.FromBody}, Description: "UUID of the item to give"},
			{Name: "itemName", From: []helpers.ParamSource{helpers.FromBody}, Description: "Name of the item to give (alternative to itemUuid)"},
			{Name: "quantity", From: []helpers.ParamSource{helpers.FromBody}, Type: helpers.TypeNumber, Description: "Quantity of items to give"},
			userIDParam(),
		},
	}
}

// Remove an item from an entity
//
// Removes an item from an entity's inventory.
// @tag Entity
// @returns Result of the remove operation
func entityRemove() helpers.APIRouteConfig {
	return helpers.APIRouteConfig{
		Type: "remove",
		OptionalParams: []helpers.ParamDef{
			clientIDParam(),
			{Name: "actorUuid", From: []helpers.ParamSource{helpers.FromBody}, Description: "UUID of the actor to remove the item from"},
			selectedParam(),
			{Name: "itemUuid", From: []helpers.ParamSource{helpers.FromBody}, Description: "UUID of the item to remove"},
			{Name: "itemName", From: []helpers.ParamSource{helpers.FromBody}, Description: "Name of the item to remove (alternative to itemUuid)"},
			{Name: "quantity", From: []helpers.ParamSource{helpers.FromBody}, Type: helpers.TypeNumber, Description: "Quantity of items to remove"},
			userIDParam(),
		},
	}
}

// Decrease an attribute
//
// Decreases a numeric attribute of an entity by the specified amount.
// @tag Entity
// @returns Result of the decrease operation
func entityDecrease() helpers.APIRouteConfig {
	return helpers.APIRouteConfig{
		Type: "decrease",
		RequiredParams: []helpers.ParamDef{
			{Name: "attribute", From: []helpers.ParamSource{helpers.FromBody}, Type: helpers.TypeString, Required: true, Description: "The attribute to decrease (e.g., hp.value)"},
			{Name: "amount", From: []helpers.ParamSource{helpers.FromBody}, Type: helpers.TypeNumber, Required: true, Description: "The amount to decrease by"},
		},
		OptionalParams: []helpers.ParamDef{clientIDParam(), uuidParam(), selectedParam(), userIDParam()},
	}
}

// Increase an attribute
//
// Increases a numeric attribute of an entity by the specified amount.
// @tag Entity
// @returns Result of the increase operation
func entityIncrease() helpers.APIRouteConfig {
	return helpers.APIRouteConfig{
		Type: "increase",
		RequiredParams: []helpers.ParamDef{
			{Name: "attribute", From: []helpers.ParamSource{helpers.FromBody}, Type: helpers.TypeString, Required: true, Description: "The attribute to increase (e.g., hp.value)"},
			{Name: "amount", From: []helpers.ParamSource{helpers.FromBody}, Type: helpers.TypeNumber, Required: true, Description: "The amount to increase by"},
		},
		OptionalParams: []helpers.ParamDef{clientIDParam(), uuidParam(), selectedParam(), userIDParam()},
	}
}

// Kill an entity
//
// Sets the entity's HP to 0.
// @tag Entity
// @returns Result of the kill operation
func entityKill() helpers.APIRouteConfig {
	return helpers.APIRouteConfig{Type: "kill", OptionalParams: []helpers.ParamDef{clientIDParam(), uuidParam(), selectedParam(), userIDParam()}}
}

// --- Search ---

// Search entities
//
// This endpoint allows searching for entities in the Foundry world based on a query string.
// Search world entities and compendiums using the native built-in search engine.
// No third-party modules required. Results are ranked by relevance: exact match,
// prefix match, substring match, word-prefix match, and subsequence match.
// @tag Search
// @returns Search results containing matching entities
func searchGet() helpers.APIRouteConfig {
	return helpers.APIRouteConfig{
		Type: "search",
		OptionalParams: []helpers.ParamDef{
			clientIDParam(),
			{Name: "query", From: []helpers.ParamSource{helpers.FromQuery}, Type: helpers.TypeString, Description: "Search query string (omit to browse all entities matching filter)"},
			{Name: "filter", From: []helpers.ParamSource{helpers.FromQuery}, Description: "Filter string — simple: filter=\"Actor\"; compound: filter=\"documentType:Item,subType:weapon\". Supported keys: documentType, subType, folder, package, resultType"},
			{Name: "excludeCompendiums", From: []helpers.ParamSource{helpers.FromQuery}, Type: helpers.TypeBoolean, Description: "Exclude compendium entries from results (default: false — compendiums are included by default)"},
			{Name: "limit", From: []helpers.ParamSource{helpers.FromQuery}, Type: helpers.TypeNumber, Description: "Maximum number of results to return (default: 200, max: 500)"},
			{Name: "minified", From: []helpers.ParamSource{helpers.FromQuery}, Type: helpers.TypeBoolean, Description: "Return minimal fields only — uuid, id, name, img, documentType (default: false)"},
			userIDParam(),
		},
		ValidateParams: func(params helpers.Params, r *http.Request) (map[string]interface{}, bool) {
			if !params.Has("query") && !params.Has("filter") {
				return map[string]interface{}{"error": "'query' or 'filter' is required"}, true
			}
			return nil, false
		},
	}
}

// --- Rolls ---

// Get recent rolls
//
// Retrieves a list of up to 20 recent rolls made in the Foundry world.
// @tag Roll
// @returns An array of recent rolls with details
func rollsGet() helpers.APIRouteConfig {
	return helpers.APIRouteConfig{Type: "rolls", OptionalParams: []helpers.ParamDef{clientIDParam(), {Name: "limit", From: []helpers.ParamSource{helpers.FromQuery}, Type: helpers.TypeNumber, Description: "Optional limit on the number of rolls to return (default is 20)"}, userIDParam()}}
}

// Get the last roll
//
// Retrieves the most recent roll made in the Foundry world.
// @tag Roll
// @returns The most recent roll with details
func lastRollGet() helpers.APIRouteConfig {
	return helpers.APIRouteConfig{Type: "last-roll", OptionalParams: []helpers.ParamDef{clientIDParam(), userIDParam()}}
}

// Make a roll
//
// Executes a roll with the specified formula.
// @tag Roll
// @returns Result of the roll operation
func rollPost() helpers.APIRouteConfig {
	return helpers.APIRouteConfig{
		Type: "roll",
		RequiredParams: []helpers.ParamDef{
			{Name: "formula", From: []helpers.ParamSource{helpers.FromBody}, Type: helpers.TypeString, Required: true, Description: "The roll formula to evaluate (e.g., \"1d20 + 5\")"},
		},
		OptionalParams: []helpers.ParamDef{
			clientIDParam(),
			{Name: "flavor", From: []helpers.ParamSource{helpers.FromBody}, Description: "Optional flavor text for the roll"},
			{Name: "createChatMessage", From: []helpers.ParamSource{helpers.FromBody}, Type: helpers.TypeBoolean, Description: "Whether to create a chat message for the roll"},
			{Name: "speaker", From: []helpers.ParamSource{helpers.FromBody}, Description: "The speaker for the roll"},
			{Name: "whisper", From: []helpers.ParamSource{helpers.FromBody}, Type: helpers.TypeArray, Description: "Users to whisper the roll result to"},
			userIDParam(),
		},
	}
}

// --- Encounters ---

// Get all active encounters
//
// Retrieves a list of all currently active encounters in the Foundry world.
// @tag Encounter
// @returns An array of active encounters with details
func encountersGet() helpers.APIRouteConfig {
	return helpers.APIRouteConfig{Type: "encounters", OptionalParams: []helpers.ParamDef{clientIDParam(), userIDParam()}}
}

// Start a new encounter
//
// Initiates a new encounter in the Foundry world.
// @tag Encounter
// @returns Details of the started encounter
func startEncounter() helpers.APIRouteConfig {
	return helpers.APIRouteConfig{
		Type: "start-encounter",
		OptionalParams: []helpers.ParamDef{
			clientIDParam(),
			{Name: "tokens", From: []helpers.ParamSource{helpers.FromBody}, Type: helpers.TypeArray, Description: "Array of token UUIDs to include in the encounter"},
			{Name: "startWithSelected", From: []helpers.ParamSource{helpers.FromBody}, Type: helpers.TypeBoolean, Description: "Whether to start with selected tokens"},
			{Name: "startWithPlayers", From: []helpers.ParamSource{helpers.FromBody}, Type: helpers.TypeBoolean, Description: "Whether to start with players"},
			{Name: "rollNPC", From: []helpers.ParamSource{helpers.FromBody}, Type: helpers.TypeBoolean, Description: "Whether to roll for NPCs"},
			{Name: "rollAll", From: []helpers.ParamSource{helpers.FromBody}, Type: helpers.TypeBoolean, Description: "Whether to roll for all tokens"},
			{Name: "name", From: []helpers.ParamSource{helpers.FromBody}, Description: "The name of the encounter"},
			userIDParam(),
		},
	}
}

var encounterParam = helpers.ParamDef{Name: "encounter", From: bqParam, Description: "The ID of the encounter (optional, defaults to current encounter)"}

// Advance to the next turn in the encounter
//
// Moves the encounter to the next turn.
// @tag Encounter
// @returns Details of the next turn
func nextTurn() helpers.APIRouteConfig { return encSimple("next-turn") }

// Advance to the next round in the encounter
//
// Moves the encounter to the next round.
// @tag Encounter
// @returns Details of the next round
func nextRound() helpers.APIRouteConfig { return encSimple("next-round") }

// Go back to the last turn in the encounter
//
// Moves the encounter back to the last turn.
// @tag Encounter
// @returns Details of the last turn
func lastTurn() helpers.APIRouteConfig { return encSimple("last-turn") }

// Go back to the last round in the encounter
//
// Moves the encounter back to the last round.
// @tag Encounter
// @returns Details of the last round
func lastRound() helpers.APIRouteConfig { return encSimple("last-round") }

// End an encounter
//
// Ends the current encounter in the Foundry world.
// @tag Encounter
// @returns Details of the ended encounter
func endEncounter() helpers.APIRouteConfig { return encSimple("end-encounter") }

func encSimple(t string) helpers.APIRouteConfig {
	return helpers.APIRouteConfig{Type: t, OptionalParams: []helpers.ParamDef{clientIDParam(), encounterParam, userIDParam()}}
}

// Add tokens to an encounter
//
// Adds selected tokens or specified UUIDs to the current encounter.
// @tag Encounter
// @returns Details of the updated encounter
func addToEncounter() helpers.APIRouteConfig {
	return helpers.APIRouteConfig{Type: "add-to-encounter", OptionalParams: []helpers.ParamDef{
		clientIDParam(), encounterParam, selectedParam(),
		{Name: "uuids", From: []helpers.ParamSource{helpers.FromBody}, Type: helpers.TypeArray, Description: "The UUIDs of the tokens to add"},
		{Name: "rollInitiative", From: []helpers.ParamSource{helpers.FromBody}, Type: helpers.TypeBoolean, Description: "Whether to roll initiative for the added tokens"},
		userIDParam(),
	}}
}

// Remove tokens from an encounter
//
// Removes selected tokens or specified UUIDs from the current encounter.
// @tag Encounter
// @returns Details of the updated encounter
func removeFromEncounter() helpers.APIRouteConfig {
	return helpers.APIRouteConfig{Type: "remove-from-encounter", OptionalParams: []helpers.ParamDef{
		clientIDParam(), encounterParam, selectedParam(),
		{Name: "uuids", From: []helpers.ParamSource{helpers.FromBody}, Type: helpers.TypeArray, Description: "The UUIDs of the tokens to remove"},
		userIDParam(),
	}}
}

// --- Macros ---

// Get all macros
//
// Retrieves a list of all macros available in the Foundry world.
// @tag Macro
// @returns An array of macros with details
func macrosGet() helpers.APIRouteConfig {
	return helpers.APIRouteConfig{Type: "macros", OptionalParams: []helpers.ParamDef{clientIDParam(), userIDParam()}}
}

// Execute a macro by UUID
//
// Executes a specific macro in the Foundry world by its UUID.
// @tag Macro
// @returns Result of the macro execution
func macroExecute(disp *service.Dispatcher, db *database.DB) helpers.APIRouteConfig {
	return helpers.APIRouteConfig{
		Type: "macro-execute",
		RequiredParams: []helpers.ParamDef{
			{Name: "uuid", From: []helpers.ParamSource{helpers.FromParams}, Type: helpers.TypeString, Required: true, Description: "UUID of the macro to execute"},
		},
		OptionalParams: []helpers.ParamDef{
			clientIDParam(),
			{Name: "args", From: []helpers.ParamSource{helpers.FromBody}, Type: helpers.TypeObject, Description: "Optional arguments to pass to the macro execution"},
			userIDParam(),
		},
		AfterForward: func(reqCtx *helpers.RequestContext, clientID string, params helpers.Params) {
			user, ok := reqCtx.GetUser()
			if !ok {
				return
			}
			uuid := params.GetString("uuid")
			var apiKeyID int64
			if reqCtx.ScopedKey != nil {
				apiKeyID = reqCtx.ScopedKey.ID
			}

			disp.Dispatch(service.NotificationContext{
				Event:       service.EventMacroExecute,
				UserID:      user.ID,
				ApiKeyID:    apiKeyID,
				ClientID:    clientID,
				Description: fmt.Sprintf("**Macro UUID:** `%s`", uuid),
			})

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			_ = db.ModuleEventLogStore().Create(ctx, &model.ModuleEventLog{
				UserID:      user.ID,
				ClientID:    clientID,
				EventType:   "macro-execute",
				Description: uuid,
			})
		},
	}
}

// --- Structure ---

// Get the structure of the Foundry world
//
// Retrieves the folder and compendium structure for the specified Foundry world.
// @tag Structure
// @returns The folder and compendium structure
func structureGet() helpers.APIRouteConfig {
	return helpers.APIRouteConfig{Type: "structure", OptionalParams: []helpers.ParamDef{
		clientIDParam(),
		{Name: "includeEntityData", From: []helpers.ParamSource{helpers.FromQuery}, Type: helpers.TypeBoolean, Description: "Whether to include full entity data or just UUIDs and names"},
		{Name: "path", From: []helpers.ParamSource{helpers.FromQuery}, Description: "Path to read structure from (null = root)"},
		{Name: "recursive", From: []helpers.ParamSource{helpers.FromQuery}, Type: helpers.TypeBoolean, Description: "Whether to read down the folder tree"},
		{Name: "recursiveDepth", From: []helpers.ParamSource{helpers.FromQuery}, Type: helpers.TypeNumber, Description: "Depth to recurse into folders (default 5)"},
		{Name: "types", From: []helpers.ParamSource{helpers.FromQuery}, Description: "Types to return (Scene/Actor/Item/JournalEntry/RollTable/Cards/Macro/Playlist), can be comma-separated or JSON array"},
		userIDParam(),
	}}
}

// This route is deprecated
//
// Use /structure with the path query parameter instead.
// @tag Structure
// @returns Error message directing to use /structure endpoint
func contentsDeprecated(w http.ResponseWriter, r *http.Request) {
	helpers.WriteJSON(w, http.StatusBadRequest, map[string]string{
		"error":   "This endpoint has been deprecated.",
		"message": "Use GET /structure with the path query parameter instead.",
		"example": "/structure?path=your/path&recursive=true",
	})
}

// Get a specific folder by name
//
// @tag Structure
// @returns The folder information and its contents
func getFolderRoute() helpers.APIRouteConfig {
	return helpers.APIRouteConfig{
		Type:           "get-folder",
		RequiredParams: []helpers.ParamDef{{Name: "name", From: bqParam, Type: helpers.TypeString, Required: true, Description: "Name of the folder to retrieve"}},
		OptionalParams: []helpers.ParamDef{clientIDParam(), userIDParam()},
	}
}

// Create a new folder
//
// @tag Structure
// @returns The created folder information
func createFolderRoute() helpers.APIRouteConfig {
	return helpers.APIRouteConfig{
		Type: "create-folder",
		RequiredParams: []helpers.ParamDef{
			{Name: "name", From: bqParam, Type: helpers.TypeString, Required: true, Description: "Name of the new folder"},
			{Name: "folderType", From: bqParam, Type: helpers.TypeString, Required: true, Description: "Type of folder (Scene, Actor, Item, JournalEntry, RollTable, Cards, Macro, Playlist)"},
		},
		OptionalParams: []helpers.ParamDef{clientIDParam(), {Name: "parentFolderId", From: bqParam, Description: "ID of the parent folder (optional for root level)"}, userIDParam()},
	}
}

// Delete a folder
//
// @tag Structure
// @returns Confirmation of deletion
func deleteFolderRoute() helpers.APIRouteConfig {
	return helpers.APIRouteConfig{
		Type:           "delete-folder",
		RequiredParams: []helpers.ParamDef{{Name: "folderId", From: bqParam, Type: helpers.TypeString, Required: true, Description: "ID of the folder to delete"}},
		OptionalParams: []helpers.ParamDef{clientIDParam(), {Name: "deleteAll", From: bqParam, Type: helpers.TypeBoolean, Description: "Whether to delete all entities in the folder"}, userIDParam()},
	}
}

// --- Utility ---

// Select token(s)
//
// Selects one or more tokens in the Foundry VTT client.
// @tag Utility
// @returns The selected token(s)
func selectPost() helpers.APIRouteConfig {
	return helpers.APIRouteConfig{
		Type: "select",
		OptionalParams: []helpers.ParamDef{
			clientIDParam(),
			{Name: "uuids", From: []helpers.ParamSource{helpers.FromBody}, Type: helpers.TypeArray, Description: "Array of UUIDs to select"},
			{Name: "name", From: []helpers.ParamSource{helpers.FromBody}, Description: "Name of the token(s) to select"},
			{Name: "data", From: []helpers.ParamSource{helpers.FromBody}, Type: helpers.TypeObject, Description: "Data to match for selection (e.g., \"data.attributes.hp.value\": 20)"},
			{Name: "overwrite", From: []helpers.ParamSource{helpers.FromBody}, Type: helpers.TypeBoolean, Description: "Whether to overwrite existing selection"},
			{Name: "all", From: []helpers.ParamSource{helpers.FromBody}, Type: helpers.TypeBoolean, Description: "Whether to select all tokens on the canvas"},
			userIDParam(),
		},
		ValidateParams: func(params helpers.Params, r *http.Request) (map[string]interface{}, bool) {
			if !params.Has("uuids") && !params.Has("name") && !params.Has("data") && !params.Has("all") {
				return map[string]interface{}{"error": "At least one of 'uuids', 'name', 'data', or 'all' is required"}, true
			}
			return nil, false
		},
	}
}

// Get selected token(s)
//
// Retrieves the currently selected token(s) in the Foundry VTT client.
// @tag Utility
// @returns The selected token(s)
func selectedGet() helpers.APIRouteConfig {
	return helpers.APIRouteConfig{Type: "selected", OptionalParams: []helpers.ParamDef{clientIDParam(), userIDParam()}}
}

// Get players/users
//
// Retrieves a list of all users configured in the Foundry VTT world.
// Useful for discovering valid userId values for permission-scoped API calls.
// @tag Utility
// @returns List of users with their IDs, names, roles, and active status
func playersGet() helpers.APIRouteConfig {
	return helpers.APIRouteConfig{Type: "players", OptionalParams: []helpers.ParamDef{clientIDParam(), userIDParam()}}
}

// Execute JavaScript
//
// Executes a JavaScript script in the Foundry VTT client.
// @tag Utility
// @returns The result of the executed script
func executeJsPost(disp *service.Dispatcher, db *database.DB) helpers.APIRouteConfig {
	return helpers.APIRouteConfig{
		Type: "execute-js", Timeout: 20 * time.Second,
		OptionalParams: []helpers.ParamDef{clientIDParam(), {Name: "script", From: []helpers.ParamSource{helpers.FromBody}, Description: "JavaScript script to execute"}, userIDParam()},
		ValidateParams: func(params helpers.Params, r *http.Request) (map[string]interface{}, bool) {
			if s := params.GetString("script"); s != "" && !helpers.ValidateScript(s) {
				return map[string]interface{}{"error": "Script contains forbidden patterns"}, true
			}
			return nil, false
		},
		AfterForward: func(reqCtx *helpers.RequestContext, clientID string, params helpers.Params) {
			user, ok := reqCtx.GetUser()
			if !ok {
				return
			}
			snippet := params.GetString("script")
			if len(snippet) > 100 {
				snippet = snippet[:100] + "…"
			}
			var apiKeyID int64
			if reqCtx.ScopedKey != nil {
				apiKeyID = reqCtx.ScopedKey.ID
			}

			// Burst detection
			if alerts.Track(fmt.Sprintf("execjs:%d", user.ID), 50, 5*time.Minute, 30*time.Minute) {
				alerts.Fire(alerts.Event{
					Type:     alerts.TypeExecuteJsBurst,
					Severity: "warning",
					Message:  "execute-js called 50+ times in 5 minutes",
					Details:  map[string]interface{}{"userId": user.ID, "clientId": clientID},
				})
			}

			disp.Dispatch(service.NotificationContext{
				Event:       service.EventExecuteJs,
				UserID:      user.ID,
				ApiKeyID:    apiKeyID,
				ClientID:    clientID,
				Description: "```js\n" + snippet + "\n```",
			})

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			_ = db.ModuleEventLogStore().Create(ctx, &model.ModuleEventLog{
				UserID:      user.ID,
				ClientID:    clientID,
				EventType:   "execute-js",
				Description: snippet,
			})
		},
	}
}

// --- Scenes ---

// Get scene(s)
//
// Retrieves one or more scenes by ID, name, active status, viewed status, or all.
// @tag Scene
// @returns Scene data
func sceneGet() helpers.APIRouteConfig {
	return helpers.APIRouteConfig{
		Type: "get-scene",
		OptionalParams: []helpers.ParamDef{
			clientIDParam(),
			{Name: "sceneId", From: []helpers.ParamSource{helpers.FromQuery}, Description: "ID of a specific scene to retrieve"},
			{Name: "name", From: []helpers.ParamSource{helpers.FromQuery}, Description: "Name of the scene to retrieve"},
			{Name: "active", From: []helpers.ParamSource{helpers.FromQuery}, Type: helpers.TypeBoolean, Description: "Set to true to get the currently active scene"},
			{Name: "viewed", From: []helpers.ParamSource{helpers.FromQuery}, Type: helpers.TypeBoolean, Description: "Set to true to get the currently viewed scene"},
			{Name: "all", From: []helpers.ParamSource{helpers.FromQuery}, Type: helpers.TypeBoolean, Description: "Set to true to get all scenes"},
			userIDParam(),
		},
		ValidateParams: func(params helpers.Params, r *http.Request) (map[string]interface{}, bool) {
			if !params.Has("sceneId") && !params.Has("name") && !params.Has("active") && !params.Has("viewed") && !params.Has("all") {
				return map[string]interface{}{"error": "At least one of 'sceneId', 'name', 'active', 'viewed', or 'all' is required"}, true
			}
			return nil, false
		},
	}
}

// Create a new scene
//
// @tag Scene
// @returns Created scene data
func sceneCreate() helpers.APIRouteConfig {
	return helpers.APIRouteConfig{
		Type:           "create-scene",
		Timeout:        30 * time.Second,
		RequiredParams: []helpers.ParamDef{{Name: "data", From: []helpers.ParamSource{helpers.FromBody}, Type: helpers.TypeObject, Required: true, Description: "Scene data object (name, width, height, grid, etc.)"}},
		OptionalParams: []helpers.ParamDef{clientIDParam(), userIDParam()},
	}
}

// Update an existing scene
//
// @tag Scene
// @returns Updated scene data
func sceneUpdate() helpers.APIRouteConfig {
	return helpers.APIRouteConfig{
		Type: "update-scene", Timeout: 30 * time.Second,
		RequiredParams: []helpers.ParamDef{{Name: "data", From: []helpers.ParamSource{helpers.FromBody}, Type: helpers.TypeObject, Required: true, Description: "Object containing the scene fields to update"}},
		OptionalParams: []helpers.ParamDef{
			clientIDParam(),
			{Name: "sceneId", From: bqParam, Description: "ID of the scene to update"},
			{Name: "name", From: bqParam, Description: "Name of the scene to update (alternative to sceneId)"},
			{Name: "active", From: bqParam, Type: helpers.TypeBoolean, Description: "Set to true to target the active scene"},
			userIDParam(),
		},
		ValidateParams: func(params helpers.Params, r *http.Request) (map[string]interface{}, bool) {
			if !params.Has("sceneId") && !params.Has("name") && !params.Has("active") {
				return map[string]interface{}{"error": "At least one of 'sceneId', 'name', or 'active' is required"}, true
			}
			return nil, false
		},
	}
}

// Delete a scene
//
// @tag Scene
// @returns Deletion result
func sceneDelete() helpers.APIRouteConfig {
	return helpers.APIRouteConfig{
		Type: "delete-scene",
		OptionalParams: []helpers.ParamDef{
			clientIDParam(),
			{Name: "sceneId", From: []helpers.ParamSource{helpers.FromQuery}, Description: "ID of the scene to delete"},
			{Name: "name", From: []helpers.ParamSource{helpers.FromQuery}, Description: "Name of the scene to delete (alternative to sceneId)"},
			userIDParam(),
		},
		ValidateParams: func(params helpers.Params, r *http.Request) (map[string]interface{}, bool) {
			if !params.Has("sceneId") && !params.Has("name") {
				return map[string]interface{}{"error": "At least one of 'sceneId' or 'name' is required"}, true
			}
			return nil, false
		},
	}
}

// Switch the active scene
//
// @tag Scene
// @returns Result of the scene switch
func switchScene() helpers.APIRouteConfig {
	return helpers.APIRouteConfig{
		Type: "switch-scene", Timeout: 30 * time.Second,
		OptionalParams: []helpers.ParamDef{
			clientIDParam(),
			{Name: "sceneId", From: bqParam, Description: "ID of the scene to activate"},
			{Name: "name", From: bqParam, Description: "Name of the scene to activate (alternative to sceneId)"},
			userIDParam(),
		},
		ValidateParams: func(params helpers.Params, r *http.Request) (map[string]interface{}, bool) {
			if !params.Has("sceneId") && !params.Has("name") {
				return map[string]interface{}{"error": "At least one of 'sceneId' or 'name' is required"}, true
			}
			return nil, false
		},
	}
}

// --- Canvas ---
var validDocTypes = map[string]string{
	"tokens": "Token", "tiles": "Tile", "drawings": "Drawing",
	"lights": "AmbientLight", "sounds": "AmbientSound", "notes": "Note",
	"templates": "MeasuredTemplate", "walls": "Wall",
}

func validateDocType(params helpers.Params, r *http.Request) (map[string]interface{}, bool) {
	dt := params.GetString("documentType")
	if _, ok := validDocTypes[dt]; !ok {
		return map[string]interface{}{
			"error":      "Invalid documentType. Must be one of: tokens, tiles, drawings, lights, sounds, notes, templates, walls",
			"validTypes": []string{"tokens", "tiles", "drawings", "lights", "sounds", "notes", "templates", "walls"},
		}, true
	}
	return nil, false
}
func addClassName(params helpers.Params, r *http.Request) map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range params {
		if k != "clientId" && k != "type" {
			result[k] = v
		}
	}
	result["className"] = validDocTypes[params.GetString("documentType")]
	return result
}

var dtParam = helpers.ParamDef{Name: "documentType", From: []helpers.ParamSource{helpers.FromParams}, Type: helpers.TypeString, Required: true, Description: "Type of canvas document (tokens, tiles, drawings, lights, sounds, notes, templates, walls)"}

// Get canvas embedded documents
//
// @tag Canvas
// @returns Array of embedded documents
func canvasGet() helpers.APIRouteConfig {
	return helpers.APIRouteConfig{
		Type:           "get-canvas-documents",
		RequiredParams: []helpers.ParamDef{dtParam},
		OptionalParams: []helpers.ParamDef{
			clientIDParam(),
			{Name: "sceneId", From: []helpers.ParamSource{helpers.FromQuery}, Description: "Scene ID to query (defaults to the active scene)"},
			{Name: "documentId", From: []helpers.ParamSource{helpers.FromQuery}, Description: "Specific document ID to retrieve"},
			userIDParam(),
		},
		ValidateParams: validateDocType,
		BuildPayload:   addClassName,
	}
}

// Create canvas embedded document(s)
//
// @tag Canvas
// @returns Created document(s)
func canvasCreate() helpers.APIRouteConfig {
	return helpers.APIRouteConfig{
		Type:    "create-canvas-document",
		Timeout: 30 * time.Second,
		RequiredParams: []helpers.ParamDef{
			dtParam,
			{Name: "data", From: []helpers.ParamSource{helpers.FromBody}, Type: helpers.TypeObject, Required: true, Description: "Document data object or array of objects to create"},
		},
		OptionalParams: []helpers.ParamDef{
			clientIDParam(),
			{Name: "sceneId", From: bqParam, Description: "Scene ID to create in (defaults to the active scene)"},
			userIDParam(),
		},
		ValidateParams: validateDocType,
		BuildPayload:   addClassName,
	}
}

// Update a canvas embedded document
//
// @tag Canvas
// @returns Updated document
func canvasUpdate() helpers.APIRouteConfig {
	return helpers.APIRouteConfig{
		Type:    "update-canvas-document",
		Timeout: 30 * time.Second,
		RequiredParams: []helpers.ParamDef{
			dtParam,
			{Name: "documentId", From: bqParam, Type: helpers.TypeString, Required: true, Description: "ID of the document to update"},
			{Name: "data", From: []helpers.ParamSource{helpers.FromBody}, Type: helpers.TypeObject, Required: true, Description: "Object containing the fields to update"},
		},
		OptionalParams: []helpers.ParamDef{
			clientIDParam(),
			{Name: "sceneId", From: bqParam, Description: "Scene ID containing the document (defaults to the active scene)"},
			userIDParam(),
		},
		ValidateParams: validateDocType,
		BuildPayload:   addClassName,
	}
}

// Delete a canvas embedded document
//
// @tag Canvas
// @returns Deletion result
func canvasDelete() helpers.APIRouteConfig {
	return helpers.APIRouteConfig{
		Type: "delete-canvas-document",
		RequiredParams: []helpers.ParamDef{
			dtParam,
			{Name: "documentId", From: []helpers.ParamSource{helpers.FromQuery}, Type: helpers.TypeString, Required: true, Description: "ID of the document to delete"},
		},
		OptionalParams: []helpers.ParamDef{
			clientIDParam(),
			{Name: "sceneId", From: []helpers.ParamSource{helpers.FromQuery}, Description: "Scene ID containing the document (defaults to the active scene)"},
			userIDParam(),
		},
		ValidateParams: validateDocType,
		BuildPayload:   addClassName,
	}
}

// --- Chat ---

// Get chat messages
//
// Retrieves chat messages from the Foundry world with optional pagination and filtering.
// @tag Chat
// @returns Paginated list of chat messages
func chatGet() helpers.APIRouteConfig {
	return helpers.APIRouteConfig{Type: "chat-messages", OptionalParams: []helpers.ParamDef{
		clientIDParam(),
		{Name: "limit", From: []helpers.ParamSource{helpers.FromQuery}, Type: helpers.TypeNumber, Description: "Maximum number of messages to return (default: 10)"},
		{Name: "offset", From: []helpers.ParamSource{helpers.FromQuery}, Type: helpers.TypeNumber, Description: "Number of messages to skip for pagination"},
		userIDParam(),
		{Name: "chatType", From: []helpers.ParamSource{helpers.FromQuery}, Type: helpers.TypeNumber, Description: "Foundry chat message type (0=OOC, 1=IC, 2=Emote, 3=Whisper, 4=Roll). Named chatType to avoid collision with WS message type field."},
		{Name: "speaker", From: []helpers.ParamSource{helpers.FromQuery}, Description: "Filter messages by speaker name or actor ID"},
	}}
}

// Send a chat message
//
// Creates a new chat message in the Foundry world.
// @tag Chat
// @returns The created chat message
func chatSend() helpers.APIRouteConfig {
	return helpers.APIRouteConfig{Type: "chat-send",
		RequiredParams: []helpers.ParamDef{
			{Name: "content", From: []helpers.ParamSource{helpers.FromBody}, Type: helpers.TypeString, Required: true, Description: "The message content (supports HTML)"},
		},
		OptionalParams: []helpers.ParamDef{
			clientIDParam(),
			{Name: "whisper", From: []helpers.ParamSource{helpers.FromBody}, Type: helpers.TypeArray, Description: "Array of user IDs to whisper the message to"},
			{Name: "speaker", From: []helpers.ParamSource{helpers.FromBody}, Description: "Actor ID to use as the message speaker"},
			{Name: "alias", From: []helpers.ParamSource{helpers.FromBody}, Description: "Display name alias for the speaker"},
			{Name: "chatType", From: []helpers.ParamSource{helpers.FromBody}, Type: helpers.TypeNumber, Description: "Foundry chat message type (0=OOC, 1=IC, 2=Emote, 3=Whisper, 4=Roll). Named chatType to avoid collision with WS message type field."},
			{Name: "flavor", From: []helpers.ParamSource{helpers.FromBody}, Description: "Flavor text displayed above the message content"},
			userIDParam(),
		},
	}
}

// Delete a specific chat message
//
// Deletes a chat message by its ID. Only the message author or a GM can delete messages.
// @tag Chat
// @returns Success confirmation
func chatDeleteMsg() helpers.APIRouteConfig {
	return helpers.APIRouteConfig{
		Type:           "chat-delete",
		RequiredParams: []helpers.ParamDef{{Name: "messageId", From: []helpers.ParamSource{helpers.FromParams}, Type: helpers.TypeString, Required: true, Description: "ID of the chat message to delete"}},
		OptionalParams: []helpers.ParamDef{clientIDParam(), userIDParam()},
	}
}

// Clear all chat messages
//
// Flushes all chat message history. Only GMs can perform this action.
// @tag Chat
// @returns Success confirmation
func chatFlush() helpers.APIRouteConfig {
	return helpers.APIRouteConfig{Type: "chat-flush", OptionalParams: []helpers.ParamDef{clientIDParam(), userIDParam()}}
}

// --- Token Movement ---

// Move a token to specific coordinates
//
// Moves a token on the canvas to the specified x,y position, optionally
// animating through waypoints. Token can be identified by UUID or name.
// @tag Canvas
// @returns Result of the token movement including new position
func tokenMove() helpers.APIRouteConfig {
	return helpers.APIRouteConfig{
		Type: "move-token",
		RequiredParams: []helpers.ParamDef{
			{Name: "x", From: bqParam, Type: helpers.TypeNumber, Required: true, Description: "Target x coordinate"},
			{Name: "y", From: bqParam, Type: helpers.TypeNumber, Required: true, Description: "Target y coordinate"},
		},
		OptionalParams: []helpers.ParamDef{
			clientIDParam(),
			{Name: "uuid", From: bqParam, Type: helpers.TypeString, Description: "UUID of the token to move (optional if name provided)"},
			{Name: "name", From: bqParam, Type: helpers.TypeString, Description: "Name of the token to move (optional if uuid provided)"},
			{Name: "waypoints", From: bqParam, Type: helpers.TypeArray, Description: "Array of waypoint objects with x and y coordinates to animate through before reaching final position"},
			{Name: "animate", From: bqParam, Type: helpers.TypeBoolean, Description: "Whether to animate the movement (default: true)"},
			{Name: "sceneId", From: bqParam, Type: helpers.TypeString, Description: "Scene ID (defaults to active scene)"},
			userIDParam(),
		},
	}
}

// --- Measure Distance ---

// Measure the distance between two points or tokens
//
// Calculates the distance between two positions on the canvas, respecting
// the grid type and measurement rules. Points can be specified as coordinates
// or by referencing tokens by UUID or name.
// @tag Canvas
// @returns Distance measurement including units and grid spaces
func measureDistance() helpers.APIRouteConfig {
	return helpers.APIRouteConfig{
		Type: "measure-distance",
		OptionalParams: []helpers.ParamDef{
			clientIDParam(),
			{Name: "originX", From: bqParam, Type: helpers.TypeNumber, Description: "Origin x coordinate (optional if originUuid/originName provided)"},
			{Name: "originY", From: bqParam, Type: helpers.TypeNumber, Description: "Origin y coordinate"},
			{Name: "targetX", From: bqParam, Type: helpers.TypeNumber, Description: "Target x coordinate (optional if targetUuid/targetName provided)"},
			{Name: "targetY", From: bqParam, Type: helpers.TypeNumber, Description: "Target y coordinate"},
			{Name: "originUuid", From: bqParam, Type: helpers.TypeString, Description: "UUID of the origin token"},
			{Name: "originName", From: bqParam, Type: helpers.TypeString, Description: "Name of the origin token"},
			{Name: "targetUuid", From: bqParam, Type: helpers.TypeString, Description: "UUID of the target token"},
			{Name: "targetName", From: bqParam, Type: helpers.TypeString, Description: "Name of the target token"},
			{Name: "sceneId", From: bqParam, Type: helpers.TypeString, Description: "Scene ID (defaults to active scene)"},
			userIDParam(),
		},
	}
}

// --- Scene Image ---

// Get the raw background image of a scene
//
// Returns the scene's background image file without any tokens, lights,
// or other canvas elements rendered on it.
// @tag Scene
// @returns The raw scene background image
func sceneRawImage() helpers.APIRouteConfig {
	return helpers.APIRouteConfig{
		Type:    "scene-raw-image",
		Timeout: 30 * time.Second,
		OptionalParams: []helpers.ParamDef{
			clientIDParam(),
			{Name: "sceneId", From: bqParam, Type: helpers.TypeString, Description: "Scene ID (defaults to viewed/active scene)"},
			{Name: "active", From: bqParam, Type: helpers.TypeBoolean, Description: "If true, explicitly use the player-facing active scene instead of the viewed scene"},
			userIDParam(),
		},
		BuildPendingRequest: func(params helpers.Params) *ws.PendingRequest {
			return &ws.PendingRequest{Format: "binary"}
		},
	}
}

// --- World Info ---

// Get comprehensive world information
//
// Returns a single object with world name, game system, Foundry version,
// all modules (with active status), all users (with online status), and
// the active scene. Useful for API clients to discover the world state.
// @tag Utility
// @returns World information object
func worldInfoGet() helpers.APIRouteConfig {
	return helpers.APIRouteConfig{
		Type:           "world-info",
		OptionalParams: []helpers.ParamDef{clientIDParam(), userIDParam()},
	}
}

// --- Effects ---

// List all available status effects
//
// Returns all status effects defined by the game system's configuration.
// Useful for discovering valid statusId values for the add/remove effect endpoints.
// @tag Effects
// @returns Array of available status effects with id, name, and icon
func statusEffectsGet() helpers.APIRouteConfig {
	return helpers.APIRouteConfig{
		Type:           "get-status-effects",
		OptionalParams: []helpers.ParamDef{clientIDParam(), userIDParam()},
	}
}

// Get all active effects on an actor or token
//
// Returns the collection of ActiveEffect documents currently applied
// to the specified actor or token.
// @tag Effects
// @returns Array of active effects
func effectsGet() helpers.APIRouteConfig {
	return helpers.APIRouteConfig{
		Type:           "get-effects",
		RequiredParams: []helpers.ParamDef{{Name: "uuid", From: bqParam, Type: helpers.TypeString, Required: true, Description: "UUID of the actor or token to query"}},
		OptionalParams: []helpers.ParamDef{clientIDParam(), userIDParam()},
	}
}

// Add an active effect to an actor or token
//
// Adds a status condition (by statusId) or a custom ActiveEffect
// (via effectData) to the specified actor or token.
// @tag Effects
// @returns Result of the add operation
func effectsAdd() helpers.APIRouteConfig {
	return helpers.APIRouteConfig{
		Type:           "add-effect",
		RequiredParams: []helpers.ParamDef{{Name: "uuid", From: bqParam, Type: helpers.TypeString, Required: true, Description: "UUID of the actor or token to add the effect to"}},
		OptionalParams: []helpers.ParamDef{
			clientIDParam(),
			{Name: "statusId", From: bqParam, Description: "Standard status condition ID (e.g., \"poisoned\", \"blinded\", \"prone\")"},
			{Name: "effectData", From: bqParam, Type: helpers.TypeObject, Description: "Custom ActiveEffect data object (name, icon, duration, changes)"},
			userIDParam(),
		},
	}
}

// Remove an active effect from an actor or token
//
// Removes an effect by its document ID (effectId) or by status condition
// identifier (statusId).
// @tag Effects
// @returns Result of the remove operation
func effectsRemove() helpers.APIRouteConfig {
	return helpers.APIRouteConfig{
		Type:           "remove-effect",
		RequiredParams: []helpers.ParamDef{{Name: "uuid", From: bqParam, Type: helpers.TypeString, Required: true, Description: "UUID of the actor or token to remove the effect from"}},
		OptionalParams: []helpers.ParamDef{
			clientIDParam(),
			{Name: "effectId", From: bqParam, Description: "The ActiveEffect document ID to remove"},
			{Name: "statusId", From: bqParam, Description: "Standard status condition ID to remove (e.g., \"poisoned\")"},
			userIDParam(),
		},
	}
}

// --- File System ---

// Get file system structure
//
// @tag FileSystem
// @returns File system structure with files and directories
func fileSystemGet() helpers.APIRouteConfig {
	return helpers.APIRouteConfig{Type: "file-system", Timeout: 15 * time.Second, OptionalParams: []helpers.ParamDef{
		clientIDParam(),
		{Name: "path", From: []helpers.ParamSource{helpers.FromQuery}, Description: "The path to retrieve (relative to source)"},
		{Name: "source", From: []helpers.ParamSource{helpers.FromQuery}, Description: "The source directory to use (data, systems, modules, etc.)"},
		{Name: "recursive", From: []helpers.ParamSource{helpers.FromQuery}, Type: helpers.TypeBoolean, Description: "Whether to recursively list all subdirectories"},
		userIDParam(),
	}}
}

// Download a file from Foundry's file system
//
// @tag FileSystem
// @returns File contents in the requested format
func downloadFileGet() helpers.APIRouteConfig {
	return helpers.APIRouteConfig{
		Type: "download-file", Timeout: 45 * time.Second,
		OptionalParams: []helpers.ParamDef{
			clientIDParam(),
			{Name: "path", From: []helpers.ParamSource{helpers.FromQuery}, Description: "The full path to the file to download"},
			{Name: "source", From: []helpers.ParamSource{helpers.FromQuery}, Description: "The source directory to use (data, systems, modules, etc.)"},
			{Name: "format", From: []helpers.ParamSource{helpers.FromQuery}, Description: "The format to return the file in (binary, base64)"},
			userIDParam(),
		},
		BuildPendingRequest: func(params helpers.Params) *ws.PendingRequest {
			return &ws.PendingRequest{Format: params.GetString("format")}
		},
	}
}

// Get all connected clients for the authenticated API key
//
// Returns a list of all currently connected Foundry VTT clients associated with
// the provided API key, including their connection details and world information.
// @tag Clients
// @param {string} x-api-key [header] API key for authentication
// @returns Object containing total count and array of connected client details
func clientsHandler(mgr *ws.ClientManager, cfg *config.Config, db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqCtx := helpers.GetRequestContext(r)
		if reqCtx == nil {
			helpers.WriteError(w, http.StatusUnauthorized, "Authentication required")
			return
		}

		// Build a map of live clients keyed by clientID
		liveClients := mgr.GetConnectedClientInfos(reqCtx.MasterAPIKey)
		liveMap := make(map[string]ws.ClientInfo, len(liveClients))
		for _, c := range liveClients {
			liveMap[c.ID] = c
		}

		type clientResponse struct {
			ws.ClientInfo
			IsOnline bool `json:"isOnline"`
		}

		var results []clientResponse

		// Add all live clients with isOnline: true
		for _, c := range liveClients {
			results = append(results, clientResponse{ClientInfo: c, IsOnline: true})
		}

		// Merge known (offline) clients from the database
		if user, ok := reqCtx.User.(*model.User); ok && user != nil {
			known, err := db.KnownClientStore().FindAllByUser(r.Context(), user.ID)
			if err == nil {
				for _, kc := range known {
					if _, online := liveMap[kc.ClientID]; !online {
						results = append(results, clientResponse{
							ClientInfo: ws.ClientInfo{
								ID:             kc.ClientID,
								WorldID:        kc.WorldID.String,
								WorldTitle:     kc.WorldTitle.String,
								FoundryVersion: kc.FoundryVersion.String,
								SystemID:       kc.SystemID.String,
								SystemTitle:    kc.SystemTitle.String,
								SystemVersion:  kc.SystemVersion.String,
								CustomName:     kc.CustomName.String,
							},
							IsOnline: false,
						})
					}
				}
			}
		}

		// Apply scoped client filtering
		if reqCtx.ScopedKey != nil && reqCtx.ScopedKey.ScopedClientID != "" {
			var filtered []clientResponse
			for _, c := range results {
				if c.ID == reqCtx.ScopedKey.ScopedClientID {
					filtered = append(filtered, c)
				}
			}
			results = filtered
		}

		if results == nil {
			results = []clientResponse{}
		}
		helpers.WriteJSON(w, http.StatusOK, map[string]interface{}{"total": len(results), "clients": results})
	}
}
