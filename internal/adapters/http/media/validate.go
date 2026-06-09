package media

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/gabriel-vasile/mimetype"
	mediastorage "github.com/rulzi/hexa-go/internal/adapters/storage/media"
)

const (
	// MaxUploadSize is the maximum allowed upload size (10 MiB).
	MaxUploadSize = 10 << 20
)

var extensionMIMEs = map[string][]string{
	".jpg":  {"image/jpeg"},
	".jpeg": {"image/jpeg"},
	".png":  {"image/png"},
	".gif":  {"image/gif"},
	".webp": {"image/webp"},
	".pdf":  {"application/pdf"},
}

// validateUpload checks file size, sanitizes the filename, and verifies MIME type.
func validateUpload(file *multipart.FileHeader, src io.Reader) (string, io.Reader, error) {
	if file.Size <= 0 {
		return "", nil, fmt.Errorf("file is empty")
	}
	if file.Size > MaxUploadSize {
		return "", nil, fmt.Errorf("file exceeds maximum size of %d bytes", MaxUploadSize)
	}

	sanitized, err := mediastorage.SanitizeFilename(file.Filename)
	if err != nil {
		return "", nil, err
	}

	content, err := io.ReadAll(src)
	if err != nil {
		return "", nil, fmt.Errorf("failed to read file content: %w", err)
	}
	if len(content) == 0 {
		return "", nil, fmt.Errorf("file is empty")
	}

	detected := mimetype.Detect(content)
	if !isAllowedMIME(sanitized, detected.String()) {
		return "", nil, fmt.Errorf("file type %q is not allowed", detected.String())
	}

	return sanitized, bytes.NewReader(content), nil
}

func isAllowedMIME(filename, mime string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	allowed, ok := extensionMIMEs[ext]
	if !ok {
		return false
	}

	for _, allowedMIME := range allowed {
		if mime == allowedMIME {
			return true
		}
	}
	return false
}
