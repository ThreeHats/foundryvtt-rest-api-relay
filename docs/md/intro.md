---
id: intro
title: Getting Started
sidebar_position: 1
---

# Getting Started with FoundryVTT REST API Relay

Welcome to the documentation for the FoundryVTT REST API Relay. This project allows you to interact with your Foundry Virtual Tabletop instance through a RESTful API, opening up possibilities for integrations, automations, and external tools.

This documentation will guide you through setting up the relay server, configuring the Foundry VTT module, and making your first API calls.

## How It Works

The project consists of two main parts:

1. **The Relay Server:** A Go application that you can host yourself (or use the public one). It acts as a bridge, managing WebSocket connections from your Foundry VTT instance(s) and exposing a secure HTTP/WebSocket API for external apps.
2. **The Foundry VTT Module:** A module you install in your Foundry VTT setup. It connects to the relay server via WebSocket, authenticates with a per-browser connection token, and listens for API commands to execute within your world.

## Navigation

- **[Installation](./installation):** Step-by-step guide to get the relay server running.
- **[Foundry Module Setup](./foundry-module):** How to install, pair, and configure the module in Foundry VTT.
- **[Authentication & Security Model](./authentication):** Credential types, threat model, and defense in depth.
- **[Your First API Call](./first-api-call):** A simple tutorial to verify your setup is working.
- **[Scoped API Keys](./scoped-keys):** Production-grade integration credentials with narrow scopes and rate limits.
- **[Building Cross-World Foundry Modules](./cross-world-modules):** How to build a Foundry module that talks to OTHER Foundry worlds via the WebSocket tunnel.
- **[API Reference](/api):** Detailed documentation for all available API endpoints.
