// Command docgen generates API documentation JSON files (api-docs.json,
// openapi.json, asyncapi.json) by parsing the Go handler source files.
//
// Usage:
//
//	go run ./cmd/docgen
//
// Run from the go-relay directory. Output is written to ../public/.
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// ---------------------------------------------------------------------------
// Version
// ---------------------------------------------------------------------------

const apiVersion = "2.2.1"

// ---------------------------------------------------------------------------
// Data types that mirror the runtime helpers.ParamDef / helpers.APIRouteConfig
// ---------------------------------------------------------------------------

type paramDef struct {
	Name        string
	Type        string   // "string", "number", "boolean", "array", "object"
	From        []string // "query", "body", "params"
	Required    bool
	Description string
}

type routeInfo struct {
	Method      string
	Path        string
	Summary     string // first line of doc comment
	Description string // subsequent lines of doc comment
	Tag         string
	Returns     string
	MsgType     string // WebSocket message type (from APIRouteConfig.Type)
	Required    []paramDef
	Optional    []paramDef
	IsSSE       bool
	IsManual    bool // not built with CreateAPIRoute/h()
}

// ---------------------------------------------------------------------------
// Main
// ---------------------------------------------------------------------------

func main() {
	baseDir := "."
	if _, err := os.Stat("internal/handler"); os.IsNotExist(err) {
		fmt.Fprintln(os.Stderr, "Run this command from the go-relay directory")
		os.Exit(1)
	}

	handlerDir := filepath.Join(baseDir, "internal", "handler")

	// 1. Parse entity.go for config-builder functions & manual handlers
	entitySrc := mustRead(filepath.Join(handlerDir, "entity.go"))
	dnd5eSrc := mustRead(filepath.Join(handlerDir, "dnd5e.go"))
	routesSrc := mustRead(filepath.Join(handlerDir, "routes.go"))
	sheetSrc := mustRead(filepath.Join(handlerDir, "sheet.go"))
	fsSrc := mustRead(filepath.Join(handlerDir, "filesystem.go"))
	sessionSrc := mustRead(filepath.Join(handlerDir, "session.go"))

	// Build function-name -> parsed config map from entity.go
	configFuncs := parseConfigFunctions(entitySrc)

	// Parse routes.go to get HTTP method + path + handler function name
	routes := parseRoutes(routesSrc)

	// Parse dnd5e.go routes (inline CreateAPIRoute calls)
	dnd5eRoutes := parseDnd5eRoutes(dnd5eSrc)

	// Parse manual handler doc comments
	sheetRoutes := parseManualHandlers(sheetSrc)
	fsRoutes := parseManualHandlers(fsSrc)
	sessionRoutes := parseManualHandlers(sessionSrc)
	sseRoutes := parseManualHandlers(routesSrc)      // SSE handlers are in routes.go
	entityManual := parseManualHandlers(entitySrc)    // clientsHandler is in entity.go

	// Build a map of manual handler func name -> routeInfo
	manualMap := make(map[string]*routeInfo)
	for i := range sheetRoutes {
		manualMap[sheetRoutes[i].funcName] = &sheetRoutes[i].info
	}
	for i := range fsRoutes {
		manualMap[fsRoutes[i].funcName] = &fsRoutes[i].info
	}
	for i := range sessionRoutes {
		manualMap[sessionRoutes[i].funcName] = &sessionRoutes[i].info
	}
	for i := range sseRoutes {
		manualMap[sseRoutes[i].funcName] = &sseRoutes[i].info
	}
	for i := range entityManual {
		manualMap[entityManual[i].funcName] = &entityManual[i].info
	}

	// Resolve routes: merge route registrations with parsed configs
	var allRoutes []routeInfo
	for _, rt := range routes {
		if rt.handlerFunc == "contentsDeprecated" {
			// Include the deprecated route with a deprecation notice
			allRoutes = append(allRoutes, routeInfo{
				Method:      rt.method,
				Path:        rt.path,
				Summary:     "This route is deprecated",
				Description: "Use /structure with the path query parameter instead.",
				Tag:         "Structure",
				Returns:     "Error message directing to use /structure endpoint",
				IsManual:    true,
			})
			continue
		}
		if cfg, ok := configFuncs[rt.handlerFunc]; ok {
			ri := routeInfo{
				Method:      rt.method,
				Path:        rt.path,
				Summary:     cfg.summary,
				Description: cfg.description,
				Tag:         cfg.tag,
				Returns:     cfg.returns,
				MsgType:     cfg.msgType,
				Required:    cfg.required,
				Optional:    cfg.optional,
			}
			allRoutes = append(allRoutes, ri)
		} else if manual, ok := manualMap[rt.handlerFunc]; ok {
			manual.Method = rt.method
			manual.Path = rt.path
			manual.IsManual = true
			allRoutes = append(allRoutes, *manual)
		}
	}

	// Add dnd5e routes with /dnd5e prefix
	for _, ri := range dnd5eRoutes {
		ri.Path = "/dnd5e" + ri.Path
		allRoutes = append(allRoutes, ri)
	}

	// Add auth routes
	allRoutes = append(allRoutes, buildAuthRoutes()...)

	// Write output files
	outDir := filepath.Join(baseDir, "..", "public")
	os.MkdirAll(outDir, 0o755)

	writeJSON(filepath.Join(outDir, "api-docs.json"), buildAPIDocs(allRoutes))
	writeJSON(filepath.Join(outDir, "openapi.json"), buildOpenAPI(allRoutes))
	writeJSON(filepath.Join(outDir, "asyncapi.json"), buildAsyncAPI(allRoutes))

	// Generate markdown documentation
	mdDir := filepath.Join(baseDir, "..", "docs", "md", "api")
	os.MkdirAll(mdDir, 0o755)
	generateMarkdown(allRoutes, mdDir)

	fmt.Printf("Generated %d routes into %s\n", len(allRoutes), outDir)
}

// ---------------------------------------------------------------------------
// Route registration parser (routes.go)
// ---------------------------------------------------------------------------

type rawRoute struct {
	method      string
	path        string
	handlerFunc string
}

func parseRoutes(src string) []rawRoute {
	var routes []rawRoute

	// Pattern: r.Get("/path", h(mgr, pending, funcName))
	reH := regexp.MustCompile(`r\.(Get|Post|Put|Delete)\("([^"]+)",\s*h\(mgr,\s*pending,\s*(\w+)\)\)`)
	for _, m := range reH.FindAllStringSubmatch(src, -1) {
		routes = append(routes, rawRoute{method: strings.ToUpper(m[1]), path: m[2], handlerFunc: m[3]})
	}

	// Pattern: r.Get("/path", handlerFunc(mgr, pending)) or handlerFunc(mgr, pending, ...)
	// Manual handlers like sheetGetHandler, uploadHandler, sessionHandshakeHandler, etc.
	reManual := regexp.MustCompile(`r\.(Get|Post|Put|Delete)\("([^"]+)",\s*(\w+)\(`)
	for _, m := range reManual.FindAllStringSubmatch(src, -1) {
		funcName := m[3]
		if funcName == "h" {
			continue // already handled above
		}
		routes = append(routes, rawRoute{method: strings.ToUpper(m[1]), path: m[2], handlerFunc: funcName})
	}

	// Pattern: r.Get("/path", funcRef) — direct function reference (no parens), e.g. contentsDeprecated
	reDirect := regexp.MustCompile(`r\.(Get|Post|Put|Delete)\("([^"]+)",\s*(\w+)\)`)
	for _, m := range reDirect.FindAllStringSubmatch(src, -1) {
		funcName := m[3]
		if funcName == "h" {
			continue
		}
		// Check it wasn't already captured
		found := false
		for _, existing := range routes {
			if existing.path == m[2] && existing.method == strings.ToUpper(m[1]) {
				found = true
				break
			}
		}
		if !found {
			routes = append(routes, rawRoute{method: strings.ToUpper(m[1]), path: m[2], handlerFunc: funcName})
		}
	}

	return routes
}

