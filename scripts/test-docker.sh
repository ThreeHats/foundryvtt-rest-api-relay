#!/usr/bin/env bash
# Run the integration suite against a Docker-managed Foundry + relay stack.
#
# Model: quicksilvervtt-lite runs the Foundry server (with the REST API module
# installed and enabled); the relay launches a headless browser against it via
# /start-session, and the module connects back to the relay. The Jest suite then
# drives the relay's HTTP API.
#
# Usage:
#   scripts/test-docker.sh [major-version]      # 11 | 12 | 13 | 14 (default 13)
#
# Credentials are read from an env file (default: .env.test.docker; override with
# ENV_FILE=path). It needs your foundryvtt.com account — the SAME variables as a
# quicksilvervtt-lite .env, so you can point ENV_FILE at one you already have:
#   FOUNDRY_USERNAME=you@example.com
#   FOUNDRY_PASSWORD=your-foundryvtt-password
#
# Optional env (file or shell):
#   MODULE_PATH  — path to the foundryvtt-rest-api checkout (default: ../foundryvtt-rest-api)
#   KEEP_RUNNING — if set, leave the stack up after tests (for debugging)

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"
MAJOR="${1:-13}"

# --- Load env file (same format as a quicksilvervtt-lite .env) ---
ENV_FILE="${ENV_FILE:-$REPO_DIR/.env.test.docker}"
if [ -f "$ENV_FILE" ]; then
  echo "Loading env from $ENV_FILE"
  set -a; . "$ENV_FILE"; set +a
fi

# The quicksilvervtt-lite image is selected by major version (its tag tracks that
# major's latest Foundry build, with the matching Node version).
case "$MAJOR" in
  11|12|13|14) ;;
  *) echo "ERROR: Unsupported Foundry version: $MAJOR (supported: 11, 12, 13, 14)" >&2; exit 1 ;;
esac

# foundryvtt.com credentials (for the Foundry download) — distinct from the world
# GM login the relay uses to join the world.
: "${FOUNDRY_USERNAME:?Set FOUNDRY_USERNAME (foundryvtt.com account) in $ENV_FILE or the shell}"
: "${FOUNDRY_PASSWORD:?Set FOUNDRY_PASSWORD (foundryvtt.com account) in $ENV_FILE or the shell}"

# World GM password the relay uses to join. The configs ship a passwordless GM,
# so this defaults to empty. To test a password-protected GM, set GM_PASSWORD and
# add a matching "password" to the GM in the per-version config.
GM_PASSWORD="${GM_PASSWORD-}"

MODULE_PATH="${MODULE_PATH:-$(realpath "$REPO_DIR/../foundryvtt-rest-api" 2>/dev/null || echo "")}"
if [ -z "$MODULE_PATH" ] || [ ! -d "$MODULE_PATH" ]; then
  echo "ERROR: Foundry REST API module not found at '$MODULE_PATH'" >&2
  echo "  Set MODULE_PATH to your foundryvtt-rest-api checkout." >&2
  exit 1
fi

RELAY_HOST_PORT="${RELAY_HOST_PORT:-3010}"   # override to coexist with a dev server on 3010
FOUNDRY_TEST_CONFIG_FILE="$(realpath "$REPO_DIR/test/foundry-configs/v${MAJOR}.json")"

echo "=== Docker Integration Tests: Foundry v${MAJOR} ==="
echo "  Relay source:  $REPO_DIR"
echo "  Module source: $MODULE_PATH"
echo "  Relay port:    $RELAY_HOST_PORT (host) -> 3010 (container)"
echo "  Foundry:       internal only (relay reaches it at foundry:30000)"
echo ""

# --- Build the REST API module (CI=1 -> output to dist/ for the dev-module mount) ---
echo "[1/4] Building Foundry REST API module..."
(cd "$MODULE_PATH" && CI=1 npm run build)

