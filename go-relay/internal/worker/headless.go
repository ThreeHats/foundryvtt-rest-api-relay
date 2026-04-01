package worker

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	cdp "github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/cdproto/target"
	"github.com/chromedp/chromedp"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/config"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/database"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/service"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/ws"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

const (
	headlessInactiveTimeout = 10 * time.Minute
	pendingSessionTimeout   = 5 * time.Minute
	clientPollInterval      = 500 * time.Millisecond
	clientPollTimeout       = 5 * time.Minute
	browserNavigateTimeout  = 180 * time.Second
	gameLoadTimeout         = 30 * time.Second
)

// HeadlessSession represents an active headless Foundry browser session.
type HeadlessSession struct {
	SessionID     string
	ClientID      string
	APIKey        string
	FoundryURL    string
	Username      string
	WorldName     string
	ContextCancel context.CancelFunc // Cancels the browser context (tab), not the whole browser
	StartedAt     time.Time
	LastActivity  time.Time
}

// PendingHeadless represents a browser context launched but not yet connected.
type PendingHeadless struct {
	SessionID        string
	ExpectedClientID string
	APIKey           string
	ContextCancel    context.CancelFunc
	StartTime        time.Time
}

// SessionInfo is the public view of a session for the GET /session endpoint.
type SessionInfo struct {
	SessionID    string `json:"sessionId"`
	ClientID     string `json:"clientId"`
	FoundryURL   string `json:"foundryUrl"`
	Username     string `json:"username"`
	WorldName    string `json:"worldName,omitempty"`
	StartedAt    int64  `json:"startedAt"`
	LastActivity int64  `json:"lastActivity"`
}

// minCreateTarget sends Target.createTarget with only the essential fields.
// Chrome 146 rejects CreateTarget with extra fields (forTab, hidden) when browserContextId is set.
type minCreateTarget struct {
	URL string                   `json:"url"`
	BCI cdp.BrowserContextID     `json:"browserContextId,omitempty"`
}

func (p *minCreateTarget) Do(ctx context.Context) (target.ID, error) {
	var res struct{ TargetID target.ID `json:"targetId"` }
	err := cdp.Execute(ctx, "Target.createTarget", p, &res)
	return res.TargetID, err
}

// HeadlessManager manages a shared Chrome browser with isolated contexts per session.
type HeadlessManager struct {
	mu             sync.RWMutex
	sessions       map[string]*HeadlessSession // sessionID -> session
	pending        map[string]*PendingHeadless // sessionID -> pending
	clientManager  *ws.ClientManager
	redis          *config.RedisClient
	maxSessions    int
	chromePath     string
	resolvedChrome string

	// Shared browser
	browserCtx    context.Context
	browserCancel context.CancelFunc
	allocCancel   context.CancelFunc
	browserReady  bool
	xvfbCmd       *exec.Cmd
}

// NewHeadlessManager creates a new headless session manager.
func NewHeadlessManager(clientManager *ws.ClientManager, redis *config.RedisClient, cfg *config.Config) *HeadlessManager {
	return &HeadlessManager{
		sessions:      make(map[string]*HeadlessSession),
		pending:       make(map[string]*PendingHeadless),
		clientManager: clientManager,
		redis:         redis,
		maxSessions:   cfg.MaxHeadlessSessions,
		chromePath:    cfg.ChromePath,
	}
}

// getChromePath resolves and caches the Chrome binary path.
func (m *HeadlessManager) getChromePath() string {
	if m.resolvedChrome != "" {
		return m.resolvedChrome
	}
	p := m.chromePath
	if p == "" {
		p = findChromeBinary()
	}
	m.resolvedChrome = p
	return p
}