// ---------------------------------------------------------------------------
// Config function parser (entity.go) — funcs that return helpers.APIRouteConfig
// ---------------------------------------------------------------------------

type parsedConfig struct {
	summary     string
	description string
	tag         string
	returns     string
	msgType     string
	required    []paramDef
	optional    []paramDef
}

func parseConfigFunctions(src string) map[string]*parsedConfig {
	configs := make(map[string]*parsedConfig)
	lines := strings.Split(src, "\n")

	// Pre-parse known helper param functions
	helperParams := parseHelperParamFuncs(lines)

	// Pre-parse package-level var params
	pkgParams := parsePkgLevelParams(lines)

	for i := 0; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])

		// Look for func declarations that return helpers.APIRouteConfig
		if !strings.HasPrefix(line, "func ") {
			continue
		}
		if !strings.Contains(line, "helpers.APIRouteConfig") {
			continue
		}

		// Extract function name
		re := regexp.MustCompile(`func (\w+)\(\)\s*helpers\.APIRouteConfig`)
		m := re.FindStringSubmatch(line)
		if m == nil {
			continue
		}
		funcName := m[1]

		// Look backwards for doc comment
		summary, desc, tag, returns := extractDocComment(lines, i)

		// Extract the function body
		bodyStart := i
		body := extractFuncBody(lines, bodyStart)

		// Parse the config from the body
		cfg := parseConfigBody(body, src, helperParams, pkgParams)
		cfg.summary = summary
		cfg.description = desc
		cfg.tag = tag
		cfg.returns = returns

		configs[funcName] = cfg
	}

	return configs
}

// extractDocComment reads the comment block above line index i.
// Go doc comments must be contiguous with the declaration (no blank lines).
func extractDocComment(lines []string, i int) (summary, desc, tag, returns string) {
	var commentLines []string
	for j := i - 1; j >= 0; j-- {
		trimmed := strings.TrimSpace(lines[j])
		if strings.HasPrefix(trimmed, "//") {
			commentLines = append([]string{trimmed}, commentLines...)
		} else {
			break
		}
	}

	var descLines []string
	for _, cl := range commentLines {
		text := strings.TrimPrefix(cl, "//")
		text = strings.TrimSpace(text)

		if strings.HasPrefix(text, "@tag ") {
			tag = strings.TrimPrefix(text, "@tag ")
			continue
		}
		if strings.HasPrefix(text, "@returns ") {
			returns = strings.TrimPrefix(text, "@returns ")
			continue
		}
		if strings.HasPrefix(text, "@param ") {
			continue // handled separately for manual handlers
		}

		if summary == "" && text != "" {
			summary = text
		} else if text != "" {
			descLines = append(descLines, text)
		}
	}
	desc = strings.Join(descLines, " ")
	return
}

// extractFuncBody returns the full text of the function body (between { and matching }).
func extractFuncBody(lines []string, startLine int) string {
	var buf strings.Builder
	depth := 0
	started := false
	for i := startLine; i < len(lines); i++ {
		line := lines[i]
		for _, ch := range line {
			if ch == '{' {
				depth++
				started = true
			}
			if ch == '}' {
				depth--
			}
		}
		buf.WriteString(line)
		buf.WriteByte('\n')
		if started && depth == 0 {
			break
		}
	}
	return buf.String()
}

// parseConfigBody extracts Type, RequiredParams, OptionalParams from a function body.
// fullSrc is the complete source file, used to resolve variables outside the immediate body.
func parseConfigBody(body, fullSrc string, helperParams map[string]paramDef, pkgParams map[string]paramDef) *parsedConfig {
	cfg := &parsedConfig{}

	// Extract Type
	reType := regexp.MustCompile(`Type:\s*"([^"]+)"`)
	if m := reType.FindStringSubmatch(body); m != nil {
		cfg.msgType = m[1]
	}

	// Check for encSimple delegation
	reEncSimple := regexp.MustCompile(`return encSimple\("([^"]+)"\)`)
	if m := reEncSimple.FindStringSubmatch(body); m != nil {
		cfg.msgType = m[1]
		// encSimple routes have: clientIDParam(), encounterParam, userIDParam()
		cfg.optional = []paramDef{
			helperParams["clientIDParam"],
			pkgParams["encounterParam"],
			helperParams["userIDParam"],
		}
		return cfg
	}

	// Extract RequiredParams
	cfg.required = extractParamSlice(body, fullSrc, "RequiredParams", helperParams, pkgParams)
	cfg.optional = extractParamSlice(body, fullSrc, "OptionalParams", helperParams, pkgParams)

	return cfg
}

// extractParamSlice extracts a []helpers.ParamDef{...} block from source.
// fullSrc is the complete file source, used to resolve variable references outside body.
func extractParamSlice(body, fullSrc, fieldName string, helperParams map[string]paramDef, pkgParams map[string]paramDef) []paramDef {
	// Find the field start
	idx := strings.Index(body, fieldName+":")
	if idx < 0 {
		return nil
	}

	// Find the opening bracket of the slice literal
	rest := body[idx:]
	brStart := strings.Index(rest, "[]helpers.ParamDef{")
	if brStart < 0 {
		return nil
	}
	rest = rest[brStart+len("[]helpers.ParamDef{"):]

	// Also handle append patterns like: append([]helpers.ParamDef{...}, rollOptions...)
	// Scope the search to just this field (up to the next field or closing brace)
	fieldStart := body[idx:]
	// Truncate at the next field assignment to avoid matching append in a different field
	for _, nextField := range []string{"RequiredParams:", "OptionalParams:", "ValidateParams:", "BuildPayload:", "BuildPendingRequest:", "Timeout:"} {
		if nextField == fieldName+":" {
			continue
		}
		if nextIdx := strings.Index(fieldStart[len(fieldName)+1:], nextField); nextIdx >= 0 {
			fieldStart = fieldStart[:len(fieldName)+1+nextIdx]
			break
		}
	}
	isAppend := false
	appendVarName := ""
	reAppend := regexp.MustCompile(`append\(\[\]helpers\.ParamDef\{([^}]*)\},\s*(\w+)\.\.\.\)`)
	if m := reAppend.FindStringSubmatch(fieldStart); m != nil {
		isAppend = true
		appendVarName = m[2]
	}

	// Find the matching closing brace for the slice literal
	depth := 1
	var paramText strings.Builder
	for _, ch := range rest {
		if ch == '{' {
			depth++
		}
		if ch == '}' {
			depth--
			if depth == 0 {
				break
			}
		}
		paramText.WriteRune(ch)
	}

	params := parseParamEntries(paramText.String(), helperParams, pkgParams)

	if isAppend && appendVarName != "" {
		// Handle known append variables — search full file source for the variable
		appendParams := resolveParamVar(fullSrc, appendVarName, helperParams, pkgParams)
		if appendParams == nil {
			// Fallback: try the local body
			appendParams = resolveParamVar(body, appendVarName, helperParams, pkgParams)
		}
		params = append(params, appendParams...)
	}

	return params
}

// parseParamEntries parses comma-separated ParamDef entries or function calls.
func parseParamEntries(text string, helperParams map[string]paramDef, pkgParams map[string]paramDef) []paramDef {
	var params []paramDef

	// Split by top-level commas (respecting nested braces)
	entries := splitTopLevel(text)

	for _, entry := range entries {
		entry = strings.TrimSpace(entry)
		if entry == "" {
			continue
		}

		// Check for helper function calls: clientIDParam(), userIDParam(), etc.
		reFuncCall := regexp.MustCompile(`^(\w+)\(\)$`)
		if m := reFuncCall.FindStringSubmatch(entry); m != nil {
			if p, ok := helperParams[m[1]]; ok {
				params = append(params, p)
			}
			continue
		}

		// Check for package-level var references: encounterParam, dtParam, actorUuid, optActorUuid
		if p, ok := pkgParams[entry]; ok {
			params = append(params, p)
			continue
		}

		// Check for local variable references (single word, no braces)
		if !strings.Contains(entry, "{") && !strings.Contains(entry, ":") {
			// Could be a local var like actorUuid, optActorUuid
			trimmed := strings.TrimSpace(entry)
			if p, ok := pkgParams[trimmed]; ok {
				params = append(params, p)
			}
			continue
		}

		// Parse inline struct literal: {Name: "x", From: ...}
		p := parseInlineParamDef(entry)
		if p.Name != "" {
			params = append(params, p)
		}
	}

	return params
}

