#!/usr/bin/env bash
#
# Benchmark: Go Relay vs Node.js Relay
#
# Usage:
#   ./benchmark/compare.sh [--skip-jest] [--rounds N]
#
# Prerequisites:
#   - Go server source at go-relay/
#   - Node.js server source at ../../ (main repo)
#   - Foundry VTT running with module connected to whichever server is on port 3010
#
# NOTE: This script does NOT kill Firefox or any browser processes.

set -euo pipefail

export PATH="$PATH:/usr/local/go/bin:$HOME/go/bin"

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
WORKTREE_DIR="$(dirname "$SCRIPT_DIR")"
GO_DIR="$WORKTREE_DIR/go-relay"
NODE_DIR="$(cd "$WORKTREE_DIR/../.." && pwd)"
RESULTS_DIR="$SCRIPT_DIR/results"
ROUNDS=${ROUNDS:-50}
SKIP_JEST=false
GO_PORT=3011
NODE_PORT=3012
TEST_PORT=3010  # Port that Foundry module connects to

for arg in "$@"; do
  case $arg in
    --skip-jest) SKIP_JEST=true ;;
    --rounds) shift; ROUNDS=$1 ;;
    --rounds=*) ROUNDS="${arg#*=}" ;;
  esac
done

mkdir -p "$RESULTS_DIR"

TIMESTAMP=$(date +%Y%m%d_%H%M%S)
REPORT="$RESULTS_DIR/benchmark_${TIMESTAMP}.md"

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m'

log() { echo -e "${BLUE}[benchmark]${NC} $1"; }
header() { echo -e "\n${GREEN}=== $1 ===${NC}"; }

# ─── Helpers ──────────────────────────────────────────────

get_pid_on_port() {
  lsof -ti:"$1" -sTCP:LISTEN 2>/dev/null || true
}

wait_for_server() {
  local port=$1 max=100
  for i in $(seq 1 $max); do
    if curl -s "http://localhost:$port/api/status" >/dev/null 2>&1; then
      return 0
    fi
    sleep 0.05
  done
  echo "Server on port $port failed to start" >&2
  return 1
}

get_mem_rss_kb() {
  local pid=$1
  if [ -f "/proc/$pid/status" ]; then
    grep VmRSS "/proc/$pid/status" 2>/dev/null | awk '{print $2}' || echo "0"
  else
    ps -o rss= -p "$pid" 2>/dev/null | tr -d ' ' || echo "0"
  fi
}

get_mem_rss_mb() {
  local kb
  kb=$(get_mem_rss_kb "$1")
  echo "scale=1; $kb / 1024" | bc 2>/dev/null || echo "?"
}

# Kill only the server we started (by saved PID), never anything else
safe_kill() {
  local pid=$1
  if [ -n "$pid" ] && kill -0 "$pid" 2>/dev/null; then
    kill "$pid" 2>/dev/null
    wait "$pid" 2>/dev/null || true
  fi
}

# ─── 1. Build sizes ──────────────────────────────────────

header "Build Sizes"

log "Building Go binary..."
cd "$GO_DIR"
go build -o /tmp/relay-bench ./cmd/server/
GO_BIN_SIZE=$(du -h /tmp/relay-bench | cut -f1)
GO_BIN_BYTES=$(stat -c%s /tmp/relay-bench 2>/dev/null || stat -f%z /tmp/relay-bench)
log "Go binary: $GO_BIN_SIZE"

log "Checking Node.js install size..."
# Use du with --max-depth=0 for speed
NODE_MODULES_SIZE=$(du -sh --max-depth=0 "$NODE_DIR/node_modules" 2>/dev/null | cut -f1 || echo "?")
[ -z "$NODE_MODULES_SIZE" ] && NODE_MODULES_SIZE="?"
log "node_modules: $NODE_MODULES_SIZE"

# ─── 2. Startup time ─────────────────────────────────────

header "Startup Time"

# Go startup
log "Measuring Go startup time..."
GO_STARTUP_TIMES=()
for i in $(seq 1 5); do
  start_ns=$(date +%s%N)
  DB_TYPE=memory PORT=$GO_PORT /tmp/relay-bench &>/dev/null &
  GO_PID=$!
  wait_for_server $GO_PORT
  end_ns=$(date +%s%N)
  elapsed_ms=$(( (end_ns - start_ns) / 1000000 ))
  GO_STARTUP_TIMES+=($elapsed_ms)
  safe_kill $GO_PID
  sleep 0.3
done
GO_STARTUP_AVG=$(echo "${GO_STARTUP_TIMES[@]}" | tr ' ' '\n' | awk '{s+=$1} END {printf "%.0f", s/NR}')
log "Go avg startup: ${GO_STARTUP_AVG}ms"