# --- Ephemeral relay secrets (override by setting them in the env file/shell) ---
export CREDENTIALS_ENCRYPTION_KEY="${CREDENTIALS_ENCRYPTION_KEY:-$(openssl rand -hex 32)}"
export ADMIN_JWT_SECRET="${ADMIN_JWT_SECRET:-$(openssl rand -base64 32)}"

# --- Write .env.test for the Jest run ---
# FOUNDRY_V{n}_URL is passed to the relay as x-foundry-url, so it must be the
# relay's container-network view of Foundry (foundry:30000), not localhost.
echo "[2/4] Writing .env.test..."
cat > "$REPO_DIR/.env.test" <<EOF
TEST_BASE_URL=http://localhost:${RELAY_HOST_PORT}
FOUNDRY_V${MAJOR}_URL=http://foundry:30000
FOUNDRY_V${MAJOR}_WORLD=test-world
TEST_FOUNDRY_VERSIONS=${MAJOR}
USE_EXISTING_SESSION=false
FOUNDRY_USERNAME=Gamemaster
FOUNDRY_PASSWORD=${GM_PASSWORD}
CAPTURE_BROWSER_CONSOLE=warn
EOF

# --- Start the stack ---
echo "[3/4] Starting Docker stack (Foundry v${MAJOR} + relay)..."
cd "$REPO_DIR"
# Export the compose vars so both `up` here and `down` in the teardown trap can
# interpolate the file (otherwise teardown fails to parse the volume specs).
export FOUNDRY_MAJOR="$MAJOR"
export RELAY_HOST_PORT FOUNDRY_TEST_CONFIG_FILE
export MODULE_SOURCE_PATH="$MODULE_PATH/dist"
docker compose -f docker-compose.test.yml up -d --build

teardown() {
  if [ -n "${KEEP_RUNNING:-}" ]; then
    echo ""
    echo "KEEP_RUNNING set — stack left up. Relay: http://localhost:${RELAY_HOST_PORT}  (Foundry: internal; \`docker compose -f docker-compose.test.yml exec foundry ...\`)"
    echo "Stop with: docker compose -f docker-compose.test.yml down -v"
  else
    echo "Tearing down..."
    docker compose -f docker-compose.test.yml down -v
  fi
}
trap teardown EXIT

# --- Wait for relay ---
echo "  Waiting for relay..."
until curl -fsS "http://localhost:${RELAY_HOST_PORT}/api/health" >/dev/null 2>&1; do sleep 2; done
echo "  Relay healthy."

# --- Wait for Foundry to finish downloading + load the world ---
echo "  Waiting for Foundry world (first run downloads Foundry + dnd5e, can take minutes)..."
WAITED=0; MAX_WAIT=600
# Check via the relay container, which shares the Docker network with Foundry —
# no Foundry host port needed.
until docker compose -f docker-compose.test.yml exec -T relay \
      curl -fsS http://foundry:30000/api/status 2>/dev/null | grep -q '"active":true'; do
  sleep 5; WAITED=$((WAITED + 5))
  if [ "$WAITED" -ge "$MAX_WAIT" ]; then
    echo "ERROR: Foundry world not active within ${MAX_WAIT}s. Recent logs:" >&2
    docker compose -f docker-compose.test.yml logs --tail=60 foundry >&2
    exit 1
  fi
  [ $((WAITED % 30)) -eq 0 ] && echo "  ...still waiting (${WAITED}s)"
done
echo "  Foundry world active."

# --- Run the suite ---
# Force the world GM login here so the foundryvtt.com FOUNDRY_USERNAME/PASSWORD
# loaded from the env file can't leak into the test process.
echo "[4/4] Running integration tests..."
set +e
FOUNDRY_USERNAME=Gamemaster FOUNDRY_PASSWORD="$GM_PASSWORD" pnpm test
TEST_EXIT_CODE=$?
set -e

exit $TEST_EXIT_CODE