// splitTopLevel splits text by commas at brace depth 0.
func splitTopLevel(text string) []string {
	var parts []string
	depth := 0
	var current strings.Builder
	for _, ch := range text {
		switch ch {
		case '{':
			depth++
			current.WriteRune(ch)
		case '}':
			depth--
			current.WriteRune(ch)
		case ',':
			if depth == 0 {
				parts = append(parts, current.String())
				current.Reset()
			} else {
				current.WriteRune(ch)
			}
		default:
			current.WriteRune(ch)
		}
	}
	if s := current.String(); strings.TrimSpace(s) != "" {
		parts = append(parts, s)
	}
	return parts
}

// parseInlineParamDef parses a single {Name: "x", Type: ..., ...} struct literal.
func parseInlineParamDef(text string) paramDef {
	var p paramDef

	// Remove outer braces
	text = strings.TrimSpace(text)
	if strings.HasPrefix(text, "{") {
		text = text[1:]
	}
	if strings.HasSuffix(text, "}") {
		text = text[:len(text)-1]
	}

	// Extract Name
	reName := regexp.MustCompile(`Name:\s*"([^"]+)"`)
	if m := reName.FindStringSubmatch(text); m != nil {
		p.Name = m[1]
	}

	// Extract Description — handle escaped quotes in Go string literals
	if idx := strings.Index(text, `Description: "`); idx >= 0 {
		rest := text[idx+len(`Description: "`):]
		// Find the closing quote that isn't escaped
		var desc strings.Builder
		for i := 0; i < len(rest); i++ {
			if rest[i] == '\\' && i+1 < len(rest) && rest[i+1] == '"' {
				desc.WriteByte('"')
				i++ // skip the escaped quote
			} else if rest[i] == '"' {
				break // unescaped closing quote
			} else {
				desc.WriteByte(rest[i])
			}
		}
		p.Description = desc.String()
	}

	// Extract Required
	if strings.Contains(text, "Required: true") {
		p.Required = true
	}

	// Extract Type
	reParamType := regexp.MustCompile(`Type:\s*helpers\.(Type\w+)`)
	if m := reParamType.FindStringSubmatch(text); m != nil {
		p.Type = mapParamType(m[1])
	}
	if p.Type == "" {
		p.Type = "string" // default
	}

	// Extract From sources
	p.From = extractFromSources(text)

	return p
}

func mapParamType(t string) string {
	switch t {
	case "TypeString":
		return "string"
	case "TypeNumber":
		return "number"
	case "TypeBoolean":
		return "boolean"
	case "TypeArray":
		return "array"
	case "TypeObject":
		return "object"
	}
	return "string"
}

func extractFromSources(text string) []string {
	var sources []string
	if strings.Contains(text, "FromQuery") {
		sources = append(sources, "query")
	}
	if strings.Contains(text, "FromBody") {
		sources = append(sources, "body")
	}
	if strings.Contains(text, "FromParams") {
		sources = append(sources, "params")
	}
	// If "bqParam" or "bq" references
	if strings.Contains(text, "bqParam") || (strings.Contains(text, "From: bq") && !strings.Contains(text, "FromBody") && !strings.Contains(text, "FromQuery")) {
		if len(sources) == 0 {
			sources = []string{"body", "query"}
		}
	}
	if len(sources) == 0 {
		sources = []string{"query"} // default
	}
	return sources
}

// parseHelperParamFuncs extracts helper functions like clientIDParam(), userIDParam(), etc.
func parseHelperParamFuncs(lines []string) map[string]paramDef {
	helpers := make(map[string]paramDef)
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		reFn := regexp.MustCompile(`^func (\w+Param)\(\)\s*helpers\.ParamDef`)
		if m := reFn.FindStringSubmatch(trimmed); m != nil {
			body := extractFuncBody(lines, i)
			p := parseReturnedParamDef(body)
			if p.Name != "" {
				helpers[m[1]] = p
			}
		}
	}
	return helpers
}

// parsePkgLevelParams extracts package-level var declarations like:
//
//	var encounterParam = helpers.ParamDef{...}
//	var dtParam = helpers.ParamDef{...}
func parsePkgLevelParams(lines []string) map[string]paramDef {
	params := make(map[string]paramDef)
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)

		// var name = helpers.ParamDef{...}
		reVar := regexp.MustCompile(`^var\s+(\w+)\s*=\s*helpers\.ParamDef\{`)
		if m := reVar.FindStringSubmatch(trimmed); m != nil {
			body := extractFuncBody(lines, i)
			p := parseInlineParamDef(body[strings.Index(body, "{"):])
			if p.Name != "" {
				params[m[1]] = p
			}
		}

		// Local variable assignment: actorUuid := helpers.ParamDef{...}
		reLocal := regexp.MustCompile(`^\s*(\w+)\s*:=\s*helpers\.ParamDef\{`)
		if m := reLocal.FindStringSubmatch(line); m != nil {
			rest := line[strings.Index(line, "{"):]
			// May span multiple lines
			if strings.Count(rest, "{") > strings.Count(rest, "}") {
				// multi-line - gather until balanced
				for j := i + 1; j < len(lines); j++ {
					rest += "\n" + lines[j]
					if strings.Count(rest, "{") <= strings.Count(rest, "}") {
						break
					}
				}
			}
			p := parseInlineParamDef(rest)
			if p.Name != "" {
				params[m[1]] = p
			}
		}
	}
	return params
}

// parseReturnedParamDef parses a function body like:
//
//	func clientIDParam() helpers.ParamDef {
//	    return helpers.ParamDef{Name: "clientId", ...}
//	}
func parseReturnedParamDef(body string) paramDef {
	idx := strings.Index(body, "helpers.ParamDef{")
	if idx < 0 {
		return paramDef{}
	}
	return parseInlineParamDef(body[idx+len("helpers.ParamDef"):])
}

// ---------------------------------------------------------------------------
// resolveParamVar resolves references to local []ParamDef variables like rollOptions, abilityParams
// ---------------------------------------------------------------------------

func resolveParamVar(fullBody, varName string, helperParams map[string]paramDef, pkgParams map[string]paramDef) []paramDef {
	// Find variable declaration in the full body
	// Pattern: varName := []helpers.ParamDef{...}
	re := regexp.MustCompile(varName + `\s*:=\s*\[\]helpers\.ParamDef\{`)
	idx := re.FindStringIndex(fullBody)
	if idx == nil {
		return nil
	}

	// Find the matching brace
	rest := fullBody[idx[1]:]
	depth := 1
	var buf strings.Builder
	for _, ch := range rest {
		if ch == '{' {
			depth++
		}
		if ch == '}' {
			depth--
			if depth == 0 {
				break
			}
		}
		buf.WriteRune(ch)
	}

	return parseParamEntries(buf.String(), helperParams, pkgParams)
}

// ---------------------------------------------------------------------------
// D&D 5e route parser
// ---------------------------------------------------------------------------

