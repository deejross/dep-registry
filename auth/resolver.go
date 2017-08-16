package auth

import (
	"errors"
	"strings"
)

// Resolve the given connection string to a specific Auth implementation.
func Resolve(path string, tm *TokenManager) (Auth, error) {
	parts := strings.SplitN(path, "://", 2)
	if len(parts) == 1 {
		return nil, errors.New("Invalid DB path: " + path)
	}

	switch parts[0] {
	case "userpass":
		return NewUserPassAuth(path, tm)
	default:
		return nil, errors.New("Unknown backend: " + parts[0])
	}
}
