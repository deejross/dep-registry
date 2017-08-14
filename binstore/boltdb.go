package binstore

import (
	"bytes"
	"io"
	"io/ioutil"

	"github.com/boltdb/bolt"
	"github.com/deejross/dep-registry/models"
	"github.com/deejross/dep-registry/util"
)

var boltBinBucket = []byte("dep-reg-binstore")

// BoltDB store.
type BoltDB struct {
	db *bolt.DB
}

// NewBoltBinStore creates a new BoltDB interface.
func NewBoltBinStore(address string) (BinStore, error) {
	db, err := bolt.Open(address, 0600, nil)
	if err != nil {
		return nil, err
	}

	if err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(boltBinBucket)
		return err
	}); err != nil {
		return nil, err
	}

	return &BoltDB{
		db: db,
	}, nil
}

// Add a new version to the BinStore.
func (s *BoltDB) Add(v *models.Version, reader io.Reader) error {
	key := []byte(v.BinID)

	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(boltBinBucket)
		if b.Get(key) != nil {
			return util.ErrAlreadyExists
		}

		val, err := ioutil.ReadAll(reader)
		if err != nil {
			return err
		}

		return b.Put(key, val)
	})
}

// Get a Version from the BinStore.
func (s *BoltDB) Get(v *models.Version) (io.Reader, error) {
	var buf *bytes.Buffer
	key := []byte(v.BinID)

	if err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(boltBinBucket)
		val := b.Get(key)

		if val == nil {
			return util.ErrNotFound
		}

		buf = bytes.NewBuffer(val)
		return nil
	}); err != nil {
		return nil, err
	}

	return buf, nil
}

// Delete a Version from the BinStore.
func (s *BoltDB) Delete(v *models.Version) error {
	key := []byte(v.BinID)

	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(boltBinBucket)
		return b.Delete(key)
	})
}
