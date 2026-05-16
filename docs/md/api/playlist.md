---
tag: playlist
---
import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';


import ApiTester from '@site/src/components/ApiTester';

# Playlist

## GET /playlists

Get all playlists

Returns all playlists in the world with their tracks/sounds, including playing status, mode, and volume information.

**Required scope:** `playlist:control`

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| clientId | string |  | query | Client ID for the Foundry world |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**array** - Array of playlists with their sounds

### Try It Out

<ApiTester
  method="GET"
  path="/playlists"
  parameters={[{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

---

## POST /playlist/play

Play a playlist or specific sound

Starts playback of an entire playlist or a specific sound within it. The playlist can be identified by ID or name. Optionally specify a specific sound/track to play within the playlist.

**Required scope:** `playlist:control`

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| clientId | string |  | query | Client ID for the Foundry world |
| playlistId | string |  | body, query | ID of the playlist (optional if playlistName provided) |
| playlistName | string |  | body, query | Name of the playlist (optional if playlistId provided) |
| soundId | string |  | body, query | ID of a specific sound to play within the playlist |
| soundName | string |  | body, query | Name of a specific sound to play (optional if soundId provided) |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**object** - Playback status confirmation

### Try It Out

<ApiTester
  method="POST"
  path="/playlist/play"
  parameters={[{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"playlistId","type":"string","required":false,"source":"body"},{"name":"playlistName","type":"string","required":false,"source":"body"},{"name":"soundId","type":"string","required":false,"source":"body"},{"name":"soundName","type":"string","required":false,"source":"body"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

---

## POST /playlist/stop

Stop a playlist

Stops playback of the specified playlist.

**Required scope:** `playlist:control`

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| clientId | string |  | query | Client ID for the Foundry world |
| playlistId | string |  | body, query | ID of the playlist (optional if playlistName provided) |
| playlistName | string |  | body, query | Name of the playlist (optional if playlistId provided) |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**object** - Stop confirmation

### Try It Out

<ApiTester
  method="POST"
  path="/playlist/stop"
  parameters={[{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"playlistId","type":"string","required":false,"source":"body"},{"name":"playlistName","type":"string","required":false,"source":"body"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

---

## POST /playlist/next

Skip to next track in a playlist

Advances to the next sound/track in the specified playlist.

**Required scope:** `playlist:control`

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| clientId | string |  | query | Client ID for the Foundry world |
| playlistId | string |  | body, query | ID of the playlist (optional if playlistName provided) |
| playlistName | string |  | body, query | Name of the playlist (optional if playlistId provided) |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**object** - Next track information

### Try It Out

<ApiTester
  method="POST"
  path="/playlist/next"
  parameters={[{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"playlistId","type":"string","required":false,"source":"body"},{"name":"playlistName","type":"string","required":false,"source":"body"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

---

## POST /playlist/volume

Set volume for a playlist or specific sound

Adjusts the volume of an entire playlist or a specific sound within it. Volume is specified as a float between 0 (silent) and 1 (full volume).

**Required scope:** `playlist:control`

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| volume | number | ✓ | body, query | Volume level from 0.0 (silent) to 1.0 (full volume) |
| clientId | string |  | query | Client ID for the Foundry world |
| playlistId | string |  | body, query | ID of the playlist (optional if playlistName provided) |
| playlistName | string |  | body, query | Name of the playlist (optional if playlistId provided) |
| soundId | string |  | body, query | ID of a specific sound to adjust volume for |
| soundName | string |  | body, query | Name of a specific sound to adjust volume for (optional if soundId provided) |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**object** - Updated volume level

### Try It Out

<ApiTester
  method="POST"
  path="/playlist/volume"
  parameters={[{"name":"volume","type":"number","required":true,"source":"body"},{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"playlistId","type":"string","required":false,"source":"body"},{"name":"playlistName","type":"string","required":false,"source":"body"},{"name":"soundId","type":"string","required":false,"source":"body"},{"name":"soundName","type":"string","required":false,"source":"body"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

---

## POST /stop-sound

Play a one-shot sound effect

Triggers playback of an audio file by its path. Useful for sound effects, ambient sounds, or any audio that should play once without being part of a playlist. Stop a playing sound Stops playback of a currently playing sound by its source path. If no src is provided, stops all currently playing sounds.

**Required scope:** `playlist:control`

### Parameters

| Name | Type | Required | Source | Description |
|------|------|----------|--------|-------------|
| clientId | string |  | query | Client ID for the Foundry world |
| src | string |  | body, query | Path to the audio file to stop (omit to stop all sounds) |
| userId | string |  | query, body | Foundry user ID or username to scope permissions (omit for GM-level access) |

### Returns

**object** - Stop confirmation

### Try It Out

<ApiTester
  method="POST"
  path="/stop-sound"
  parameters={[{"name":"clientId","type":"string","required":false,"source":"query"},{"name":"src","type":"string","required":false,"source":"body"},{"name":"userId","type":"string","required":false,"source":"query"}]}
/>

