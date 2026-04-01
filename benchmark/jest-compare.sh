#!/usr/bin/env bash
#
# Jest Integration Test Performance Comparison: Go vs Node.js
#
# Runs the full Jest test suite against each server and generates a
# comprehensive comparison report with per-file timings, memory, and CPU.
#
# Usage:
#   ./benchmark/jest-compare.sh [--go-only] [--node-only] [--skip-build]
#
# Prerequisites:
#   - Foundry VTT instances running (v12 on 30012, v13 on 30013)
#   - Foundry module installed and configured to connect to port 3010
#   - .env.test configured in the repo root
#   - Go source at go-relay/, Node dist at dist/
#
# NOTE: This script NEVER kills Firefox or any browser processes.

set -euo pipefail

export PATH="$PATH:/usr/local/go/bin:$HOME/go/bin"

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
WORKTREE_DIR="$(dirname "$SCRIPT_DIR")"
GO_DIR="$WORKTREE_DIR/go-relay"
RESULTS_DIR="$SCRIPT_DIR/results"
PORT=3010

RUN_GO=true
RUN_NODE=true
SKIP_BUILD=false

for arg in "$@"; do
  case $arg in
    --go-only)   RUN_NODE=false ;;
    --node-only) RUN_GO=false ;;
    --skip-build) SKIP_BUILD=true ;;
  esac
done

mkdir -p "$RESULTS_DIR"

TIMESTAMP=$(date +%Y%m%d_%H%M%S)
REPORT="$RESULTS_DIR/jest_compare_${TIMESTAMP}.md"

# Temp files for background polling
PEAK_RSS_FILE=$(mktemp)
RSS_POLL_PID=""
trap 'rm -f "$PEAK_RSS_FILE"; [ -n "$RSS_POLL_PID" ] && kill "$RSS_POLL_PID" 2>/dev/null || true' EXIT

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BOLD='\033[1m'
NC='\033[0m'

log()    { echo -e "${BLUE}[bench]${NC} $1"; }
header() { echo -e "\n${GREEN}${BOLD}=== $1 ===${NC}"; }
warn()   { echo -e "${YELLOW}[warn]${NC} $1"; }
err()    { echo -e "${RED}[error]${NC} $1"; }

# ─── Helpers ──────────────────────────────────────────────

get_mem_rss_kb() {
  local pid=$1
  if [ -f "/proc/$pid/status" ]; then
    grep VmRSS "/proc/$pid/status" 2>/dev/null | awk '{print $2}' || echo "0"
  else
    ps -o rss= -p "$pid" 2>/dev/null | tr -d ' ' || echo "0"
  fi
}

get_cpu_ticks() {
  # Returns utime+stime from /proc/pid/stat (in clock ticks)
  local pid=$1
  if [ -f "/proc/$pid/stat" ]; then
    awk '{print $14 + $15}' "/proc/$pid/stat" 2>/dev/null || echo "0"
  else
    echo "0"
  fi
}

safe_kill() {
  local pid=$1
  if [ -n "$pid" ] && kill -0 "$pid" 2>/dev/null; then
    kill "$pid" 2>/dev/null
    wait "$pid" 2>/dev/null || true
  fi
}

wait_for_server() {
  local port=$1 max=200  # 10s max
  for i in $(seq 1 $max); do
    if curl -s "http://localhost:$port/api/status" >/dev/null 2>&1; then
      return 0
    fi
    sleep 0.05
  done
  err "Server on port $port failed to start"
  return 1
}


start_rss_polling() {
  local pid=$1
  echo "0" > "$PEAK_RSS_FILE"
  (
    while kill -0 "$pid" 2>/dev/null; do
      local rss
      rss=$(get_mem_rss_kb "$pid")
      local current_peak
      current_peak=$(cat "$PEAK_RSS_FILE" 2>/dev/null || echo "0")
      if [ "$rss" -gt "$current_peak" ] 2>/dev/null; then
        echo "$rss" > "$PEAK_RSS_FILE"
      fi
      sleep 1
    done
  ) &
  RSS_POLL_PID=$!
}