# Node startup
log "Measuring Node.js startup time..."
NODE_STARTUP_TIMES=()
for i in $(seq 1 5); do
  start_ns=$(date +%s%N)
  bash -c "cd '$NODE_DIR' && DB_TYPE=memory PORT=$NODE_PORT node dist/index.js" &>/dev/null &
  NODE_BASH_PID=$!
  wait_for_server $NODE_PORT
  end_ns=$(date +%s%N)
  elapsed_ms=$(( (end_ns - start_ns) / 1000000 ))
  NODE_STARTUP_TIMES+=($elapsed_ms)
  # Kill the node process (child of bash)
  pkill -P $NODE_BASH_PID 2>/dev/null || true
  safe_kill $NODE_BASH_PID
  sleep 0.5
done
NODE_STARTUP_AVG=$(echo "${NODE_STARTUP_TIMES[@]}" | tr ' ' '\n' | awk '{s+=$1} END {printf "%.0f", s/NR}')
log "Node avg startup: ${NODE_STARTUP_AVG}ms"

# ─── 3. Memory usage (idle) ──────────────────────────────

header "Memory Usage (Idle)"

cd "$GO_DIR"
DB_TYPE=memory PORT=$GO_PORT /tmp/relay-bench &>/dev/null &
wait_for_server $GO_PORT
sleep 2
GO_PID=$(lsof -ti:$GO_PORT -sTCP:LISTEN 2>/dev/null | head -1)
GO_MEM_IDLE=$(get_mem_rss_mb $GO_PID)
log "Go idle RSS: ${GO_MEM_IDLE} MB"

bash -c "cd '$NODE_DIR' && DB_TYPE=memory PORT=$NODE_PORT node dist/index.js" &>/dev/null &
NODE_BASH_PID=$!
wait_for_server $NODE_PORT
sleep 2
NODE_PID=$(lsof -ti:$NODE_PORT -sTCP:LISTEN 2>/dev/null | head -1)
[ -z "$NODE_PID" ] && NODE_PID=$NODE_BASH_PID
NODE_MEM_IDLE=$(get_mem_rss_mb $NODE_PID)
log "Node idle RSS: ${NODE_MEM_IDLE} MB"

# ─── 4. API Response Times ───────────────────────────────

header "API Response Times ($ROUNDS rounds each)"

run_bench() {
  local port=$1 label=$2
  local times_status=() times_search=() times_create=()

  for i in $(seq 1 $ROUNDS); do
    # GET /api/status (no auth needed in memory mode)
    t=$(curl -o /dev/null -s -w '%{time_total}' "http://localhost:$port/api/status")
    times_status+=($t)

    # GET /search (requires client, will get error but measures middleware overhead)
    t=$(curl -o /dev/null -s -w '%{time_total}' "http://localhost:$port/search?query=test")
    times_search+=($t)

    # GET /api/health
    t=$(curl -o /dev/null -s -w '%{time_total}' "http://localhost:$port/api/health")
    times_create+=($t)
  done

  # Calculate averages in ms
  avg_status=$(printf '%s\n' "${times_status[@]}" | awk '{s+=$1} END {printf "%.2f", s/NR*1000}')
  avg_search=$(printf '%s\n' "${times_search[@]}" | awk '{s+=$1} END {printf "%.2f", s/NR*1000}')
  avg_health=$(printf '%s\n' "${times_create[@]}" | awk '{s+=$1} END {printf "%.2f", s/NR*1000}')

  # p95
  p95_status=$(printf '%s\n' "${times_status[@]}" | sort -n | awk -v n=$ROUNDS 'NR==int(n*0.95){printf "%.2f", $1*1000}')
  p95_search=$(printf '%s\n' "${times_search[@]}" | sort -n | awk -v n=$ROUNDS 'NR==int(n*0.95){printf "%.2f", $1*1000}')
  p95_health=$(printf '%s\n' "${times_create[@]}" | sort -n | awk -v n=$ROUNDS 'NR==int(n*0.95){printf "%.2f", $1*1000}')

  echo "${label}_STATUS_AVG=$avg_status"
  echo "${label}_STATUS_P95=$p95_status"
  echo "${label}_SEARCH_AVG=$avg_search"
  echo "${label}_SEARCH_P95=$p95_search"
  echo "${label}_HEALTH_AVG=$avg_health"
  echo "${label}_HEALTH_P95=$p95_health"
}

log "Benchmarking Go server..."
eval "$(run_bench $GO_PORT GO)"
log "  /api/status  avg=${GO_STATUS_AVG}ms  p95=${GO_STATUS_P95}ms"
log "  /search      avg=${GO_SEARCH_AVG}ms  p95=${GO_SEARCH_P95}ms"
log "  /api/health  avg=${GO_HEALTH_AVG}ms  p95=${GO_HEALTH_P95}ms"

log "Benchmarking Node.js server..."
eval "$(run_bench $NODE_PORT NODE)"
log "  /api/status  avg=${NODE_STATUS_AVG}ms  p95=${NODE_STATUS_P95}ms"
log "  /search      avg=${NODE_SEARCH_AVG}ms  p95=${NODE_SEARCH_P95}ms"
log "  /api/health  avg=${NODE_HEALTH_AVG}ms  p95=${NODE_HEALTH_P95}ms"