func parseDnd5eRoutes(src string) []routeInfo {
	var routes []routeInfo
	lines := strings.Split(src, "\n")

	// Pre-parse known helpers from dnd5e.go context (it uses helpers from entity.go too)
	// We need to parse local variables: actorUuid, optActorUuid, abilityParams, rollOptions
	helperParams := map[string]paramDef{
		"clientIDParam": {Name: "clientId", Type: "string", From: []string{"query"}, Description: "Client ID for the Foundry world"},
		"userIDParam":   {Name: "userId", Type: "string", From: []string{"query", "body"}, Description: "Foundry user ID or username to scope permissions (omit for GM-level access)"},
		"selectedParam": {Name: "selected", Type: "boolean", From: []string{"query", "body"}, Description: "Whether to get the selected entity"},
	}
	pkgParams := parsePkgLevelParams(lines)

	// Parse inline route registrations
	for i := 0; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])

		// Pattern: r.Get("/path", helpers.CreateAPIRoute(mgr, pending, helpers.APIRouteConfig{
		// or r.Post("/path", helpers.CreateAPIRoute(mgr, pending, helpers.APIRouteConfig{
		reRoute := regexp.MustCompile(`r\.(Get|Post|Put|Delete)\("([^"]+)",\s*helpers\.CreateAPIRoute\(mgr,\s*pending,\s*helpers\.APIRouteConfig\{`)
		m := reRoute.FindStringSubmatch(line)
		if m == nil {
			continue
		}

		method := strings.ToUpper(m[1])
		path := m[2]

		// Extract doc comment above this line
		summary, desc, tag, returns := extractDocComment(lines, i)

		// Extract the config body
		bodyStart := i
		body := extractFuncBody(lines, bodyStart)

		// Handle the for-range loop for use-ability etc.
		if strings.Contains(lines[max(0, i-3)], "for _, route := range") {
			// These are generated in the loop; we handle them separately below
			continue
		}

		cfg := parseConfigBody(body, src, helperParams, pkgParams)

		ri := routeInfo{
			Method:      method,
			Path:        path,
			Summary:     summary,
			Description: desc,
			Tag:         tag,
			Returns:     returns,
			MsgType:     cfg.msgType,
			Required:    cfg.required,
			Optional:    cfg.optional,
		}
		routes = append(routes, ri)
	}

	// Handle the for-range loop for use-ability, use-feature, use-spell, use-item
	routes = append(routes, parseDnd5eForLoop(src, helperParams, pkgParams)...)

	return routes
}

func parseDnd5eForLoop(src string, helperParams map[string]paramDef, pkgParams map[string]paramDef) []routeInfo {
	// Find the for-range block
	reLoop := regexp.MustCompile(`for _, route := range \[\]struct \{[\s\S]*?\}\{([\s\S]*?)\} \{`)
	m := reLoop.FindStringSubmatch(src)
	if m == nil {
		return nil
	}

	// Extract path/type pairs
	rePair := regexp.MustCompile(`\{"([^"]+)",\s*"([^"]+)"\}`)
	pairs := rePair.FindAllStringSubmatch(m[1], -1)

	// The shared config uses: RequiredParams: actorUuid, OptionalParams: abilityParams
	actorUuid := pkgParams["actorUuid"]
	if actorUuid.Name == "" {
		actorUuid = paramDef{Name: "actorUuid", Type: "string", From: []string{"body", "query"}, Required: true, Description: "UUID of the actor"}
	}

	// Parse abilityParams from source
	abilityParams := resolveParamVar(src, "abilityParams", helperParams, pkgParams)

	// Doc comments for use-* routes
	useComments := map[string][2]string{
		"use-ability": {"Use an ability", "Activates a specific ability for an actor, optionally targeting another entity"},
		"use-feature": {"Use a feature", "Activates a specific feature for an actor, optionally targeting another entity"},
		"use-spell":   {"Use a spell", "Casts a specific spell for an actor, optionally targeting another entity"},
		"use-item":    {"Use an item", "Uses a specific item for an actor, optionally targeting another entity"},
	}

	var routes []routeInfo
	for _, pair := range pairs {
		path := pair[1]
		reqType := pair[2]

		summary := "Use " + strings.TrimPrefix(reqType, "use-")
		desc := ""
		if c, ok := useComments[reqType]; ok {
			summary = c[0]
			desc = c[1]
		}

		ri := routeInfo{
			Method:   "POST",
			Path:     path,
			Summary:  summary,
			Description: desc,
			Tag:      "Dnd5e",
			Returns:  "Result of the " + strings.ReplaceAll(reqType, "-", " ") + " operation",
			MsgType:  reqType,
			Required: []paramDef{actorUuid},
			Optional: abilityParams,
		}
		routes = append(routes, ri)
	}

	return routes
}

// ---------------------------------------------------------------------------
// Manual handler doc-comment parser (for sheet.go, filesystem.go, session.go, routes.go)
// ---------------------------------------------------------------------------

type manualHandler struct {
	funcName string
	info     routeInfo
}

func parseManualHandlers(src string) []manualHandler {
	var handlers []manualHandler
	lines := strings.Split(src, "\n")

	for i := 0; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])

		// Look for func declarations that are NOT config builders
		// (they don't return helpers.APIRouteConfig)
		if !strings.HasPrefix(line, "func ") {
			continue
		}
		if strings.Contains(line, "helpers.APIRouteConfig") {
			continue
		}
		// Must be an exported-ish handler func (lowercase is fine, we match by name)
		reFn := regexp.MustCompile(`^func (\w+Handler|\w+Subscribe)\(`)
		m := reFn.FindStringSubmatch(line)
		if m == nil {
			continue
		}
		funcName := m[1]

		// Extract doc comment
		summary, desc, tag, returns := extractDocComment(lines, i)
		if tag == "" {
			continue // no doc comment, skip
		}

		// Extract @param annotations from comments
		required, optional := extractParamAnnotations(lines, i)

		isSSE := strings.Contains(funcName, "Subscribe")

		handlers = append(handlers, manualHandler{
			funcName: funcName,
			info: routeInfo{
				Summary:     summary,
				Description: desc,
				Tag:         tag,
				Returns:     returns,
				Required:    required,
				Optional:    optional,
				IsSSE:       isSSE,
				IsManual:    true,
			},
		})
	}

	return handlers
}

// extractParamAnnotations parses @param annotations from doc comments.
// Format: @param {type} name [location] description
func extractParamAnnotations(lines []string, funcLine int) (required, optional []paramDef) {
	reParam := regexp.MustCompile(`@param\s+\{(\w+)\}\s+(\S+)\s+\[([^\]]+)\]\s+(.+)`)

	for j := funcLine - 1; j >= 0; j-- {
		trimmed := strings.TrimSpace(lines[j])
		if !strings.HasPrefix(trimmed, "//") {
			if trimmed == "" {
				continue
			}
			break
		}
		text := strings.TrimPrefix(trimmed, "//")
		text = strings.TrimSpace(text)

		if m := reParam.FindStringSubmatch(text); m != nil {
			pType := m[1]
			name := m[2]
			location := m[3]
			description := m[4]

			isRequired := strings.Contains(strings.ToLower(location), "required")
			fromSources := parseLocationToSources(location)

			p := paramDef{
				Name:        name,
				Type:        pType,
				From:        fromSources,
				Required:    isRequired,
				Description: description,
			}

			if isRequired {
				required = append(required, p)
			} else {
				optional = append(optional, p)
			}
		}
	}

	// Reverse since we read bottom-up
	for i, j := 0, len(required)-1; i < j; i, j = i+1, j-1 {
		required[i], required[j] = required[j], required[i]
	}
	for i, j := 0, len(optional)-1; i < j; i, j = i+1, j-1 {
		optional[i], optional[j] = optional[j], optional[i]
	}

	return required, optional
}

func parseLocationToSources(loc string) []string {
	loc = strings.ToLower(loc)
	// Remove "required" flag before parsing locations
	loc = strings.Replace(loc, "required", "", -1)
	loc = strings.Replace(loc, ",,", ",", -1)
	loc = strings.Trim(loc, ", ")

	var sources []string
	if strings.Contains(loc, "query") {
		sources = append(sources, "query")
	}
	if strings.Contains(loc, "body") {
		sources = append(sources, "body")
	}
	if strings.Contains(loc, "params") {
		sources = append(sources, "params")
	}
	if strings.Contains(loc, "header") {
		sources = append(sources, "header")
	}
	if len(sources) == 0 {
		sources = []string{"query"}
	}
	return sources
}

// ---------------------------------------------------------------------------
// Output 1: api-docs.json
// ---------------------------------------------------------------------------