stop_rss_polling() {
  if [ -n "$RSS_POLL_PID" ]; then
    kill "$RSS_POLL_PID" 2>/dev/null || true
    wait "$RSS_POLL_PID" 2>/dev/null || true
    RSS_POLL_PID=""
  fi
  cat "$PEAK_RSS_FILE" 2>/dev/null || echo "0"
}

kb_to_mb() {
  echo "scale=1; $1 / 1024" | bc 2>/dev/null || echo "?"
}

ticks_to_seconds() {
  # Clock ticks to seconds (usually 100 ticks/sec on Linux)
  local hz
  hz=$(getconf CLK_TCK 2>/dev/null || echo 100)
  echo "scale=2; $1 / $hz" | bc 2>/dev/null || echo "?"
}

parse_junit_times() {
  # Extract per-file durations from junit.xml
  # Output: filename<TAB>time_seconds
  local xml_file=$1
  grep '<testsuite ' "$xml_file" | while read -r line; do
    local name time
    name=$(echo "$line" | sed -n 's/.*name="\([^"]*\)".*/\1/p')
    time=$(echo "$line" | sed -n 's/.*time="\([^"]*\)".*/\1/p')
    # Clean up name: just get the test file basename
    name=$(basename "$name" .test.ts 2>/dev/null || echo "$name")
    name=$(echo "$name" | sed 's|.*/||; s|\.test\.ts$||; s|tests/integration/||')
    if [ -n "$name" ] && [ -n "$time" ]; then
      echo -e "${name}\t${time}"
    fi
  done
}

parse_junit_summary() {
  # Extract total tests, failures, time from junit.xml
  local xml_file=$1
  local tests failures time_total
  tests=$(grep '<testsuites' "$xml_file" | sed -n 's/.*tests="\([^"]*\)".*/\1/p' || echo "?")
  failures=$(grep '<testsuites' "$xml_file" | sed -n 's/.*failures="\([^"]*\)".*/\1/p' || echo "?")
  time_total=$(grep '<testsuites' "$xml_file" | sed -n 's/.*time="\([^"]*\)".*/\1/p' || echo "?")
  echo "$tests $failures $time_total"
}

# ─── Build ────────────────────────────────────────────────

header "Build"

if [ "$SKIP_BUILD" = false ]; then
  if [ "$RUN_GO" = true ]; then
    log "Building Go binary..."
    cd "$GO_DIR"
    go build -o /tmp/relay-bench ./cmd/server/
    log "  Done."
  fi

  if [ "$RUN_NODE" = true ]; then
    if [ ! -f "$WORKTREE_DIR/dist/index.js" ]; then
      log "Building Node.js server..."
      cd "$WORKTREE_DIR"
      pnpm run build 2>&1 | tail -3
    else
      log "Node.js dist/index.js exists."
    fi
  fi
fi

GO_BIN_SIZE="N/A"
GO_BIN_BYTES=0
if [ -f /tmp/relay-bench ]; then
  GO_BIN_SIZE=$(du -h /tmp/relay-bench | cut -f1)
  GO_BIN_BYTES=$(stat -c%s /tmp/relay-bench 2>/dev/null || stat -f%z /tmp/relay-bench)
fi

NODE_MODULES_SIZE="N/A"
if [ -d "$WORKTREE_DIR/node_modules" ]; then
  NODE_MODULES_SIZE=$(du -sh "$WORKTREE_DIR/node_modules" 2>/dev/null | cut -f1)
fi

log "Go binary: $GO_BIN_SIZE | Node modules: $NODE_MODULES_SIZE"

# ─── Results storage ──────────────────────────────────────

declare -A R  # associative array for all results

# ─── Run test cycle ───────────────────────────────────────

