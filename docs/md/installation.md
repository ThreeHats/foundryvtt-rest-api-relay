---
id: installation
title: Installation
sidebar_position: 2
---

# Installation

There are two primary ways to install the FoundryVTT REST API Relay server: using Docker (recommended for ease of use and deployment) or manually.

## Recommended: Docker Installation

Using Docker and Docker Compose is the simplest way to get the relay server running.

1.  **Clone the repository:**
    ```bash
    git clone https://github.com/JustAnotherIdea/foundryvtt-rest-api-relay.git
    cd foundryvtt-rest-api-relay
    ```

Instead of cloning the repository you also just might download the docker-compose.yml file and use that (which will be the latest working version).

Or copy the following and create your own docker-compose.yml
```bash
services:
  relay:
    image: threehats/foundryvtt-rest-api-relay:latest
    container_name: foundryvtt-rest-api-relay
    ports:
      - "3010:3010"
    environment:
      - NODE_ENV=production
      - PORT=3010
      - DB_TYPE=sqlite
      # Optional: Configure connection handling (defaults shown)
      - WEBSOCKET_PING_INTERVAL_MS=20000  # (20 seconds)
      - CLIENT_CLEANUP_INTERVAL_MS=15000  # (15 seconds)
    volumes:
      - ./data:/app/data
    command: pnpm local:sqlite
    restart: unless-stopped
```
By adding:   
```bash
# increase monthly request limit
FREE_API_REQUEST_LIMIT=100000
# increase daily request limit
DAILY_REQUEST_LIMIT=3000   
```
in the enironment section you can increase the amount of API requests for your local installation.

2.  **Start the server:**
    ```bash
    docker-compose up -d
    ```
    This command will pull the latest Docker image and start the relay server in the background. The server will be available at `http://localhost:3010`.

3.  **Database Initialization:**
    The default Docker setup uses an SQLite database for persistence, which is stored in the `data` directory. When you first start the server, a default admin user is created.
    - **Email:** `admin@example.com`
    - **Password:** `admin123`
    
    You can log into the web interface to get your API key.

4.  **Stopping the server:**
    ```bash
    docker-compose down
    ```

4.  **Updating the server:**
-   **[Updating the docker image](./update-docker-image):** Commands to update your docker image.

### Using PostgreSQL
If you prefer to use PostgreSQL for your database, you can use the provided `docker-compose.postgres.yml` file. See the [PostgreSQL Setup Guide](/postgres-setup) for more details.

### Relay + Foundry + duckDNS
For an in depth guide for a full setup using duckDNS see [Relay + App + DNS Example](/relay-app-duckdns-example)

## Manual Installation

If you prefer not to use Docker, you can run the server directly using Node.js.

1.  **Prerequisites:**
    - Node.js (v18 or later)
    - pnpm package manager (`npm install -g pnpm`)

2.  **Clone the repository:**
    ```bash
    git clone https://github.com/JustAnotherIdea/foundryvtt-rest-api-relay.git
    cd foundryvtt-rest-api-relay
    ```

3.  **Install dependencies:**
    ```bash
    pnpm install
    ```

4.  **Build SQLite native module (required for local:sqlite mode):**
    ```bash
    cd node_modules/.pnpm/sqlite3@5.1.7/node_modules/sqlite3 && npm run install && cd -
    ```
    
    > **Note:** This step is necessary to compile the SQLite native bindings for your system. If you skip this step, you'll get a "Could not locate the bindings file" error when running with SQLite.

5.  **Run the server:**
    - **For development (with auto-reloading):**
      ```bash
      pnpm dev
      ```
    - **For production:**
      First, build the project:
      ```bash
      pnpm build
      ```
      Then, start the server using SQLite:
      ```bash
      pnpm local:sqlite
      ```
      Or with an in-memory database (not recommended for production):
      ```bash
      pnpm local
      ```

The server will be running at `http://localhost:3010`.
