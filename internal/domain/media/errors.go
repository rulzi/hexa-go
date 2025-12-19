package media

import "errors"

var (
	// ErrMediaNotFound is returned when a media is not found
	ErrMediaNotFound = errors.New("media not found")
	// ErrNameRequired is returned when media name is missing
	ErrNameRequired = errors.New("name is required")
	// ErrPathRequired is returned when media path is missing
	ErrPathRequired = errors.New("path is required")
)
