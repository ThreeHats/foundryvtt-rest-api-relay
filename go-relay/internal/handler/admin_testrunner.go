package handler

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/config"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/database"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/handler/helpers"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

// TestRunEvent is one item in the streaming output of a test run.
type TestRunEvent struct {
	Type      string    `json:"type"` // "stdout" | "stderr" | "result" | "complete" | "error"
	Data      string    `json:"data"`
	Timestamp time.Time `json:"timestamp"`
}

// TestRun captures the state and output of a single spawned process (test suite or doc generation).
type TestRun struct {
	ID         string
	StartedAt  time.Time
	FinishedAt time.Time
	Status     string // "running" | "completed" | "failed" | "cancelled"
	Cmd        *exec.Cmd
	cancel     context.CancelFunc

	mu     sync.RWMutex
	events []TestRunEvent
	subs   []chan TestRunEvent
}

func (r *TestRun) appendEvent(ev TestRunEvent) {
	r.mu.Lock()
	r.events = append(r.events, ev)
	for _, sub := range r.subs {
		select {
		case sub <- ev:
		default: // drop if subscriber is slow
		}
	}
	r.mu.Unlock()
}

func (r *TestRun) finalize(status string) {
	r.mu.Lock()
	r.Status = status
	r.FinishedAt = time.Now()
	for _, sub := range r.subs {
		close(sub)
	}
	r.subs = nil
	r.mu.Unlock()
}

// Snapshot returns a copy of all events accumulated so far.
func (r *TestRun) Snapshot() []TestRunEvent {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]TestRunEvent, len(r.events))
	copy(out, r.events)
	return out
}

// Subscribe registers a channel for new events. The channel is closed when the run finishes.
func (r *TestRun) Subscribe() chan TestRunEvent {
	ch := make(chan TestRunEvent, 64)
	r.mu.Lock()
	if r.Status == "running" {
		r.subs = append(r.subs, ch)
	} else {
		// already finished — close immediately so caller doesn't block
		close(ch)
	}
	r.mu.Unlock()
	return ch
}

// SnapshotAndSubscribe atomically returns all accumulated events and registers
// a subscription channel. This prevents the race between snapshot and subscribe
// where events appended in the gap (including the complete event) would be lost.
func (r *TestRun) SnapshotAndSubscribe() ([]TestRunEvent, chan TestRunEvent) {
	r.mu.Lock()
	defer r.mu.Unlock()
	events := make([]TestRunEvent, len(r.events))
	copy(events, r.events)
	ch := make(chan TestRunEvent, 2048)
	if r.Status == "running" {
		r.subs = append(r.subs, ch)
	} else {
		// Already finished — snapshot has everything; close channel immediately.
		close(ch)
	}
	return events, ch
}

// testRunRegistry holds all active and recent process runs.
var testRunRegistry = struct {
	mu   sync.Mutex
	runs map[string]*TestRun
}{runs: make(map[string]*TestRun)}

// operationMutex ensures only one operation (test run or doc generation) runs at a time.
var operationMutex sync.Mutex

// curatedEnvKeys is the whitelist of .env.test keys exposed via the admin env editor.
var curatedEnvKeys = []string{
	"TEST_BASE_URL",
	"TEST_API_KEY",
	"TEST_FOUNDRY_VERSIONS",
	"FOUNDRY_V12_URL",
	"FOUNDRY_V13_URL",
	"FOUNDRY_V12_WORLD",
	"FOUNDRY_V13_WORLD",
	"TEST_DEFAULT_WORLD",
	"USE_EXISTING_SESSION",
	"TEST_CLIENT_ID_V12",
	"TEST_CLIENT_ID_V13",
	"FOUNDRY_USERNAME",
	"FOUNDRY_PASSWORD",
	"TEST_USER_EMAIL",
	"TEST_USER_PASSWORD",
	"TEST_ADMIN_EMAIL",
	"TEST_ADMIN_PASSWORD",
	"TEST_PLAYER_USER_ID_V12",
	"TEST_PLAYER_USER_ID_V13",
	"CAPTURE_BROWSER_CONSOLE",
}

