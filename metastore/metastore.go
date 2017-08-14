package metastore

import "github.com/deejross/dep-registry/models"

// MetaStore represents a metadata store.
type MetaStore interface {
	// AddUpdateImport adds an Import if it doesn't exist.
	AddImportIfNotExists(m *models.Import) error

	// UpdateImport updates an import.
	UpdateImport(m *models.Import) error

	// AddVersion adds a Version to an import.
	AddVersion(v *models.Version) error

	// GetImport gets an Import.
	GetImport(url string) (*models.Import, error)

	// GetVersions gets a list of Versions for an Import.
	GetVersions(m *models.Import) ([]*models.Version, error)

	// DeleteImport deletes an import and all its versions.
	DeleteImport(url string) error

	// DeleteVersion deletes a version.
	DeleteVersion(m *models.Import, v *models.Version) error
}
