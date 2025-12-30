---
id: authentication
title: Authentication
sidebar_position: 4
---

# Authentication

Most endpoints in the FoundryVTT REST API require authentication to ensure that only authorized users can access and modify your Foundry VTT world data.

## API Key

Authentication is handled via an API key. You must include your API key in the `x-api-key` header for every request to a protected endpoint.

### How to Get Your API Key

1.  **Start the Relay Server:**
    Follow the [Installation](./installation) guide to get your relay server running.

2.  **Access the Web Interface:**
    Open your browser and navigate to the server's address (e.g., `http://localhost:3010`).

3.  **Create an Account:**
    
    **Option A: Register a new account**
    - Click **Sign Up** and create a new account with your email and password
    - Log in with your new credentials

    :::warning Security Notice
    If your server is accessible over a network, or in a production deployment:
    - Use a strong, unique password
    - Run behind HTTPS (e.g., with a reverse proxy like nginx or Caddy)
    - Use `wss://` instead of `ws://` for WebSocket connections
    :::

4.  **Retrieve Your API Key:**
    After logging in, you'll see the **Dashboard** with your account information:
    - Your **API Key** is displayed with a **Copy** button — click it to copy the key to your clipboard
    - You can click **Regenerate** to create a new API key (this invalidates the old one)

## Making Authenticated Requests

When making a request to the API, include the `x-api-key` header with your copied key.

**Example using `curl`:**
```bash
curl -X GET http://localhost:3010/structure \
  -H "x-api-key: YOUR_API_KEY_HERE"
```

**Example using Postman:**
You can download the [Postman Collection](https://github.com/JustAnotherIdea/foundryvtt-rest-api-relay/blob/main/Foundry%20REST%20API%20Documentation.postman_collection.json) to easily test the API. In Postman, you can set the `x-api-key` header for the entire collection or for individual requests.

## Unauthenticated Endpoints

A few endpoints do not require authentication, such as the health check endpoint (`/api/health`) or the main web interface. All data-related endpoints are protected.