// startXvfb starts a virtual X display for GPU-accelerated Chrome.
// Returns the display string (e.g. ":99") and a cleanup function.
func startXvfb() (string, *exec.Cmd, error) {
	// Find an available display number
	for display := 99; display < 120; display++ {
		displayStr := fmt.Sprintf(":%d", display)
		// Check if display is already in use
		if _, err := os.Stat(fmt.Sprintf("/tmp/.X11-unix/X%d", display)); err == nil {
			continue // Already in use
		}
		cmd := exec.Command("Xvfb", displayStr, "-screen", "0", "1920x1080x24", "-ac", "-nolisten", "tcp")
		if err := cmd.Start(); err != nil {
			continue
		}
		// Give Xvfb a moment to start
		time.Sleep(200 * time.Millisecond)
		// Verify it started
		if _, err := os.Stat(fmt.Sprintf("/tmp/.X11-unix/X%d", display)); err != nil {
			cmd.Process.Kill()
			continue
		}
		log.Info().Str("display", displayStr).Msg("Xvfb virtual display started")
		return displayStr, cmd, nil
	}
	return "", nil, fmt.Errorf("could not find available display for Xvfb")
}

// ensureBrowser starts the shared browser with GPU acceleration via Xvfb.
func (m *HeadlessManager) ensureBrowser() error {
	if m.browserReady {
		return nil
	}

	// Try to start Xvfb for GPU-accelerated non-headless mode
	useGPU := false
	xvfbDisplay, xvfbCmd, err := startXvfb()
	if err != nil {
		log.Warn().Err(err).Msg("Xvfb not available, falling back to headless mode (no GPU)")
	} else {
		useGPU = true
		os.Setenv("DISPLAY", xvfbDisplay)
		m.xvfbCmd = xvfbCmd
	}

	opts := []chromedp.ExecAllocatorOption{
		chromedp.NoFirstRun,
		chromedp.NoDefaultBrowserCheck,
		chromedp.NoSandbox,
		chromedp.Flag("disable-setuid-sandbox", true),
		chromedp.Flag("disable-dev-shm-usage", true),
		chromedp.Flag("disable-extensions", true),
		chromedp.Flag("disable-web-security", true),
		chromedp.Flag("disable-background-networking", true),
		chromedp.Flag("disable-background-timer-throttling", true),
		chromedp.Flag("disable-backgrounding-occluded-windows", true),
		chromedp.Flag("disable-renderer-backgrounding", true),
		chromedp.Flag("disable-sync", true),
		chromedp.Flag("disable-breakpad", true),
		chromedp.Flag("disable-default-apps", true),
		chromedp.Flag("disable-infobars", true),
		chromedp.Flag("disable-popup-blocking", true),
		chromedp.Flag("disable-translate", true),
		chromedp.Flag("metrics-recording-only", true),
		chromedp.Flag("mute-audio", true),
		chromedp.Flag("window-size", "1920,1080"),
		chromedp.WindowSize(1920, 1080),
	}

	if useGPU {
		// Non-headless mode with real GPU
		opts = append(opts,
			chromedp.Flag("headless", false),
			chromedp.Flag("disable-gpu", false),
			chromedp.Flag("enable-gpu-rasterization", true),
			chromedp.Flag("enable-oop-rasterization", true),
			chromedp.Flag("enable-webgl", true),
			chromedp.Flag("ignore-gpu-blocklist", true),
		)
		log.Info().Msg("GPU acceleration enabled via Xvfb")
	} else {
		// Standard headless
		opts = append(opts,
			chromedp.Headless,
			chromedp.Flag("enable-gpu-rasterization", true),
			chromedp.Flag("enable-oop-rasterization", true),
		)
	}

	chromePath := m.getChromePath()
	if chromePath != "" {
		opts = append(opts, chromedp.ExecPath(chromePath))
	}

	allocCtx, allocCancel := chromedp.NewExecAllocator(context.Background(), opts...)
	browserCtx, browserCancel := chromedp.NewContext(allocCtx)

	if err := chromedp.Run(browserCtx, chromedp.Navigate("about:blank")); err != nil {
		browserCancel()
		allocCancel()
		if xvfbCmd != nil {
			xvfbCmd.Process.Kill()
		}
		return fmt.Errorf("start browser: %w", err)
	}

	m.browserCtx = browserCtx
	m.browserCancel = browserCancel
	m.allocCancel = allocCancel
	m.browserReady = true

	mode := "headless (no GPU)"
	if useGPU {
		mode = "GPU-accelerated (Xvfb)"
	}
	log.Info().Str("chrome", chromePath).Str("mode", mode).Msg("Shared Chrome browser started")
	return nil
}

