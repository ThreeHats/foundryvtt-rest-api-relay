package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type Endpoint struct {
	Name    string
	Method  string
	Path    string
	Body    interface{}
	Headers map[string]string
}

type Result struct {
	Endpoint string
	Duration time.Duration
	Status   int
	Error    error
	BodySize int
}

type Stats struct {
	Endpoint   string
	Count      int
	Errors     int
	Avg        time.Duration
	P50        time.Duration
	P95        time.Duration
	P99        time.Duration
	Min        time.Duration
	Max        time.Duration
	TotalBytes int64
	StatusMap  map[int]int
}

func main() {
	baseURL := flag.String("url", "http://localhost:3010", "Base URL")
	apiKey := flag.String("key", "", "API key (or TEST_API_KEY env)")
	clientsFlag := flag.String("client", "", "Client IDs, comma-separated (or TEST_CLIENT_ID_V13 env)")
	concurrency := flag.Int("c", 20, "Concurrent workers")
	duration := flag.Int("d", 15, "Duration in seconds")
	rampUp := flag.Int("ramp", 2, "Ramp-up seconds")
	flag.Parse()

	if *apiKey == "" {
		*apiKey = os.Getenv("TEST_API_KEY")
	}
	if *clientsFlag == "" {
		*clientsFlag = os.Getenv("TEST_CLIENT_ID_V13")
	}
	if *apiKey == "" || *clientsFlag == "" {
		fmt.Println("ERROR: --key and --client required (or set TEST_API_KEY and TEST_CLIENT_ID_V13 env vars)")
		os.Exit(1)
	}

	clientIDs := strings.Split(*clientsFlag, ",")
	for i := range clientIDs {
		clientIDs[i] = strings.TrimSpace(clientIDs[i])
	}

	// Verify server is up
	fmt.Printf("Target: %s\n", *baseURL)
	resp, err := doRequest(*baseURL+"/api/status", "GET", *apiKey, nil)
	if err != nil {
		fmt.Printf("ERROR: Server not reachable: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Server: %s\n", resp.Body)

	// Auto-discover clients if "auto" is specified
	if len(clientIDs) == 1 && clientIDs[0] == "auto" {
		fmt.Println("Auto-discovering connected clients...")
		clientIDs = nil // Clear "auto" placeholder
		resp, err = doRequest(*baseURL+"/clients", "GET", *apiKey, nil)
		if err != nil {
			fmt.Printf("ERROR: Failed to reach /clients: %v\n", err)
			os.Exit(1)
		}
		if resp.Status != 200 {
			fmt.Printf("ERROR: /clients returned %d: %s\n", resp.Status, resp.Body)
			os.Exit(1)
		}
		var clientsResp map[string]interface{}
		if json.Unmarshal([]byte(resp.Body), &clientsResp) == nil {
			if clients, ok := clientsResp["clients"].([]interface{}); ok {
				for _, c := range clients {
					if cm, ok := c.(map[string]interface{}); ok {
						if id, ok := cm["id"].(string); ok {
							clientIDs = append(clientIDs, id)
						}
					}
				}
			}
		}
		if len(clientIDs) == 0 {
			fmt.Println("ERROR: No connected clients found. Connect Foundry instances first.")
			os.Exit(1)
		}
	}

	// Verify each client is connected
	fmt.Printf("Clients: %d\n", len(clientIDs))
	for _, cid := range clientIDs {
		resp, err = doRequest(*baseURL+"/clients?clientId="+cid, "GET", *apiKey, nil)
		if err != nil || resp.Status != 200 {
			fmt.Printf("  WARNING: %s - may not be connected (status %d)\n", cid, resp.Status)
		} else {
			fmt.Printf("  %s - connected\n", cid)
		}
	}

	// Discover real actor UUIDs per client
	type clientInfo struct {
		ID       string
		ActorUUID string
	}
	var clients []clientInfo

	fmt.Println("Discovering world entities...")
	for _, cid := range clientIDs {
		actorUUID := discoverUUID(*baseURL, *apiKey, cid, "Actor")
		if actorUUID == "" {
			actorUUID = "Actor.nonexistent"
			fmt.Printf("  %s: no actors found\n", cid)
		} else {
			fmt.Printf("  %s: %s\n", cid, actorUUID)
		}
		clients = append(clients, clientInfo{ID: cid, ActorUUID: actorUUID})
	}

	// Build endpoints — for each client, create the full set of WS endpoints
	// Non-WS endpoints are shared (no clientId needed)
	var endpoints []Endpoint

	// Shared endpoints (no WS round-trip)
	endpoints = append(endpoints,
		Endpoint{Name: "GET /api/status", Method: "GET", Path: "/api/status"},
		Endpoint{Name: "GET /api/health", Method: "GET", Path: "/api/health"},
		Endpoint{Name: "GET /clients", Method: "GET", Path: "/clients"},
	)

	// Per-client WS endpoints
	for _, c := range clients {
		suffix := ""
		if len(clients) > 1 {
			// Add client index to name for multi-client runs
			suffix = " [" + c.ID[:8] + "]"
		}
		cid := c.ID
		aUUID := c.ActorUUID

		endpoints = append(endpoints,
			Endpoint{Name: "GET /get (entity)" + suffix, Method: "GET", Path: "/get?uuid=" + aUUID + "&clientId=" + cid},
			Endpoint{Name: "GET /search" + suffix, Method: "GET", Path: "/search?query=goblin&clientId=" + cid},
			Endpoint{Name: "GET /rolls" + suffix, Method: "GET", Path: "/rolls?clientId=" + cid},
			Endpoint{Name: "GET /lastroll" + suffix, Method: "GET", Path: "/lastroll?clientId=" + cid},
			Endpoint{Name: "GET /macros" + suffix, Method: "GET", Path: "/macros?clientId=" + cid},
			Endpoint{Name: "GET /encounters" + suffix, Method: "GET", Path: "/encounters?clientId=" + cid},
			Endpoint{Name: "GET /structure" + suffix, Method: "GET", Path: "/structure?clientId=" + cid},
			Endpoint{Name: "GET /selected" + suffix, Method: "GET", Path: "/selected?clientId=" + cid},
			Endpoint{Name: "GET /players" + suffix, Method: "GET", Path: "/players?clientId=" + cid},
			Endpoint{Name: "GET /effects" + suffix, Method: "GET", Path: "/effects?uuid=" + aUUID + "&clientId=" + cid},
			Endpoint{Name: "GET /file-system" + suffix, Method: "GET", Path: "/file-system?clientId=" + cid},
			Endpoint{Name: "GET /scene (all)" + suffix, Method: "GET", Path: "/scene?all=true&clientId=" + cid},
			Endpoint{Name: "GET /chat" + suffix, Method: "GET", Path: "/chat?limit=5&clientId=" + cid},
			Endpoint{Name: "POST /roll" + suffix, Method: "POST", Path: "/roll?clientId=" + cid,
				Body: map[string]interface{}{"formula": "1d20", "createChatMessage": false}},
		)
	}

	for i := range endpoints {
		endpoints[i].Headers = map[string]string{"x-api-key": *apiKey}
	}

	// Get server PID for memory tracking
	serverPID := getServerPID(*baseURL)

	memBefore := getMemRSS(serverPID)
	cpuBefore := getCPUTime(serverPID)

	fmt.Printf("\nLoad test: %d concurrent workers, %ds duration, %ds ramp-up\n", *concurrency, *duration, *rampUp)
	fmt.Printf("Endpoints: %d\n", len(endpoints))
	fmt.Printf("Server PID: %d, Memory before: %s\n", serverPID, formatBytes(memBefore))
	fmt.Println(strings.Repeat("─", 80))

	// Run load test
	results := runLoadTest(endpoints, *baseURL, *concurrency, time.Duration(*duration)*time.Second, time.Duration(*rampUp)*time.Second)

	memAfter := getMemRSS(serverPID)
	cpuAfter := getCPUTime(serverPID)

	// Calculate stats per endpoint
	statsByEndpoint := calcStats(results)

	// Print results
	fmt.Println()
	fmt.Println(strings.Repeat("─", 110))
	fmt.Printf("%-25s %7s %6s %8s %8s %8s %8s %8s %6s\n",
		"ENDPOINT", "REQS", "ERRS", "AVG", "P50", "P95", "P99", "MAX", "RPS")
	fmt.Println(strings.Repeat("─", 110))

	var totalReqs, totalErrors int
	var totalDuration time.Duration
	for _, s := range statsByEndpoint {
		totalReqs += s.Count
		totalErrors += s.Errors
		totalDuration += s.Avg * time.Duration(s.Count)

		errStr := fmt.Sprintf("%d", s.Errors)
		if s.Errors > 0 {
			errStr = fmt.Sprintf("\033[31m%d\033[0m", s.Errors)
		}

		fmt.Printf("%-25s %7d %6s %8s %8s %8s %8s %8s %6.0f\n",
			truncate(s.Endpoint, 25),
			s.Count, errStr,
			fmtDur(s.Avg), fmtDur(s.P50), fmtDur(s.P95), fmtDur(s.P99), fmtDur(s.Max),
			float64(s.Count)/float64(*duration))
	}

	fmt.Println(strings.Repeat("─", 110))

	overallAvg := totalDuration / time.Duration(totalReqs)
	fmt.Printf("%-25s %7d %6d %8s %8s %8s %8s %8s %6.0f\n",
		"TOTAL", totalReqs, totalErrors,
		fmtDur(overallAvg), "", "", "", "",
		float64(totalReqs)/float64(*duration))

	fmt.Println()
	fmt.Printf("Memory: %s → %s (Δ %s)\n",
		formatBytes(memBefore), formatBytes(memAfter), formatBytes(memAfter-memBefore))
	if cpuBefore >= 0 && cpuAfter >= 0 {
		fmt.Printf("CPU time: %.2fs → %.2fs (Δ %.2fs)\n", cpuBefore, cpuAfter, cpuAfter-cpuBefore)
	}
	fmt.Printf("Goroutines/Threads: %d\n", runtime.NumGoroutine())
}

// discoverUUID fetches the structure endpoint and finds the first world actor UUID.
// World actor UUIDs look like "Actor.abc123" — NOT "Compendium.xxx".
func discoverUUID(baseURL, apiKey, clientID, entityType string) string {
	// First try: get entity directly (returns world entities)
	resp, err := doRequest(
		baseURL+"/get?clientId="+clientID,
		"GET", apiKey, nil)
	if err == nil && resp.Status == 200 {
		var data map[string]interface{}
		if json.Unmarshal([]byte(resp.Body), &data) == nil {
			// Response may contain a list of actors or a single entity
			if uuid := extractWorldUUID(data, entityType); uuid != "" {
				return uuid
			}
		}
		// Try as array
		var arr []interface{}
		if json.Unmarshal([]byte(resp.Body), &arr) == nil {
			for _, item := range arr {
				if m, ok := item.(map[string]interface{}); ok {
					if uuid := extractWorldUUID(m, entityType); uuid != "" {
						return uuid
					}
				}
			}
		}
	}

	// Second try: structure endpoint with actor type
	resp, err = doRequest(
		baseURL+"/structure?types="+entityType+"&includeEntityData=true&clientId="+clientID,
		"GET", apiKey, nil)
	if err != nil || resp.Status != 200 {
		return ""
	}

	var data interface{}
	if err := json.Unmarshal([]byte(resp.Body), &data); err != nil {
		return ""
	}

	return findWorldUUID(data, entityType)
}

// extractWorldUUID gets a UUID from an entity map if it's a world entity (not compendium).
func extractWorldUUID(m map[string]interface{}, entityType string) string {
	uuid, _ := m["uuid"].(string)
	if uuid != "" && strings.HasPrefix(uuid, entityType+".") && !strings.HasPrefix(uuid, "Compendium.") {
		return uuid
	}
	return ""
}

// findWorldUUID recursively searches for a world entity UUID (Actor.xxx, not Compendium.xxx).
func findWorldUUID(data interface{}, entityType string) string {
	switch v := data.(type) {
	case map[string]interface{}:
		if uuid, ok := v["uuid"].(string); ok && uuid != "" {
			// Only accept world-level UUIDs like "Actor.abc123"
			if strings.HasPrefix(uuid, entityType+".") && !strings.Contains(uuid, "Compendium") {
				return uuid
			}
		}
		// Also check _id and construct UUID
		if id, ok := v["_id"].(string); ok && id != "" {
			if _, hasName := v["name"]; hasName {
				// This looks like an entity document — check it's not a compendium
				candidateUUID := entityType + "." + id
				if uuid, ok := v["uuid"].(string); ok {
					candidateUUID = uuid
				}
				if !strings.Contains(candidateUUID, "Compendium") {
					return candidateUUID
				}
			}
		}
		// Recurse into known container keys
		for _, key := range []string{"children", "entities", "contents", "data", "actors", "items"} {
			if child, ok := v[key]; ok {
				if result := findWorldUUID(child, entityType); result != "" {
					return result
				}
			}
		}
	case []interface{}:
		for _, item := range v {
			if result := findWorldUUID(item, entityType); result != "" {
				return result
			}
		}
	}
	return ""
}

func runLoadTest(endpoints []Endpoint, baseURL string, concurrency int, duration, rampUp time.Duration) []Result {
	var mu sync.Mutex
	var results []Result
	var running int32

	stop := make(chan struct{})
	var wg sync.WaitGroup

	// Ramp up workers
	workerDelay := rampUp / time.Duration(concurrency)

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()

			// Stagger start
			time.Sleep(workerDelay * time.Duration(workerID))
			atomic.AddInt32(&running, 1)

			client := &http.Client{Timeout: 30 * time.Second}

			for {
				select {
				case <-stop:
					return
				default:
				}

				// Pick a random endpoint (round-robin by worker)
				ep := endpoints[workerID%len(endpoints)]

				start := time.Now()
				var body io.Reader
				if ep.Body != nil {
					b, _ := json.Marshal(ep.Body)
					body = bytes.NewReader(b)
				}

				req, _ := http.NewRequest(ep.Method, baseURL+ep.Path, body)
				for k, v := range ep.Headers {
					req.Header.Set(k, v)
				}
				if ep.Body != nil {
					req.Header.Set("Content-Type", "application/json")
				}

				resp, err := client.Do(req)
				elapsed := time.Since(start)

				r := Result{
					Endpoint: ep.Name,
					Duration: elapsed,
					Error:    err,
				}

				if err == nil {
					r.Status = resp.StatusCode
					bodyBytes, _ := io.ReadAll(resp.Body)
					r.BodySize = len(bodyBytes)
					resp.Body.Close()
				}

				mu.Lock()
				results = append(results, r)
				mu.Unlock()
			}
		}(i)
	}

	// Wait for ramp-up then run for duration
	time.Sleep(rampUp)
	fmt.Printf("All %d workers active, running for %s...\n", concurrency, duration-rampUp)
	time.Sleep(duration - rampUp)
	close(stop)
	wg.Wait()

	return results
}

