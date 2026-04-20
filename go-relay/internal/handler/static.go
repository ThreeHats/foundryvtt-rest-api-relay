package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ThreeHats/foundryvtt-rest-api-relay/go-relay/internal/handler/helpers"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

// findFile checks multiple paths and returns the first that exists.
func findFile(paths ...string) string {
	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}
	return ""
}

// ServeStaticPage serves a static HTML page from multiple fallback paths.
func ServeStaticPage(paths ...string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p := findFile(paths...)
		if p == "" {
			http.NotFound(w, r)
			return
		}
		http.ServeFile(w, r, p)
	}
}

// OpenAPIHandler serves OpenAPI spec with dynamic server URL injection.
func OpenAPIHandler(specPath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := os.ReadFile(specPath)
		if err != nil {
			http.NotFound(w, r)
			return
		}

		var spec map[string]interface{}
		if err := json.Unmarshal(data, &spec); err != nil {
			http.Error(w, "Invalid spec", http.StatusInternalServerError)
			return
		}

		// Inject server URL from request
		scheme := "http"
		if r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https" {
			scheme = "https"
		}
		serverURL := scheme + "://" + r.Host

		if servers, ok := spec["servers"].([]interface{}); ok && len(servers) > 0 {
			if s, ok := servers[0].(map[string]interface{}); ok {
				s["url"] = serverURL
			}
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "public, max-age=300")
		json.NewEncoder(w).Encode(spec)
	}
}

// AsyncAPIHandler serves AsyncAPI spec with dynamic WebSocket URL injection.
func AsyncAPIHandler(specPath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := os.ReadFile(specPath)
		if err != nil {
			http.NotFound(w, r)
			return
		}

		var spec map[string]interface{}
		if err := json.Unmarshal(data, &spec); err != nil {
			http.Error(w, "Invalid spec", http.StatusInternalServerError)
			return
		}

		wsScheme := "ws"
		if r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https" {
			wsScheme = "wss"
		}
		wsURL := wsScheme + "://" + r.Host + "/ws/api"

		if servers, ok := spec["servers"].(map[string]interface{}); ok {
			if prod, ok := servers["production"].(map[string]interface{}); ok {
				prod["url"] = wsURL
			}
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "public, max-age=300")
		json.NewEncoder(w).Encode(spec)
	}
}

// APIDocsHandler serves api-docs.json with dynamic baseUrl.
func APIDocsHandler(docsPath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := os.ReadFile(docsPath)
		if err != nil {
			http.NotFound(w, r)
			return
		}

		var docs map[string]interface{}
		if err := json.Unmarshal(data, &docs); err != nil {
			http.Error(w, "Invalid docs", http.StatusInternalServerError)
			return
		}

		scheme := "http"
		if r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https" {
			scheme = "https"
		}
		docs["baseUrl"] = scheme + "://" + r.Host

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "public, max-age=300")
		json.NewEncoder(w).Encode(docs)
	}
}

