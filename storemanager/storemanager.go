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
func NewStoreManager(bin binstore.BinStore, meta metastore.MetaStore) *StoreManager {
	return &StoreManager{
		bin:  bin,
		meta: meta,
	}
}

// Add a new Version.
func (s *StoreManager) Add(m *models.Import, v *models.Version, reader io.Reader) error {
	if err := s.meta.AddImportIfNotExists(m); err != nil {
		return err
	}
	if err := s.meta.AddVersion(v); err != nil {
		return err
	}
	return s.bin.Add(v, reader)
}

// Get an Import.
func (s *StoreManager) Get(url string) (*models.Import, error) {
	return s.meta.GetImport(url)
}

// GetVersions gets a list of Versions.
func (s *StoreManager) GetVersions(url string) ([]*models.Version, error) {
	m, err := s.meta.GetImport(url)
	if err != nil {
		return nil, err
	}

	return s.meta.GetVersions(m)
}

// GetVersion gets a Version.
func (s *StoreManager) GetVersion(url string, versionName string) (*models.Version, error) {
	m, err := s.meta.GetImport(url)
	if err != nil {
		return nil, err
	}

	versions, err := s.meta.GetVersions(m)
	if err != nil {
		return nil, err
	}

	if len(versionName) == 0 {
		return versions[len(versions)-1], nil
	}

	for _, v := range versions {
		if v.Name == versionName {
			return v, nil
		}
	}

	return nil, models.ErrVersionNotFound
}

// GetVersionBinary downloads the binary for the version.
func (s *StoreManager) GetVersionBinary(v *models.Version) (io.Reader, error) {
	return s.bin.Get(v)
}

// DisableImport disables an import and all its versions.
func (s *StoreManager) DisableImport(url string) error {
	return s.meta.DisableImport(url)
}

// DisableVersion disables a version.
func (s *StoreManager) DisableVersion(url, version string) error {
	m, err := s.meta.GetImport(url)
	if err != nil {
		return err
	}

	versions, err := s.meta.GetVersions(m)
	if err != nil {
		return err
	}

	var v *models.Version
	for _, ver := range versions {
		if ver.Name == version {
			v = ver
			break
		}
	}

	if v != nil {
		return s.meta.DisableVersion(m, v)
	}

	return nil
}

// EnableImport enables an import and all its versions.
func (s *StoreManager) EnableImport(url string) error {
	return s.meta.EnableImport(url)
}

// EnableVersion enables a version.
func (s *StoreManager) EnableVersion(url, version string) error {
	m, err := s.meta.GetImport(url)
	if err != nil {
		return err
	}

	versions, err := s.meta.GetVersions(m)
	if err != nil {
		return err
	}

	var v *models.Version
	for _, ver := range versions {
		if ver.Name == version {
			v = ver
			break
		}
	}

	if v != nil {
		return s.meta.EnableVersion(m, v)
	}

	return nil
}

// DeleteImport deletes an import and all its versions.
func (s *StoreManager) DeleteImport(url string) error {
	versions, err := s.GetVersions(url)
	if err != nil {
		return err
	}

	if err := s.meta.DeleteImport(url); err != nil {
		return err
	}

	for _, v := range versions {
		s.bin.Delete(v)
	}

	return nil
}

// DeleteVersion deletes a version.
func (s *StoreManager) DeleteVersion(m *models.Import, v *models.Version) error {
	if err := s.meta.DeleteVersion(m, v); err != nil {
		return err
	}
	return s.bin.Delete(v)
}