run_test_cycle() {
  local label=$1        # "go" or "node"
  local start_cmd=$2    # command to start server
  local server_pid=""

  header "Testing $label server"

  # Start server
  log "Starting $label server on port $PORT..."
  local start_ns=$(date +%s%N)

  eval "$start_cmd" &>/dev/null &
  server_pid=$!

  wait_for_server $PORT
  local ready_ns=$(date +%s%N)
  local startup_ms=$(( (ready_ns - start_ns) / 1000000 ))
  R["${label}_startup"]=$startup_ms
  log "  Server ready in ${startup_ms}ms (PID: $server_pid)"

  # Get the actual server PID (might differ from bash wrapper PID)
  local actual_pid
  actual_pid=$(lsof -ti:$PORT -sTCP:LISTEN 2>/dev/null | head -1 || echo "$server_pid")

  # Let the server fully initialize (DB migrations, cron, etc.)
  sleep 10

  # Pre-test metrics
  local pre_rss=$(get_mem_rss_kb "$actual_pid")
  local pre_cpu=$(get_cpu_ticks "$actual_pid")
  R["${label}_pre_rss"]=$pre_rss

  # Start peak RSS polling
  start_rss_polling "$actual_pid"

  # Run Jest tests
  log "Running Jest test suite..."
  local test_start_ns=$(date +%s%N)
  cd "$WORKTREE_DIR"

  local test_exit=0
  pnpm test 2>&1 | tee "$RESULTS_DIR/${label}_jest_${TIMESTAMP}.log" || test_exit=$?

  local test_end_ns=$(date +%s%N)
  local test_wall_ms=$(( (test_end_ns - test_start_ns) / 1000000 ))
  local test_wall_s=$(echo "scale=1; $test_wall_ms / 1000" | bc)
  R["${label}_test_wall"]=$test_wall_s
  R["${label}_test_exit"]=$test_exit

  # Post-test metrics
  local post_rss=$(get_mem_rss_kb "$actual_pid")
  local post_cpu=$(get_cpu_ticks "$actual_pid")
  local peak_rss
  stop_rss_polling
  peak_rss=$(cat "$PEAK_RSS_FILE" 2>/dev/null || echo "0")

  R["${label}_post_rss"]=$post_rss
  R["${label}_peak_rss"]=$peak_rss

  # CPU time delta
  local cpu_delta=$(( post_cpu - pre_cpu ))
  R["${label}_cpu_ticks"]=$cpu_delta
  R["${label}_cpu_seconds"]=$(ticks_to_seconds $cpu_delta)

  # Copy junit.xml
  if [ -f "$WORKTREE_DIR/test-results/junit.xml" ]; then
    cp "$WORKTREE_DIR/test-results/junit.xml" "$RESULTS_DIR/${label}_junit_${TIMESTAMP}.xml"
    R["${label}_junit"]="$RESULTS_DIR/${label}_junit_${TIMESTAMP}.xml"

    # Parse summary
    read tests failures time_total <<< $(parse_junit_summary "$RESULTS_DIR/${label}_junit_${TIMESTAMP}.xml")
    R["${label}_tests"]=$tests
    R["${label}_failures"]=$failures
    R["${label}_junit_time"]=$time_total
  else
    warn "junit.xml not found!"
    R["${label}_junit"]=""
    R["${label}_tests"]="?"
    R["${label}_failures"]="?"
    R["${label}_junit_time"]="?"
  fi

  # Log summary
  log "  Wall time: ${test_wall_s}s | Exit: $test_exit"
  log "  RSS: pre=$(kb_to_mb $pre_rss)MB → peak=$(kb_to_mb $peak_rss)MB → post=$(kb_to_mb $post_rss)MB"
  log "  CPU time: ${R["${label}_cpu_seconds"]}s"
  log "  Tests: ${R["${label}_tests"]} | Failures: ${R["${label}_failures"]}"

  # Stop server
  log "Stopping $label server..."
  safe_kill $server_pid
  # Also kill the actual process if different
  if [ "$actual_pid" != "$server_pid" ]; then
    safe_kill "$actual_pid"
  fi
  sleep 2

  log "$label cycle complete."
}

# ─── Go test cycle ────────────────────────────────────────

if [ "$RUN_GO" = true ]; then
  # The Go server reads .env from its own directory
  run_test_cycle "go" "cd '$GO_DIR' && DB_TYPE=sqlite FREE_API_REQUESTS_LIMIT=999999999 DAILY_REQUEST_LIMIT=999999999 PORT=$PORT /tmp/relay-bench"
