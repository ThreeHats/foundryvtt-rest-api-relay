---
id: api-reference
title: API Reference
sidebar_label: API Reference
---

# FoundryVTT REST API Reference

This is the API reference for the FoundryVTT REST API Relay server.

## Available Endpoints

The API is divided into several resource groups:

- [Auth](/api/auth/) - Registration, login, API key and connection token management
- [Canvas](/api/canvas/) - Manipulate canvas embedded documents (tokens, walls, lights, etc.)
- [Chat](/api/chat/) - Send, retrieve, and subscribe to chat messages (includes SSE streaming)
- [Clients](/api/clients/) - List connected Foundry VTT worlds
- [DnD5e](/api/dnd5e/) - D&D 5th Edition specific endpoints
- [Effects](/api/effects/) - Read and manage active effects on actors and tokens
- [Encounter](/api/encounter/) - Combat encounter management
- [Entity](/api/entity/) - General entity manipulation (actors, items, etc.)
- [Events](/api/events/) - Subscribe to real-time Foundry hook events via SSE
- [File System](/api/fileSystem/) - File management operations
- [Macro](/api/macro/) - Execute and manage macros
- [Playlist](/api/playlist/) - Control playlists and audio tracks
- [Roll](/api/roll/) - Perform dice rolls and subscribe to roll events (includes SSE streaming)
- [Scene](/api/scene/) - Create, update, delete, and switch scenes
- [Search](/api/search/) - Search the Foundry database
- [Session](/api/session/) - Manage headless Foundry sessions
- [Sheet](/api/sheet/) - Interact with character sheets
- [Structure](/api/structure/) - Get information about world structure
- [User](/api/user/) - List and manage Foundry VTT users
- [Utility](/api/utility/) - Miscellaneous utility endpoints
- [WebSocket](/api/websocket/) - Bidirectional real-time communication via WebSocket

Each endpoint is documented with its path, HTTP method, required and optional parameters, and expected response format.

## Real-Time Streaming (SSE)

Some endpoints support **Server-Sent Events (SSE)** for real-time streaming:

- `GET /chat/subscribe` - Stream chat messages as they are created, updated, or deleted
- `GET /rolls/subscribe` - Stream dice roll events as they occur
- `GET /hooks/subscribe` - Stream Foundry hook events (combat, actor, scene, and more)

These endpoints return a persistent `text/event-stream` connection. See the individual endpoint documentation for code examples in JavaScript, TypeScript, Python, and cURL.

## Authentication

Most API endpoints require authentication with an API key. Provide your API key in the `x-api-key` header with each request.
