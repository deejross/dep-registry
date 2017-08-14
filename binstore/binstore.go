package binstore

import (
	"io"

	"github.com/deejross/dep-registry/models"
)

// BinStore represents a binary store.
type BinStore interface {
	// Add to the binary store.
	Add(v *models.Version, reader io.Reader) error

	// Get a binary from the store.
	Get(v *models.Version) (io.Reader, error)

	// Delete a binary from the store.
	Delete(v *models.Version) error
}
