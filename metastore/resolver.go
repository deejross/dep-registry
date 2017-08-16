package metastore

import (
	"errors"
	"strings"
)

// Resolve the given connection string to a specific MetaStore implementation.
func Resolve(path string) (MetaStore, error) {
	parts := strings.SplitN(path, "://", 2)
	if len(parts) == 1 {
		return nil, errors.New("Invalid DB path: " + path)
	}

	switch parts[0] {
	case "boltdb":
		return NewBoltMetaStore(parts[1])
	default:
		return nil, errors.New("Unknown backend: " + parts[0])
	}
}
