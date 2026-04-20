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

	// --- Concentration Tracking ---

	actorNameParam := helpers.ParamDef{Name: "actorName", From: bq, Type: helpers.TypeString, Description: "Name of the actor (optional if actorUuid provided)"}

	// Check if an actor is concentrating on a spell
	//
	// Returns whether the actor currently has a concentration effect active,
	// and if so, what spell they are concentrating on.
	// @tag Dnd5e
	// @returns Concentration status with effect details and spell name
	r.Get("/concentration", helpers.CreateAPIRoute(mgr, pending, helpers.APIRouteConfig{
		Type:           "get-concentration",
		OptionalParams: []helpers.ParamDef{clientIDParam(), optActorUuid, actorNameParam, userIDParam()},
	}))

	// Break an actor's concentration
	//
	// Removes the concentration effect from the actor, ending any
	// spell that requires concentration.
	// @tag Dnd5e
	// @returns Confirmation that concentration was broken
	r.Post("/break-concentration", helpers.CreateAPIRoute(mgr, pending, helpers.APIRouteConfig{
		Type:           "break-concentration",
		OptionalParams: []helpers.ParamDef{clientIDParam(), optActorUuid, actorNameParam, userIDParam()},
	}))

	// Roll a concentration saving throw
	//
	// Rolls a Constitution saving throw to maintain concentration after taking damage.
	// The DC is calculated as max(10, floor(damage/2)). Returns the roll result
	// and whether concentration was maintained or broken.
	// @tag Dnd5e
	// @returns Roll result and concentration maintained status
	r.Post("/concentration-save", helpers.CreateAPIRoute(mgr, pending, helpers.APIRouteConfig{
		Type: "concentration-save",
		RequiredParams: []helpers.ParamDef{
			{Name: "damage", From: bq, Type: helpers.TypeNumber, Required: true, Description: "Amount of damage taken (used to calculate DC = max(10, floor(damage/2)))"},
		},
		OptionalParams: append([]helpers.ParamDef{clientIDParam(), optActorUuid, actorNameParam, userIDParam()}, rollOptions...),
	}))

	// --- Inventory Management ---

	itemUuidParam := helpers.ParamDef{Name: "itemUuid", From: bq, Type: helpers.TypeString, Description: "UUID of the item (optional if itemName provided)"}
	itemNameParam := helpers.ParamDef{Name: "itemName", From: bq, Type: helpers.TypeString, Description: "Name of the item (optional if itemUuid provided)"}

	// Equip or unequip an item
	//
	// Changes the equipped status of an item in an actor's inventory.
	// @tag Dnd5e
	// @returns Updated equipment status
	r.Post("/equip-item", helpers.CreateAPIRoute(mgr, pending, helpers.APIRouteConfig{
		Type: "equip-item",
		RequiredParams: []helpers.ParamDef{
			{Name: "equipped", From: bq, Type: helpers.TypeBoolean, Required: true, Description: "Whether the item should be equipped (true) or unequipped (false)"},
		},
		OptionalParams: []helpers.ParamDef{clientIDParam(), optActorUuid, actorNameParam, itemUuidParam, itemNameParam, userIDParam()},
	}))

	// Attune or unattune an item
	//
	// Changes the attunement status of a magic item in an actor's inventory.
	// @tag Dnd5e
	// @returns Updated attunement status
	r.Post("/attune-item", helpers.CreateAPIRoute(mgr, pending, helpers.APIRouteConfig{
		Type: "attune-item",
		RequiredParams: []helpers.ParamDef{
			{Name: "attuned", From: bq, Type: helpers.TypeBoolean, Required: true, Description: "Whether the item should be attuned (true) or unattuned (false)"},
		},
		OptionalParams: []helpers.ParamDef{clientIDParam(), optActorUuid, actorNameParam, itemUuidParam, itemNameParam, userIDParam()},
	}))

	// Transfer currency between actors
	//
	// Moves currency from one actor to another. Validates that the source
	// actor has sufficient funds before transferring.
	// @tag Dnd5e
	// @returns Transfer result with updated balances
	r.Post("/transfer-currency", helpers.CreateAPIRoute(mgr, pending, helpers.APIRouteConfig{
		Type: "transfer-currency",
		RequiredParams: []helpers.ParamDef{
			{Name: "currency", From: bq, Type: helpers.TypeObject, Required: true, Description: "Currency amounts to transfer, e.g. pp, gp, ep, sp, cp denomination keys with numeric values"},
		},
		OptionalParams: []helpers.ParamDef{
			clientIDParam(),
			{Name: "sourceActorUuid", From: bq, Type: helpers.TypeString, Description: "UUID of the source actor (optional if sourceActorName provided)"},
			{Name: "sourceActorName", From: bq, Type: helpers.TypeString, Description: "Name of the source actor"},
			{Name: "targetActorUuid", From: bq, Type: helpers.TypeString, Description: "UUID of the target actor (optional if targetActorName provided)"},
			{Name: "targetActorName", From: bq, Type: helpers.TypeString, Description: "Name of the target actor"},
			userIDParam(),
		},
	}))

	return r
}
