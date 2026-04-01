package helpers

import "regexp"

// Forbidden script patterns for defense-in-depth validation.
// Module-side permission enforcement is the primary security boundary.
var forbiddenPatterns = []*regexp.Regexp{
	regexp.MustCompile(`localStorage`),
	regexp.MustCompile(`sessionStorage`),
	regexp.MustCompile(`document\.cookie`),
	regexp.MustCompile(`eval\(`),
	regexp.MustCompile(`new Worker\(`),
	regexp.MustCompile(`new SharedWorker\(`),
	regexp.MustCompile(`__proto__`),
	regexp.MustCompile(`atob\(`),
	regexp.MustCompile(`btoa\(`),
	regexp.MustCompile(`crypto\.`),
	regexp.MustCompile(`Intl\.`),
	regexp.MustCompile(`postMessage\(`),
	regexp.MustCompile(`XMLHttpRequest`),
	regexp.MustCompile(`importScripts\(`),
	regexp.MustCompile(`apiKey`),
	regexp.MustCompile(`privateKey`),
	regexp.MustCompile(`password`),
	regexp.MustCompile(`Function\(`),
	regexp.MustCompile(`Function\.constructor`),
	regexp.MustCompile(`globalThis`),
	regexp.MustCompile(`game\.settings\.set`),
	regexp.MustCompile(`Reflect\.`),
	regexp.MustCompile(`Proxy`),
	regexp.MustCompile(`import\(`),
}

// ValidateScript checks a script for forbidden patterns.
// Returns true if the script is safe, false if it contains forbidden patterns.
func ValidateScript(script string) bool {
	for _, pattern := range forbiddenPatterns {
		if pattern.MatchString(script) {
			return false
		}
	}
	return true
}
