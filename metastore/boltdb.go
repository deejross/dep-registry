package metastore

import (
	"encoding/json"
	"strings"

	"github.com/boltdb/bolt"
	"github.com/deejross/dep-registry/models"
	"github.com/deejross/dep-registry/util"
)

var boltMetaBucket = []byte("dep-reg-metastore")

// BoltDB MetaStore implementation.
type BoltDB struct {
	db *bolt.DB
}

// NewBoltMetaStore creates a new BoltDB interface.
func NewBoltMetaStore(address string) (MetaStore, error) {
	db, err := bolt.Open(strings.Replace(address, "boltdb://", "", 1), 0600, nil)
	if err != nil {
		return nil, err
	}

	if err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(boltMetaBucket)
		return err
	}); err != nil {
		return nil, err
	}

	return &BoltDB{
		db: db,
	}, nil
}

// AddImportIfNotExists adds an Import if it doesn't exist.
func (s *BoltDB) AddImportIfNotExists(m *models.Import) error {
	key := []byte(m.ImportURL)

	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(boltMetaBucket)
		if b.Get(key) != nil {
			return nil
		}

		val, err := json.Marshal(m)
		if err != nil {
			return err
		}

		return b.Put(key, val)
	})
}

// UpdateImport updates an existing Import.
func (s *BoltDB) UpdateImport(m *models.Import) error {
	key := []byte(m.ImportURL)

	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(boltMetaBucket)
		val, err := json.Marshal(m)
		if err != nil {
			return err
		}

		return b.Put(key, val)
	})
}

// AddVersion adds a new version to an import.
func (s *BoltDB) AddVersion(v *models.Version) error {
	key := []byte(v.ImportURL + ":versions")

	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(boltMetaBucket)
		versionsB := b.Get(key)
		var versions []*models.Version

		if versionsB == nil {
			versions = []*models.Version{}
		} else {
			err := json.Unmarshal(versionsB, versions)
			if err != nil {
				return err
			}
		}

		for _, ver := range versions {
			if ver.Name == v.Name {
				return util.ErrAlreadyExists
			}
		}

		versions = append(versions, v)

		val, err := json.Marshal(versions)
		if err != nil {
			return err
		}

		return b.Put(key, val)
	})
}

// GetImport gets an Import.
func (s *BoltDB) GetImport(url string) (*models.Import, error) {
	key := []byte(url)
	var imp *models.Import

	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(boltMetaBucket)
		val := b.Get(key)

		if val == nil {
			return util.ErrNotFound
		}

		err := json.Unmarshal(val, imp)
		return err
	})

	return imp, err
}

// GetVersions gets a list of Versions for an Import.
func (s *BoltDB) GetVersions(m *models.Import) ([]*models.Version, error) {
	key := []byte(m.ImportURL + ":versions")
	var v []*models.Version

	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(boltMetaBucket)
		val := b.Get(key)

		if val == nil {
			return util.ErrNotFound
		}

		err := json.Unmarshal(val, v)
		return err
	})

	return v, err
}

// DeleteImport deletes an import and all its versions.
func (s *BoltDB) DeleteImport(url string) error {
	key := []byte(url)

	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(boltMetaBucket)
		if err := b.Delete(key); err != nil {
			return err
		}

		key = append(key, []byte(":versions")...)
		return b.Delete(key)
	})
}

// DeleteVersion deletes a version.
func (s *BoltDB) DeleteVersion(m *models.Import, v *models.Version) error {
	key := []byte(v.ImportURL + ":versions")

	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(boltMetaBucket)
		versionsB := b.Get(key)
		var versions []*models.Version

		if versionsB == nil {
			versions = []*models.Version{}
		} else {
			err := json.Unmarshal(versionsB, versions)
			if err != nil {
				return err
			}
		}

		newVersions := []*models.Version{}
		for _, ver := range versions {
			if ver.Name == v.Name {
				continue
			}
			newVersions = append(newVersions, ver)
		}

		val, err := json.Marshal(newVersions)
		if err != nil {
			return err
		}

		return b.Put(key, val)
	})
}