func calcStats(results []Result) []Stats {
	byEndpoint := make(map[string][]Result)
	for _, r := range results {
		byEndpoint[r.Endpoint] = append(byEndpoint[r.Endpoint], r)
	}

	var stats []Stats
	for name, rs := range byEndpoint {
		s := Stats{
			Endpoint:  name,
			Count:     len(rs),
			StatusMap: make(map[int]int),
		}

		var durations []time.Duration
		for _, r := range rs {
			if r.Error != nil {
				s.Errors++
				continue
			}
			durations = append(durations, r.Duration)
			s.StatusMap[r.Status]++
			s.TotalBytes += int64(r.BodySize)
			if r.Status >= 500 {
				s.Errors++
			}
		}

		if len(durations) == 0 {
			stats = append(stats, s)
			continue
		}

		sort.Slice(durations, func(i, j int) bool { return durations[i] < durations[j] })

		var total time.Duration
		for _, d := range durations {
			total += d
		}

		s.Avg = total / time.Duration(len(durations))
		s.Min = durations[0]
		s.Max = durations[len(durations)-1]
		s.P50 = percentile(durations, 50)
		s.P95 = percentile(durations, 95)
		s.P99 = percentile(durations, 99)

		stats = append(stats, s)
	}

	// Sort by name for consistent output
	sort.Slice(stats, func(i, j int) bool { return stats[i].Endpoint < stats[j].Endpoint })
	return stats
}

