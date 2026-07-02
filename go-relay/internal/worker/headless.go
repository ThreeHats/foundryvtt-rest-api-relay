package worker

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	goruntime "runtime"
	"strings"
	"sync"
	"time"

	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/alerts"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/config"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/database"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/model"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/service"
	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/ws"
	cdp "github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/cdproto/target"
	"github.com/chromedp/chromedp"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

const (
	pendingSessionTimeout   = 5 * time.Minute
	clientPollInterval      = 500 * time.Millisecond
	clientPollTimeout       = 5 * time.Minute
	autoStartPollTimeout    = 60 * time.Second // tighter poll for auto-start path
	autoStartCooldownOnFail = 2 * time.Minute  // backoff after a failed auto-start
	browserNavigateTimeout  = 180 * time.Second
	gameLoadTimeout         = 30 * time.Second
)

// webGLSpoof is injected via Page.addScriptToEvaluateOnNewDocument before Foundry
// loads. It handles two cases:
//  1. Real WebGL context exists (SwiftShader/Mesa) — overrides renderer strings so
//     Foundry's hardware-acceleration check doesn't see "SwiftShader"/"llvmpipe".
//  2. WebGL context is null (unavailable in this container) — returns a full stub
//     context so PIXI.js doesn't crash before the login form renders.
//
// Foundry checks gl.RENDERER (0x1F01), gl.VENDOR (0x1F00), and the
// WEBGL_debug_renderer_info unmasked values (0x9246 / 0x9245) for known software
// renderer substrings ("SwiftShader", "llvmpipe", "softpipe", etc.).
const webGLSpoof = `
	(function() {
		const origGetContext = HTMLCanvasElement.prototype.getContext;
		HTMLCanvasElement.prototype.getContext = function(type, attrs) {
			var ctx = origGetContext.call(this, type, attrs);
			if (type === 'webgl' || type === 'webgl2') {
				if (!ctx) {
					ctx = {
						canvas: this,
						drawingBufferWidth: this.width || 1280,
						drawingBufferHeight: this.height || 720,
						getParameter: function(p) {
							switch (p) {
								case 0x9246: return 'ANGLE (Intel, Intel(R) UHD Graphics 620 Direct3D11 vs_5_0 ps_5_0, D3D11-27.20.100.8681)';
								case 0x9245: return 'Google Inc. (Intel)';
								case 0x1F01: return 'ANGLE (Intel, Intel(R) UHD Graphics 620)';
								case 0x1F00: return 'Intel';
								case 0x8B8C: return 16; // MAX_VERTEX_ATTRIBS
								case 0x8869: return 16; // MAX_VERTEX_UNIFORM_VECTORS
								case 0x8872: return 16; // MAX_TEXTURE_IMAGE_UNITS
								case 0x8B8D: return 16; // MAX_FRAGMENT_UNIFORM_VECTORS
							}
							return null;
						},
						getExtension: function() { return null; },
						getShaderPrecisionFormat: function() { return { rangeMin: 127, rangeMax: 127, precision: 23 }; },
						getSupportedExtensions: function() { return []; },
						createShader: function() { return {}; },
						shaderSource: function() {}, compileShader: function() {},
						getShaderParameter: function() { return true; },
						getShaderInfoLog: function() { return ''; },
						createProgram: function() { return {}; },
						attachShader: function() {}, linkProgram: function() {},
						getProgramParameter: function() { return true; },
						getProgramInfoLog: function() { return ''; },
						useProgram: function() {},
						getAttribLocation: function() { return 0; },
						getUniformLocation: function() { return {}; },
						uniform1f: function() {}, uniform2f: function() {}, uniform3f: function() {}, uniform4f: function() {},
						uniform1i: function() {}, uniform2i: function() {}, uniform3i: function() {}, uniform4i: function() {},
						uniformMatrix2fv: function() {}, uniformMatrix3fv: function() {}, uniformMatrix4fv: function() {},
						bindBuffer: function() {}, bufferData: function() {}, bufferSubData: function() {},
						enableVertexAttribArray: function() {}, disableVertexAttribArray: function() {},
						vertexAttribPointer: function() {},
						drawArrays: function() {}, drawElements: function() {},
						enable: function() {}, disable: function() {}, blendFunc: function() {},
						clearColor: function() {}, clear: function() {},
						viewport: function() {}, scissor: function() {},
						createTexture: function() { return {}; }, deleteTexture: function() {},
						bindTexture: function() {}, texImage2D: function() {},
						texParameteri: function() {}, pixelStorei: function() {},
						createFramebuffer: function() { return {}; }, bindFramebuffer: function() {},
						framebufferTexture2D: function() {},
						checkFramebufferStatus: function() { return 0x8CD5; }, // FRAMEBUFFER_COMPLETE
						readPixels: function() {}, getError: function() { return 0; },
						activeTexture: function() {}, generateMipmap: function() {},
						createBuffer: function() { return {}; }, deleteBuffer: function() {},
						createRenderbuffer: function() { return {}; }, bindRenderbuffer: function() {},
						renderbufferStorage: function() {}, framebufferRenderbuffer: function() {},
						deleteFramebuffer: function() {}, deleteRenderbuffer: function() {},
						deleteShader: function() {}, deleteProgram: function() {},
						depthMask: function() {}, depthFunc: function() {}, cullFace: function() {},
						stencilFunc: function() {}, stencilOp: function() {}, colorMask: function() {},
						lineWidth: function() {}, polygonOffset: function() {},
						isContextLost: function() { return false; },
					};
					if (type === 'webgl2') {
						ctx.UNSIGNED_INT_24_8 = 0x84FA;
						ctx.READ_FRAMEBUFFER = 0x8CA8;
						ctx.DRAW_FRAMEBUFFER = 0x8CA9;
						ctx.DEPTH_COMPONENT24 = 0x81A6;
						ctx.DEPTH_STENCIL = 0x84F9;
						ctx.blitFramebuffer = function() {};
						ctx.renderbufferStorageMultisample = function() {};
						ctx.bindVertexArray = function() {};
						ctx.createVertexArray = function() { return {}; };
						ctx.deleteVertexArray = function() {};
					}
				} else {
					const origGetParam = ctx.getParameter.bind(ctx);
					ctx.getParameter = function(param) {
						switch (param) {
							case 0x9246: return 'ANGLE (Intel, Intel(R) UHD Graphics 620 Direct3D11 vs_5_0 ps_5_0, D3D11-27.20.100.8681)'; // UNMASKED_RENDERER_WEBGL
							case 0x9245: return 'Google Inc. (Intel)';                          // UNMASKED_VENDOR_WEBGL
							case 0x1F01: return 'ANGLE (Intel, Intel(R) UHD Graphics 620)';    // RENDERER
							case 0x1F00: return 'Intel';                                        // VENDOR
						}
						return origGetParam(param);
					};
				}
			}
			return ctx;
		};
	})();
`

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
	URL string               `json:"url"`
	BCI cdp.BrowserContextID `json:"browserContextId,omitempty"`
}

