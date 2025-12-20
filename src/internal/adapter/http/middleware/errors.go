package middleware

import "strings"

// isUUIDParsingError checks if an error is a PostgreSQL UUID parsing error
// PostgreSQL returns SQLSTATE 22P02 for invalid UUID format
func isUUIDParsingError(err error) bool {
	if err == nil {
		return false
	}
	errStr := err.Error()
	return strings.Contains(errStr, "invalid input syntax for type uuid") ||
		strings.Contains(errStr, "22P02") ||
		strings.Contains(errStr, "invalid UUID format")
}