// canvasDocumentTypes lists the valid canvas document types to expand in api-docs.json.
var canvasDocumentTypes = []string{"tokens", "tiles", "drawings", "lights", "sounds", "notes", "templates", "walls"}

// chiPathToExpress converts Chi path format to Express format for api-docs.json.
// {param} -> :param, /* -> /:path
func chiPathToExpress(path string) string {
	re := regexp.MustCompile(`\{(\w+)\}`)
	path = re.ReplaceAllString(path, ":$1")
	path = strings.Replace(path, "/*", "/:path", 1)
	return path
}

func buildAPIDocs(routes []routeInfo) map[string]interface{} {
	var endpoints []map[string]interface{}

	for _, r := range routes {
		// Expand canvas {documentType} routes into individual routes
		if strings.Contains(r.Path, "{documentType}") {
			for _, dt := range canvasDocumentTypes {
				expanded := strings.Replace(r.Path, "{documentType}", dt, 1)
				ep := map[string]interface{}{
					"method":      r.Method,
					"path":        expanded,
					"description": r.Summary,
				}

				var reqParams []map[string]interface{}
				for _, p := range r.Required {
					if p.Name == "documentType" {
						continue // already baked into the path
					}
					reqParams = append(reqParams, paramToAPIDoc(p))
				}
				var optParams []map[string]interface{}
				for _, p := range r.Optional {
					if p.Name == "documentType" {
						continue
					}
					optParams = append(optParams, paramToAPIDoc(p))
				}
				if reqParams == nil {
					reqParams = []map[string]interface{}{}
				}
				if optParams == nil {
					optParams = []map[string]interface{}{}
				}
				ep["requiredParameters"] = reqParams
				ep["optionalParameters"] = optParams
				endpoints = append(endpoints, ep)
			}
			continue
		}

		ep := map[string]interface{}{
			"method":      r.Method,
			"path":        chiPathToExpress(r.Path),
			"description": r.Summary,
		}

		var reqParams []map[string]interface{}
		for _, p := range r.Required {
			reqParams = append(reqParams, paramToAPIDoc(p))
		}
		var optParams []map[string]interface{}
		for _, p := range r.Optional {
			optParams = append(optParams, paramToAPIDoc(p))
		}

		if reqParams == nil {
			reqParams = []map[string]interface{}{}
		}
		if optParams == nil {
			optParams = []map[string]interface{}{}
		}

		ep["requiredParameters"] = reqParams
		ep["optionalParameters"] = optParams
		endpoints = append(endpoints, ep)
	}

	return map[string]interface{}{
		"version": apiVersion,
		"baseUrl": "https://your-relay-server.com",
		"authentication": map[string]interface{}{
			"required":    true,
			"headerName":  "x-api-key",
			"description": "API key must be included in the x-api-key header for all endpoints except /api/status and /api/docs",
		},
		"endpoints": endpoints,
	}
}

func paramToAPIDoc(p paramDef) map[string]interface{} {
	location := strings.Join(p.From, "/")
	// Filter out non-HTTP locations
	location = strings.Replace(location, "params", "path", 1)

	return map[string]interface{}{
		"name":        p.Name,
		"type":        p.Type,
		"description": p.Description,
		"location":    location,
	}
}

// ---------------------------------------------------------------------------
// Output 2: openapi.json
// ---------------------------------------------------------------------------

func buildOpenAPI(routes []routeInfo) map[string]interface{} {
	tags := collectTags(routes)
	paths := buildOpenAPIPaths(routes)

	return map[string]interface{}{
		"openapi": "3.0.3",
		"info": map[string]interface{}{
			"title":       "FoundryVTT REST API",
			"description": "REST API relay server for accessing Foundry VTT data remotely. Provides WebSocket connectivity and HTTP endpoints to interact with Foundry VTT worlds.",
			"version":     apiVersion,
			"license":     map[string]string{"name": "MIT"},
		},
		"servers": []map[string]interface{}{
			{"url": "http://localhost:3010", "description": "Replaced dynamically at /openapi.json"},
		},
		"security": []map[string]interface{}{
			{"apiKey": []string{}},
		},
		"tags":       tags,
		"paths":      paths,
		"components": openAPIComponents(),
	}
}

func collectTags(routes []routeInfo) []map[string]string {
	tagSet := make(map[string]bool)
	for _, r := range routes {
		if r.Tag != "" {
			tagSet[r.Tag] = true
		}
	}
	var tagNames []string
	for t := range tagSet {
		tagNames = append(tagNames, t)
	}
	sort.Strings(tagNames)

	var tags []map[string]string
	for _, t := range tagNames {
		tags = append(tags, map[string]string{"name": t})
	}
	return tags
}

func buildOpenAPIPaths(routes []routeInfo) map[string]interface{} {
	paths := make(map[string]interface{})

	for _, r := range routes {
		// Convert chi path params to OpenAPI: {documentType} stays as-is
		oaPath := r.Path
		// Convert wildcard to named param for OpenAPI
		oaPath = strings.ReplaceAll(oaPath, "/*", "/{path}")

		method := strings.ToLower(r.Method)

		operationID := fmt.Sprintf("%s_%s", method, strings.ReplaceAll(strings.TrimPrefix(oaPath, "/"), "/", "_"))
		operationID = strings.ReplaceAll(operationID, "{", "")
		operationID = strings.ReplaceAll(operationID, "}", "")
		operationID = strings.ReplaceAll(operationID, "-", "_")

		operation := map[string]interface{}{
			"summary":     r.Summary,
			"operationId": operationID,
			"responses":   buildResponses(r),
		}

		if r.Tag != "" {
			operation["tags"] = []string{r.Tag}
		}

		// Build parameters
		var params []map[string]interface{}
		allParamDefs := append(r.Required, r.Optional...)
		for _, p := range allParamDefs {
			for _, from := range p.From {
				if from == "body" {
					continue // body params go in requestBody
				}
				in := from
				if in == "params" {
					in = "path"
				}
				param := map[string]interface{}{
					"name":     p.Name,
					"in":       in,
					"required": p.Required,
					"description": p.Description,
					"schema":   map[string]string{"type": openAPIType(p.Type)},
				}
				params = append(params, param)
				break // only add once per param
			}
		}

		if len(params) > 0 {
			operation["parameters"] = params
		}

		// Build requestBody for POST/PUT/PATCH/DELETE with body params
		if method == "post" || method == "put" || method == "patch" || method == "delete" {
			bodyProps := make(map[string]interface{})
			var requiredProps []string
			for _, p := range allParamDefs {
				hasBody := false
				for _, f := range p.From {
					if f == "body" {
						hasBody = true
						break
					}
				}
				if !hasBody {
					continue
				}
				prop := map[string]interface{}{"type": openAPIType(p.Type)}
				if p.Description != "" {
					prop["description"] = p.Description
				}
				bodyProps[p.Name] = prop
				if p.Required {
					requiredProps = append(requiredProps, p.Name)
				}
			}
			if len(bodyProps) > 0 {
				schema := map[string]interface{}{
					"type":       "object",
					"properties": bodyProps,
				}
				if len(requiredProps) > 0 {
					schema["required"] = requiredProps
				}
				operation["requestBody"] = map[string]interface{}{
					"content": map[string]interface{}{
						"application/json": map[string]interface{}{
							"schema": schema,
						},
					},
				}
			}
		}

		if _, ok := paths[oaPath]; !ok {
			paths[oaPath] = make(map[string]interface{})
		}
		paths[oaPath].(map[string]interface{})[method] = operation
	}

	return paths
}

