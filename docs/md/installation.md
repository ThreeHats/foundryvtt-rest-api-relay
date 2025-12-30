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
```yaml
services:
  relay:
    # For production stability, pin to a specific version tag instead of 'latest'
    # See available tags at: https://github.com/ThreeHats/foundryvtt-rest-api-relay/tags
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

:::tip Version Pinning for Production
For production deployments, replace `latest` with a specific version tag (e.g., `threehats/foundryvtt-rest-api-relay:2.1.0`) to avoid unexpected breaking changes from updates.
:::

By adding:   
```yaml
# increase monthly request limit
- FREE_API_REQUESTS_LIMIT=100000
# increase daily request limit
- DAILY_REQUEST_LIMIT=3000   
```
in the environment section you can increase the amount of API requests for your local installation.

2.  **Start the server:**
    ```bash
    docker-compose up -d
    ```
    This command will pull the latest Docker image and start the relay server in the background. The server will be available at `http://localhost:3010`.

3.  **Create Your Account:**
    The default Docker setup uses an SQLite database for persistence, stored in the `data` directory.
    - Open `http://localhost:3010` in your browser
    - Click **Sign Up** and create an account
    - Your API key will be displayed on the dashboard after logging in

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
    - Node.js v18 or later (v20 recommended)
    - pnpm package manager (`npm install -g pnpm`)

2.  **Clone the repository:**
    ```bash
    git clone https://github.com/ThreeHats/foundryvtt-rest-api-relay.git
    cd foundryvtt-rest-api-relay
    ```

3.  **Install dependencies:**
    ```bash
    pnpm install
    ```

4.  **Build SQLite native module (required for local:sqlite mode):**
    
    **Linux/macOS:**
    ```bash
    cd node_modules/.pnpm/sqlite3@*/node_modules/sqlite3 && npm run install && cd -
    ```
    
    **Windows (PowerShell):**
    ```powershell
    # Find the sqlite3 directory (version may vary)
    cd (Get-ChildItem -Path node_modules/.pnpm -Filter "sqlite3@*" -Directory).FullName
    cd node_modules/sqlite3
    npm run install
    cd ../../../..
    ```
    
    > **Note:** This step compiles the SQLite native bindings for your system. If you skip this step, you'll get a "Could not locate the bindings file" error when running with SQLite. Windows users may need [Visual Studio Build Tools](https://visualstudio.microsoft.com/visual-cpp-build-tools/) installed.

5.  **Install Chrome for Puppeteer (required for headless Foundry sessions):**
    ```bash
    npx puppeteer browsers install chrome
    ```
    
    > **Note:** Puppeteer is used to create headless browser sessions for interacting with Foundry VTT. If you skip this step, you'll get a "Could not find Chrome" error when the server tries to start headless sessions.

6.  **Run the server:**
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