// newIsolatedTab creates a new tab in an isolated browser context (separate cookies/storage).
func (m *HeadlessManager) newIsolatedTab() (ctx context.Context, cancel context.CancelFunc, err error) {
	if err := m.ensureBrowser(); err != nil {
		return nil, nil, err
	}

	c := chromedp.FromContext(m.browserCtx)
	browserExec := cdp.WithExecutor(m.browserCtx, c.Browser)

	// Create isolated browser context via browser-level connection
	bcID, err := target.CreateBrowserContext().WithDisposeOnDetach(true).Do(browserExec)
	if err != nil {
		return nil, nil, fmt.Errorf("create browser context: %w", err)
	}

	// Create target using minimal params (Chrome 146 compat)
	tid, err := (&minCreateTarget{URL: "about:blank", BCI: bcID}).Do(browserExec)
	if err != nil {
		target.DisposeBrowserContext(bcID).Do(browserExec)
		return nil, nil, fmt.Errorf("create target: %w", err)
	}

	// Attach chromedp context to the new target
	tabCtx, tabCancel := chromedp.NewContext(m.browserCtx, chromedp.WithTargetID(tid))

	// Wrap cancel to also dispose the browser context
	combinedCancel := func() {
		tabCancel()
		target.DisposeBrowserContext(bcID).Do(browserExec)
	}

	return tabCtx, combinedCancel, nil
}

// OnClientDisconnected cleans up any headless session associated with the disconnected client.
func (m *HeadlessManager) OnClientDisconnected(clientID string) {
	m.mu.Lock()
	for id, s := range m.sessions {
		if s.ClientID == clientID {
			log.Info().Str("sessionId", id).Str("clientId", clientID).Msg("Cleaning up headless session for disconnected client")
			s.ContextCancel()
			delete(m.sessions, id)
			if m.redis != nil && m.redis.IsConnected() {
				ctx := context.Background()
				m.redis.SafeDel(ctx, fmt.Sprintf("headless_session:%s", id))
				m.redis.SafeDel(ctx, fmt.Sprintf("headless_client:%s", clientID))
			}
		}
	}
	m.mu.Unlock()
}