func buildResponses(r routeInfo) map[string]interface{} {
	successDesc := "Successful response"
	if r.Returns != "" {
		successDesc = r.Returns
	}

	responses := map[string]interface{}{
		"200": map[string]interface{}{
			"description": successDesc,
			"content": map[string]interface{}{
				"application/json": map[string]interface{}{
					"schema": map[string]string{"type": "object"},
				},
			},
		},
		"400": map[string]interface{}{
			"description": "Bad request - missing or invalid parameters",
			"content": map[string]interface{}{
				"application/json": map[string]interface{}{
					"schema": map[string]interface{}{
						"type": "object",
						"properties": map[string]interface{}{
							"error": map[string]string{"type": "string"},
						},
					},
					"example": map[string]string{
						"error":    "Client ID is required",
						"howToUse": "Add ?clientId=yourClientId to your request",
					},
				},
			},
		},
		"401": map[string]interface{}{
			"description": "Unauthorized - invalid or missing API key",
			"content": map[string]interface{}{
				"application/json": map[string]interface{}{
					"schema": map[string]interface{}{
						"type": "object",
						"properties": map[string]interface{}{
							"error": map[string]string{"type": "string"},
						},
					},
				},
			},
		},
	}

	// SSE endpoints return event-stream
	if r.IsSSE {
		responses["200"] = map[string]interface{}{
			"description": successDesc,
			"content": map[string]interface{}{
				"text/event-stream": map[string]interface{}{
					"schema": map[string]string{"type": "string"},
				},
			},
		}
	}

	return responses
}

func openAPIType(t string) string {
	switch t {
	case "number":
		return "number"
	case "boolean":
		return "boolean"
	case "array":
		return "array"
	case "object":
		return "object"
	default:
		return "string"
	}
}

func openAPIComponents() map[string]interface{} {
	return map[string]interface{}{
		"securitySchemes": map[string]interface{}{
			"apiKey": map[string]interface{}{
				"type": "apiKey",
				"in":   "header",
				"name": "x-api-key",
			},
		},
	}
}

// ---------------------------------------------------------------------------
// Output 3: asyncapi.json
// ---------------------------------------------------------------------------

func buildAsyncAPI(routes []routeInfo) map[string]interface{} {
	channels := make(map[string]interface{})

	// Deduplicate by message type
	seen := make(map[string]bool)

	for _, r := range routes {
		msgType := r.MsgType
		if msgType == "" || r.IsSSE || r.IsManual {
			continue // SSE and manual handlers don't have WS message types (except those with MsgType set)
		}
		if seen[msgType] {
			continue
		}
		seen[msgType] = true

		channelName := "request/" + msgType

		// Build publish (send) payload properties
		props := map[string]interface{}{
			"type": map[string]interface{}{
				"type":        "string",
				"const":       msgType,
				"description": "Message type",
			},
			"requestId": map[string]interface{}{
				"type":        "string",
				"description": "Unique request ID for correlation",
			},
		}
		var requiredFields []string
		requiredFields = append(requiredFields, "type", "requestId")

		allParams := append(r.Required, r.Optional...)
		for _, p := range allParams {
			if p.Name == "clientId" {
				continue // clientId is part of the WS connection, not the message
			}
			prop := map[string]interface{}{
				"type": openAPIType(p.Type),
			}
			if p.Description != "" {
				prop["description"] = p.Description
			}
			props[p.Name] = prop
			if p.Required {
				requiredFields = append(requiredFields, p.Name)
			}
		}

		// Result type
		resultType := msgType + "-result"

		channel := map[string]interface{}{
			"description": r.Summary,
			"publish": map[string]interface{}{
				"operationId": "send_" + msgType,
				"summary":     "Send " + msgType + " request",
				"message": map[string]interface{}{
					"name":    msgType,
					"summary": r.Summary,
					"payload": map[string]interface{}{
						"type":       "object",
						"properties": props,
						"required":   requiredFields,
					},
				},
			},
			"subscribe": map[string]interface{}{
				"operationId": "receive_" + msgType + "_result",
				"summary":     "Receive " + msgType + " response",
				"message": map[string]interface{}{
					"name": resultType,
					"payload": map[string]interface{}{
						"type": "object",
						"properties": map[string]interface{}{
							"type": map[string]interface{}{
								"type":        "string",
								"const":       resultType,
								"description": "Response message type",
							},
							"requestId": map[string]interface{}{
								"type":        "string",
								"description": "Echoed request ID",
							},
							"clientId": map[string]interface{}{
								"type":        "string",
								"description": "Foundry client ID",
							},
							"error": map[string]interface{}{
								"type":        "string",
								"description": "Error message if request failed",
							},
						},
						"required": []string{"type", "requestId"},
					},
				},
			},
		}

		channels[channelName] = channel
	}

	// Add event channels
	channels["events/chat"] = buildChatEventChannel()
	channels["events/roll"] = buildRollEventChannel()
	channels["control/subscribe"] = buildSubscribeControlChannel()
	channels["control/unsubscribe"] = buildUnsubscribeControlChannel()

	return map[string]interface{}{
		"asyncapi": "2.6.0",
		"info": map[string]interface{}{
			"title":       "FoundryVTT WebSocket API",
			"description": "Client-facing WebSocket API for bidirectional communication with Foundry VTT through the relay server. Connect to /ws/api with a token and clientId to send requests and receive real-time events.",
			"version":     apiVersion,
			"license":     map[string]string{"name": "MIT"},
		},
		"servers": map[string]interface{}{
			"production": map[string]interface{}{
				"url":         "ws://localhost:3010/ws/api",
				"protocol":    "ws",
				"description": "Replaced dynamically at /asyncapi.json",
				"variables": map[string]interface{}{
					"token":    map[string]string{"description": "API key for authentication"},
					"clientId": map[string]string{"description": "ID of the connected Foundry instance to target"},
				},
			},
		},
		"channels": channels,
		"components": map[string]interface{}{
			"schemas": map[string]interface{}{
				"ErrorMessage": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"type":      map[string]interface{}{"type": "string", "const": "error"},
						"error":     map[string]interface{}{"type": "string"},
						"requestId": map[string]interface{}{"type": "string"},
					},
					"required": []string{"type", "error"},
				},
				"ConnectedMessage": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"type":     map[string]interface{}{"type": "string", "const": "connected"},
						"clientId": map[string]interface{}{"type": "string"},
						"supportedTypes": map[string]interface{}{
							"type":  "array",
							"items": map[string]string{"type": "string"},
						},
						"eventChannels": map[string]interface{}{
							"type":  "array",
							"items": map[string]string{"type": "string"},
						},
					},
				},
			},
		},
	}
}

func buildChatEventChannel() map[string]interface{} {
	return map[string]interface{}{
		"description": "Real-time chat message events from Foundry. Subscribe with { type: \"subscribe\", channel: \"chat-events\" }.",
		"subscribe": map[string]interface{}{
			"operationId": "receive_chat_event",
			"summary":     "Receive chat events",
			"message": map[string]interface{}{
				"name": "chat-event",
				"payload": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"type":      map[string]interface{}{"type": "string", "const": "chat-event"},
						"event":     map[string]interface{}{"type": "string", "enum": []string{"create", "update", "delete"}, "description": "Event type"},
						"message":   map[string]interface{}{"type": "object", "description": "Chat message data"},
						"clientId":  map[string]interface{}{"type": "string"},
						"timestamp": map[string]interface{}{"type": "number"},
					},
					"required": []string{"type", "event"},
				},
			},
		},
	}
}

func buildRollEventChannel() map[string]interface{} {
	return map[string]interface{}{
		"description": "Real-time dice roll events from Foundry. Subscribe with { type: \"subscribe\", channel: \"roll-events\" }.",
		"subscribe": map[string]interface{}{
			"operationId": "receive_roll_event",
			"summary":     "Receive roll events",
			"message": map[string]interface{}{
				"name": "roll-event",
				"payload": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"type":      map[string]interface{}{"type": "string", "const": "roll-event"},
						"roll":      map[string]interface{}{"type": "object", "description": "Roll data"},
						"clientId":  map[string]interface{}{"type": "string"},
						"timestamp": map[string]interface{}{"type": "number"},
					},
					"required": []string{"type", "roll"},
				},
			},
		},
	}
}

