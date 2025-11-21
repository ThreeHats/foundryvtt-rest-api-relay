# Foundry REST API
Join the [discord](https://discord.gg/U634xNGRAC) server for updates, questions, and discussions

This project consists of two main components:

- [Relay Server](https://github.com/ThreeHats/foundryvtt-rest-api-relay): A WebSocket server that facilitates communication between Foundry VTT and external applications.
- [Foundry Module](https://github.com/ThreeHats/foundryvtt-rest-api): A Foundry VTT module that connects to the relay server and provides access to Foundry data through a REST API.

## Foundry REST API Relay Server
The server provides WebSocket connectivity and a REST API to access Foundry VTT data remotely.

### Features
- [Documentation](https://github.com/ThreeHats/foundryvtt-rest-api/wiki)
- WebSocket relay to connect Foundry clients with external applications
- REST API endpoints for searching Foundry content and retrieving entity data
- Client management for tracking Foundry connections
- Data storage and search results
- [Roadmap](https://github.com/users/ThreeHats/projects/7)

### Installation

#### Using Docker Compose (Recommended)
The easiest way to run the relay server is using Docker Compose:

```bash
# Clone the repository
git clone https://github.com/ThreeHats/foundryvtt-rest-api-relay.git
cd foundryvtt-rest-api-relay

# Start the server
docker-compose up -d

# To stop the server
docker-compose down
```

The server will be available at http://localhost:3010 and will automatically restart unless manually stopped.

#### Manual Installation
```bash
### Install dependencies
pnpm install

### Build SQLite native module (required for local:sqlite mode)
cd node_modules/.pnpm/sqlite3@5.1.7/node_modules/sqlite3 && npm run install && cd -

### Build for production
pnpm build

### Run with SQLite database (recommended)
pnpm run local:sqlite

### Run in development mode (requires postgres)
pnpm dev

### Start production server (requires postgres)
pnpm start
```

### Configuration

The server can be configured using environment variables:

- `PORT`: The port the server listens on (default: `3010`).
- `NODE_ENV`: Set to `production` for production deployments.
- `WEBSOCKET_PING_INTERVAL_MS`: Interval in milliseconds for sending WebSocket protocol pings to keep connections alive (default: `20000`).
- `CLIENT_CLEANUP_INTERVAL_MS`: Interval in milliseconds for checking and removing inactive clients (default: `15000`).
- `REDIS_URL`: Connection URL for Redis (optional, used for multi-instance deployments and session storage).

When using Docker Compose, you can set these in the `environment` section of the `docker-compose.yml` file.

### Documentation

This project uses TypeDoc and Docusaurus for comprehensive API documentation. The documentation is automatically generated from TypeScript source code and includes both manual documentation and auto-generated API references.

#### Development

To work with the documentation:

```bash
# Install documentation dependencies
pnpm docs:install

# Generate API documentation from TypeScript source
pnpm docs:generate

# Start the documentation development server
pnpm docs:dev
```

The documentation will be available at [http://localhost:3000](http://localhost:3000) and includes:

- **Manual Documentation**: Getting started guides, installation instructions, and usage examples
- **API Reference**: Auto-generated from TypeScript source code using TypeDoc
- **Interactive Navigation**: Browse the codebase structure and find specific functions, classes, and types

#### Building for Production

```bash
# Build static documentation files
pnpm docs:build

# Serve the built documentation
pnpm docs:serve
```

The documentation system automatically:
- Generates markdown files from TypeScript comments and type definitions
- Creates sidebar navigation based on code structure
- Links to source code on GitHub
- Updates when source code changes

## Foundry REST API Module
A Foundry VTT module that connects to the relay server and provides access to Foundry data.

### Features
- WebSocket connection to relay server
- Integration with Foundry's QuickInsert for powerful search capabilities
- Entity retrieval by UUID
- Configurable WebSocket relay URL and token

### Installation
1. Install the module with the latest manifest link [https://github.com/ThreeHats/foundryvtt-rest-api/releases/latest/download/module.json]([https://github.com/ThreeHats/foundryvtt-rest-api/releases/latest/download/module.json](https://github.com/ThreeHats/foundryvtt-rest-api/releases/latest/download/module.json))
2. Configure the WebSocket relay URL in module settings
3. Set your relay token (defaults to your world ID)

### Configuration
After installing the module, go to the module settings to configure:

- WebSocket Relay URL - URL for the WebSocket relay server (default: ws://localhost:3010)
- WebSocket Relay Token - Token for grouping users together (default: your world ID)

### Technical Details
#### Server Architecture
- Express.js - HTTP server framework
- WebSocket - For real-time communication
- Data Store - In-memory storage for entities and search results
- Client Manager - Handles client connections and message routing

#### Module Architecture
- Built with TypeScript for Foundry VTT
- Integrates with Foundry's QuickInsert for powerful search capabilities
- Provides WebSocket relay functionality for external applications

