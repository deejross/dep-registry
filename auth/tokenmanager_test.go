package auth

import (
	"testing"
	"time"
)

func TestTokens(t *testing.T) {
	theUser := "the-username"

	tm := NewTokenManager([]byte("super-secret-key"), time.Minute)
	token, err := tm.Generate(theUser)
	if err != nil {
		t.Fatal(err)
	}

	username, err := tm.Validate(token)
	if err != nil {
		t.Fatal(err)
	}

	if username != theUser {
		t.Fatal("Token validation failed", username, "does not equal", theUser)
	}
}