func buildSubscribeControlChannel() map[string]interface{} {
	return map[string]interface{}{
		"description": "Subscribe to event channels",
		"publish": map[string]interface{}{
			"operationId": "subscribe",
			"summary":     "Subscribe to an event channel",
			"message": map[string]interface{}{
				"name": "subscribe",
				"payload": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"type":      map[string]interface{}{"type": "string", "const": "subscribe"},
						"channel":   map[string]interface{}{"type": "string", "enum": []string{"chat-events", "roll-events"}},
						"requestId": map[string]interface{}{"type": "string"},
					},
					"required": []string{"type", "channel"},
				},
			},
		},
	}
}

func buildUnsubscribeControlChannel() map[string]interface{} {
	return map[string]interface{}{
		"description": "Unsubscribe from event channels",
		"publish": map[string]interface{}{
			"operationId": "unsubscribe",
			"summary":     "Unsubscribe from an event channel",
			"message": map[string]interface{}{
				"name": "unsubscribe",
				"payload": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"type":      map[string]interface{}{"type": "string", "const": "unsubscribe"},
						"channel":   map[string]interface{}{"type": "string", "enum": []string{"chat-events", "roll-events"}, "description": "Omit to unsubscribe from all channels"},
						"requestId": map[string]interface{}{"type": "string"},
					},
					"required": []string{"type"},
				},
			},
		},
	}
}

// ---------------------------------------------------------------------------
// Auth routes (defined inline in auth.go, hardcoded here for documentation)
// ---------------------------------------------------------------------------

func buildAuthRoutes() []routeInfo {
	return []routeInfo{
		{
			Method:      "POST",
			Path:        "/auth/register",
			Summary:     "Register a new user account",
			Description: "Creates a new user with the provided email and password. Returns the new user's API key.",
			Tag:         "Auth",
			Returns:     "User object with API key",
			IsManual:    true,
			Required: []paramDef{
				{Name: "email", Type: "string", From: []string{"body"}, Required: true, Description: "The user's email address"},
				{Name: "password", Type: "string", From: []string{"body"}, Required: true, Description: "The user's password (min 8 chars, must include uppercase, lowercase, and number)"},
			},
		},
		{
			Method:      "POST",
			Path:        "/auth/login",
			Summary:     "Log in with email and password",
			Description: "Authenticates a user and returns their account details including API key.",
			Tag:         "Auth",
			Returns:     "User object with API key",
			IsManual:    true,
			Required: []paramDef{
				{Name: "email", Type: "string", From: []string{"body"}, Required: true, Description: "The user's email address"},
				{Name: "password", Type: "string", From: []string{"body"}, Required: true, Description: "The user's password"},
			},
		},
		{
			Method:      "POST",
			Path:        "/auth/regenerate-key",
			Summary:     "Regenerate API key",
			Description: "Generates a new API key for the authenticated user. The old key is invalidated.",
			Tag:         "Auth",
			Returns:     "New API key",
			IsManual:    true,
			Required: []paramDef{
				{Name: "email", Type: "string", From: []string{"body"}, Required: true, Description: "The user's email address"},
				{Name: "password", Type: "string", From: []string{"body"}, Required: true, Description: "The user's password"},
			},
		},
		{
			Method:      "GET",
			Path:        "/auth/user-data",
			Summary:     "Get user data",
			Description: "Retrieves the authenticated user's account details including usage statistics and subscription status.",
			Tag:         "Auth",
			Returns:     "User data object",
			IsManual:    true,
		},
		{
			Method:      "GET",
			Path:        "/auth/export-data",
			Summary:     "Export user data",
			Description: "Exports all stored user data for GDPR data portability compliance.",
			Tag:         "Auth",
			Returns:     "Complete user data export",
			IsManual:    true,
		},
		{
			Method:      "DELETE",
			Path:        "/auth/account",
			Summary:     "Delete user account",
			Description: "Permanently deletes the user's account and all associated data.",
			Tag:         "Auth",
			Returns:     "Confirmation of deletion",
			IsManual:    true,
			Required: []paramDef{
				{Name: "confirmEmail", Type: "string", From: []string{"body"}, Required: true, Description: "The user's email address (must match account email)"},
				{Name: "password", Type: "string", From: []string{"body"}, Required: true, Description: "The user's password for verification"},
			},
		},
		{
			Method:      "POST",
			Path:        "/auth/change-password",
			Summary:     "Change password while logged in",
			Description: "Allows an authenticated user to change their password.",
			Tag:         "Auth",
			Returns:     "Success message",
			IsManual:    true,
			Required: []paramDef{
				{Name: "currentPassword", Type: "string", From: []string{"body"}, Required: true, Description: "The user's current password"},
				{Name: "newPassword", Type: "string", From: []string{"body"}, Required: true, Description: "The new password (min 8 chars, must include uppercase, lowercase, and number)"},
			},
		},
		{
			Method:      "POST",
			Path:        "/auth/forgot-password",
			Summary:     "Request a password reset",
			Description: "Sends a password reset email if the account exists.",
			Tag:         "Auth",
			Returns:     "Generic success message",
			IsManual:    true,
			Required: []paramDef{
				{Name: "email", Type: "string", From: []string{"body"}, Required: true, Description: "The email address associated with the account"},
			},
		},
		{
			Method:      "POST",
			Path:        "/auth/reset-password",
			Summary:     "Reset password with token",
			Description: "Resets the user's password using a valid reset token.",
			Tag:         "Auth",
			Returns:     "Success message",
			IsManual:    true,
			Required: []paramDef{
				{Name: "token", Type: "string", From: []string{"body"}, Required: true, Description: "The password reset token from the email"},
				{Name: "password", Type: "string", From: []string{"body"}, Required: true, Description: "The new password (min 8 chars, must include uppercase, lowercase, and number)"},
			},
		},
		{
			Method:      "GET",
			Path:        "/auth/validate-reset-token/{token}",
			Summary:     "Validate a password reset token",
			Description: "Checks whether a password reset token is still valid before showing the reset form.",
			Tag:         "Auth",
			Returns:     "Token validity status",
			IsManual:    true,
			Required: []paramDef{
				{Name: "token", Type: "string", From: []string{"params"}, Required: true, Description: "The password reset token to validate"},
			},
		},
		{
			Method:      "POST",
			Path:        "/auth/api-keys",
			Summary:     "Create a new scoped API key",
			Description: "Creates a sub-key under the authenticated user's master key with optional restrictions.",
			Tag:         "Auth",
			Returns:     "New scoped API key details",
			IsManual:    true,
			Required: []paramDef{
				{Name: "name", Type: "string", From: []string{"body"}, Required: true, Description: "Friendly name for the key"},
			},
			Optional: []paramDef{
				{Name: "scopedClientId", Type: "string", From: []string{"body"}, Description: "Lock to specific Foundry client ID"},
				{Name: "scopedUserId", Type: "string", From: []string{"body"}, Description: "Lock to specific Foundry user ID"},
				{Name: "dailyLimit", Type: "string", From: []string{"body"}, Description: "Per-key daily request cap"},
				{Name: "expiresAt", Type: "string", From: []string{"body"}, Description: "Expiry timestamp (ISO 8601)"},
				{Name: "foundryUrl", Type: "string", From: []string{"body"}, Description: "Foundry instance URL for headless sessions"},
				{Name: "foundryUsername", Type: "string", From: []string{"body"}, Description: "Foundry login username"},
				{Name: "foundryPassword", Type: "string", From: []string{"body"}, Description: "Foundry login password (encrypted at rest)"},
			},
		},
		{
			Method:      "GET",
			Path:        "/auth/api-keys",
			Summary:     "List all scoped API keys",
			Description: "Returns all scoped keys for the authenticated user.",
			Tag:         "Auth",
			Returns:     "Array of scoped API keys",
			IsManual:    true,
		},
		{
			Method:      "DELETE",
			Path:        "/auth/api-keys/{id}",
			Summary:     "Delete a scoped API key",
			Description: "Permanently deletes a scoped key.",
			Tag:         "Auth",
			Returns:     "Success message",
			IsManual:    true,
			Required: []paramDef{
				{Name: "id", Type: "string", From: []string{"params"}, Required: true, Description: "The scoped key ID"},
			},
		},
	}
}

