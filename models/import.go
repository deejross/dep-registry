package models

import "errors"

var (
	// ErrImportNotFound indicates that the given Import was not found.
	ErrImportNotFound = errors.New("Import not found")

	// ErrVersionNotFound indicates that the given version was not found for the Import.
	ErrVersionNotFound = errors.New("Version not found")
)

// Import object.
type Import struct {
	ImportURL   string
	Name        string
	Description string
	ProjectURL  string
}

// Version object.
type Version struct {
	ImportURL   string
	Name        string
	BinID       string
	ArchiveType string
}