func (p *minCreateTarget) Do(ctx context.Context) (target.ID, error) {
	var res struct {
		TargetID target.ID `json:"targetId"`
	}
	err := cdp.Execute(ctx, "Target.createTarget", p, &res)
	return res.TargetID, err
}

// HeadlessDeps groups the dependencies the headless manager needs for the
// AutoStartForKnownClient flow (the new remote-request auto-start path).
// Wired by the server during initialization to avoid the import cycle the
// worker package would have if it imported database directly at the
// constructor level.
type HeadlessDeps struct {
	DB  *database.DB
	Cfg *config.Config
}

// HeadlessManager manages a shared Chrome browser with isolated contexts per session.
type HeadlessManager struct {
	mu              sync.RWMutex
	sessions        map[string]*HeadlessSession // sessionID -> session
	pending         map[string]*PendingHeadless // sessionID -> pending
	clientManager   *ws.ClientManager
	redis           *config.RedisClient
	maxSessions     int
	inactiveTimeout time.Duration
	chromePath      string
	resolvedChrome  string
	dataDir         string // absolute path to the data directory (for browser logs)
	userDataDir     string // persistent Chrome profile dir (V8 bytecode + HTTP cache)
	jsHeapMB        int    // max V8 old-space heap in MB
	windowWidth     int    // viewport width
	windowHeight    int    // viewport height
	enableSHM       bool   // allow Chrome to use /dev/shm
	renderMode      string // configured render mode (auto|gpu|xvfb|swiftshader)

	// headlessDeps is set after construction by SetDeps. Used by
	// AutoStartForKnownClient (the remote-request auto-start path) to look
	// up users, credentials, and persist headless connection tokens.
	headlessDeps *HeadlessDeps

	// Shared browser
	browserCtx    context.Context
	browserCancel context.CancelFunc
	allocCancel   context.CancelFunc
	browserReady  bool
	xvfbCmd       *exec.Cmd

	// Request queuing: tracks in-flight launches per clientID
	launchQueues   map[string]*launchQueue // clientID -> queue
	launchQueuesMu sync.Mutex

	// Failure cooldowns: blocks rapid retries after a failed auto-start
	failureCooldowns   map[string]autoStartCooldown
	failureCooldownsMu sync.Mutex
}

type autoStartCooldown struct {
	Until  time.Time
	Reason string // original failure message, surfaced in subsequent "on cooldown" errors
}

// SetDeps wires the database + config dependencies needed by
// AutoStartForKnownClient. Called by the server after constructing the
// HeadlessManager and the database.
func (m *HeadlessManager) SetDeps(deps *HeadlessDeps) {
	m.headlessDeps = deps
}

// launchQueue tracks an in-flight headless launch and waiters.
type launchQueue struct {
	done     chan string // closed when the launch resolves (success or failure)
	errVal   error       // set before done is closed on failure; safe for all waiters to read after <-done
	clientID string      // set before done is closed on success
}

// NewHeadlessManager creates a new headless session manager.
func NewHeadlessManager(clientManager *ws.ClientManager, redis *config.RedisClient, cfg *config.Config) *HeadlessManager {
	userDataDir := cfg.ChromeUserDataDir
	if userDataDir == "" {
		userDataDir = filepath.Join(cfg.DataDir, "chrome-profile")
	}
	// 0 means sessions never time out due to inactivity.
	var inactiveTimeout time.Duration
	if cfg.HeadlessInactiveTimeoutSecs > 0 {
		inactiveTimeout = time.Duration(cfg.HeadlessInactiveTimeoutSecs) * time.Second
	}
	return &HeadlessManager{
		sessions:         make(map[string]*HeadlessSession),
		pending:          make(map[string]*PendingHeadless),
		clientManager:    clientManager,
		redis:            redis,
		maxSessions:      cfg.MaxHeadlessSessions,
		inactiveTimeout:  inactiveTimeout,
		chromePath:       cfg.ChromePath,
		dataDir:          cfg.DataDir,
		userDataDir:      userDataDir,
		jsHeapMB:         cfg.ChromeJSHeapMB,
		windowWidth:      cfg.ChromeWindowWidth,
		windowHeight:     cfg.ChromeWindowHeight,
		enableSHM:        cfg.ChromeEnableSHM,
		renderMode:       cfg.ChromeGPUMode,
		launchQueues:     make(map[string]*launchQueue),
		failureCooldowns: make(map[string]autoStartCooldown),
	}
}

// IsLaunching returns true if a headless session is currently being launched for the given client.
func (m *HeadlessManager) IsLaunching(clientID string) bool {
	m.launchQueuesMu.Lock()
	defer m.launchQueuesMu.Unlock()
	_, exists := m.launchQueues[clientID]
	return exists
}

// tryRegisterLaunch atomically registers a new launch slot for clientID.
// Returns (true, nil) when this caller owns the slot.
// Returns (false, q) when a launch is already in flight — caller should wait on q.done.
func (m *HeadlessManager) tryRegisterLaunch(clientID string) (bool, *launchQueue) {
	m.launchQueuesMu.Lock()
	defer m.launchQueuesMu.Unlock()
	if q, exists := m.launchQueues[clientID]; exists {
		return false, q
	}
	q := &launchQueue{
		done: make(chan string, 1),
	}
	m.launchQueues[clientID] = q
	return true, nil
}

// CompleteLaunch signals that a headless launch completed successfully.
func (m *HeadlessManager) CompleteLaunch(clientID string, resultClientID string) {
	m.launchQueuesMu.Lock()
	q, exists := m.launchQueues[clientID]
	if exists {
		q.clientID = resultClientID
		close(q.done) // unblock all waiters
		delete(m.launchQueues, clientID)
	}
	m.launchQueuesMu.Unlock()
}

// FailLaunch signals that a headless launch failed.
func (m *HeadlessManager) FailLaunch(clientID string, err error) {
	m.launchQueuesMu.Lock()
	q, exists := m.launchQueues[clientID]
	if exists {
		q.errVal = err // set before close so every waiter sees it
		close(q.done)
		delete(m.launchQueues, clientID)
	}
	m.launchQueuesMu.Unlock()
}