// ---------------------------------------------------------------------------
// Output 4: Markdown documentation
// ---------------------------------------------------------------------------

// tagToFilename maps tag names to their markdown filenames.
func tagToFilename(tag string) string {
	if tag == "" {
		return ""
	}
	// Convert first char to lowercase, keep rest as-is
	r := []rune(tag)
	r[0] = []rune(strings.ToLower(string(r[0])))[0]
	return string(r)
}

func generateMarkdown(routes []routeInfo, outDir string) {
	// Group routes by tag
	tagRoutes := make(map[string][]routeInfo)
	var tagOrder []string
	tagSeen := make(map[string]bool)

	for _, r := range routes {
		tag := r.Tag
		if tag == "" {
			continue
		}
		if !tagSeen[tag] {
			tagSeen[tag] = true
			tagOrder = append(tagOrder, tag)
		}
		tagRoutes[tag] = append(tagRoutes[tag], r)
	}

	sort.Strings(tagOrder)

	for _, tag := range tagOrder {
		rts := tagRoutes[tag]
		filename := tagToFilename(tag) + ".md"
		fpath := filepath.Join(outDir, filename)

		var buf strings.Builder

		// Frontmatter + imports
		buf.WriteString("---\ntag: " + strings.ToLower(tag) + "\n---\n")
		buf.WriteString("import Tabs from '@theme/Tabs';\nimport TabItem from '@theme/TabItem';\n\n\n")
		buf.WriteString("import ApiTester from '@site/src/components/ApiTester';\n\n")
		buf.WriteString("# " + tag + "\n")

		for idx, r := range rts {
			buf.WriteString("\n## " + r.Method + " " + chiPathToExpress(r.Path) + "\n\n")

			// Description
			if r.Summary != "" {
				buf.WriteString(r.Summary + "\n")
			}
			if r.Description != "" {
				buf.WriteString("\n" + r.Description + "\n")
			}

			// Parameters
			allParams := append(r.Required, r.Optional...)
			if len(allParams) > 0 {
				buf.WriteString("\n### Parameters\n\n")
				buf.WriteString("| Name | Type | Required | Source | Description |\n")
				buf.WriteString("|------|------|----------|--------|-------------|\n")
				for _, p := range allParams {
					reqMark := ""
					if p.Required {
						reqMark = "\u2713"
					}
					source := strings.Join(p.From, ", ")
					buf.WriteString(fmt.Sprintf("| %s | %s | %s | %s | %s |\n", p.Name, p.Type, reqMark, source, p.Description))
				}
			}

			// Returns
			if r.Returns != "" {
				buf.WriteString("\n### Returns\n\n")
				// Determine return type keyword
				retType := "object"
				retDesc := r.Returns
				lowerRet := strings.ToLower(r.Returns)
				if strings.Contains(lowerRet, "binary") || strings.Contains(lowerRet, "image") {
					retType = "binary"
				} else if strings.Contains(lowerRet, "sse") || strings.Contains(lowerRet, "event stream") {
					retType = "SSE stream"
				} else if strings.Contains(lowerRet, "array") {
					retType = "array"
				}
				buf.WriteString(fmt.Sprintf("**%s** - %s\n", retType, retDesc))
			}

			// Try It Out (ApiTester component)
			buf.WriteString("\n### Try It Out\n\n")
			buf.WriteString("<ApiTester\n")
			buf.WriteString(fmt.Sprintf("  method=\"%s\"\n", r.Method))
			buf.WriteString(fmt.Sprintf("  path=\"%s\"\n", chiPathToExpress(r.Path)))
			var paramObjs []string
			for _, p := range allParams {
				source := ""
				if len(p.From) > 0 {
					source = p.From[0]
				}
				paramObjs = append(paramObjs, fmt.Sprintf(
					`{"name":"%s","type":"%s","required":%v,"source":"%s"}`,
					p.Name, p.Type, p.Required, source))
			}
			buf.WriteString(fmt.Sprintf("  parameters={[%s]}\n", strings.Join(paramObjs, ",")))
			buf.WriteString("/>\n")

			// Separator between endpoints (not after last)
			if idx < len(rts)-1 {
				buf.WriteString("\n---\n")
			}
		}

		buf.WriteString("\n")

		if err := os.WriteFile(fpath, []byte(buf.String()), 0o644); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to write %s: %v\n", fpath, err)
			continue
		}
		fmt.Printf("  Wrote %s\n", fpath)
	}

	// Generate websocket.md
	generateWebSocketMarkdown(routes, outDir)
}

func generateWebSocketMarkdown(routes []routeInfo, outDir string) {
	var buf strings.Builder

	buf.WriteString(`---
tag: WebSocket
---

# WebSocket API

The WebSocket API provides bidirectional communication with Foundry VTT through the relay server. It supports the same operations as the REST API, plus real-time event subscriptions.

## Connection

Connect to the WebSocket endpoint with your API key and target Foundry client ID:

` + "```" + `
ws://<host>/ws/api?token=<apiKey>&clientId=<clientId>
` + "```" + `

On successful connection, you will receive a ` + "`connected`" + ` message listing all supported message types and event channels.

## Message Format

All messages are JSON objects with a ` + "`type`" + ` field. Request messages must also include a ` + "`requestId`" + ` for correlation.

### Request

` + "```json" + `
{
  "type": "search",
  "requestId": "my-unique-id",
  "query": "dragon"
}
` + "```" + `

### Response

` + "```json" + `
{
  "type": "search-result",
  "requestId": "my-unique-id",
  "clientId": "abc123",
  "results": [...]
}
` + "```" + `

## Event Subscriptions

Subscribe to real-time events from Foundry:

` + "```json" + `
{
  "type": "subscribe",
  "channel": "chat-events",
  "filters": { "speaker": "GM" }
}
` + "```" + `

Available channels: ` + "`chat-events`" + `, ` + "`roll-events`" + `

To unsubscribe:

` + "```json" + `
{
  "type": "unsubscribe",
  "channel": "chat-events"
}
` + "```" + `

## Supported Message Types

| Type | Description | Required Params |
|------|-------------|-----------------|
`)

	// Collect WS message types from routes
	seen := make(map[string]bool)
	for _, r := range routes {
		if r.MsgType == "" || r.IsSSE {
			continue
		}
		if seen[r.MsgType] {
			continue
		}
		seen[r.MsgType] = true

		desc := r.Summary
		if r.Description != "" {
			desc += " " + r.Description
		}

		var reqNames []string
		for _, p := range r.Required {
			if p.Name != "clientId" {
				reqNames = append(reqNames, "`"+p.Name+"`")
			}
		}
		reqStr := strings.Join(reqNames, ", ")
		if reqStr == "" {
			reqStr = "\u2014"
		}

		buf.WriteString(fmt.Sprintf("| `%s` | %s | %s |\n", r.MsgType, desc, reqStr))
	}

	buf.WriteString(`
## AsyncAPI Spec

The full AsyncAPI specification is available at ` + "`/asyncapi.json`" + `.
`)

	wsPath := filepath.Join(outDir, "websocket.md")
	if err := os.WriteFile(wsPath, []byte(buf.String()), 0o644); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to write %s: %v\n", wsPath, err)
		return
	}
	fmt.Printf("  Wrote %s\n", wsPath)
}

// ---------------------------------------------------------------------------
// Utilities
// ---------------------------------------------------------------------------

func mustRead(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read %s: %v\n", path, err)
		os.Exit(1)
	}
	return string(data)
}

func writeJSON(path string, data interface{}) {
	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to marshal JSON: %v\n", err)
		os.Exit(1)
	}
	b = append(b, '\n')
	if err := os.WriteFile(path, b, 0o644); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to write %s: %v\n", path, err)
		os.Exit(1)
	}
	fmt.Printf("  Wrote %s (%d bytes)\n", path, len(b))
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