// LaunchSession creates a new isolated browser context, logs into Foundry, and waits for the client to connect.
func (m *HeadlessManager) LaunchSession(apiKey, foundryURL, username, password, worldName string) (sessionID, clientID string, err error) {
	// Clean up stale sessions
	m.mu.Lock()
	for id, s := range m.sessions {
		if m.clientManager.GetClient(s.ClientID) == nil {
			log.Info().Str("sessionId", id).Str("clientId", s.ClientID).Msg("Removing stale headless session")
			s.ContextCancel()
			delete(m.sessions, id)
		}
	}
	count := 0
	for _, s := range m.sessions {
		if s.APIKey == apiKey {
			count++
		}
	}
	m.mu.Unlock()

	if count >= m.maxSessions {
		return "", "", fmt.Errorf("maximum headless sessions (%d) reached for this API key", m.maxSessions)
	}

	// Create a new isolated tab (shared browser, isolated cookies)
	tabCtx, tabCancel, err := m.newIsolatedTab()
	if err != nil {
		return "", "", fmt.Errorf("create isolated tab: %w", err)
	}

	// Set up console capture for this context
	logFile := setupBrowserConsoleCapture(tabCtx, username, foundryURL)
	if logFile != "" {
		log.Info().Str("logFile", logFile).Msg("Browser console logging enabled")
	}

	log.Info().Str("url", foundryURL).Str("username", username).Msg("Launching headless session (new browser context)")

	// Inject WebGL override before navigation
	chromedp.Run(tabCtx, chromedp.ActionFunc(func(ctx context.Context) error {
		p := page.AddScriptToEvaluateOnNewDocument(`
			const origGetContext = HTMLCanvasElement.prototype.getContext;
			HTMLCanvasElement.prototype.getContext = function(type, attrs) {
				const ctx = origGetContext.call(this, type, attrs);
				if (ctx && (type === 'webgl' || type === 'webgl2')) {
					const origGetParam = ctx.getParameter.bind(ctx);
					ctx.getParameter = function(param) {
						if (param === 0x9246) return 'ANGLE (Google, Vulkan 1.3.0, SwiftShader)';
						if (param === 0x9245) return 'Google Inc. (Google)';
						return origGetParam(param);
					};
				}
				return ctx;
			};
		`)
		_, err := p.Do(ctx)
		return err
	}))

	// Set viewport
	chromedp.Run(tabCtx, chromedp.EmulateViewport(1920, 1080))

	// Navigate
	navCtx, navCancel := context.WithTimeout(tabCtx, browserNavigateTimeout)
	defer navCancel()

	if err := chromedp.Run(navCtx, chromedp.Navigate(foundryURL)); err != nil {
		tabCancel()
		return "", "", fmt.Errorf("navigate to Foundry: %w", err)
	}

	// Dismiss notifications (quick — GPU mode loads fast)
	time.Sleep(300 * time.Millisecond)
	chromedp.Run(tabCtx, chromedp.Evaluate(`
		document.querySelectorAll('.notification .close, .notification a.close, #notifications .notification .close, .notification-pip, .notification .notification-close').forEach(el => el.click());
	`, nil))

	// World selection — only if we're on the world list page, not the login page
	if worldName != "" {
		// Check which page we're on: world list or login
		pageType := detectPage(tabCtx)
		if pageType == "worldList" {
			log.Info().Str("world", worldName).Msg("Selecting world")
			if err := selectWorld(tabCtx, worldName); err != nil {
				log.Warn().Err(err).Msg("World selection failed, continuing")
			}
		} else {
			log.Info().Str("pageType", pageType).Msg("Already past world selection, skipping")
		}
	}

	// Snapshot clients before login
	existingClients := make(map[string]bool)
	for _, cid := range m.clientManager.GetConnectedClients(apiKey) {
		existingClients[cid] = true
	}

	// Login
	log.Info().Str("username", username).Msg("Logging in")
	userID, err := loginToFoundry(tabCtx, username, password)
	if err != nil {
		tabCancel()
		return "", "", fmt.Errorf("login failed: %w", err)
	}

	// Wait for game canvas
	log.Info().Msg("Waiting for game canvas")
	loadCtx, loadCancel := context.WithTimeout(tabCtx, gameLoadTimeout)
	defer loadCancel()

	err = waitForAnySelector(loadCtx, []string{"#ui-left", "#sidebar", "#game", ".vtt"})
	if err != nil {
		// Debug screenshot
		var screenshot []byte
		if screenshotErr := chromedp.Run(tabCtx, chromedp.CaptureScreenshot(&screenshot)); screenshotErr == nil && len(screenshot) > 0 {
			debugPath := "data/headless-debug.png"
			os.WriteFile(debugPath, screenshot, 0644)
			log.Warn().Str("screenshot", debugPath).Msg("Saved debug screenshot")
		}
		var pageURL, pageTitle string
		chromedp.Run(tabCtx, chromedp.Location(&pageURL), chromedp.Title(&pageTitle))
		log.Warn().Str("url", pageURL).Str("title", pageTitle).Msg("Browser state at timeout")

		tabCancel()
		return "", "", fmt.Errorf("game canvas did not load: %w", err)
	}

	// Generate session ID
	sessionID = uuid.New().String()
	expectedClientID := fmt.Sprintf("foundry-%s", userID)

	// Register pending session
	m.mu.Lock()
	m.pending[sessionID] = &PendingHeadless{
		SessionID:        sessionID,
		ExpectedClientID: expectedClientID,
		APIKey:           apiKey,
		ContextCancel:    tabCancel,
		StartTime:        time.Now(),
	}
	m.mu.Unlock()

	log.Info().Str("sessionId", sessionID).Str("expectedClient", expectedClientID).Int("existingClients", len(existingClients)).Msg("Polling for client connection")

	// Poll for client connection
	pollCtx, pollCancel := context.WithTimeout(context.Background(), clientPollTimeout)
	defer pollCancel()

	ticker := time.NewTicker(clientPollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Priority 1: exact match
			if client := m.clientManager.GetClient(expectedClientID); client != nil && client.APIKey() == apiKey {
				connectedClientID := expectedClientID
				m.mu.Lock()
				delete(m.pending, sessionID)
				m.sessions[sessionID] = &HeadlessSession{
					SessionID: sessionID, ClientID: connectedClientID, APIKey: apiKey,
					FoundryURL: foundryURL, Username: username, WorldName: worldName,
					ContextCancel: tabCancel, StartedAt: time.Now(), LastActivity: time.Now(),
				}
				m.mu.Unlock()
				m.registerInRedis(sessionID, connectedClientID)
				log.Info().Str("sessionId", sessionID).Str("clientId", connectedClientID).Msg("Headless session established (exact match)")
				return sessionID, connectedClientID, nil
			}

			// Priority 2: new client whose ID contains the userId
			currentClients := m.clientManager.GetConnectedClients(apiKey)
			for _, cid := range currentClients {
				if !existingClients[cid] && strings.Contains(cid, userID) {
					connectedClientID := cid
					m.mu.Lock()
					delete(m.pending, sessionID)
					m.sessions[sessionID] = &HeadlessSession{
						SessionID: sessionID, ClientID: connectedClientID, APIKey: apiKey,
						FoundryURL: foundryURL, Username: username, WorldName: worldName,
						ContextCancel: tabCancel, StartedAt: time.Now(), LastActivity: time.Now(),
					}
					m.mu.Unlock()
					m.registerInRedis(sessionID, connectedClientID)
					log.Info().Str("sessionId", sessionID).Str("clientId", connectedClientID).Msg("Headless session established (userId match)")
					return sessionID, connectedClientID, nil
				}
			}
		case <-pollCtx.Done():
			m.mu.Lock()
			delete(m.pending, sessionID)
			m.mu.Unlock()
			tabCancel()
			return "", "", fmt.Errorf("client connection timed out after %s", clientPollTimeout)
		}
	}
}

