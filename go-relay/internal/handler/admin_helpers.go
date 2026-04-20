package handler

import "strconv"

// parseAdminSubject parses the JWT subject claim (admin user ID) into int64.
func parseAdminSubject(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}
