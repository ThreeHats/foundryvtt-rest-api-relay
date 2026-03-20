---
id: api-reference
title: API Reference
sidebar_label: API Reference
---

# FoundryVTT REST API Reference

This is the API reference for the FoundryVTT REST API Relay server.

## Interactive API Testing

Each endpoint page includes a **"Try it out"** widget that lets you send real requests directly from the documentation. Configure your server URL and API key once, and they'll be saved across pages.

An <a href="/openapi.json">OpenAPI 3.0 spec</a> is available at `/openapi.json` for use with external tools like Swagger UI, Postman, or Insomnia.

## Available Endpoints

The API is divided into several resource groups:

- [Clients](/api/clients/) - List connected Foundry VTT clients
- [DnD5e](/api/dnd5e/) - D&D 5th Edition specific endpoints
- [Encounter](/api/encounter/) - Combat encounter management
- [Entity](/api/entity/) - General entity manipulation (actors, items, etc.)
- [File System](/api/fileSystem/) - File management operations
- [Macro](/api/macro/) - Execute and manage macros
- [Roll](/api/roll/) - Perform dice rolls
- [Search](/api/search/) - Search the Foundry database
- [Session](/api/session/) - Manage headless Foundry sessions
- [Sheet](/api/sheet/) - Interact with character sheets
- [Structure](/api/structure/) - Get information about world structure
- [Utility](/api/utility/) - Miscellaneous utility endpoints

Each endpoint is documented with its path, HTTP method, required and optional parameters, and expected response format.

## Authentication

Most API endpoints require authentication with an API key. Provide your API key in the `x-api-key` header with each request.