func (m *HeadlessManager) registerInRedis(sessionID, clientID string) {
	if m.redis != nil && m.redis.IsConnected() {
		ctx := context.Background()
		ttl := 3 * time.Hour
		m.redis.SafeSet(ctx, fmt.Sprintf("headless_session:%s", sessionID), clientID, ttl)
		m.redis.SafeSet(ctx, fmt.Sprintf("headless_client:%s", clientID), sessionID, ttl)
	}
}

// EndSession closes an isolated browser context for a session.
func (m *HeadlessManager) EndSession(sessionID string) error {
	m.mu.Lock()

	if s, ok := m.sessions[sessionID]; ok {
		delete(m.sessions, sessionID)
		m.mu.Unlock()
		s.ContextCancel() // Close the tab/context, not the whole browser
		m.clientManager.RemoveClient(s.ClientID)
		if m.redis != nil && m.redis.IsConnected() {
			ctx := context.Background()
			m.redis.SafeDel(ctx, fmt.Sprintf("headless_session:%s", sessionID))
			m.redis.SafeDel(ctx, fmt.Sprintf("headless_client:%s", s.ClientID))
		}
		log.Info().Str("sessionId", sessionID).Msg("Headless session ended")
		return nil
	}

	if p, ok := m.pending[sessionID]; ok {
		delete(m.pending, sessionID)
		m.mu.Unlock()
		p.ContextCancel()
		log.Info().Str("sessionId", sessionID).Msg("Pending headless session ended")
		return nil
	}

	m.mu.Unlock()
	return fmt.Errorf("session not found: %s", sessionID)
}

