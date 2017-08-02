package metastore

import "github.com/deejross/dep-registry/models"

// MetaStore represents a metadata store.
type MetaStore interface {
	// AddUpdateImport adds or updates an Import.
	AddUpdateImport(m *models.Import) error

	// AddVersion adds a Version to an import.
	AddVersion(v *models.Version) error

	// GetImport gets an Import.
	GetImport(url string) (*models.Import, error)

	// GetVersions gets a list of Versions for an Import.
	GetVersions(m *models.Import) ([]*models.Version, error)
}