// ProxyAssetHandler proxies Foundry assets with CDN fallbacks.
func ProxyAssetHandler() http.HandlerFunc {
	client := &http.Client{Timeout: 30 * time.Second}

	return func(w http.ResponseWriter, r *http.Request) {
		assetPath := chi.URLParam(r, "*")
		if assetPath == "" {
			http.NotFound(w, r)
			return
		}

		// Font Awesome → CDN redirect
		if strings.Contains(assetPath, "/webfonts/fa-") || strings.Contains(assetPath, "/fonts/fontawesome/") || strings.Contains(assetPath, "/fonts/fa-") {
			filename := filepath.Base(assetPath)
			http.Redirect(w, r, "https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.5.2/webfonts/"+filename, http.StatusFound)
			return
		}

		// The Forge assets → placeholder
		if strings.Contains(assetPath, "forgevtt-module") || strings.Contains(assetPath, "forge-vtt.com") {
			if strings.HasSuffix(assetPath, ".css") {
				w.Header().Set("Content-Type", "text/css")
				w.Write([]byte("/* Placeholder for The Forge CSS */"))
			} else if strings.HasSuffix(assetPath, ".js") {
				w.Header().Set("Content-Type", "application/javascript")
				w.Write([]byte("// Placeholder for The Forge JS"))
			} else {
				// 1x1 transparent PNG
				w.Header().Set("Content-Type", "image/png")
				transparentPNG := []byte{0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a, 0x00, 0x00, 0x00, 0x0d, 0x49, 0x48, 0x44, 0x52, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01, 0x08, 0x06, 0x00, 0x00, 0x00, 0x1f, 0x15, 0xc4, 0x89, 0x00, 0x00, 0x00, 0x0a, 0x49, 0x44, 0x41, 0x54, 0x78, 0x9c, 0x62, 0x00, 0x00, 0x00, 0x02, 0x00, 0x01, 0xe2, 0x21, 0xbc, 0x33, 0x00, 0x00, 0x00, 0x00, 0x49, 0x45, 0x4e, 0x44, 0xae, 0x42, 0x60, 0x82}
				w.Write(transparentPNG)
			}
			return
		}

		// Texture files → GitHub fallback
		if strings.Contains(assetPath, "texture1.webp") || strings.Contains(assetPath, "texture2.webp") || strings.Contains(assetPath, "parchment.jpg") {
			http.Redirect(w, r, "https://raw.githubusercontent.com/foundryvtt/dnd5e/master/ui/parchment.jpg", http.StatusFound)
			return
		}

		// General proxy — try to fetch from Foundry
		foundryURL := "http://localhost:30000/" + assetPath
		resp, err := client.Get(foundryURL)
		if err != nil {
			log.Debug().Err(err).Str("path", assetPath).Msg("Asset proxy failed")
			http.NotFound(w, r)
			return
		}
		defer resp.Body.Close()

		// Copy headers
		for key, values := range resp.Header {
			for _, v := range values {
				w.Header().Add(key, v)
			}
		}
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(resp.StatusCode)
		io.Copy(w, resp.Body)
	}
}

// DocsFileServer serves Docusaurus documentation with SPA fallback.
// For any path that doesn't match a real file, it tries path/index.html
// (Docusaurus generates /foo/index.html for route /foo), then falls back
// to the root index.html for client-side routing.
func DocsFileServer(docsDir string) http.HandlerFunc {
	fsys := http.Dir(docsDir)

	return func(w http.ResponseWriter, r *http.Request) {
		// Strip /docs prefix
		path := strings.TrimPrefix(r.URL.Path, "/docs")
		if path == "" || path == "/" {
			http.ServeFile(w, r, filepath.Join(docsDir, "index.html"))
			return
		}

		// Trailing slash redirect
		if strings.HasSuffix(path, "/") {
			http.Redirect(w, r, strings.TrimSuffix(r.URL.Path, "/"), http.StatusMovedPermanently)
			return
		}

		// 1. Try exact file (CSS, JS, images, etc.)
		if f, err := fsys.Open(path); err == nil {
			stat, _ := f.Stat()
			f.Close()
			if stat != nil && !stat.IsDir() {
				http.ServeFile(w, r, filepath.Join(docsDir, path))
				return
			}
		}

		// 2. Try path/index.html (Docusaurus generates these for each route)
		indexPath := filepath.Join(path, "index.html")
		if f, err := fsys.Open(indexPath); err == nil {
			f.Close()
			http.ServeFile(w, r, filepath.Join(docsDir, indexPath))
			return
		}

		// 3. Try path.html
		htmlPath := path + ".html"
		if f, err := fsys.Open(htmlPath); err == nil {
			f.Close()
			http.ServeFile(w, r, filepath.Join(docsDir, htmlPath))
			return
		}

		// 4. SPA fallback — serve root index.html for client-side routing
		http.ServeFile(w, r, filepath.Join(docsDir, "index.html"))
	}
}

// unused import guard
var _ = helpers.WriteJSON
