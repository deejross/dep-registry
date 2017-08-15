package gate

import (
	"errors"
	"io"

	"github.com/deejross/dep-registry/auth"
	"github.com/deejross/dep-registry/models"
	"github.com/deejross/dep-registry/storemanager"
)

// ErrNotAuthorized indicates the user does not have permission to perform the requested action.
var ErrNotAuthorized = errors.New("Not authorized")

// Gate validates and enforces the proper logic when interacting with the stores.
type Gate struct {
	a  auth.Auth
	sm *storemanager.StoreManager
	tm *auth.TokenManager
}

// NewGate returns a new Gate object.
func NewGate(a auth.Auth, sm *storemanager.StoreManager, tm *auth.TokenManager) *Gate {
	return &Gate{
		a:  a,
		sm: sm,
		tm: tm,
	}
}

// Login generates a token on successful login.
func (g *Gate) Login(username, password string) (string, error) {
	return g.a.Login(username, password)
}

// ParseToken parses the auth token and returns the associated User object if valid.
func (g *Gate) ParseToken(token string) (*auth.User, error) {
	if len(token) == 0 {
		return nil, nil
	}

	username, err := g.tm.Validate(token)
	if err != nil {
		return nil, err
	}

	return g.a.GetUser(username)
}

// CanUser determines if a user can perform an action, returns nil if successful.
func (g *Gate) CanUser(user *auth.User, m *models.Import, write bool) error {
	if user.Disabled {
		return ErrNotAuthorized
	}

	if user == nil {
		if write || m.Private {
			return ErrNotAuthorized
		}
	} else if user.Admin {
		return nil
	} else if write {
		for _, name := range m.Owners {
			if user.Username == name {
				return nil
			}
		}
	} else if m.Private {
		for _, name := range m.Owners {
			if user.Username == name {
				return nil
			}
		}
		for _, name := range m.Readers {
			if user.Username == name {
				return nil
			}
		}
	} else {
		return nil
	}

	return ErrNotAuthorized
}

// Add a new Version.
func (g *Gate) Add(token string, m *models.Import, v *models.Version, reader io.Reader) error {
	user, err := g.ParseToken(token)
	if err != nil {
		return err
	}

	if err := g.CanUser(user, m, true); err != nil {
		return err
	}

	return g.sm.Add(m, v, reader)
}

// Get an Import.
func (g *Gate) Get(token, url string) (*models.Import, error) {
	user, err := g.ParseToken(token)
	if err != nil {
		return nil, err
	}

	m, err := g.sm.Get(url)
	if err != nil {
		return nil, err
	}

	if err := g.CanUser(user, m, false); err != nil {
		return nil, err
	}

	return m, nil
}

// GetVersions gets a list of Versions.
func (g *Gate) GetVersions(token, url string) ([]*models.Version, error) {
	user, err := g.ParseToken(token)
	if err != nil {
		return nil, err
	}

	m, err := g.sm.Get(url)
	if err != nil {
		return nil, err
	}

	if err := g.CanUser(user, m, false); err != nil {
		return nil, err
	}

	return g.sm.GetVersions(url)
}

// GetVersion gets a Version.
func (g *Gate) GetVersion(token, url, versionName string) (*models.Version, error) {
	user, err := g.ParseToken(token)
	if err != nil {
		return nil, err
	}

	m, err := g.sm.Get(url)
	if err != nil {
		return nil, err
	}

	if err := g.CanUser(user, m, false); err != nil {
		return nil, err
	}

	return g.sm.GetVersion(url, versionName)
}

// GetVersionBinary downloads the binary for the version.
func (g *Gate) GetVersionBinary(token, url, versionName string) (io.Reader, error) {
	user, err := g.ParseToken(token)
	if err != nil {
		return nil, err
	}

	m, err := g.sm.Get(url)
	if err != nil {
		return nil, err
	}

	if err := g.CanUser(user, m, false); err != nil {
		return nil, err
	}

	v, err := g.sm.GetVersion(url, versionName)
	if err != nil {
		return nil, err
	}

	return g.sm.GetVersionBinary(v)
}

// DeleteImport deletes an import and all its versions.
func (g *Gate) DeleteImport(token, url string) error {
	user, err := g.ParseToken(token)
	if err != nil {
		return err
	}

	m, err := g.sm.Get(url)
	if err != nil {
		return err
	}

	if err := g.CanUser(user, m, true); err != nil {
		return err
	}

	return g.sm.DeleteImport(url)
}

// DeleteVersion deletes a version.
func (g *Gate) DeleteVersion(token, url, versionName string) error {
	user, err := g.ParseToken(token)
	if err != nil {
		return err
	}

	m, err := g.sm.Get(url)
	if err != nil {
		return err
	}

	if err := g.CanUser(user, m, false); err != nil {
		return err
	}

	v, err := g.sm.GetVersion(url, versionName)
	if err != nil {
		return err
	}

	return g.sm.DeleteVersion(m, v)
}