// AdminTestRunnerRouter exposes endpoints to trigger, stream, and cancel process runs,
// serve the HTML test report, and read/write .env.test.
func AdminTestRunnerRouter(db *database.DB, cfg *config.Config) chi.Router {
	r := chi.NewRouter()

	// POST /admin/api/tests/run — start the full integration test suite
	r.Post("/run", func(w http.ResponseWriter, req *http.Request) {
		if !operationMutex.TryLock() {
			helpers.WriteError(w, http.StatusConflict, "An operation is already in progress")
			return
		}

		run, err := startTestRun()
		if err != nil {
			operationMutex.Unlock()
			helpers.WriteError(w, http.StatusInternalServerError, err.Error())
			return
		}

		go func(r *TestRun) { <-waitForRunComplete(r); operationMutex.Unlock() }(run)

		auditAdmin(req, db, "tests.run", "test_run", run.ID, "full")
		helpers.WriteJSON(w, http.StatusAccepted, map[string]interface{}{
			"runId":  run.ID,
			"status": run.Status,
		})
	})

	// POST /admin/api/tests/docs/run — regenerate documentation (pnpm run docs:full)
	r.Post("/docs/run", func(w http.ResponseWriter, req *http.Request) {
		if !operationMutex.TryLock() {
			helpers.WriteError(w, http.StatusConflict, "An operation is already in progress")
			return
		}

		run, err := startDocsRun()
		if err != nil {
			operationMutex.Unlock()
			helpers.WriteError(w, http.StatusInternalServerError, err.Error())
			return
		}

		go func(r *TestRun) { <-waitForRunComplete(r); operationMutex.Unlock() }(run)

		auditAdmin(req, db, "tests.docs", "docs_run", run.ID, "full")
		helpers.WriteJSON(w, http.StatusAccepted, map[string]interface{}{
			"runId":  run.ID,
			"status": run.Status,
		})
	})

	// GET /admin/api/tests/runs — list recent runs
	r.Get("/runs", func(w http.ResponseWriter, req *http.Request) {
		testRunRegistry.mu.Lock()
		defer testRunRegistry.mu.Unlock()
		out := make([]map[string]interface{}, 0, len(testRunRegistry.runs))
		for _, run := range testRunRegistry.runs {
			out = append(out, map[string]interface{}{
				"id":         run.ID,
				"status":     run.Status,
				"startedAt":  run.StartedAt,
				"finishedAt": run.FinishedAt,
			})
		}
		sort.Slice(out, func(i, j int) bool {
			return out[i]["startedAt"].(time.Time).After(out[j]["startedAt"].(time.Time))
		})
		helpers.WriteJSON(w, http.StatusOK, map[string]interface{}{"runs": out})
	})

	// GET /admin/api/tests/run/{id}/stream — SSE stream of run output
	r.Get("/run/{id}/stream", func(w http.ResponseWriter, req *http.Request) {
		id := chi.URLParam(req, "id")
		testRunRegistry.mu.Lock()
		run, ok := testRunRegistry.runs[id]
		testRunRegistry.mu.Unlock()
		if !ok {
			helpers.WriteError(w, http.StatusNotFound, "Run not found")
			return
		}

		flusher, ok := w.(http.Flusher)
		if !ok {
			helpers.WriteError(w, http.StatusInternalServerError, "Streaming not supported")
			return
		}
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		// Atomic snapshot + subscribe — no events can be missed in the gap.
		events, ch := run.SnapshotAndSubscribe()

		// On EventSource reconnect the browser sends Last-Event-ID; skip events
		// the client already received to avoid duplicate terminal lines.
		startFrom := 0
		if lastID := req.Header.Get("Last-Event-ID"); lastID != "" {
			if n, err := strconv.Atoi(lastID); err == nil && n >= 0 {
				startFrom = n + 1
			}
		}
		for i := startFrom; i < len(events); i++ {
			writeSSEWithID(w, events[i], i)
		}
		flusher.Flush()

		// Stream new events; return immediately if the client disconnects.
		ctx := req.Context()
		baseIdx := len(events)
		for i := 0; ; i++ {
			select {
			case ev, ok := <-ch:
				if !ok {
					return // channel closed = run finished
				}
				writeSSEWithID(w, ev, baseIdx+i)
				flusher.Flush()
			case <-ctx.Done():
				return // client disconnected
			}
		}
	})


	// POST /admin/api/tests/run/{id}/cancel
	r.Post("/run/{id}/cancel", func(w http.ResponseWriter, req *http.Request) {
		id := chi.URLParam(req, "id")
		testRunRegistry.mu.Lock()
		run, ok := testRunRegistry.runs[id]
		testRunRegistry.mu.Unlock()
		if !ok {
			helpers.WriteError(w, http.StatusNotFound, "Run not found")
			return
		}
		if run.Status != "running" {
			helpers.WriteError(w, http.StatusConflict, "Run is not running")
			return
		}
		if run.cancel != nil {
			run.cancel()
		}
		if run.Cmd != nil && run.Cmd.Process != nil {
			_ = run.Cmd.Process.Signal(syscall.SIGTERM)
			time.AfterFunc(5*time.Second, func() {
				_ = run.Cmd.Process.Kill()
			})
		}
		auditAdmin(req, db, "tests.cancel", "test_run", id, "")
		helpers.WriteJSON(w, http.StatusOK, map[string]string{"message": "Cancellation requested"})
	})

	// GET /admin/api/tests/env — return curated .env.test values
	r.Get("/env", func(w http.ResponseWriter, req *http.Request) {
		root, err := projectRoot()
		if err != nil {
			helpers.WriteError(w, http.StatusInternalServerError, err.Error())
			return
		}
		parsed := parseDotEnv(filepath.Join(root, ".env.test"))
		result := make(map[string]string, len(curatedEnvKeys))
		for _, key := range curatedEnvKeys {
			result[key] = parsed[key]
		}
		helpers.WriteJSON(w, http.StatusOK, result)
	})

	// POST /admin/api/tests/env — update curated .env.test values
	r.Post("/env", func(w http.ResponseWriter, req *http.Request) {
		var body map[string]string
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			helpers.WriteError(w, http.StatusBadRequest, "Invalid JSON")
			return
		}
		// Whitelist: only write keys in curatedEnvKeys
		updates := make(map[string]string)
		for _, key := range curatedEnvKeys {
			if val, ok := body[key]; ok {
				updates[key] = val
			}
		}
		root, err := projectRoot()
		if err != nil {
			helpers.WriteError(w, http.StatusInternalServerError, err.Error())
			return
		}
		if err := writeDotEnvKeys(filepath.Join(root, ".env.test"), updates); err != nil {
			helpers.WriteError(w, http.StatusInternalServerError, "Failed to write .env.test: "+err.Error())
			return
		}
		auditAdmin(req, db, "tests.env.update", "env_file", ".env.test", "")
		helpers.WriteJSON(w, http.StatusOK, map[string]string{"message": "Saved"})
	})

	return r
}

