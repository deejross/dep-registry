package storemanager

import (
	"io"

	"github.com/deejross/dep-registry/binstore"
	"github.com/deejross/dep-registry/metastore"
	"github.com/deejross/dep-registry/models"
)

// StoreManager is the high-level manager of BinStore and MetaStore and provides transactional operations.
type StoreManager struct {
	bin  binstore.BinStore
	meta metastore.MetaStore
}

// NewStoreManager creates a new StoreManager.
func NewStoreManager(bin binstore.BinStore, meta metastore.MetaStore) StoreManager {
	return StoreManager{
		bin:  bin,
		meta: meta,
	}
}

// Add a new Version.
func (s StoreManager) Add(m *models.Import, v *models.Version, reader io.Reader) error {
	if err := s.meta.AddUpdateImport(m); err != nil {
		return err
	}
	if err := s.meta.AddVersion(v); err != nil {
		return err
	}
	return s.bin.Add(v, reader)
}

// Get an Import.
func (s StoreManager) Get(url string) (*models.Import, error) {
	return s.meta.GetImport(url)
}

// GetVersions gets a list of Versions.
func (s StoreManager) GetVersions(url string) ([]*models.Version, error) {
	m, err := s.meta.GetImport(url)
	if err != nil {
		return nil, err
	}

	return s.meta.GetVersions(m)
}

// GetVersion gets a Version.
func (s StoreManager) GetVersion(url string, versionName string) (*models.Version, error) {
	m, err := s.meta.GetImport(url)
	if err != nil {
		return nil, err
	}

	versions, err := s.meta.GetVersions(m)
	if err != nil {
		return nil, err
	}

	for _, v := range versions {
		if v.Name == versionName {
			return v, nil
		}
	}

	return nil, models.ErrVersionNotFound
}

// GetVersionBinary downloads the binary for the version.
func (s StoreManager) GetVersionBinary(v *models.Version) (io.Reader, error) {
	return s.bin.Get(v)
}