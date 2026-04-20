package ws

import (
	"github.com/rs/zerolog/log"
)

// ModuleNotifyEvent is the parsed module-notify event.
type ModuleNotifyEvent struct {
	ClientID  string
	APIKey    string
	EventName string // "settings-change" | "execute-js" | "macro-execute"
	Details   string // Free-form details (e.g., script content for execute-js)
	Actor     string // GM/user that triggered the event
	World     string // World title
}

// ModuleNotifyHandler is the callback the server provides to handle these events.
// We use a callback to avoid an import cycle (this package can't import database).
type ModuleNotifyHandler func(event ModuleNotifyEvent)

// RegisterModuleNotifyHandler registers the WebSocket handler for the
// `module-notify` message type. The Foundry module sends these messages when
// in-Foundry events occur (settings changes, execute-js, macro-execute, etc.)
// — events the relay cannot observe directly.
//
// Message format from module:
//
//	{
//	  "type": "module-notify",
//	  "event": "settings-change" | "execute-js" | "macro-execute",
//	  "details": "...",   // free-form description
//	  "actor": "GM Name", // who triggered the event (optional)
//	  "world": "World Title" // optional
//	}
//
// The handler callback (set by server.go) looks up the user by API key and
// dispatches via the unified notification system.
func RegisterModuleNotifyHandler(manager *ClientManager, handler ModuleNotifyHandler) {
	manager.OnMessageType("module-notify", func(client *Client, data map[string]interface{}) {
		eventName, _ := data["event"].(string)
		details, _ := data["details"].(string)
		actor, _ := data["actor"].(string)
		world, _ := data["world"].(string)

		if eventName == "" {
			log.Warn().Str("clientId", client.ID()).Msg("module-notify missing event field")
			return
		}

		if handler == nil {
			return
		}

		handler(ModuleNotifyEvent{
			ClientID:  client.ID(),
			APIKey:    client.APIKey(),
			EventName: eventName,
			Details:   details,
			Actor:     actor,
			World:     world,
		})
	})

	log.Info().Msg("Registered module-notify handler")
}