func writeSSE(w http.ResponseWriter, ev TestRunEvent) {
	data, _ := json.Marshal(ev)
	fmt.Fprintf(w, "event: %s\ndata: %s\n\n", ev.Type, data)
}

func writeSSEWithID(w http.ResponseWriter, ev TestRunEvent, id int) {
	data, _ := json.Marshal(ev)
	fmt.Fprintf(w, "id: %d\nevent: %s\ndata: %s\n\n", id, ev.Type, data)
}

// TestReportHandler serves the last generated HTML test report.
// Intentionally unauthenticated — the report contains no secrets, and the
// /admin prefix is already behind the IP allowlist at the server level.
func TestReportHandler(w http.ResponseWriter, r *http.Request) {
	root, err := projectRoot()
	if err != nil {
		helpers.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	reportPath := filepath.Join(root, "test-results", "test-report.html")
	data, err := os.ReadFile(reportPath)
	if err != nil {
		helpers.WriteError(w, http.StatusNotFound, "No test report found — run the test suite first")
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(data)
}

// projectRoot returns the monorepo root (the directory that contains tests/).
// When the server runs from go-relay/, the root is the parent directory.
func projectRoot() (string, error) {
	cwd, _ := os.Getwd()
	parent := filepath.Dir(cwd)
	if _, err := os.Stat(filepath.Join(parent, "tests")); err == nil {
		return parent, nil
	}
	if _, err := os.Stat(filepath.Join(cwd, "tests")); err == nil {
		return cwd, nil
	}
	return "", fmt.Errorf("cannot locate project root from %s", cwd)
}

// launchProcess spawns cmd with the given args, wires up stdout/stderr streaming,
// registers the run in testRunRegistry, and returns immediately.
func launchProcess(cmd string, args ...string) (*TestRun, error) {
	id := uuid.New().String()

	root, err := projectRoot()
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())
	command := exec.CommandContext(ctx, cmd, args...)
	command.Dir = root

	stdout, err := command.StdoutPipe()
	if err != nil {
		cancel()
		return nil, err
	}
	stderr, err := command.StderrPipe()
	if err != nil {
		cancel()
		return nil, err
	}

	run := &TestRun{
		ID:        id,
		StartedAt: time.Now(),
		Status:    "running",
		Cmd:       command,
		cancel:    cancel,
	}

	testRunRegistry.mu.Lock()
	testRunRegistry.runs[id] = run
	testRunRegistry.mu.Unlock()

	if err := command.Start(); err != nil {
		run.appendEvent(TestRunEvent{Type: "error", Data: err.Error(), Timestamp: time.Now()})
		run.finalize("failed")
		cancel()
		return nil, err
	}

	go pumpReader(run, "stdout", stdout)
	go pumpReader(run, "stderr", stderr)

	go func() {
		err := command.Wait()
		status := "completed"
		if err != nil {
			status = "failed"
			run.appendEvent(TestRunEvent{Type: "error", Data: err.Error(), Timestamp: time.Now()})
		}
		run.appendEvent(TestRunEvent{Type: "complete", Data: status, Timestamp: time.Now()})
		run.finalize(status)
		log.Info().Str("runId", id).Str("cmd", cmd).Str("status", status).Msg("admin process finished")
	}()

	return run, nil
}

// startTestRun runs the full integration test suite via `pnpm run test`.
// pnpm run test invokes scripts/runTestsInOrder.ts which builds the ordered file
// list from TEST_ORDER and passes --testSequencer for correct lifecycle ordering.
func startTestRun() (*TestRun, error) {
	return launchProcess("pnpm", "run", "test")
}

// startDocsRun regenerates the documentation via `pnpm run docs:full`.
func startDocsRun() (*TestRun, error) {
	return launchProcess("pnpm", "run", "docs:full")
}

func pumpReader(run *TestRun, kind string, r interface{ Read([]byte) (int, error) }) {
	scanner := bufio.NewScanner(struct{ readerWrap }{readerWrap{r}})
	scanner.Buffer(make([]byte, 1024*1024), 1024*1024) // 1 MB lines
	for scanner.Scan() {
		run.appendEvent(TestRunEvent{
			Type:      kind,
			Data:      scanner.Text(),
			Timestamp: time.Now(),
		})
	}
}

type readerWrap struct {
	r interface{ Read([]byte) (int, error) }
}

func (rw readerWrap) Read(p []byte) (int, error) { return rw.r.Read(p) }

// waitForRunComplete returns a channel that closes when the run reaches a terminal state.
func waitForRunComplete(run *TestRun) <-chan struct{} {
	done := make(chan struct{})
	go func() {
		ticker := time.NewTicker(200 * time.Millisecond)
		defer ticker.Stop()
		for range ticker.C {
			run.mu.RLock()
			s := run.Status
			run.mu.RUnlock()
			if s != "running" {
				close(done)
				return
			}
		}
	}()
	return done
}

// parseDotEnv reads a dotenv file and returns a map of key→value pairs.
// Blank lines and lines starting with # are ignored. Values may be quoted.
// Returns an empty map (not an error) if the file does not exist.
func parseDotEnv(path string) map[string]string {
	result := map[string]string{}
	data, err := os.ReadFile(path)
	if err != nil {
		return result
	}
	scanner := bufio.NewScanner(bytes.NewReader(data))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		idx := strings.IndexByte(line, '=')
		if idx < 0 {
			continue
		}
		key := strings.TrimSpace(line[:idx])
		val := strings.TrimSpace(line[idx+1:])
		// Strip surrounding single or double quotes
		if len(val) >= 2 &&
			((val[0] == '"' && val[len(val)-1] == '"') ||
				(val[0] == '\'' && val[len(val)-1] == '\'')) {
			val = val[1 : len(val)-1]
		}
		result[key] = val
	}
	return result
}

