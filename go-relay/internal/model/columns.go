package model

import "strings"

// Col converts a snake_case column name to a quoted camelCase column name for
// use in SQL queries. Sequelize creates camelCase columns in both SQLite and
// PostgreSQL, so any snake_case name passed here is converted algorithmically.
//
// All-lowercase names with no underscores (id, email, status, etc.) are
// returned as-is — they are identical in both conventions. CamelCase names
// MUST be quoted: Postgres case-folds unquoted identifiers to lowercase and
// the column is not found, while SQLite's case-insensitive identifiers mask
// the bug locally (this silently broke key-request approval on Postgres).
func Col(dbType, name string) string {
	if !strings.Contains(name, "_") {
		if name == strings.ToLower(name) {
			return name
		}
		return `"` + name + `"`
	}
	parts := strings.Split(name, "_")
	result := parts[0]
	for _, part := range parts[1:] {
		if len(part) > 0 {
			result += strings.ToUpper(part[:1]) + part[1:]
		}
	}
	return `"` + result + `"`
}

// NormalizeColumnName strips underscores and lowercases for sqlx mapper matching.
// This allows db:"api_key" tags to match both "api_key" (Go-created) and "apiKey" (Sequelize-created) columns.
func NormalizeColumnName(s string) string {
	return strings.ToLower(strings.ReplaceAll(s, "_", ""))
}
