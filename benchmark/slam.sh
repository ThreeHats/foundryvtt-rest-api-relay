#!/usr/bin/env bash
#
# Slam test: Hammers both Go and Node.js servers with real API traffic
#
# Usage:
#   ./benchmark/slam.sh [-c concurrency] [-d duration] [--go-only] [--node-only]
#
# This test runs each server on port 3010 (one at a time) so the Foundry module
# connects to it. Requires Foundry VTT running with the module installed.
#
# NOTE: Does NOT kill Firefox or any browser processes.

set -eo pipefail
export PATH="$PATH:/usr/local/go/bin:$HOME/go/bin"

# Cleanup on exit — always stop any server we started
cleanup() {
  local pid
  pid=$(lsof -ti:3010 -sTCP:LISTEN 2>/dev/null || true)
  if [ -n "$pid" ]; then
    kill "$pid" 2>/dev/null || true
  fi
}
trap cleanup EXIT

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
WORKTREE_DIR="$(dirname "$SCRIPT_DIR")"
GO_DIR="$WORKTREE_DIR/go-relay"
NODE_DIR="$(cd "$WORKTREE_DIR/../../.." && pwd)"
RESULTS_DIR="$SCRIPT_DIR/results"

PORT=3010
CONCURRENCY=20
DURATION=15
RUN_GO=true
RUN_NODE=true
CLIENT_IDS=""

# Load test env
ENV_FILE="$NODE_DIR/.env.test"
if [ -f "$ENV_FILE" ]; then
  export TEST_API_KEY=$(grep "^TEST_API_KEY=" "$ENV_FILE" | cut -d= -f2)
  export TEST_CLIENT_ID_V13=$(grep "^TEST_CLIENT_ID_V13=" "$ENV_FILE" | cut -d= -f2)
fi

if [ -z "$TEST_API_KEY" ]; then
  echo "ERROR: TEST_API_KEY must be set in $ENV_FILE"
  exit 1
fi

while [ $# -gt 0 ]; do
  case $1 in
    -c) CONCURRENCY=$2; shift 2 ;;
    -d) DURATION=$2; shift 2 ;;
    --clients) CLIENT_IDS=$2; shift 2 ;;
    --go-only) RUN_NODE=false; shift ;;
    --node-only) RUN_GO=false; shift ;;
    *) shift ;;
  esac
done

# Default to auto-discovery if no clients specified
if [ -z "$CLIENT_IDS" ]; then
  CLIENT_IDS="auto"
fi

mkdir -p "$RESULTS_DIR"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)

GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

log() { echo -e "${BLUE}[slam]${NC} $1"; }

# Build the load test tool
log "Building load test tool..."
cd "$SCRIPT_DIR"
go build -o /tmp/loadtest ./loadtest.go

# Build Go server
log "Building Go server..."
cd "$GO_DIR"
go build -o /tmp/relay-go ./cmd/server/

wait_for_foundry() {
  local max=90 logfile=${1:-/dev/null}
  log "Waiting for Foundry module to connect (up to ${max}s)..."
  for i in $(seq 1 $max); do
    # Check server log for "Client connected" (works regardless of auth/rate limits)
    if grep -q "Client connected" "$logfile" 2>/dev/null; then
      log "Foundry connected!"
      return 0
    fi
    # Also try the API as a backup
    result=$(curl -s -H "x-api-key: $TEST_API_KEY" "http://localhost:$PORT/clients" 2>/dev/null || echo '{}')
    total=$(echo "$result" | python3 -c "import sys,json; print(json.load(sys.stdin).get('total',0))" 2>/dev/null || echo "0")
    if [ "$total" -ge 1 ] 2>/dev/null; then
      log "Foundry connected!"
      return 0
    fi
    sleep 1
  done
  log "${YELLOW}WARNING: Foundry module did not connect within ${max}s${NC}"
  return 1
}

stop_server() {
  local pid
  pid=$(lsof -ti:$PORT -sTCP:LISTEN 2>/dev/null || true)
  if [ -n "$pid" ]; then
    kill "$pid" 2>/dev/null || true
    # Wait for port to be free
    for i in $(seq 1 20); do
      lsof -ti:$PORT -sTCP:LISTEN &>/dev/null || break
      sleep 0.2
    done
  fi
}

# ─── Go Server ────────────────────────────────────────────

if [ "$RUN_GO" = true ]; then
  echo ""
  echo -e "${GREEN}╔══════════════════════════════════════════════════════════════╗${NC}"
  echo -e "${GREEN}║                    GO SERVER LOAD TEST                       ║${NC}"
  echo -e "${GREEN}╚══════════════════════════════════════════════════════════════╝${NC}"

  stop_server
  sleep 1

  # Reset rate limit counters before test
  sqlite3 "$GO_DIR/data/relay.db" "UPDATE Users SET requests_today = 0, requests_this_month = 0;" 2>/dev/null || true

  log "Starting Go server on port $PORT..."
  DB_TYPE=sqlite FREE_API_REQUESTS_LIMIT=999999999 DAILY_REQUEST_LIMIT=999999999 PORT=$PORT /tmp/relay-go &>/tmp/slam-go.log &
  GO_PID=$!
  sleep 1

  if ! curl -s "http://localhost:$PORT/api/status" >/dev/null 2>&1; then
    log "${RED}Go server failed to start${NC}"
    cat /tmp/slam-go.log
    exit 1
  fi

  wait_for_foundry /tmp/slam-go.log

  log "Running load test: ${CONCURRENCY} workers, ${DURATION}s..."
  echo ""
  /tmp/loadtest -url "http://localhost:$PORT" -key "$TEST_API_KEY" -client "$CLIENT_IDS" \
    -c "$CONCURRENCY" -d "$DURATION" 2>&1 | tee "$RESULTS_DIR/slam_go_${TIMESTAMP}.txt"

  stop_server
  sleep 2
fi

# ─── Node.js Server ──────────────────────────────────────

if [ "$RUN_NODE" = true ]; then
  echo ""
  echo -e "${GREEN}╔══════════════════════════════════════════════════════════════╗${NC}"
  echo -e "${GREEN}║                  NODE.JS SERVER LOAD TEST                    ║${NC}"
  echo -e "${GREEN}╚══════════════════════════════════════════════════════════════╝${NC}"

  stop_server
  sleep 1

  log "Starting Node.js server on port $PORT..."
  cd "$NODE_DIR"
  DB_TYPE=memory PORT=$PORT FREE_API_REQUESTS_LIMIT=999999999 DAILY_REQUEST_LIMIT=999999999 \
    node dist/index.js &>/tmp/slam-node.log &
  NODE_PID=$!
  sleep 2

  if ! curl -s "http://localhost:$PORT/api/status" >/dev/null 2>&1; then
    log "${RED}Node.js server failed to start${NC}"
    cat /tmp/slam-node.log
    exit 1
  fi

  wait_for_foundry /tmp/slam-node.log

  log "Running load test: ${CONCURRENCY} workers, ${DURATION}s..."
  echo ""
  /tmp/loadtest -url "http://localhost:$PORT" -key "$TEST_API_KEY" -client "$CLIENT_IDS" \
    -c "$CONCURRENCY" -d "$DURATION" 2>&1 | tee "$RESULTS_DIR/slam_node_${TIMESTAMP}.txt"

  stop_server
fi

echo ""
echo -e "${GREEN}Results saved to: ${RESULTS_DIR}/${NC}"
echo -e "  Go:   slam_go_${TIMESTAMP}.txt"
echo -e "  Node: slam_node_${TIMESTAMP}.txt"
