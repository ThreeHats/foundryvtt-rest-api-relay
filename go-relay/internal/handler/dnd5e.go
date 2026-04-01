package handler

import (
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/handler/helpers"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/ws"
	"github.com/go-chi/chi/v5"
)

// Dnd5eRouter creates D&D 5e system-specific routes.
func Dnd5eRouter(mgr *ws.ClientManager, pending *ws.PendingRequests) chi.Router {
	r := chi.NewRouter()

	bq := []helpers.ParamSource{helpers.FromBody, helpers.FromQuery}

	actorUuid := helpers.ParamDef{Name: "actorUuid", From: bq, Type: helpers.TypeString, Required: true, Description: "UUID of the actor"}
	optActorUuid := helpers.ParamDef{Name: "actorUuid", From: bq, Type: helpers.TypeString, Description: "UUID of the actor (optional if selected is true)"}
	abilityParams := []helpers.ParamDef{
		clientIDParam(),
		{Name: "abilityUuid", From: bq, Type: helpers.TypeString, Description: "The UUID of the specific ability (optional if abilityName provided)"},
		{Name: "abilityName", From: bq, Type: helpers.TypeString, Description: "The name of the ability if UUID not provided (optional if abilityUuid provided)"},
		{Name: "targetUuid", From: bq, Type: helpers.TypeString, Description: "The UUID of the target for the ability (optional)"},
		{Name: "targetName", From: bq, Type: helpers.TypeString, Description: "The name of the target if UUID not provided (optional)"},
		userIDParam(),
	}
	rollOptions := []helpers.ParamDef{
		{Name: "advantage", From: bq, Type: helpers.TypeBoolean, Description: "Roll with advantage"},
		{Name: "disadvantage", From: bq, Type: helpers.TypeBoolean, Description: "Roll with disadvantage"},
		{Name: "bonus", From: bq, Type: helpers.TypeString, Description: "Extra bonus formula to add (e.g., \"1d4\", \"+2\")"},
		{Name: "createChatMessage", From: bq, Type: helpers.TypeBoolean, Description: "Whether to post the roll to chat (default: true)"},
	}

	// Get detailed information for a specific D&D 5e actor
	//
	// Retrieves comprehensive details about an actor including stats, inventory,
	// spells, features, and other character information based on the requested details array.
	// @tag Dnd5e
	// @returns Actor details object containing requested information
	r.Get("/get-actor-details", helpers.CreateAPIRoute(mgr, pending, helpers.APIRouteConfig{
		Type: "get-actor-details",
		RequiredParams: []helpers.ParamDef{
			actorUuid,
			{Name: "details", From: bq, Type: helpers.TypeArray, Required: true, Description: "Array of detail types to retrieve (e.g., [\"resources\", \"items\", \"spells\", \"features\"])"},
		},
		OptionalParams: []helpers.ParamDef{clientIDParam(), userIDParam()},
	}))

	// Modify the charges for a specific item owned by an actor
	//
	// Increases or decreases the charges/uses of an item in an actor's inventory.
	// Useful for consumable items like potions, scrolls, or charged magic items.
	// @tag Dnd5e
	// @returns Result of the charge modification operation
	r.Post("/modify-item-charges", helpers.CreateAPIRoute(mgr, pending, helpers.APIRouteConfig{
		Type: "modify-item-charges",
		RequiredParams: []helpers.ParamDef{
			actorUuid,
			{Name: "amount", From: bq, Type: helpers.TypeNumber, Required: true, Description: "The amount to modify charges by (positive or negative)"},
		},
		OptionalParams: []helpers.ParamDef{
			clientIDParam(),
			{Name: "itemUuid", From: bq, Type: helpers.TypeString, Description: "The UUID of the specific item (optional if itemName provided)"},
			{Name: "itemName", From: bq, Type: helpers.TypeString, Description: "The name of the item if UUID not provided (optional if itemUuid provided)"},
			userIDParam(),
		},
	}))

	// use-ability, use-feature, use-spell, use-item all share the same shape
	for _, route := range []struct {
		path    string
		reqType string
	}{
		{"/use-ability", "use-ability"},
		{"/use-feature", "use-feature"},
		{"/use-spell", "use-spell"},
		{"/use-item", "use-item"},
	} {
		rt := route
		r.Post(rt.path, helpers.CreateAPIRoute(mgr, pending, helpers.APIRouteConfig{
			Type:           rt.reqType,
			RequiredParams: []helpers.ParamDef{actorUuid},
			OptionalParams: abilityParams,
		}))
	}

	// Perform a short rest for an actor
	//
	// Triggers the D&D 5e short rest workflow including hit dice recovery,
	// class feature resets, and HP recovery.
	// @tag Dnd5e
	// @returns Result of the short rest operation
	r.Post("/short-rest", helpers.CreateAPIRoute(mgr, pending, helpers.APIRouteConfig{
		Type: "short-rest",
		OptionalParams: []helpers.ParamDef{
			clientIDParam(), optActorUuid, selectedParam(),
			{Name: "autoHD", From: bq, Type: helpers.TypeBoolean, Description: "Automatically spend hit dice during short rest"},
			{Name: "autoHDThreshold", From: bq, Type: helpers.TypeNumber, Description: "HP threshold below which to auto-spend hit dice (0-1 as fraction of max HP)"},
			userIDParam(),
		},
	}))

	// Perform a long rest for an actor
	//
	// Triggers the D&D 5e long rest workflow including full HP recovery,
	// spell slot restoration, hit dice recovery, and feature resets.
	// @tag Dnd5e
	// @returns Result of the long rest operation
	r.Post("/long-rest", helpers.CreateAPIRoute(mgr, pending, helpers.APIRouteConfig{
		Type: "long-rest",
		OptionalParams: []helpers.ParamDef{
			clientIDParam(), optActorUuid, selectedParam(),
			{Name: "newDay", From: bq, Type: helpers.TypeBoolean, Description: "Whether the long rest marks a new day (default: true)"},
			userIDParam(),
		},
	}))

	// Roll a skill check for an actor
	//
	// Rolls a D&D 5e skill check with all applicable modifiers including
	// proficiency, expertise, Jack of All Trades, and conditional bonuses.
	// @tag Dnd5e
	// @returns Result of the skill check roll
	r.Post("/skill-check", helpers.CreateAPIRoute(mgr, pending, helpers.APIRouteConfig{
		Type: "skill-check",
		RequiredParams: []helpers.ParamDef{
			actorUuid,
			{Name: "skill", From: bq, Type: helpers.TypeString, Required: true, Description: "Skill abbreviation (e.g., \"acr\", \"ath\", \"ste\", \"prc\")"},
		},
		OptionalParams: append([]helpers.ParamDef{clientIDParam(), userIDParam()}, rollOptions...),
	}))

	// Roll an ability saving throw for an actor
	//
	// Rolls a D&D 5e ability saving throw with all applicable modifiers.
	// @tag Dnd5e
	// @returns Result of the saving throw roll
	r.Post("/ability-save", helpers.CreateAPIRoute(mgr, pending, helpers.APIRouteConfig{
		Type: "ability-save",
		RequiredParams: []helpers.ParamDef{
			actorUuid,
			{Name: "ability", From: bq, Type: helpers.TypeString, Required: true, Description: "Ability abbreviation (e.g., \"str\", \"dex\", \"con\", \"int\", \"wis\", \"cha\")"},
		},
		OptionalParams: append([]helpers.ParamDef{clientIDParam(), userIDParam()}, rollOptions...),
	}))

	// Roll an ability check for an actor
	//
	// Rolls a D&D 5e ability check (raw ability test, not a skill check)
	// with all applicable modifiers.
	// @tag Dnd5e
	// @returns Result of the ability check roll
	r.Post("/ability-check", helpers.CreateAPIRoute(mgr, pending, helpers.APIRouteConfig{
		Type: "ability-check",
		RequiredParams: []helpers.ParamDef{
			actorUuid,
			{Name: "ability", From: bq, Type: helpers.TypeString, Required: true, Description: "Ability abbreviation (e.g., \"str\", \"dex\", \"con\", \"int\", \"wis\", \"cha\")"},
		},
		OptionalParams: append([]helpers.ParamDef{clientIDParam(), userIDParam()}, rollOptions...),
	}))

	// Roll a death saving throw for an actor
	//
	// Rolls a D&D 5e death saving throw, handling DC 10 CON save,
	// three successes/failures tracking, nat 20 healing, and nat 1 double failure.
	// @tag Dnd5e
	// @returns Result of the death saving throw
	r.Post("/death-save", helpers.CreateAPIRoute(mgr, pending, helpers.APIRouteConfig{
		Type:           "death-save",
		RequiredParams: []helpers.ParamDef{actorUuid},
		OptionalParams: []helpers.ParamDef{
			clientIDParam(),
			{Name: "advantage", From: bq, Type: helpers.TypeBoolean, Description: "Roll with advantage"},
			{Name: "createChatMessage", From: bq, Type: helpers.TypeBoolean, Description: "Whether to post the roll to chat (default: true)"},
			userIDParam(),
		},
	}))

	// Modify the experience points for a specific actor
	//
	// Adds or removes experience points from an actor.
	// @tag Dnd5e
	// @returns Result of the experience modification operation
	r.Post("/modify-experience", helpers.CreateAPIRoute(mgr, pending, helpers.APIRouteConfig{
		Type: "modify-experience",
		RequiredParams: []helpers.ParamDef{
			{Name: "amount", From: bq, Type: helpers.TypeNumber, Required: true, Description: "The amount of experience to add (can be negative)"},
		},
		OptionalParams: []helpers.ParamDef{
			clientIDParam(), optActorUuid, selectedParam(), userIDParam(),
		},
	}))

	return r
}
