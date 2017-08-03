package binstore

import "github.com/boltdb/bolt"

// BoltDB store.
type BoltDB struct {
	db bolt.DB
}
