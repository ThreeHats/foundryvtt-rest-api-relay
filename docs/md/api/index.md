---
id: api-reference
title: API Reference
sidebar_label: API Reference
---

# FoundryVTT REST API Reference

This is the API reference for the FoundryVTT REST API Relay server.

## Available Endpoints

The API is divided into several resource groups:

- [Canvas](/api/canvas/) - Manipulate canvas embedded documents (tokens, walls, lights, etc.)
- [Chat](/api/chat/) - Send, retrieve, and subscribe to chat messages (includes SSE streaming)
- [DnD5e](/api/dnd5e/) - D&D 5th Edition specific endpoints
- [Encounter](/api/encounter/) - Combat encounter management
- [Entity](/api/entity/) - General entity manipulation (actors, items, etc.)
- [File System](/api/fileSystem/) - File management operations
- [Macro](/api/macro/) - Execute and manage macros
- [Roll](/api/roll/) - Perform dice rolls and subscribe to roll events (includes SSE streaming)
- [Scene](/api/scene/) - Create, update, delete, and switch scenes
- [Search](/api/search/) - Search the Foundry database
- [Session](/api/session/) - Manage headless Foundry sessions
- [Sheet](/api/sheet/) - Interact with character sheets
- [Structure](/api/structure/) - Get information about world structure
- [Utility](/api/utility/) - Miscellaneous utility endpoints

Each endpoint is documented with its path, HTTP method, required and optional parameters, and expected response format.

## Real-Time Streaming (SSE)

Some endpoints support **Server-Sent Events (SSE)** for real-time streaming:

- `GET /chat/subscribe` - Stream chat messages as they are created, updated, or deleted
- `GET /rolls/subscribe` - Stream dice roll events as they occur

These endpoints return a persistent `text/event-stream` connection. See the individual endpoint documentation for code examples in JavaScript, TypeScript, Python, and cURL.

## Authentication

Most API endpoints require authentication with an API key. Provide your API key in the `x-api-key` header with each request.
