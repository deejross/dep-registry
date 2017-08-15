package auth

import (
	"golang.org/x/crypto/bcrypt"
)

// User object.
type User struct {
	Username string `json:"username,omitempty"`
	Disabled bool   `json:"disabled,omitempty"`
	Admin    bool   `json:"admin,omitempty"`
}

// Auth provides authentication.
type Auth interface {
	// Login validates the given credentials and if successful, generates a token.
	Login(username, password string) (string, error)

	// AddUser adds a new user.
	AddUser(user *User) error

	// UpdateUser updates an existing user.
	UpdateUser(username string, user *User) error

	// SetPassword sets a password for a user.
	SetPassword(username, password string) error

	// GetUser gets a User object.
	GetUser(username string) (*User, error)

	// DeleteUser deletes a User.
	DeleteUser(username string) error
}

// HashPassword creates a secure hash of a password for storage.
func HashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

// ValidatePassword checks the given plaintext password against the stored hash.
func ValidatePassword(hash []byte, password string) error {
	return bcrypt.CompareHashAndPassword(hash, []byte(password))
}
