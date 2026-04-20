package handler

import (
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/handler/helpers"
)

// --- Playlist Control ---

// Get all playlists
//
// Returns all playlists in the world with their tracks/sounds,
// including playing status, mode, and volume information.
// @tag Playlist
// @returns Array of playlists with their sounds
func playlistsGet() helpers.APIRouteConfig {
	return helpers.APIRouteConfig{
		Type:           "get-playlists",
		OptionalParams: []helpers.ParamDef{clientIDParam(), userIDParam()},
	}
}

// Play a playlist or specific sound
//
// Starts playback of an entire playlist or a specific sound within it.
// The playlist can be identified by ID or name. Optionally specify a
// specific sound/track to play within the playlist.
// @tag Playlist
// @returns Playback status confirmation
func playlistPlay() helpers.APIRouteConfig {
	bq := []helpers.ParamSource{helpers.FromBody, helpers.FromQuery}
	return helpers.APIRouteConfig{
		Type: "playlist-play",
		OptionalParams: []helpers.ParamDef{
			clientIDParam(),
			{Name: "playlistId", From: bq, Type: helpers.TypeString, Description: "ID of the playlist (optional if playlistName provided)"},
			{Name: "playlistName", From: bq, Type: helpers.TypeString, Description: "Name of the playlist (optional if playlistId provided)"},
			{Name: "soundId", From: bq, Type: helpers.TypeString, Description: "ID of a specific sound to play within the playlist"},
			{Name: "soundName", From: bq, Type: helpers.TypeString, Description: "Name of a specific sound to play (optional if soundId provided)"},
			userIDParam(),
		},
	}
}

// Stop a playlist
//
// Stops playback of the specified playlist.
// @tag Playlist
// @returns Stop confirmation
func playlistStop() helpers.APIRouteConfig {
	bq := []helpers.ParamSource{helpers.FromBody, helpers.FromQuery}
	return helpers.APIRouteConfig{
		Type: "playlist-stop",
		OptionalParams: []helpers.ParamDef{
			clientIDParam(),
			{Name: "playlistId", From: bq, Type: helpers.TypeString, Description: "ID of the playlist (optional if playlistName provided)"},
			{Name: "playlistName", From: bq, Type: helpers.TypeString, Description: "Name of the playlist (optional if playlistId provided)"},
			userIDParam(),
		},
	}
}

// Skip to next track in a playlist
//
// Advances to the next sound/track in the specified playlist.
// @tag Playlist
// @returns Next track information
func playlistNext() helpers.APIRouteConfig {
	bq := []helpers.ParamSource{helpers.FromBody, helpers.FromQuery}
	return helpers.APIRouteConfig{
		Type: "playlist-next",
		OptionalParams: []helpers.ParamDef{
			clientIDParam(),
			{Name: "playlistId", From: bq, Type: helpers.TypeString, Description: "ID of the playlist (optional if playlistName provided)"},
			{Name: "playlistName", From: bq, Type: helpers.TypeString, Description: "Name of the playlist (optional if playlistId provided)"},
			userIDParam(),
		},
	}
}

// Set volume for a playlist or specific sound
//
// Adjusts the volume of an entire playlist or a specific sound within it.
// Volume is specified as a float between 0 (silent) and 1 (full volume).
// @tag Playlist
// @returns Updated volume level
func playlistVolume() helpers.APIRouteConfig {
	bq := []helpers.ParamSource{helpers.FromBody, helpers.FromQuery}
	return helpers.APIRouteConfig{
		Type: "playlist-volume",
		RequiredParams: []helpers.ParamDef{
			{Name: "volume", From: bq, Type: helpers.TypeNumber, Required: true, Description: "Volume level from 0.0 (silent) to 1.0 (full volume)"},
		},
		OptionalParams: []helpers.ParamDef{
			clientIDParam(),
			{Name: "playlistId", From: bq, Type: helpers.TypeString, Description: "ID of the playlist (optional if playlistName provided)"},
			{Name: "playlistName", From: bq, Type: helpers.TypeString, Description: "Name of the playlist (optional if playlistId provided)"},
			{Name: "soundId", From: bq, Type: helpers.TypeString, Description: "ID of a specific sound to adjust volume for"},
			{Name: "soundName", From: bq, Type: helpers.TypeString, Description: "Name of a specific sound to adjust volume for (optional if soundId provided)"},
			userIDParam(),
		},
	}
}

// Play a one-shot sound effect
//
// Triggers playback of an audio file by its path. Useful for sound effects,
// ambient sounds, or any audio that should play once without being part of a playlist.
// @tag Playlist
// @returns Playback confirmation
// Stop a playing sound
//
// Stops playback of a currently playing sound by its source path.
// If no src is provided, stops all currently playing sounds.
// @tag Playlist
// @returns Stop confirmation
func soundStop() helpers.APIRouteConfig {
	bq := []helpers.ParamSource{helpers.FromBody, helpers.FromQuery}
	return helpers.APIRouteConfig{
		Type: "stop-sound",
		OptionalParams: []helpers.ParamDef{
			clientIDParam(),
			{Name: "src", From: bq, Type: helpers.TypeString, Description: "Path to the audio file to stop (omit to stop all sounds)"},
			userIDParam(),
		},
	}
}

func soundPlay() helpers.APIRouteConfig {
	bq := []helpers.ParamSource{helpers.FromBody, helpers.FromQuery}
	return helpers.APIRouteConfig{
		Type: "play-sound",
		RequiredParams: []helpers.ParamDef{
			{Name: "src", From: bq, Type: helpers.TypeString, Required: true, Description: "Path to the audio file (e.g., \"sounds/effect.mp3\")"},
		},
		OptionalParams: []helpers.ParamDef{
			clientIDParam(),
			{Name: "volume", From: bq, Type: helpers.TypeNumber, Description: "Volume from 0.0 to 1.0 (default: 0.5)"},
			{Name: "loop", From: bq, Type: helpers.TypeBoolean, Description: "Whether to loop the sound (default: false)"},
			userIDParam(),
		},
	}
}