// ListSessions returns info about all active sessions.
func (m *HeadlessManager) ListSessions() []SessionInfo {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var infos []SessionInfo
	for _, s := range m.sessions {
		infos = append(infos, SessionInfo{
			SessionID: s.SessionID, ClientID: s.ClientID, FoundryURL: s.FoundryURL,
			Username: s.Username, WorldName: s.WorldName,
			StartedAt: s.StartedAt.UnixMilli(), LastActivity: s.LastActivity.UnixMilli(),
		})
	}
	return infos
}

// ValidateHeadlessSession checks if a client ID belongs to a valid headless session.
func (m *HeadlessManager) ValidateHeadlessSession(clientID, token string) (bool, error) {
	if !strings.HasPrefix(clientID, "foundry-") {
		return true, nil
	}
	if m.redis != nil && m.redis.IsConnected() {
		ctx := context.Background()
		sessionID, err := m.redis.SafeGet(ctx, fmt.Sprintf("headless_client:%s", clientID))
		if err != nil || sessionID == "" {
			return true, nil
		}
		ttl := 3 * time.Hour
		m.redis.SafeExpire(ctx, fmt.Sprintf("headless_session:%s", sessionID), ttl)
		m.redis.SafeExpire(ctx, fmt.Sprintf("headless_client:%s", clientID), ttl)
		return true, nil
	}
	return true, nil
}

// StartHeadlessWithStoredCredentials launches a session using encrypted credentials.
func (m *HeadlessManager) StartHeadlessWithStoredCredentials(scopedKeyID int64, masterAPIKey string, db *database.DB, cfg *config.Config, worldName string) (string, error) {
	ctx := context.Background()
	key, err := db.ApiKeyStore().FindByID(ctx, scopedKeyID)
	if err != nil || key == nil {
		return "", fmt.Errorf("scoped key not found")
	}
	if !key.HasStoredCredentials() {
		return "", fmt.Errorf("scoped key has no stored credentials")
	}

	password, err := service.Decrypt(
		key.EncryptedFoundryPassword.String, key.PasswordIV.String,
		key.PasswordAuthTag.String, cfg.CredentialsEncryptionKey,
	)
	if err != nil {
		return "", fmt.Errorf("decrypt credentials: %w", err)
	}

	// Check if already connected
	clients := m.clientManager.GetConnectedClients(masterAPIKey)
	for _, cid := range clients {
		if m.clientManager.GetClient(cid) != nil {
			return cid, nil
		}
	}

	_, clientID, err := m.LaunchSession(masterAPIKey, key.FoundryURL.String, key.FoundryUsername.String, password, worldName)
	if err != nil {
		return "", err
	}
	return clientID, nil
}

// Shutdown closes all sessions and the shared browser.
func (m *HeadlessManager) Shutdown() {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, s := range m.sessions {
		s.ContextCancel()
	}
	for _, p := range m.pending {
		p.ContextCancel()
	}
	m.sessions = make(map[string]*HeadlessSession)
	m.pending = make(map[string]*PendingHeadless)
	if m.browserReady {
		m.browserCancel()
		m.allocCancel()
		m.browserReady = false
	}
	if m.xvfbCmd != nil {
		m.xvfbCmd.Process.Kill()
		m.xvfbCmd = nil
	}
	log.Info().Msg("All headless sessions and browser stopped")
}

