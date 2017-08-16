package binstore

import (
	"errors"
	"strings"
)

// Resolve the given connection string to a specific BinStore implementation.
func Resolve(path string) (BinStore, error) {
	parts := strings.SplitN(path, "://", 2)
	if len(parts) == 1 {
		return nil, errors.New("Invalid DB path: " + path)
	}

	switch parts[0] {
	case "boltdb":
		return NewBoltBinStore(parts[1])
	default:
		return nil, errors.New("Unknown backend: " + parts[0])
	}
}
