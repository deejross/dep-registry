package auth

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/boltdb/bolt"
)

var (
	// ErrUserAlreadyExists indicates the user already exists.
	ErrUserAlreadyExists = errors.New("User already exists")

	// ErrUserDoesNotExist indicates the given username does not exist.
	ErrUserDoesNotExist = errors.New("User does not exist")

	// ErrUsernameEmpty indicates the given username was empty.
	ErrUsernameEmpty = errors.New("Username cannot be empty")

	// ErrPasswordTooShort indicates the given password was too short.
	ErrPasswordTooShort = errors.New("Password must be at least 6 characters long")

	boltAuthBucket = []byte("dep-reg-auth")
	passSuffix     = ":pass"
)

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
	if len(username) == 0 {
		return "", ErrUsernameEmpty
	}
	if len(password) < 6 {
		return "", ErrPasswordTooShort
	}

	key := []byte(username + passSuffix)
	token := ""

	err := a.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(boltAuthBucket)
		passHash := b.Get(key)
		if passHash == nil {
			return ErrUserDoesNotExist
		}

		if err := ValidatePassword(passHash, password); err != nil {
			return err
		}

		var err error
		token, err = a.tm.Generate(username)
		return err
	})

	return token, err
}

// AddUser adds a new user.
func (a *UserPassAuth) AddUser(user *User) error {
	if len(user.Username) == 0 {
		return ErrUsernameEmpty
	}

	key := []byte(user.Username)

	return a.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(boltAuthBucket)
		if b.Get(key) != nil {
			return ErrUserAlreadyExists
		}

		bs, err := json.Marshal(user)
		if err != nil {
			return err
		}

		return b.Put(key, bs)
	})
}

// UpdateUser updates an existing user.
func (a *UserPassAuth) UpdateUser(user *User) error {
	if len(user.Username) == 0 {
		return ErrUsernameEmpty
	}

	key := []byte(user.Username)

	return a.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(boltAuthBucket)
		if b.Get(key) == nil {
			return ErrUserDoesNotExist
		}

		bs, err := json.Marshal(user)
		if err != nil {
			return err
		}

		return b.Put(key, bs)
	})
}

// SetPassword sets a password for a user.
func (a *UserPassAuth) SetPassword(username, password string) error {
	if len(username) == 0 {
		return ErrUsernameEmpty
	}
	if len(password) < 6 {
		return ErrPasswordTooShort
	}

	key := []byte(username + passSuffix)

	return a.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(boltAuthBucket)
		passHash, err := HashPassword(password)
		if err != nil {
			return err
		}

		return b.Put(key, passHash)
	})
}

// GetUser gets a User object.
func (a *UserPassAuth) GetUser(username string) (*User, error) {
	if len(username) == 0 {
		return nil, ErrUsernameEmpty
	}

	key := []byte(username)
	user := &User{}

	err := a.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(boltAuthBucket)
		userB := b.Get(key)
		if userB == nil {
			return ErrUserDoesNotExist
		}

		return json.Unmarshal(userB, user)
	})

	return user, err
}

// DeleteUser deletes a User.
func (a *UserPassAuth) DeleteUser(username string) error {
	if len(username) == 0 {
		return ErrUsernameEmpty
	}

	key := []byte(username)

	return a.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(boltAuthBucket)
		return b.Delete(key)
	})
}