// WaitForLaunch blocks until an in-flight headless launch completes.
// Returns the clientID or error. Timeout is 5 minutes.
func (m *HeadlessManager) WaitForLaunch(clientID string, timeout time.Duration) (string, error) {
	m.launchQueuesMu.Lock()
	q, exists := m.launchQueues[clientID]
	m.launchQueuesMu.Unlock()

	if !exists {
		return "", fmt.Errorf("no pending launch for client %s", clientID)
	}

	select {
	case <-q.done:
		if q.clientID != "" {
			return q.clientID, nil
		}
		if q.errVal != nil {
			return "", q.errVal
		}
		return "", fmt.Errorf("launch failed for client %s", clientID)
	case <-time.After(timeout):
		return "", fmt.Errorf("timed out waiting for headless session (client %s)", clientID)
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

// startXvfb starts a virtual framebuffer display for Chrome.
// Xvfb gives Chrome a non-headless environment for rendering but goes through
// Mesa software (llvmpipe/softpipe) — it does NOT access the real GPU.
func startXvfb(width, height int) (string, *exec.Cmd, error) {
	for display := 99; display < 120; display++ {
		displayStr := fmt.Sprintf(":%d", display)
		if _, err := os.Stat(fmt.Sprintf("/tmp/.X11-unix/X%d", display)); err == nil {
			continue // already in use
		}
		screenArg := fmt.Sprintf("%dx%dx24", width, height) // 24-bit required for OpenGL/GLX contexts
		cmd := exec.Command("Xvfb", displayStr, "-screen", "0", screenArg, "-ac", "-nolisten", "tcp", "-dpi", "96")
		if err := cmd.Start(); err != nil {
			continue
		}
		socketPath := fmt.Sprintf("/tmp/.X11-unix/X%d", display)
		started := false
		for i := 0; i < 10; i++ {
			if _, err := os.Stat(socketPath); err == nil {
				started = true
				break
			}
			time.Sleep(50 * time.Millisecond)
		}
		if !started {
			cmd.Process.Kill()
			continue
		}
		log.Info().Str("display", displayStr).Msg("Xvfb virtual display started")
		return displayStr, cmd, nil
	}
	return "", nil, fmt.Errorf("could not find available display for Xvfb")
}

// findNvidiaVulkanICD returns the path to NVIDIA's Vulkan ICD JSON file.
// When set as VK_ICD_FILENAMES, the Vulkan loader uses NVIDIA's GPU exclusively,
// avoiding Mesa software devices (llvmpipe, etc.) that appear in the default ICD scan.
// Returns "" if not found.
func findNvidiaVulkanICD() string {
	for _, p := range []string{
		"/etc/vulkan/icd.d/nvidia_icd.json",
		"/usr/share/vulkan/icd.d/nvidia_icd.json",
	} {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}
	return ""
}

// hasDRIAccess returns true if at least one DRI render node (/dev/dri/renderD*)
// exists and is readable by the current process.
func hasDRIAccess() bool {
	matches, _ := filepath.Glob("/dev/dri/renderD*")
	for _, path := range matches {
		if f, err := os.Open(path); err == nil {
			f.Close()
			return true
		}
	}
	return false
}

// hasNvidiaAccess returns true if nvidia-container-toolkit has exposed an NVIDIA
// GPU to this container. /dev/nvidiactl is only present when the toolkit has
// made GPU access available.
func hasNvidiaAccess() bool {
	_, err := os.Stat("/dev/nvidiactl")
	return err == nil
}

// resolveRenderMode returns the effective render mode for ensureBrowser.
// When configured to "auto" it detects the best available option:
//
//	nvidia       if /dev/nvidiactl is present (NVIDIA via nvidia-container-toolkit)
//	gpu          if /dev/dri/renderD* is readable (Intel/AMD via Ozone headless + ANGLE GL)
//	xvfb         if the Xvfb binary is in PATH
//	swiftshader  always available as final fallback
//
// NVIDIA is checked before DRI because modern NVIDIA drivers also expose
// /dev/dri/renderD* nodes (NVIDIA DRM); Mesa cannot drive these, so a
// DRI-first check would cause GPU mode to fail and fall through to SwiftShader.
func resolveRenderMode(configured string) string {
	if configured != "" && configured != "auto" {
		return configured
	}
	if goruntime.GOOS == "linux" && hasNvidiaAccess() {
		return "nvidia"
	}
	if goruntime.GOOS == "linux" && hasDRIAccess() {
		return "gpu"
	}
	if _, err := exec.LookPath("Xvfb"); err == nil {
		return "xvfb"
	}
	return "swiftshader"
}

// ensureBrowser starts the shared browser, selecting the best available render
// mode and falling back to SwiftShader if the primary mode fails.
func (m *HeadlessManager) ensureBrowser() error {
	if m.browserReady {
		return nil
	}
	mode := resolveRenderMode(m.renderMode)
	if err := m.startBrowserWithMode(mode); err != nil {
		if mode != "swiftshader" {
			log.Warn().Err(err).Str("attempted", mode).
				Msg("Browser failed to start, falling back to SwiftShader")
			return m.startBrowserWithMode("swiftshader")
		}
		return fmt.Errorf("start browser: %w", err)
	}
	return nil
}

// WarmUpBrowser eagerly starts the shared Chrome browser during server startup
// so the first session request doesn't pay the cold-start penalty (~1–3s).
// If Chrome fails to start it logs a warning but does not fatal — the server
// remains functional for non-headless requests.
func (m *HeadlessManager) WarmUpBrowser() {
	if err := m.ensureBrowser(); err != nil {
		log.Warn().Err(err).Msg("Chrome pre-launch failed; headless sessions will cold-start on first request")
	} else {
		log.Info().Msg("Chrome browser pre-launched successfully")
	}
}

// buildBaseOpts returns the Chrome allocator options common to all render modes.
func (m *HeadlessManager) buildBaseOpts() []chromedp.ExecAllocatorOption {
	opts := []chromedp.ExecAllocatorOption{
		chromedp.NoFirstRun,
		chromedp.NoDefaultBrowserCheck,
		chromedp.NoSandbox,
		chromedp.Flag("disable-setuid-sandbox", true),
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
		chromedp.Flag("autoplay-policy", "no-user-gesture-required"),
		chromedp.Flag("window-size", fmt.Sprintf("%d,%d", m.windowWidth, m.windowHeight)),
		chromedp.WindowSize(m.windowWidth, m.windowHeight),
		chromedp.UserDataDir(m.userDataDir),
		chromedp.Flag("js-flags", fmt.Sprintf("--max-old-space-size=%d", m.jsHeapMB)),
		chromedp.Flag("disk-cache-size", "209715200"),
		chromedp.Flag("no-zygote", true),
		chromedp.Flag("renderer-process-limit", "4"),
		chromedp.Flag("bwsi", true),
		chromedp.Flag("disable-features", "MediaRouter,DialMediaRouteProvider,Translate,OptimizationHints,InterestCohortAPI"),
	}
	if !m.enableSHM {
		opts = append(opts, chromedp.Flag("disable-dev-shm-usage", true))
	}
	return opts
}

// startBrowserWithMode launches Chrome with the given render mode's flags,
// waits for a blank navigation to confirm it started, and stores the contexts.
func (m *HeadlessManager) startBrowserWithMode(mode string) error {
	opts := m.buildBaseOpts()

	switch mode {
	case "nvidia":
		// ANGLE Vulkan + native system Vulkan ICD for hardware NVIDIA rendering.
		//
		// --use-angle=vulkan: ANGLE uses Vulkan as its rendering backend.
		// --use-vulkan=native: Chrome uses the system Vulkan loader (libvulkan.so.1)
		//   instead of its bundled SwiftShader Vulkan (libvk_swiftshader.so).
		// --ozone-platform=headless: surfaceless EGL mode; no X11 display needed.
		//
		// The system Vulkan loader reads NVIDIA's ICD (nvidia_icd.json which
		// resolves to libGLX_nvidia.so.0 exporting vk_icdGetInstanceProcAddr).
		// VK_ICD_FILENAMES restricts the Vulkan loader to the NVIDIA ICD so Mesa
		// software devices (llvmpipe, etc.) are not enumerated alongside it.
		//
		// Requires: NVIDIA_DRIVER_CAPABILITIES=graphics (set in Dockerfile ENV)
		// so libGLX_nvidia.so.0 exports vk_icdGetInstanceProcAddr.
		if icd := findNvidiaVulkanICD(); icd != "" {
			os.Setenv("VK_ICD_FILENAMES", icd)
		}
		opts = append(opts,
			chromedp.Headless,
			chromedp.Flag("ozone-platform", "headless"),
			chromedp.Flag("use-gl", "angle"),
			chromedp.Flag("use-angle", "vulkan"),
			chromedp.Flag("use-vulkan", "native"),
			chromedp.Flag("ignore-gpu-blocklist", true),
			chromedp.Flag("enable-gpu-rasterization", true),
			chromedp.Flag("enable-oop-rasterization", true),
			chromedp.Flag("enable-webgl", true),
			chromedp.Flag("enable-zero-copy", true),
		)
		log.Info().Str("vk_icd", os.Getenv("VK_ICD_FILENAMES")).
			Msg("NVIDIA GPU mode: ANGLE Vulkan + native Vulkan ICD (hardware NVIDIA rendering)")

	case "gpu":
		// Ozone headless + ANGLE GL: real GPU via Mesa/DRI without a display server.
		// --ozone-platform=headless switches Chrome's graphics stack to EGL
		// surfaceless mode; --use-gl=angle --use-angle=gl routes through ANGLE
		// which handles EGL surface setup and calls Mesa DRI → GPU.
		// Old headless mode (chromedp.Headless) preserves ChromeDP automation
		// compatibility — --headless=new has different page-lifecycle semantics.
		opts = append(opts,
			chromedp.Headless,
			chromedp.Flag("ozone-platform", "headless"),
			chromedp.Flag("use-gl", "angle"),
			chromedp.Flag("use-angle", "gl"),
			chromedp.Flag("enable-gpu-rasterization", true),
			chromedp.Flag("enable-oop-rasterization", true),
			chromedp.Flag("enable-webgl", true),
			chromedp.Flag("ignore-gpu-blocklist", true),
			chromedp.Flag("enable-zero-copy", true),
		)
		log.Info().Msg("GPU mode: Ozone headless + ANGLE GL (Mesa/DRI)")

	case "xvfb":
		display, cmd, err := startXvfb(m.windowWidth, m.windowHeight)
		if err != nil {
			return err
		}
		os.Setenv("DISPLAY", display)
		m.xvfbCmd = cmd
		opts = append(opts,
			chromedp.Flag("headless", false),
			chromedp.Flag("disable-gpu", false),
			chromedp.Flag("enable-gpu-rasterization", true),
			chromedp.Flag("enable-oop-rasterization", true),
			chromedp.Flag("enable-webgl", true),
			chromedp.Flag("ignore-gpu-blocklist", true),
		)
		log.Info().Str("display", display).Msg("Xvfb mode: Mesa software rendering")

	default: // swiftshader
		opts = append(opts,
			chromedp.Headless,
			chromedp.Flag("enable-gpu-rasterization", true),
			chromedp.Flag("enable-oop-rasterization", true),
			chromedp.Flag("use-gl", "swiftshader"),
			chromedp.Flag("use-angle", "swiftshader"),
			chromedp.Flag("enable-unsafe-swiftshader", true),
			chromedp.Flag("enable-webgl", true),
			chromedp.Flag("ignore-gpu-blocklist", true),
		)
		log.Info().Msg("SwiftShader mode: software WebGL")
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
		if m.xvfbCmd != nil {
			m.xvfbCmd.Process.Kill()
			m.xvfbCmd = nil
		}
		return err
	}

	m.browserCtx = browserCtx
	m.browserCancel = browserCancel
	m.allocCancel = allocCancel
	m.browserReady = true

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
				m.redis.SafeSRem(ctx, "headless:global_sessions", id)
			}
		}
	}
	m.mu.Unlock()
}

