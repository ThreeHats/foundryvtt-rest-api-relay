package model

import "strings"

// Col converts a snake_case column name to a quoted camelCase column name for
// use in SQL queries. Sequelize creates camelCase columns in both SQLite and
// PostgreSQL, so any snake_case name passed here is converted algorithmically.
//
// Names with no underscores (id, email, password, name, etc.) are returned
// as-is — they are identical in both conventions.
func Col(dbType, name string) string {
	if !strings.Contains(name, "_") {
		return name
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