// StartCleanupLoop starts goroutines that clean up inactive and orphaned sessions.
func (m *HeadlessManager) StartCleanupLoop(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(60 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				m.cleanupInactive()
			case <-ctx.Done():
				return
			}
		}
	}()
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				m.cleanupOrphanedPending()
			case <-ctx.Done():
				return
			}
		}
	}()
}

func (m *HeadlessManager) cleanupInactive() {
	m.mu.Lock()
	now := time.Now()
	var toEnd []string
	for id, s := range m.sessions {
		if now.Sub(s.LastActivity) > headlessInactiveTimeout {
			toEnd = append(toEnd, id)
		}
	}
	m.mu.Unlock()
	for _, id := range toEnd {
		log.Info().Str("sessionId", id).Msg("Cleaning up inactive headless session")
		m.EndSession(id)
	}
}

func (m *HeadlessManager) cleanupOrphanedPending() {
	m.mu.Lock()
	now := time.Now()
	var toEnd []string
	for id, p := range m.pending {
		if now.Sub(p.StartTime) > pendingSessionTimeout {
			toEnd = append(toEnd, id)
		}
	}
	m.mu.Unlock()
	for _, id := range toEnd {
		log.Info().Str("sessionId", id).Msg("Cleaning up orphaned pending session")
		m.EndSession(id)
	}
}

// --- Helper functions ---

func setupBrowserConsoleCapture(ctx context.Context, username, foundryURL string) string {
	captureLevel := os.Getenv("CAPTURE_BROWSER_CONSOLE")
	if captureLevel == "" {
		return ""
	}
	os.MkdirAll("data/browser-logs", 0755)
	timestamp := time.Now().Format("2006-01-02T15-04-05")
	filename := fmt.Sprintf("data/browser-logs/headless_%s_%s.log", username, timestamp)
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to create browser log file")
		return ""
	}

	chromedp.ListenTarget(ctx, func(ev interface{}) {
		switch ev := ev.(type) {
		case *runtime.EventConsoleAPICalled:
			level := string(ev.Type)
			switch captureLevel {
			case "error":
				if level != "error" {
					return
				}
			case "warn":
				if level != "error" && level != "warning" {
					return
				}
			}
			var args []string
			for _, arg := range ev.Args {
				val := string(arg.Value)
				if val == "" {
					val = arg.Description
				}
				if len(val) > 1 && val[0] == '"' {
					val = val[1 : len(val)-1]
				}
				args = append(args, val)
			}
			fmt.Fprintf(f, "[%s] [%s] %s\n", time.Now().Format(time.RFC3339), level, strings.Join(args, " "))
		case *runtime.EventExceptionThrown:
			if ev.ExceptionDetails != nil {
				text := ev.ExceptionDetails.Text
				if ev.ExceptionDetails.Exception != nil {
					text += ": " + ev.ExceptionDetails.Exception.Description
				}
				fmt.Fprintf(f, "[%s] [exception] %s\n", time.Now().Format(time.RFC3339), text)
			}
		}
	})
	chromedp.Run(ctx, runtime.Enable())
	return filename
}

func waitForAnySelector(ctx context.Context, selectors []string) error {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			for _, sel := range selectors {
				var nodes []*cdp.Node
				err := chromedp.Run(ctx, chromedp.Nodes(sel, &nodes, chromedp.ByQuery, chromedp.AtLeast(0)))
				if err == nil && len(nodes) > 0 {
					log.Info().Str("selector", sel).Msg("Game canvas detected")
					return nil
				}
			}
		}
	}
}

// detectPage checks whether we're on the world list, login page, or game canvas.
func detectPage(ctx context.Context) string {
	checkCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	var result string
	err := chromedp.Run(checkCtx, chromedp.Evaluate(`
		(function() {
			if (document.querySelector('input[name="password"]')) return 'login';
			if (document.querySelector('li.package.world')) return 'worldList';
			if (document.querySelector('#ui-left, #sidebar, #game')) return 'game';
			return 'unknown';
		})()
	`, &result))
	if err != nil {
		return "unknown"
	}
	return result
}

