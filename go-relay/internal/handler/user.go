package handler

import (
	"net/http"

	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/handler/helpers"
)

// List all Foundry users
//
// Retrieves a list of all users configured in the Foundry VTT world, including their roles and online status.
// This is a GM-only operation.
// @tag User
// @returns Array of user objects with id, name, role, isGM, active, color, avatar, and character fields
func usersGet() helpers.APIRouteConfig {
	return helpers.APIRouteConfig{
		Type:           "get-users",
		OptionalParams: []helpers.ParamDef{clientIDParam(), userIDParam()},
	}
}

// Get a single Foundry user
//
// Retrieves a single user by their ID or name.
// This is a GM-only operation.
// @tag User
// @returns User object with id, name, role, isGM, active, color, avatar, and character fields
func userGet() helpers.APIRouteConfig {
	return helpers.APIRouteConfig{
		Type: "get-user",
		OptionalParams: []helpers.ParamDef{
			clientIDParam(),
			{Name: "id", From: []helpers.ParamSource{helpers.FromQuery}, Type: helpers.TypeString, Description: "ID of the user to retrieve"},
			{Name: "name", From: []helpers.ParamSource{helpers.FromQuery}, Type: helpers.TypeString, Description: "Name of the user to retrieve (alternative to id)"},
			userIDParam(),
		},
		ValidateParams: func(params helpers.Params, r *http.Request) (map[string]interface{}, bool) {
			if !params.Has("id") && !params.Has("name") {
				return map[string]interface{}{"error": "At least one of 'id' or 'name' is required"}, true
			}
			return nil, false
		},
	}
}

// Create a new Foundry user
//
// Creates a new user in the Foundry VTT world with the specified name, role, and optional password.
// This is a GM-only operation.
// @tag User
// @returns The created user object
func userCreate() helpers.APIRouteConfig {
	return helpers.APIRouteConfig{
		Type: "create-user",
		RequiredParams: []helpers.ParamDef{
			{Name: "name", From: []helpers.ParamSource{helpers.FromBody}, Type: helpers.TypeString, Required: true, Description: "Username for the new user"},
		},
		OptionalParams: []helpers.ParamDef{
			clientIDParam(),
			{Name: "role", From: []helpers.ParamSource{helpers.FromBody}, Type: helpers.TypeNumber, Description: "User role: 0=None, 1=Player, 2=Trusted, 3=Assistant, 4=GM (default: 1)"},
			{Name: "password", From: []helpers.ParamSource{helpers.FromBody}, Type: helpers.TypeString, Description: "Password for the new user"},
			userIDParam(),
		},
	}
}

// Update an existing Foundry user
//
// Updates fields on an existing user. Identify the user by id or name, then pass the fields to update in the data object.
// Cannot demote the last GM user. This is a GM-only operation.
// @tag User
// @returns The updated user object
func userUpdate() helpers.APIRouteConfig {
	return helpers.APIRouteConfig{
		Type: "update-user",
		RequiredParams: []helpers.ParamDef{
			{Name: "data", From: []helpers.ParamSource{helpers.FromBody}, Type: helpers.TypeObject, Required: true, Description: "Object containing user fields to update (name, role, password, color, avatar, character)"},
		},
		OptionalParams: []helpers.ParamDef{
			clientIDParam(),
			{Name: "id", From: bqParam, Type: helpers.TypeString, Description: "ID of the user to update"},
			{Name: "name", From: bqParam, Type: helpers.TypeString, Description: "Name of the user to update (alternative to id)"},
			userIDParam(),
		},
		ValidateParams: func(params helpers.Params, r *http.Request) (map[string]interface{}, bool) {
			if !params.Has("id") && !params.Has("name") {
				return map[string]interface{}{"error": "At least one of 'id' or 'name' is required"}, true
			}
			return nil, false
		},
	}
}

// Delete a Foundry user
//
// Permanently deletes a user from the Foundry VTT world.
// Cannot delete yourself or the last GM user. This is a GM-only operation.
// @tag User
// @returns Deletion result
func userDelete() helpers.APIRouteConfig {
	return helpers.APIRouteConfig{
		Type: "delete-user",
		OptionalParams: []helpers.ParamDef{
			clientIDParam(),
			{Name: "id", From: []helpers.ParamSource{helpers.FromQuery}, Type: helpers.TypeString, Description: "ID of the user to delete"},
			{Name: "name", From: []helpers.ParamSource{helpers.FromQuery}, Type: helpers.TypeString, Description: "Name of the user to delete (alternative to id)"},
			userIDParam(),
		},
		ValidateParams: func(params helpers.Params, r *http.Request) (map[string]interface{}, bool) {
			if !params.Has("id") && !params.Has("name") {
				return map[string]interface{}{"error": "At least one of 'id' or 'name' is required"}, true
			}
			return nil, false
		},
	}
}
