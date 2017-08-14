package util

import "errors"

var (
	// ErrNotFound Import or Version was not found.
	ErrNotFound = errors.New("Requested resource was not found")

	// ErrAlreadyExists Import or Version already exists and cannot be overwritten.
	ErrAlreadyExists = errors.New("Resource already exists and cannot be overritten")

	// ErrDisabled Import or Version has been disabled and cannot be downloaded.
	ErrDisabled = errors.New("Resource disabled")
)