func selectWorld(ctx context.Context, worldName string) error {
	selCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	if err := chromedp.Run(selCtx, chromedp.WaitVisible(`li.package.world`, chromedp.ByQuery)); err != nil {
		return fmt.Errorf("world list not found: %w", err)
	}
	js := fmt.Sprintf(`
		(function() {
			const worlds = document.querySelectorAll('li.package.world');
			for (const w of worlds) {
				const title = w.querySelector('h3.package-title, .package-title');
				if (title && title.textContent.trim().toLowerCase() === %q) {
					const playBtn = w.querySelector('a.control.play, button.control.play');
					if (playBtn) { playBtn.click(); return 'clicked'; }
				}
			}
			return 'not_found';
		})()
	`, strings.ToLower(worldName))
	var result string
	if err := chromedp.Run(selCtx, chromedp.Evaluate(js, &result)); err != nil {
		return fmt.Errorf("world selection eval: %w", err)
	}
	if result != "clicked" {
		return fmt.Errorf("world %q not found", worldName)
	}
	time.Sleep(2 * time.Second)
	return nil
}

func loginToFoundry(ctx context.Context, username, password string) (string, error) {
	loginCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	if err := chromedp.Run(loginCtx, chromedp.WaitVisible(`input[name="password"]`, chromedp.ByQuery)); err != nil {
		return "", fmt.Errorf("login form not found: %w", err)
	}

	js := fmt.Sprintf(`
		(function() {
			var userId = %q;
			const sel = document.querySelector('select[name="userid"]');
			if (sel) {
				const options = Array.from(sel.options);
				const match = options.find(o => o.textContent.trim().toLowerCase() === %q);
				if (match) {
					sel.value = match.value;
					sel.dispatchEvent(new Event('change', {bubbles: true}));
					userId = match.value;
				}
			}
			const pwInput = document.querySelector('input[name="password"]');
			if (pwInput) {
				const nativeInputValueSetter = Object.getOwnPropertyDescriptor(window.HTMLInputElement.prototype, 'value').set;
				nativeInputValueSetter.call(pwInput, %q);
				pwInput.dispatchEvent(new Event('input', {bubbles: true}));
				pwInput.dispatchEvent(new Event('change', {bubbles: true}));
			}
			const submitBtn = document.querySelector('button[type="submit"], button[name="join"], form button');
			if (submitBtn) { submitBtn.click(); }
			else { const form = document.querySelector('form'); if (form) form.submit(); }
			return userId;
		})()
	`, username, strings.ToLower(username), password)

	var userID string
	if err := chromedp.Run(loginCtx, chromedp.Evaluate(js, &userID)); err != nil {
		return "", fmt.Errorf("login eval: %w", err)
	}
	if userID == "" {
		userID = username
	}
	log.Info().Str("userId", userID).Msg("Login submitted, waiting for page transition")
	time.Sleep(1 * time.Second)
	return userID, nil
}

func findChromeBinary() string {
	candidates := []string{
		"chromium-browser", "chromium", "google-chrome", "google-chrome-stable",
		"/snap/bin/chromium",
		"/usr/bin/chromium-browser", "/usr/bin/chromium",
		"/usr/bin/google-chrome", "/usr/bin/google-chrome-stable",
		"/Applications/Google Chrome.app/Contents/MacOS/Google Chrome",
		"/Applications/Chromium.app/Contents/MacOS/Chromium",
	}
	for _, c := range candidates {
		if path, err := exec.LookPath(c); err == nil {
			log.Info().Str("path", path).Msg("Auto-detected Chrome/Chromium binary")
			return path
		}
	}
	return ""
}

// Unused import guards
var _ = rand.Read
var _ = rsa.GenerateKey
var _ = target.CreateBrowserContext