// writeDotEnvKeys updates specific key=value pairs in a dotenv file, preserving
// all other lines (comments, blank lines, unrelated keys). Keys not already
// present as uncommented lines are appended at the end of the file.
func writeDotEnvKeys(path string, updates map[string]string) error {
	// Read existing content, or start empty
	var lines []string
	if data, err := os.ReadFile(path); err == nil {
		scanner := bufio.NewScanner(bytes.NewReader(data))
		for scanner.Scan() {
			lines = append(lines, scanner.Text())
		}
	}

	written := map[string]bool{}

	// Update in-place: find non-commented KEY=... lines and replace the value
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "#") {
			continue
		}
		idx := strings.IndexByte(trimmed, '=')
		if idx < 0 {
			continue
		}
		key := strings.TrimSpace(trimmed[:idx])
		if newVal, ok := updates[key]; ok {
			lines[i] = key + "=" + newVal
			written[key] = true
		}
	}

	// Append any keys that weren't found in the existing file
	var toAppend []string
	for _, key := range curatedEnvKeys { // iterate in stable order
		if updates[key] != "" || written[key] {
			if !written[key] {
				toAppend = append(toAppend, key+"="+updates[key])
			}
		}
	}
	if len(toAppend) > 0 {
		if len(lines) > 0 && lines[len(lines)-1] != "" {
			lines = append(lines, "")
		}
		lines = append(lines, toAppend...)
	}

	content := strings.Join(lines, "\n")
	if !strings.HasSuffix(content, "\n") {
		content += "\n"
	}

	// Atomic write via temp file + rename
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, []byte(content), 0644); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}
