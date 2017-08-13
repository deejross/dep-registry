package models

import (
	"errors"

	"github.com/deejross/dep-registry/util"
)

const (
	// ArchTar is a tar archive with no compression.
	ArchTar ArchType = "tar"

	// ArchTarGz is a tar archive with gzip compression.
	ArchTarGz ArchType = "tgz"

	// ArchZip is a zip archive.
	ArchZip ArchType = "zip"
)

var (
	// ErrImportNotFound indicates that the given Import was not found.
	ErrImportNotFound = errors.New("Import not found")

	// ErrVersionNotFound indicates that the given version was not found for the Import.
	ErrVersionNotFound = errors.New("Version not found")
)

// ArchType represents an archive type.
type ArchType string

// Import object.
type Import struct {
	ImportURL   string `json:"import_url,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	ProjectURL  string `json:"project_url,omitempty"`
}

// NewImport creates a new Import object.
func NewImport(url string) *Import {
	return &Import{
		ImportURL: url,
		Name:      url,
	}
}

// Version object.
type Version struct {
	ImportURL   string   `json:"import_url,omitempty"`
	Name        string   `json:"name,omitempty"`
	BinID       string   `json:"bin_id,omitempty"`
	ArchiveType ArchType `json:"archive_type,omitempty"`
}

// NewVersion creates a new Version object.
func NewVersion(m *Import, name string, archive ArchType) *Version {
	return &Version{
		ImportURL:   m.ImportURL,
		Name:        name,
		BinID:       util.UUID4(),
		ArchiveType: archive,
	}
}
