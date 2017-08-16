package auth

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

// TokenManager object.
type TokenManager struct {
	signingKey []byte
	ttl        time.Duration
}

// NewTokenManager returns a new TokenManager.
func NewTokenManager(key []byte, ttl time.Duration) *TokenManager {
	return &TokenManager{
		signingKey: key,
		ttl:        ttl,
	}
}

// Generate a token for a username.
func (t *TokenManager) Generate(username string) (string, error) {
	claims := jwt.StandardClaims{
		Subject:   username,
		Issuer:    "dep-registry",
		ExpiresAt: time.Now().Add(t.ttl).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(t.signingKey)
}

// Validate and return the username for a token.
func (t *TokenManager) Validate(token string) (string, error) {
	tok, err := jwt.Parse(token, func(tok *jwt.Token) (interface{}, error) {
		return t.signingKey, nil
	})
	if err != nil {
		return "", err
	}
	if _, ok := tok.Method.(*jwt.SigningMethodHMAC); !ok {
		return "", jwt.ErrSignatureInvalid
	}

	if !tok.Valid {
		return "", jwt.ErrInvalidKey
	}

	return tok.Claims.(jwt.StandardClaims).Subject, nil
}