# ─── 5. Memory under load ────────────────────────────────

header "Memory Usage (After Load)"

GO_MEM_LOAD=$(get_mem_rss_mb $GO_PID)
NODE_MEM_LOAD=$(get_mem_rss_mb $NODE_PID)
log "Go after-load RSS: ${GO_MEM_LOAD} MB"
log "Node after-load RSS: ${NODE_MEM_LOAD} MB"

# ─── 6. Concurrent request throughput ─────────────────────

header "Concurrent Throughput (100 requests, 10 concurrent)"

bench_concurrent() {
  local port=$1
  local start_ns=$(date +%s%N)
  seq 1 100 | xargs -P 10 -I{} curl -s -o /dev/null "http://localhost:$port/api/status"
  local end_ns=$(date +%s%N)
  local elapsed_ms=$(( (end_ns - start_ns) / 1000000 ))
  local rps=$(echo "scale=1; 100000 / $elapsed_ms" | bc 2>/dev/null || echo "?")
  echo "$elapsed_ms $rps"
}

read GO_CONC_MS GO_CONC_RPS <<< $(bench_concurrent $GO_PORT)
log "Go: ${GO_CONC_MS}ms total, ${GO_CONC_RPS} req/s"

read NODE_CONC_MS NODE_CONC_RPS <<< $(bench_concurrent $NODE_PORT)
log "Node: ${NODE_CONC_MS}ms total, ${NODE_CONC_RPS} req/s"

# ─── Cleanup temp servers ─────────────────────────────────

safe_kill $GO_PID
pkill -P $NODE_BASH_PID 2>/dev/null || true
safe_kill $NODE_BASH_PID

# ─── 7. Jest test suite timing (optional) ────────────────

if [ "$SKIP_JEST" = false ]; then
  header "Jest Test Suite Timing"
  log "${YELLOW}Skipping Jest comparison — requires Foundry module on port 3010${NC}"
  log "To run: start each server on port 3010 and run 'pnpm test' separately"
  JEST_GO_TIME="N/A"
  JEST_NODE_TIME="N/A"
else
  JEST_GO_TIME="skipped"
  JEST_NODE_TIME="skipped"
fi

# ─── 8. Generate Report ──────────────────────────────────

header "Generating Report"

cat > "$REPORT" << REPORT_EOF
# Benchmark: Go vs Node.js Relay Server
**Date:** $(date -u +"%Y-%m-%d %H:%M:%S UTC")
**Rounds:** $ROUNDS per endpoint

## Build & Deploy

| Metric | Go | Node.js | Difference |
|--------|-----|---------|------------|
| Binary / Install Size | $GO_BIN_SIZE | $NODE_MODULES_SIZE (node_modules) | — |
| Startup Time (avg) | ${GO_STARTUP_AVG}ms | ${NODE_STARTUP_AVG}ms | $(echo "scale=1; $NODE_STARTUP_AVG / $GO_STARTUP_AVG" | bc 2>/dev/null || echo "?")x faster |

## Memory Usage

| Metric | Go | Node.js | Difference |
|--------|-----|---------|------------|
| Idle RSS | ${GO_MEM_IDLE} MB | ${NODE_MEM_IDLE} MB | $(echo "scale=1; $NODE_MEM_IDLE / $GO_MEM_IDLE" | bc 2>/dev/null || echo "?")x less |
| After Load RSS | ${GO_MEM_LOAD} MB | ${NODE_MEM_LOAD} MB | $(echo "scale=1; $NODE_MEM_LOAD / $GO_MEM_LOAD" | bc 2>/dev/null || echo "?")x less |

## Response Times (avg / p95)

| Endpoint | Go avg | Go p95 | Node avg | Node p95 |
|----------|--------|--------|----------|----------|
| GET /api/status | ${GO_STATUS_AVG}ms | ${GO_STATUS_P95}ms | ${NODE_STATUS_AVG}ms | ${NODE_STATUS_P95}ms |
| GET /search | ${GO_SEARCH_AVG}ms | ${GO_SEARCH_P95}ms | ${NODE_SEARCH_AVG}ms | ${NODE_SEARCH_P95}ms |
| GET /api/health | ${GO_HEALTH_AVG}ms | ${GO_HEALTH_P95}ms | ${NODE_HEALTH_AVG}ms | ${NODE_HEALTH_P95}ms |

## Concurrent Throughput (100 requests, 10 concurrent)

| Metric | Go | Node.js |
|--------|-----|---------|
| Total Time | ${GO_CONC_MS}ms | ${NODE_CONC_MS}ms |
| Requests/sec | ${GO_CONC_RPS} | ${NODE_CONC_RPS} |

## Test Suite

| Metric | Go | Node.js |
|--------|-----|---------|
| Jest Suite Time | ${JEST_GO_TIME} | ${JEST_NODE_TIME} |
REPORT_EOF

log "Report saved to: $REPORT"
echo ""
cat "$REPORT"
