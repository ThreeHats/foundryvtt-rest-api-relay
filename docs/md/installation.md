---
id: installation
title: Installation
sidebar_position: 2
---

# Installation

## Docker (Recommended)

**Prerequisites:** [Docker Engine](https://docs.docker.com/engine/install/) + [Docker Compose](https://docs.docker.com/compose/install/linux/) (Linux) or [Docker Desktop](https://docs.docker.com/get-started/get-docker/) (Windows/Mac).

Three commands to get started:

```bash
mkdir foundry-relay && cd foundry-relay
curl -O https://raw.githubusercontent.com/ThreeHats/foundryvtt-rest-api-relay/main/docker-compose.local.yml
docker compose -f docker-compose.local.yml up -d
```

That's it. The relay is now running at `http://localhost:3010`.

Open it in your browser, click **Sign Up**, and create your account. **Your relay key is shown exactly once** after registration — copy it into a password manager before dismissing the dialog. Then create a scoped API key from the dashboard's **API Keys** tab for all integrations.

---

## Optional Configuration

To configure SMTP, rate limits, or other settings, grab the example config and edit it before starting:

```bash
curl -O https://raw.githubusercontent.com/ThreeHats/foundryvtt-rest-api-relay/main/.env.example
cp .env.example .env
# edit .env, then:
docker compose -f docker-compose.local.yml up -d
```

See [Server Configuration](./configuration) for the full list of environment variables.

---

## GPU Acceleration (Optional)

GPU acceleration improves headless Chrome performance for screenshots, sheet rendering, and canvas-heavy operations. The relay auto-detects the best available option once the device is exposed.

- **NVIDIA (Linux):** install `nvidia-container-toolkit`, set it as the default Docker runtime, and restart Docker. No compose changes needed — the image already requests the right capabilities. See [GPU Acceleration in Docker](./configuration#gpu-acceleration-in-docker) for the full setup. Note: use Docker Engine CE, not Docker Desktop — Docker Desktop on Linux runs in a VM and can't pass through the GPU.
- **Intel / AMD (Linux):** expose `/dev/dri` in your compose file — see the commented-out example in `docker-compose.local.yml`.
- **Windows (Docker Desktop + WSL2):** NVIDIA passthrough works — follow the Linux steps inside WSL2.
- **Mac:** no GPU passthrough available; the relay falls back to software rendering automatically and still works.

:::note NVIDIA toolkit version
Requires `nvidia-container-toolkit` >= ~1.14. Some distros (Pop!_OS, Ubuntu) ship older versions. If you see `nvidia-container-runtime did not terminate successfully: exit status 2`, upgrade from the NVIDIA repo:
```bash
sudo apt-get install \
  nvidia-container-toolkit=1.19.0-1 \
  nvidia-container-toolkit-base=1.19.0-1 \
  libnvidia-container-tools=1.19.0-1 \
  libnvidia-container1=1.19.0-1
sudo systemctl restart docker
```
:::

---

## Stopping and Updating

```bash
# Stop
docker compose -f docker-compose.local.yml down

# Update to latest image
docker compose -f docker-compose.local.yml pull
docker compose -f docker-compose.local.yml up -d
```

See [Updating the Docker Image](./update-docker-image) for more detail.

---

## Other Setup Options

- **PostgreSQL** — for a more production-ready database backend: [PostgreSQL Setup Guide](./postgres-setup)
- **HTTPS + reverse proxy** — to expose the relay on the internet: [Relay + App + DNS Example](./relay-app-duckdns-example)
- **Fly.io / cloud deployment** — managed cloud without the public relay: see [Server Configuration](./configuration)

---

## Accessing the Relay From Other Devices on Your LAN

The relay binds to `0.0.0.0` by default, so it's reachable from any device on your local network. The startup logs print every reachable URL:

```
INF Local URL url=http://localhost:3010
INF LAN URL (reachable from other devices on your network) url=http://192.168.1.42:3010
```

When configuring the Foundry module, point it at the LAN IP (`http://192.168.1.42:3010`), not `localhost`, unless Foundry is on the same machine as the relay.

If it doesn't work: check your host firewall has port `3010` open, and that your compose file publishes `3010:3010`.

---

## Manual Installation (No Docker)

1. **Prerequisites:** Go 1.22+, Node.js v18+ and pnpm (for frontend build and tests), Chromium/Chrome (for headless sessions)

2. **Clone and run:**
    ```bash
    git clone https://github.com/ThreeHats/foundryvtt-rest-api-relay.git
    cd foundryvtt-rest-api-relay
    pnpm run local:sqlite
    ```

    Or build the binary directly:
    ```bash
    cd go-relay
    go build -o relay ./cmd/server/
    DB_TYPE=sqlite PORT=3010 ./relay
    ```

The server will be running at `http://localhost:3010`.
