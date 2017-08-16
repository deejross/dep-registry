package auth

import (
	"strings"

	"github.com/boltdb/bolt"
)

var boltAuthBucket = []byte("dep-reg-auth")

// UserPassAuth implements basic user password authentication using a BoltDB backend.
type UserPassAuth struct {
	db *bolt.DB
	tm *TokenManager
}

// NewUserPassAuth creates a new UserPassAuth object.
func NewUserPassAuth(address string, tm *TokenManager) (*UserPassAuth, error) {
	db, err := bolt.Open(strings.Replace(address, "userpass://", "", 1), 0600, nil)
	if err != nil {
		return nil, err
	}

	if err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(boltAuthBucket)
		return err
	}); err != nil {
		return nil, err
	}

	return &UserPassAuth{
		db: db,
		tm: tm,
	}, nil
}

// Login validates the given credentials and if successful, generates a token.
func (a *UserPassAuth) Login(username, password string) (string, error) {

}

// AddUser adds a new user.
func (a *UserPassAuth) AddUser(user *User) error {

}

// UpdateUser updates an existing user.
func (a *UserPassAuth) UpdateUser(username string, user *User) error {

}

// SetPassword sets a password for a user.
func (a *UserPassAuth) SetPassword(username, password string) error {

}

// GetUser gets a User object.
func (a *UserPassAuth) GetUser(username string) (*User, error) {

}

// DeleteUser deletes a User.
func (a *UserPassAuth) DeleteUser(username string) error {

}
