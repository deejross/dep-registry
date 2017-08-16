package auth

import (
	"os"
	"testing"
	"time"
)

var boltAddress = "auth.test.bolt"
var tm = NewTokenManager([]byte("super-secret-key"), time.Minute)
var upa *UserPassAuth

func TestNewUserPassAuth(t *testing.T) {
	os.Remove(boltAddress)

	var err error
	upa, err = NewUserPassAuth("userpass://"+boltAddress, tm)
	if err != nil {
		t.Fatal(err)
	}
}

func TestLoginUnknownUser(t *testing.T) {
	token, err := upa.Login("username", "password")
	if err == nil {
		t.Fatal("Expected an error since user does not exist")
	}

	if len(token) > 0 {
		t.Fatal("Expected empty token, got", token)
	}
}

func TestAddUser(t *testing.T) {
	user := &User{
		Username: "username",
	}

	if err := upa.AddUser(user); err != nil {
		t.Fatal(err)
	}
}

func TestGetUser(t *testing.T) {
	user, err := upa.GetUser("username")
	if err != nil {
		t.Fatal(err)
	}

	if user.Username != "username" {
		t.Fatal("Expected 'username', got", user.Username)
	}
}

func TestSetPassword(t *testing.T) {
	if err := upa.SetPassword("username", "password"); err != nil {
		t.Fatal(err)
	}
}

func TestLoginSuccess(t *testing.T) {
	token, err := upa.Login("username", "password")
	if err != nil {
		t.Fatal(err)
	}

	username, err := tm.Validate(token)
	if err != nil {
		t.Fatal("Token", token, "did not validate:", err)
	}

	if username != "username" {
		t.Fatal("Expected username, got", username)
	}
}

func TestUpdateUserFail(t *testing.T) {
	user := &User{
		Username: "no-username",
	}

	if err := upa.UpdateUser(user); err != ErrUserDoesNotExist {
		t.Fatal("Expected ErrUserDoesNotExist, got", err)
	}
}

func TestUpdateUserSuccess(t *testing.T) {
	user := &User{
		Username: "username",
		Admin:    true,
	}

	if err := upa.UpdateUser(user); err != nil {
		t.Fatal(err)
	}

	user, err := upa.GetUser("username")
	if err != nil {
		t.Fatal(err)
	}

	if !user.Admin {
		t.Fatal("Expected admin, update failed")
	}
}

func TestDeleteUser(t *testing.T) {
	if err := upa.DeleteUser("username"); err != nil {
		t.Fatal(err)
	}
}

func TestCleanup(t *testing.T) {
	if err := os.Remove(boltAddress); err != nil {
		t.Fatal(err)
	}
}