func percentile(sorted []time.Duration, p int) time.Duration {
	if len(sorted) == 0 {
		return 0
	}
	idx := int(math.Ceil(float64(p)/100.0*float64(len(sorted)))) - 1
	if idx < 0 {
		idx = 0
	}
	if idx >= len(sorted) {
		idx = len(sorted) - 1
	}
	return sorted[idx]
}

type httpResponse struct {
	Status int
	Body   string
}

func doRequest(url, method, apiKey string, body interface{}) (*httpResponse, error) {
	var bodyReader io.Reader
	if body != nil {
		b, _ := json.Marshal(body)
		bodyReader = bytes.NewReader(b)
	}
	req, _ := http.NewRequest(method, url, bodyReader)
	req.Header.Set("x-api-key", apiKey)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)
	return &httpResponse{Status: resp.StatusCode, Body: string(b)}, nil
}

func getServerPID(baseURL string) int {
	// Extract port from URL
	parts := strings.Split(baseURL, ":")
	port := parts[len(parts)-1]
	out, err := exec.Command("lsof", "-ti:"+port, "-sTCP:LISTEN").Output()
	if err != nil {
		return 0
	}
	var pid int
	fmt.Sscanf(strings.TrimSpace(string(out)), "%d", &pid)
	return pid
}