fi

# ─── Node test cycle ──────────────────────────────────────

if [ "$RUN_NODE" = true ]; then
  run_test_cycle "node" "cd '$WORKTREE_DIR' && DB_TYPE=sqlite FREE_API_REQUESTS_LIMIT=999999999 DAILY_REQUEST_LIMIT=999999999 PORT=$PORT node dist/index.js"
fi

# ─── Generate Report ──────────────────────────────────────

header "Generating Report"

# Hardware info
HW_CPU=$(grep 'model name' /proc/cpuinfo 2>/dev/null | head -1 | cut -d: -f2 | xargs || echo "unknown")
HW_RAM=$(grep MemTotal /proc/meminfo 2>/dev/null | awk '{printf "%.0f GB", $2/1024/1024}' || echo "unknown")
HW_KERNEL=$(uname -r)

# Build per-file comparison table
PER_FILE_TABLE=""
if [ -n "${R[go_junit]:-}" ] && [ -n "${R[node_junit]:-}" ]; then
  # Parse both junit files into temp files
  GO_TIMES=$(mktemp)
  NODE_TIMES=$(mktemp)
  parse_junit_times "${R[go_junit]}" | sort -k1 > "$GO_TIMES"
  parse_junit_times "${R[node_junit]}" | sort -k1 > "$NODE_TIMES"

  # Join on test file name
  PER_FILE_TABLE="| Test File | Go (s) | Node (s) | Speedup |
|-----------|--------|----------|---------|
"
  while IFS=$'\t' read -r name go_time; do
    node_time=$(grep "^${name}"$'\t' "$NODE_TIMES" | cut -f2 || echo "?")
    if [ -n "$node_time" ] && [ "$node_time" != "?" ] && [ "$go_time" != "0" ]; then
      speedup=$(echo "scale=1; $node_time / $go_time" | bc 2>/dev/null || echo "?")
      PER_FILE_TABLE+="| ${name} | ${go_time} | ${node_time} | ${speedup}x |
"
    else
      PER_FILE_TABLE+="| ${name} | ${go_time} | ${node_time:-N/A} | — |
"
    fi
  done < "$GO_TIMES"

  rm -f "$GO_TIMES" "$NODE_TIMES"
elif [ -n "${R[go_junit]:-}" ]; then
  PER_FILE_TABLE="| Test File | Go (s) |
|-----------|--------|
"
  while IFS=$'\t' read -r name time; do
    PER_FILE_TABLE+="| ${name} | ${time} |
"
  done < <(parse_junit_times "${R[go_junit]}")
elif [ -n "${R[node_junit]:-}" ]; then
  PER_FILE_TABLE="| Test File | Node (s) |
|-----------|----------|
"
  while IFS=$'\t' read -r name time; do
    PER_FILE_TABLE+="| ${name} | ${time} |
"
  done < <(parse_junit_times "${R[node_junit]}")
fi

# Compute advantages (only if both ran)
WALL_ADV="—"
MEM_IDLE_ADV="—"
MEM_PEAK_ADV="—"
CPU_ADV="—"
if [ "$RUN_GO" = true ] && [ "$RUN_NODE" = true ]; then
  WALL_ADV=$(echo "scale=1; ${R[node_test_wall]} / ${R[go_test_wall]}" | bc 2>/dev/null || echo "?")
  if [ "${R[go_pre_rss]}" -gt 0 ] 2>/dev/null; then
    MEM_IDLE_ADV=$(echo "scale=1; ${R[node_pre_rss]} / ${R[go_pre_rss]}" | bc 2>/dev/null || echo "?")
  fi
  if [ "${R[go_peak_rss]}" -gt 0 ] 2>/dev/null; then
    MEM_PEAK_ADV=$(echo "scale=1; ${R[node_peak_rss]} / ${R[go_peak_rss]}" | bc 2>/dev/null || echo "?")
  fi
  if [ "${R[go_cpu_ticks]}" -gt 0 ] 2>/dev/null; then
    CPU_ADV=$(echo "scale=1; ${R[node_cpu_ticks]} / ${R[go_cpu_ticks]}" | bc 2>/dev/null || echo "?")
  fi
