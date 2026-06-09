package media

import (
	"errors"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	// ErrInvalidFilename is returned when a filename is empty or unsafe.
	ErrInvalidFilename = errors.New("invalid filename")
	// ErrExtensionNotAllowed is returned when the file extension is missing or not permitted.
	ErrExtensionNotAllowed = errors.New("file extension not allowed")

	allowedExtensions = map[string]struct{}{
		".jpg":  {},
		".jpeg": {},
		".png":  {},
		".gif":  {},
		".webp": {},
		".pdf":  {},
	}

	invalidFilenameChars = regexp.MustCompile(`[^a-zA-Z0-9._-]`)
)

// AllowedExtensions returns the set of permitted file extensions (lowercase, with dot).
func AllowedExtensions() map[string]struct{} {
	return allowedExtensions
}

// SanitizeFilename strips path components, enforces an extension whitelist, and
// replaces disallowed characters in the basename.
func SanitizeFilename(filename string) (string, error) {
	base := filepath.Base(filename)
	if base == "" || base == "." || base == ".." {
		return "", ErrInvalidFilename
	}

	ext := strings.ToLower(filepath.Ext(base))
	if ext == "" {
		return "", ErrExtensionNotAllowed
	}
	if _, ok := allowedExtensions[ext]; !ok {
		return "", ErrExtensionNotAllowed
	}

	nameWithoutExt := strings.TrimSuffix(base, filepath.Ext(base))
	nameWithoutExt = invalidFilenameChars.ReplaceAllString(nameWithoutExt, "_")
	nameWithoutExt = strings.Trim(nameWithoutExt, "._-")
	if nameWithoutExt == "" {
		nameWithoutExt = "file"
	}

	return nameWithoutExt + ext, nil
}