// LaunchSession creates a new isolated browser context, logs into Foundry, and waits for the client to connect.
// injectConnectionTokenSeed installs a Page.addScriptToEvaluateOnNewDocument
// script that seeds the Foundry module's connection token into localStorage
// BEFORE any page scripts run. This ensures the module's `init` hook reads
// the seeded token instead of finding an empty client-scope setting and
// reporting "not paired."
//
// The localStorage key is `foundry-rest-api.connectionToken` and the value is
// JSON.stringify(rawToken) — matching Foundry v13's client-settings storage
// format (verified at client/helpers/client-settings.mjs:266 and :349).
//
// Call this AFTER creating the isolated tab context but BEFORE navigating to
// the Foundry URL.
func injectConnectionTokenSeed(tabCtx context.Context, rawToken string) error {
	if rawToken == "" {
		return nil
	}
	tokenJSON, err := json.Marshal(rawToken)
	if err != nil {
		return fmt.Errorf("marshal seed token: %w", err)
	}
	seedScript := fmt.Sprintf(`
		try {
			window.localStorage.setItem("foundry-rest-api.connectionToken", JSON.stringify(%s));
		} catch (e) {
			console.error("[REST API headless] failed to seed connection token:", e);
		}
	`, string(tokenJSON))
	return chromedp.Run(tabCtx, chromedp.ActionFunc(func(ctx context.Context) error {
		p := page.AddScriptToEvaluateOnNewDocument(seedScript)
		_, err := p.Do(ctx)
		return err
	}))
}

