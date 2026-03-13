---
id: swagger-ui-testing
title: Local Swagger UI Testing
sidebar_position: 7
---

# Local Swagger UI Testing

For manual API testing on a Windows machine, the recommended option is to run a local Swagger UI instance with the OpenAPI files from this repository.

This gives you a browser-based interface for:

- browsing the `core` and `dnd5e` API definitions
- trying requests against your local relay server
- validating request parameters and payloads against the current OpenAPI source files

## Why Swagger UI Is Recommended

Swagger UI is the primary manual testing workflow for this project because it reads the OpenAPI YAML files directly from `documentation/openapi`.

That means you are testing against the same source used to generate the API reference documentation.

## Prerequisites

Before starting Swagger UI, make sure you have:

1. Docker Desktop installed on Windows
2. A running relay server
3. A running Foundry VTT world connected to the relay
4. Your API key available

## Option 1: Use the Included Docker Compose Setup

The repository `docker-compose.yml` already contains a `swagger-ui` service.

### Step 1: Configure the OpenAPI file URLs

The Swagger UI container reads three environment variables:

- `SWAGGER_UI_CORE_OPENAPI_URL`
- `SWAGGER_UI_DND5E_OPENAPI_URL`
- `SWAGGER_UI_RELAY_SERVER_URL`

Example values are already included in `.env.test.example`:

```env
SWAGGER_UI_RELAY_SERVER_URL=http://localhost:3010
SWAGGER_UI_CORE_OPENAPI_URL=/swagger/generated-openapi-core.yaml
SWAGGER_UI_DND5E_OPENAPI_URL=/swagger/generated-openapi-dnd5e.yaml
```

`SWAGGER_UI_RELAY_SERVER_URL` is used to rewrite the OpenAPI `servers.url` value before Swagger UI serves the specs, so `Try it out` requests go to your local relay by default.

You can use them in either of these ways:

### Option A: Copy them into `.env`

Create a local `.env` file in the project root and include:

```env
SWAGGER_UI_RELAY_SERVER_URL=http://localhost:3010
SWAGGER_UI_CORE_OPENAPI_URL=/swagger/generated-openapi-core.yaml
SWAGGER_UI_DND5E_OPENAPI_URL=/swagger/generated-openapi-dnd5e.yaml
```

### Option B: Start Docker Compose with `.env.test.example`

From PowerShell:

```powershell
docker compose --env-file .env.test.example up -d swagger-ui
```

If you also want to start the relay from the same compose file:

```powershell
docker compose --env-file .env.test.example up -d
```

### Step 2: Open Swagger UI

Open:

```text
http://localhost:8081/swagger
```

You will see two API definitions:

- `core`
- `dnd5e`

### Step 3: Test requests

When trying protected endpoints:

1. expand an endpoint
2. click **Try it out**
3. provide your `x-api-key`
4. provide `clientId` when required
5. execute the request

## Option 2: Run Swagger UI Only

If the relay is already running elsewhere and you only want the external Swagger UI container:

```powershell
docker compose --env-file .env.test.example up -d swagger-ui
```

This is useful when you want a local testing UI without changing your relay setup.

## Windows Notes

- Run commands from PowerShell in the repository root.
- Docker Desktop must be running before you start the compose stack.
- If `Try it out` calls the wrong server, update `SWAGGER_UI_RELAY_SERVER_URL`.
- If the browser shows a missing spec error, verify that the two `SWAGGER_UI_*_OPENAPI_URL` variables match the generated files served by the container.

## Postman: Legacy Alternative

The Postman collection is still kept in the repository for compatibility:

[`documentation/openapi/Foundry REST API Documentation.postman_collection.json`](../../documentation/openapi/Foundry%20REST%20API%20Documentation.postman_collection.json)

However, Postman is now considered a legacy fallback and is **not the recommended manual testing workflow**.

Reasons:

- it is not the canonical source anymore
- the OpenAPI YAML files are the current reference
- Swagger UI matches the current documentation pipeline more closely

Use Postman only if you already have an existing workflow that depends on it.