fi

# Write report
cat > "$REPORT" << REPORT_EOF
# Jest Integration Test Performance: Go vs Node.js

**Date:** $(date -u +"%Y-%m-%d %H:%M:%S UTC")
**CPU:** $HW_CPU
**RAM:** $HW_RAM
**Kernel:** $HW_KERNEL

---

## Summary

| Metric | Go | Node.js | Go Advantage |
|--------|-----|---------|--------------|
| Total Suite Time | ${R[go_test_wall]:-N/A}s | ${R[node_test_wall]:-N/A}s | ${WALL_ADV}x faster |
| Tests Passed | ${R[go_tests]:-N/A} | ${R[node_tests]:-N/A} | — |
| Tests Failed | ${R[go_failures]:-N/A} | ${R[node_failures]:-N/A} | — |
| Server Startup | ${R[go_startup]:-N/A}ms | ${R[node_startup]:-N/A}ms | — |

## Memory Usage

| Metric | Go | Node.js | Go Advantage |
|--------|-----|---------|--------------|
| Idle (pre-test) | $(kb_to_mb ${R[go_pre_rss]:-0})MB | $(kb_to_mb ${R[node_pre_rss]:-0})MB | ${MEM_IDLE_ADV}x less |
| Peak (during tests) | $(kb_to_mb ${R[go_peak_rss]:-0})MB | $(kb_to_mb ${R[node_peak_rss]:-0})MB | ${MEM_PEAK_ADV}x less |
| Post-test | $(kb_to_mb ${R[go_post_rss]:-0})MB | $(kb_to_mb ${R[node_post_rss]:-0})MB | — |
| Growth (peak - idle) | $(echo "scale=1; (${R[go_peak_rss]:-0} - ${R[go_pre_rss]:-0}) / 1024" | bc 2>/dev/null || echo "?")MB | $(echo "scale=1; (${R[node_peak_rss]:-0} - ${R[node_pre_rss]:-0}) / 1024" | bc 2>/dev/null || echo "?")MB | — |

## CPU Utilization

| Metric | Go | Node.js | Go Advantage |
|--------|-----|---------|--------------|
| CPU Time (user+sys) | ${R[go_cpu_seconds]:-N/A}s | ${R[node_cpu_seconds]:-N/A}s | ${CPU_ADV}x less |
| Efficiency (tests/cpu-sec) | $(echo "scale=1; ${R[go_tests]:-0} / ${R[go_cpu_seconds]:-1}" | bc 2>/dev/null || echo "?") | $(echo "scale=1; ${R[node_tests]:-0} / ${R[node_cpu_seconds]:-1}" | bc 2>/dev/null || echo "?") | — |

## Memory Timeline

- **Go:** idle=$(kb_to_mb ${R[go_pre_rss]:-0})MB → peak=$(kb_to_mb ${R[go_peak_rss]:-0})MB → post=$(kb_to_mb ${R[go_post_rss]:-0})MB (Δ+$(echo "scale=1; (${R[go_peak_rss]:-0} - ${R[go_pre_rss]:-0}) / 1024" | bc 2>/dev/null || echo "?")MB)
- **Node:** idle=$(kb_to_mb ${R[node_pre_rss]:-0})MB → peak=$(kb_to_mb ${R[node_peak_rss]:-0})MB → post=$(kb_to_mb ${R[node_post_rss]:-0})MB (Δ+$(echo "scale=1; (${R[node_peak_rss]:-0} - ${R[node_pre_rss]:-0}) / 1024" | bc 2>/dev/null || echo "?")MB)

## Build & Deploy

| Metric | Go | Node.js |
|--------|-----|---------|
| Artifact Size | $GO_BIN_SIZE (single binary) | $NODE_MODULES_SIZE (node_modules) |

## Per-File Test Durations

$PER_FILE_TABLE

---
*Generated by \`benchmark/jest-compare.sh\`*
REPORT_EOF

log "Report saved to: $REPORT"
echo ""
cat "$REPORT"