func getMemRSS(pid int) int64 {
	if pid == 0 {
		return 0
	}
	data, err := os.ReadFile(fmt.Sprintf("/proc/%d/status", pid))
	if err != nil {
		return 0
	}
	for _, line := range strings.Split(string(data), "\n") {
		if strings.HasPrefix(line, "VmRSS:") {
			var kb int64
			fmt.Sscanf(strings.TrimPrefix(line, "VmRSS:"), "%d", &kb)
			return kb * 1024
		}
	}
	return 0
}

func getCPUTime(pid int) float64 {
	if pid == 0 {
		return -1
	}
	data, err := os.ReadFile(fmt.Sprintf("/proc/%d/stat", pid))
	if err != nil {
		return -1
	}
	fields := strings.Fields(string(data))
	if len(fields) < 15 {
		return -1
	}
	var utime, stime int64
	fmt.Sscanf(fields[13], "%d", &utime)
	fmt.Sscanf(fields[14], "%d", &stime)
	clkTck := 100.0 // sysconf(_SC_CLK_TCK), usually 100 on Linux
	return float64(utime+stime) / clkTck
}

func formatBytes(b int64) string {
	if b < 1024 {
		return fmt.Sprintf("%dB", b)
	}
	if b < 1024*1024 {
		return fmt.Sprintf("%.1fKB", float64(b)/1024)
	}
	return fmt.Sprintf("%.1fMB", float64(b)/(1024*1024))
}

func fmtDur(d time.Duration) string {
	if d == 0 {
		return "-"
	}
	if d < time.Millisecond {
		return fmt.Sprintf("%.0fµs", float64(d)/float64(time.Microsecond))
	}
	return fmt.Sprintf("%.1fms", float64(d)/float64(time.Millisecond))
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n-2] + ".."
}