// LaunchSession launches a headless Foundry session. If seedToken is non-empty,
// the connection token is injected into localStorage before navigation so the
// Foundry module can connect without a manual pairing flow.
func (m *HeadlessManager) LaunchSession(apiKey, foundryURL, username, password, worldName, seedToken string) (sessionID, clientID string, err error) {
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

	if m.maxSessions > 0 && count >= m.maxSessions {
		if alerts.Track("headless_limit", 3, 10*time.Minute, 1*time.Hour) {
			alerts.Fire(alerts.Event{
				Type:     alerts.TypeHeadlessSessionFlood,
				Severity: "warning",
				Message:  "Headless session limit hit 3+ times in 10 minutes",
				Details:  map[string]interface{}{"maxSessions": m.maxSessions},
			})
		}
		return "", "", fmt.Errorf("maximum headless sessions (%d) reached for this API key", m.maxSessions)
	}

	// Global cap check across all relay instances via Redis
	if err := m.checkGlobalSessionCap(); err != nil {
		return "", "", err
	}

	// Create a new isolated tab (shared browser, isolated cookies)
	tabCtx, tabCancel, err := m.newIsolatedTab()
	if err != nil {
		return "", "", fmt.Errorf("create isolated tab: %w", err)
	}

	// Set up console capture for this context
	logFile := setupBrowserConsoleCapture(tabCtx, m.dataDir, username, "", foundryURL)
	if logFile != "" {
		log.Info().Str("logFile", logFile).Msg("Browser console logging enabled")
	}

	log.Info().Str("url", foundryURL).Str("username", username).Msg("Launching headless session (new browser context)")

	// Inject WebGL override before navigation
	chromedp.Run(tabCtx, chromedp.ActionFunc(func(ctx context.Context) error {
		_, err := page.AddScriptToEvaluateOnNewDocument(webGLSpoof).Do(ctx)
		return err
	}))

	// Inject connection token into localStorage BEFORE navigation so the
	// Foundry module reads it from its client-scope setting on init.
	if seedToken != "" {
		if err := injectConnectionTokenSeed(tabCtx, seedToken); err != nil {
			tabCancel()
			return "", "", fmt.Errorf("inject connection token seed: %w", err)
		}
		log.Info().Msg("Seeded connection token into headless browser localStorage")
	}

	// Set viewport
	chromedp.Run(tabCtx, chromedp.EmulateViewport(int64(m.windowWidth), int64(m.windowHeight)))

	// Navigate
	navCtx, navCancel := context.WithTimeout(tabCtx, browserNavigateTimeout)
	defer navCancel()

	if err := chromedp.Run(navCtx, chromedp.Navigate(foundryURL)); err != nil {
		tabCancel()
		return "", "", fmt.Errorf("navigate to Foundry: %w", err)
	}

	// Dismiss any notifications that appeared during navigation
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
	_, err = loginToFoundry(tabCtx, username, password)
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

	// Register pending session (no predicted client ID — the Foundry module
	// always connects with its paired fvtt_* ID which we can't know in advance)
	m.mu.Lock()
	m.pending[sessionID] = &PendingHeadless{
		SessionID:     sessionID,
		APIKey:        apiKey,
		ContextCancel: tabCancel,
		StartTime:     time.Now(),
	}
	m.mu.Unlock()

	log.Info().Str("sessionId", sessionID).Int("existingClients", len(existingClients)).Msg("Polling for client connection")

	// Poll for client connection
	pollCtx, pollCancel := context.WithTimeout(context.Background(), clientPollTimeout)
	defer pollCancel()

	ticker := time.NewTicker(clientPollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Accept any new client on this API key that wasn't connected before we launched the browser
			currentClients := m.clientManager.GetConnectedClients(apiKey)
			for _, cid := range currentClients {
				if !existingClients[cid] {
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
					log.Info().Str("sessionId", sessionID).Str("clientId", connectedClientID).Msg("Headless session established")
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
		// Track in global set for multi-instance session cap
		m.redis.SafeSAdd(ctx, "headless:global_sessions", sessionID)
		m.redis.SafeExpire(ctx, "headless:global_sessions", ttl)
	}
}

// checkGlobalSessionCap returns an error if the Redis-backed global headless
// session count is at or above maxSessions. No-op when Redis is unavailable.
func (m *HeadlessManager) checkGlobalSessionCap() error {
	if m.redis == nil || !m.redis.IsConnected() {
		return nil
	}
	ctx := context.Background()
	globalCount, _ := m.redis.SafeSCard(ctx, "headless:global_sessions")
	if m.maxSessions > 0 && int(globalCount) >= m.maxSessions {
		return fmt.Errorf("maximum headless sessions (%d) reached globally across all relay instances", m.maxSessions)
	}
	return nil
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
			m.redis.SafeSRem(ctx, "headless:global_sessions", sessionID)
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
// It mints a fresh headless connection token and seeds it into the Chrome
// instance's localStorage so the Foundry module can connect without a
// manual pairing flow.
// AutoStartForKnownClient implements the ws.HeadlessAutoStarter interface.
// It launches a headless session for a specific (userId, clientId) pair, used
// by the remote-request handler when a target client is offline.
//
// Concurrent calls for the same clientId join the in-flight launch rather than
// spawning a parallel browser. A per-client cooldown (autoStartCooldownOnFail)
// blocks rapid retries after a failure.
func (m *HeadlessManager) AutoStartForKnownClient(ctx context.Context, userID int64, clientID string) (string, error) {
	// Reject immediately if we're still cooling down from a recent failure.
	m.failureCooldownsMu.Lock()
	if cd, ok := m.failureCooldowns[clientID]; ok {
		if time.Now().Before(cd.Until) {
			remaining := time.Until(cd.Until).Round(time.Second)
			reason := cd.Reason
			m.failureCooldownsMu.Unlock()
			return "", fmt.Errorf("%s (auto-start blocked for %s, retry after cooldown)", reason, remaining)
		}
		delete(m.failureCooldowns, clientID)
	}
	m.failureCooldownsMu.Unlock()

	// Deduplicate: if a launch is already in flight, wait for its result.
	registered, existingQ := m.tryRegisterLaunch(clientID)
	if !registered {
		log.Info().Str("clientId", clientID).Msg("Auto-start already in flight; joining result")
		select {
		case <-existingQ.done:
			if existingQ.clientID != "" {
				return existingQ.clientID, nil
			}
			if existingQ.errVal != nil {
				return "", existingQ.errVal
			}
			return "", fmt.Errorf("in-flight auto-start failed for %s", clientID)
		case <-ctx.Done():
			return "", ctx.Err()
		}
	}

	// We own the slot. Recover panics into an error so it's always resolved
	// below — a leaked launchQueues entry would block this client forever.
	resultID, err := func() (rid string, rerr error) {
		defer func() {
			if p := recover(); p != nil {
				rerr = fmt.Errorf("auto-start panicked: %v", p)
			}
		}()
		return m.doAutoStartForKnownClient(ctx, userID, clientID)
	}()
	if err != nil {
		m.FailLaunch(clientID, err)
		m.failureCooldownsMu.Lock()
		m.failureCooldowns[clientID] = autoStartCooldown{
			Until:  time.Now().Add(autoStartCooldownOnFail),
			Reason: err.Error(),
		}
		m.failureCooldownsMu.Unlock()
		return "", err
	}
	m.CompleteLaunch(clientID, resultID)
	return resultID, nil
}

// doAutoStartForKnownClient is the core auto-start logic.
//
// Flow:
//  1. Look up the user (FindByID) to get apiKeyHash for ClientManager registration
//  2. Look up the KnownClient row, verify ownership, get its credentialId
//  3. Decrypt the stored Foundry password
//  4. Mint a fresh ConnectionToken in the DB with source = "headless"
//  5. Spawn an isolated ChromeDP tab
//  6. Inject `Page.addScriptToEvaluateOnNewDocument` to seed localStorage
//  7. Navigate to Foundry, log in, wait for game canvas
//  8. Poll for the EXACT target clientId to appear in ClientManager
//  9. Return the connected clientId
//
// On any error, the seed token is revoked from the DB.
func (m *HeadlessManager) doAutoStartForKnownClient(ctx context.Context, userID int64, clientID string) (string, error) {
	if m.headlessDeps == nil || m.headlessDeps.DB == nil {
		return "", fmt.Errorf("headless manager has no DB wired; cannot auto-start")
	}
	db := m.headlessDeps.DB
	cfg := m.headlessDeps.Cfg

	log.Info().Str("clientId", clientID).Int64("userId", userID).Msg("Headless auto-start: resolving credentials for offline world")

	// 1. User lookup → apiKeyHash for ClientManager registration
	user, err := db.UserStore().FindByID(ctx, userID)
	if err != nil || user == nil {
		log.Warn().Err(err).Str("clientId", clientID).Int64("userId", userID).Msg("Headless auto-start blocked: owning user not found")
		return "", fmt.Errorf("user %d not found", userID)
	}
	if !user.APIKeyHash.Valid || user.APIKeyHash.String == "" {
		log.Warn().Str("clientId", clientID).Int64("userId", userID).Msg("Headless auto-start blocked: user has no relay key hash (regenerate the relay key in the dashboard)")
		return "", fmt.Errorf("user %d has no apiKeyHash (must regenerate relay key)", userID)
	}
	accountIdentifier := user.APIKeyHash.String

	// 2. KnownClient lookup + auto-start gate
	known, err := db.KnownClientStore().FindByClientID(ctx, userID, clientID)
	if err != nil || known == nil {
		log.Warn().Err(err).Str("clientId", clientID).Int64("userId", userID).Msg("Headless auto-start blocked: clientId is not a known client for this account (connect this world once so it is remembered)")
		return "", fmt.Errorf("clientId %s not found in known clients for user %d", clientID, userID)
	}
	if !bool(known.AutoStartOnRemoteRequest) {
		log.Warn().Str("clientId", clientID).Int64("userId", userID).Msg("Headless auto-start blocked: auto-start is not enabled for this client (turn it on in the Connections tab)")
		return "", fmt.Errorf("clientId %s does not have auto-start enabled", clientID)
	}

	// 3. Credential resolution — explicit credentialId required; no implicit fallback
	if !known.CredentialID.Valid {
		log.Warn().Str("clientId", clientID).Int64("userId", userID).Msg("Headless auto-start blocked: no credential assigned to this client (select one in the Connections tab)")
		return "", fmt.Errorf("clientId %s has auto-start enabled but no credential assigned; select one in Connections", clientID)
	}
	credential, err := db.CredentialStore().FindByID(ctx, known.CredentialID.Int64)
	if err != nil || credential == nil || credential.UserID != userID {
		log.Warn().Err(err).Str("clientId", clientID).Int64("userId", userID).Int64("credentialId", known.CredentialID.Int64).Msg("Headless auto-start blocked: assigned credential is missing or not owned by this account (re-create it in the Credentials tab)")
		return "", fmt.Errorf("KnownClient.credentialId references missing or unauthorized credential %d", known.CredentialID.Int64)
	}

	password, err := service.Decrypt(
		credential.EncryptedFoundryPassword, credential.PasswordIV,
		credential.PasswordAuthTag, cfg.CredentialsEncryptionKey,
	)
	if err != nil {
		log.Warn().Err(err).Str("clientId", clientID).Int64("userId", userID).Int64("credentialId", known.CredentialID.Int64).Msg("Headless auto-start blocked: could not decrypt stored password (CREDENTIALS_ENCRYPTION_KEY may have changed since the credential was saved — re-save it)")
		return "", fmt.Errorf("decrypt credential: %w", err)
	}

	// 4. Mint a fresh connection token. We hash for storage but keep the
	// raw value to inject into the headless browser's localStorage.
	rawToken, tokenHash, err := GenerateHeadlessToken()
	if err != nil {
		return "", fmt.Errorf("generate headless token: %w", err)
	}
	headlessTokenName := fmt.Sprintf("headless %s %s", clientID, time.Now().Format("2006-01-02 15:04"))
	headlessToken := &model.ConnectionToken{
		UserID:    userID,
		TokenHash: tokenHash,
		Name:      headlessTokenName,
		Source:    model.TokenSourceHeadless,
	}
	if err := db.ConnectionTokenStore().Create(ctx, headlessToken); err != nil {
		return "", fmt.Errorf("persist headless token: %w", err)
	}

	// Cleanup helper: revoke the token if anything below fails. We use a
	// flag so the success path can skip the revoke.
	revoked := false
	revokeOnFail := func() {
		if revoked {
			return
		}
		// Use a fresh background context — the caller's may be cancelled
		_ = db.ConnectionTokenStore().Delete(context.Background(), headlessToken.ID)
		revoked = true
	}
	defer revokeOnFail()

	// 5. Launch the headless session with the seed token. This is a
	// near-clone of LaunchSession with two key differences:
	//   - We inject the localStorage seed BEFORE navigation
	//   - We poll for the EXACT clientId, not the legacy foundry-{userId} format
	// World to launch on Foundry's setup screen: prefer the credential's
	// explicit default world; fall back to the world this client last connected
	// as. Empty means skip world selection (Foundry already has a world active).
	worldToLaunch := credential.World
	if worldToLaunch == "" {
		worldToLaunch = known.WorldTitle.String
	}

	resultClientID, err := m.launchHeadlessWithSeededToken(ctx, launchSeededOpts{
		AccountIdentifier: accountIdentifier,
		FoundryURL:        credential.FoundryURL,
		Username:          credential.FoundryUsername,
		Password:          password,
		WorldName:         worldToLaunch,
		ExpectedClientID:  clientID,
		ExpectedWorldID:   known.WorldID.String,
		SeedToken:         rawToken,
		FoundryVersion:    known.FoundryVersion.String,
	})
	if err != nil {
		return "", err
	}

	// 6. Success — keep the token alive for the duration of the session.
	// The session cleanup loop will revoke it (and broadcast the disconnect)
	// when the headless session times out.
	revoked = true
	log.Info().
		Str("clientId", resultClientID).
		Int64("userId", userID).
		Int64("tokenId", headlessToken.ID).
		Msg("Headless auto-start succeeded for remote-request")
	return resultClientID, nil
}

// launchSeededOpts groups the parameters for launchHeadlessWithSeededToken.
type launchSeededOpts struct {
	AccountIdentifier string // user.APIKeyHash — what ClientManager registers under
	FoundryURL        string
	Username          string
	Password          string
	WorldName         string // optional
	ExpectedClientID  string // poll for this exact clientId in ClientManager
	ExpectedWorldID   string // if set, abort immediately if Foundry is in a different world
	SeedToken         string // raw connection token to inject into localStorage
	FoundryVersion    string // known Foundry version (for log filename), empty if unknown
}

// launchHeadlessWithSeededToken is the AutoStartForKnownClient launch core.
// It mirrors LaunchSession but injects the connection token into localStorage
// before navigation and polls for an EXACT clientId rather than the legacy
// foundry-{userId} format.
func (m *HeadlessManager) launchHeadlessWithSeededToken(ctx context.Context, opts launchSeededOpts) (string, error) {
	// Stale-session cleanup (same as LaunchSession)
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
		if s.APIKey == opts.AccountIdentifier {
			count++
		}
	}
	m.mu.Unlock()

	if m.maxSessions > 0 && count >= m.maxSessions {
		if alerts.Track("headless_limit", 3, 10*time.Minute, 1*time.Hour) {
			alerts.Fire(alerts.Event{
				Type:     alerts.TypeHeadlessSessionFlood,
				Severity: "warning",
				Message:  "Headless session limit hit 3+ times in 10 minutes",
				Details:  map[string]interface{}{"maxSessions": m.maxSessions},
			})
		}
		return "", fmt.Errorf("maximum headless sessions (%d) reached for this account", m.maxSessions)
	}

	// Global cap check across all relay instances via Redis
	if err := m.checkGlobalSessionCap(); err != nil {
		return "", err
	}

	tabCtx, tabCancel, err := m.newIsolatedTab()
	if err != nil {
		return "", fmt.Errorf("create isolated tab: %w", err)
	}

	// Console capture for debugging
	logFile := setupBrowserConsoleCapture(tabCtx, m.dataDir, opts.Username, opts.FoundryVersion, opts.FoundryURL)
	if logFile != "" {
		log.Info().Str("logFile", logFile).Msg("Browser console logging enabled")
	}

	log.Info().
		Str("url", opts.FoundryURL).
		Str("username", opts.Username).
		Str("expectedClientId", opts.ExpectedClientID).
		Msg("Launching headless session with seeded connection token")

	// Inject WebGL override (same as LaunchSession)
	chromedp.Run(tabCtx, chromedp.ActionFunc(func(ctx context.Context) error {
		_, err := page.AddScriptToEvaluateOnNewDocument(webGLSpoof).Do(ctx)
		return err
	}))

	// THE KEY DIFFERENCE FROM LaunchSession: inject the connection token
	// into localStorage BEFORE Foundry's module init runs.
	//
	// Foundry stores client-scope settings in localStorage at the key
	// `<moduleId>.<settingKey>` with the value JSON-stringified (verified at
	// /home/noah/Foundry/foundry-v13/code/client/helpers/client-settings.mjs:266
	// for the key format and :349 for the JSON.stringify).
	//
	// addScriptToEvaluateOnNewDocument runs the script on every new document,
	// before any other page scripts. This guarantees the token is in
	// localStorage by the time the Foundry module's `init` hook fires and the
	// WebSocketManager constructor reads it.
	if err := injectConnectionTokenSeed(tabCtx, opts.SeedToken); err != nil {
		tabCancel()
		return "", fmt.Errorf("inject connection token seed: %w", err)
	}
	log.Info().Msg("Seeded connection token into headless browser localStorage (AutoStart path)")

	// Set viewport
	chromedp.Run(tabCtx, chromedp.EmulateViewport(int64(m.windowWidth), int64(m.windowHeight)))

	// Navigate
	navCtx, navCancel := context.WithTimeout(tabCtx, browserNavigateTimeout)
	defer navCancel()

	if err := chromedp.Run(navCtx, chromedp.Navigate(opts.FoundryURL)); err != nil {
		tabCancel()
		return "", fmt.Errorf("navigate to Foundry: %w", err)
	}

	// Dismiss notifications
	time.Sleep(300 * time.Millisecond)
	chromedp.Run(tabCtx, chromedp.Evaluate(`
		document.querySelectorAll('.notification .close, .notification a.close, #notifications .notification .close, .notification-pip, .notification .notification-close').forEach(el => el.click());
	`, nil))

	// World selection (only if needed)
	if opts.WorldName != "" {
		pageType := detectPage(tabCtx)
		if pageType == "worldList" {
			log.Info().Str("world", opts.WorldName).Msg("Selecting world")
			if err := selectWorld(tabCtx, opts.WorldName); err != nil {
				log.Warn().Err(err).Msg("World selection failed, continuing")
			}
		}
	}

	// Login
	log.Info().Str("username", opts.Username).Msg("Logging in")
	if _, err := loginToFoundry(tabCtx, opts.Username, opts.Password); err != nil {
		tabCancel()
		return "", fmt.Errorf("login failed: %w", err)
	}

	// Wait for game canvas
	log.Info().Msg("Waiting for game canvas")
	loadCtx, loadCancel := context.WithTimeout(tabCtx, gameLoadTimeout)
	defer loadCancel()
	if err := waitForAnySelector(loadCtx, []string{"#ui-left", "#sidebar", "#game", ".vtt"}); err != nil {
		tabCancel()
		return "", fmt.Errorf("game canvas did not load: %w", err)
	}

	// Wrong-world fast fail: if we know which world we need, verify it immediately
	// after canvas load rather than burning the full poll timeout on a world that
	// will never produce the expected clientId.
	if opts.ExpectedWorldID != "" {
		var loadedWorld string
		if err := chromedp.Run(tabCtx, chromedp.Evaluate(`
			(function(){ try { return game?.world?.id ?? ''; } catch(e) { return ''; } })()
		`, &loadedWorld)); err != nil {
			log.Warn().Err(err).Msg("Wrong-world check eval failed; skipping fast-fail")
		}
		if loadedWorld != "" && loadedWorld != opts.ExpectedWorldID {
			tabCancel()
			return "", fmt.Errorf("wrong world loaded: server is running %q, expected %q", loadedWorld, opts.ExpectedWorldID)
		}
	}

	// Now poll for the EXACT clientId we expect to register. The Foundry
	// module reads our seeded token from localStorage and connects with the
	// existing world clientId, so the resulting WS client should appear under
	// opts.ExpectedClientID, NOT the legacy foundry-{userId} format.
	sessionID := uuid.New().String()
	m.mu.Lock()
	m.pending[sessionID] = &PendingHeadless{
		SessionID:        sessionID,
		ExpectedClientID: opts.ExpectedClientID,
		APIKey:           opts.AccountIdentifier,
		ContextCancel:    tabCancel,
		StartTime:        time.Now(),
	}
	m.mu.Unlock()

	pollCtx, pollCancel := context.WithTimeout(context.Background(), autoStartPollTimeout)
	defer pollCancel()
	ticker := time.NewTicker(clientPollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if client := m.clientManager.GetClient(opts.ExpectedClientID); client != nil {
				if client.APIKey() == opts.AccountIdentifier {
					m.mu.Lock()
					delete(m.pending, sessionID)
					m.sessions[sessionID] = &HeadlessSession{
						SessionID:     sessionID,
						ClientID:      opts.ExpectedClientID,
						APIKey:        opts.AccountIdentifier,
						FoundryURL:    opts.FoundryURL,
						Username:      opts.Username,
						WorldName:     opts.WorldName,
						ContextCancel: tabCancel,
						StartedAt:     time.Now(),
						LastActivity:  time.Now(),
					}
					m.mu.Unlock()
					m.registerInRedis(sessionID, opts.ExpectedClientID)
					log.Info().
						Str("sessionId", sessionID).
						Str("clientId", opts.ExpectedClientID).
						Msg("Headless session established (seeded-token, exact match)")
					return opts.ExpectedClientID, nil
				}
				// Wrong account — this shouldn't happen but be defensive
				log.Warn().
					Str("clientId", opts.ExpectedClientID).
					Msg("Found client with expected ID but wrong account identifier")
			}
		case <-pollCtx.Done():
			m.mu.Lock()
			delete(m.pending, sessionID)
			m.mu.Unlock()
			tabCancel()
			return "", fmt.Errorf("seeded headless client did not register within %s", autoStartPollTimeout)
		}
	}
}

// generateHeadlessToken creates a fresh 32-byte random token for an auto-
// started headless session. Returns the raw value (sent to ChromeDP for
// localStorage injection) and its SHA-256 hash (stored in the DB).
// GenerateHeadlessToken creates a fresh 32-byte random token for an auto-
// started headless session. Returns the raw value (sent to ChromeDP for
// localStorage injection) and its SHA-256 hash (stored in the DB).
func GenerateHeadlessToken() (raw, hash string, err error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", "", err
	}
	raw = hex.EncodeToString(b)
	sum := sha256.Sum256([]byte(raw))
	hash = hex.EncodeToString(sum[:])
	return raw, hash, nil
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
		if m.inactiveTimeout > 0 && now.Sub(s.LastActivity) > m.inactiveTimeout {
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

func setupBrowserConsoleCapture(ctx context.Context, dataDir, username, foundryVersion, foundryURL string) string {
	captureLevel := os.Getenv("CAPTURE_BROWSER_CONSOLE")
	if captureLevel == "" {
		return ""
	}
	logDir := filepath.Join(dataDir, "browser-logs")
	os.MkdirAll(logDir, 0755)
	timestamp := time.Now().Format("2006-01-02T15-04-05")
	var filename string
	if foundryVersion != "" {
		filename = filepath.Join(logDir, fmt.Sprintf("headless_%s_v%s_%s.log", username, foundryVersion, timestamp))
	} else {
		filename = filepath.Join(logDir, fmt.Sprintf("headless_%s_%s.log", username, timestamp))
	}
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

// CleanBrowserLogs deletes browser log files older than maxAgeDays from the
// browser-logs subdirectory of dataDir. Safe to call when logging is disabled
// (no-op if the directory does not exist).
func CleanBrowserLogs(dataDir string, maxAgeDays int) {
	logDir := filepath.Join(dataDir, "browser-logs")
	entries, err := os.ReadDir(logDir)
	if err != nil {
		return // directory doesn't exist or unreadable — nothing to clean
	}
	cutoff := time.Now().AddDate(0, 0, -maxAgeDays)
	deleted := 0
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		info, err := entry.Info()
		if err != nil {
			continue
		}
		if info.ModTime().Before(cutoff) {
			if removeErr := os.Remove(filepath.Join(logDir, entry.Name())); removeErr == nil {
				deleted++
			}
		}
	}
	if deleted > 0 {
		log.Info().Int("count", deleted).Msg("Cleaned up old browser log files")
	}
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

// detectPage waits for one of the known page states to appear and returns
// which one won. Times out after 15 seconds and returns "unknown".
// Polling avoids a race where the world list is rendered by JavaScript after
// the initial HTML load (observed in Foundry v14+).
func detectPage(ctx context.Context) string {
	checkCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	ticker := time.NewTicker(300 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-checkCtx.Done():
			return "unknown"
		case <-ticker.C:
			var result string
			err := chromedp.Run(checkCtx, chromedp.Evaluate(`
				(function() {
					if (document.querySelector('input[name="password"]')) return 'login';
					if (document.querySelector('li.package.world')) return 'worldList';
					if (document.querySelector('#ui-left, #sidebar, #game')) return 'game';
					return 'unknown';
				})()
			`, &result))
			if err == nil && result != "unknown" {
				return result
			}
		}
	}
}

func selectWorld(ctx context.Context, worldName string) error {
	selCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	if err := chromedp.Run(selCtx, chromedp.WaitVisible(`li.package.world`, chromedp.ByQuery)); err != nil {
		return fmt.Errorf("world list not found: %w", err)
	}
	// Match on either the world's display title or its id (data-package-id slug),
	// so callers can pass whichever they have. Comparison is case-insensitive.
	js := fmt.Sprintf(`
		(function() {
			const target = %q;
			const worlds = document.querySelectorAll('li.package.world');
			for (const w of worlds) {
				const title = w.querySelector('h3.package-title, .package-title');
				const titleText = title ? title.textContent.trim().toLowerCase() : '';
				const id = (w.getAttribute('data-package-id') || w.dataset.packageId || '').trim().toLowerCase();
				if (titleText === target || (id && id === target)) {
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

	// Pre-submit: if a user-select dropdown is present, verify the target user
	// exists and is not already active (disabled/greyed = already logged in).
	var preCheck string
	if err := chromedp.Run(loginCtx, chromedp.Evaluate(fmt.Sprintf(`
		(function() {
			const sel = document.querySelector('select[name="userid"]');
			if (!sel) return 'ok'; // single-user or password-only form
			const opts = Array.from(sel.options);
			const match = opts.find(o => o.textContent.trim().toLowerCase() === %q);
			if (!match) return 'not-found';
			if (match.disabled) return 'already-active';
			return 'ok';
		})()
	`, strings.ToLower(username)), &preCheck)); err != nil {
		log.Warn().Err(err).Msg("Pre-submit user check eval failed; proceeding without pre-check")
	}
	switch preCheck {
	case "not-found":
		return "", fmt.Errorf("configured Foundry user not found on this server (check credential settings)")
	case "already-active":
		return "", fmt.Errorf("configured Foundry user is already logged in (another active session may be holding the slot)")
	}

	// POST /join with an explicit password instead of filling the form: Foundry's
	// form drops an empty password field on submit, so a passwordless GM (stored as
	// createPassword("")) can never log in through it. Posting ourselves guarantees
	// the password — including "" — reaches the server.
	js := fmt.Sprintf(`
		(async function() {
			var userId = %q;
			const sel = document.querySelector('select[name="userid"]');
			if (sel) {
				const match = Array.from(sel.options).find(o => o.textContent.trim().toLowerCase() === %q);
				if (match) userId = match.value;
			}
			try {
				const resp = await fetch('/join', {
					method: 'POST',
					headers: { 'Content-Type': 'application/json' },
					body: JSON.stringify({ action: 'join', userid: userId, password: %q })
				});
				if (resp.ok) {
					let data = {};
					try { data = await resp.json(); } catch (e) {}
					if (!data || data.status !== 'success') return 'fail:' + ((data && data.message) || 'join was not accepted');
					setTimeout(function() { window.location.href = data.redirect || '/game'; }, 0);
					return 'ok:' + userId;
				}
				const t = (await resp.text()).replace(/\s+/g, ' ').trim().substring(0, 120);
				return 'fail:' + t;
			} catch (e) {
				return 'fail:' + (e && e.message ? e.message : 'join request failed');
			}
		})()
	`, username, strings.ToLower(username), password)

	var result string
	if err := chromedp.Run(loginCtx, chromedp.Evaluate(js, &result,
		func(p *runtime.EvaluateParams) *runtime.EvaluateParams { return p.WithAwaitPromise(true) },
	)); err != nil {
		return "", fmt.Errorf("login eval: %w", err)
	}
	if strings.HasPrefix(result, "fail:") {
		return "", fmt.Errorf("login rejected by Foundry: %s", strings.TrimPrefix(result, "fail:"))
	}
	userID := strings.TrimPrefix(result, "ok:")
	if userID == "" {
		userID = username
	}
	// The POST already returned Foundry's authoritative verdict; the caller waits
	// for the game canvas to confirm the world finished loading.
	log.Info().Str("userId", userID).Msg("Headless login succeeded via POST /join")
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
