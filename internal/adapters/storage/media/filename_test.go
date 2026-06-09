package media

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSanitizeFilename_AllowedExtensions(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{"photo.jpg", "photo.jpg"},
		{"photo.JPEG", "photo.jpeg"},
		{"image.png", "image.png"},
		{"animation.gif", "animation.gif"},
		{"modern.webp", "modern.webp"},
		{"document.pdf", "document.pdf"},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			got, err := SanitizeFilename(tc.input)
			require.NoError(t, err)
			assert.Equal(t, tc.expected, got)
		})
	}
}

func TestSanitizeFilename_RejectsDisallowedExtensions(t *testing.T) {
	testCases := []string{
		"cmd.sh",
		"index.php",
		"script.js",
		"data.json",
		"video.mp4",
		"archive.zip",
	}

	for _, filename := range testCases {
		t.Run(filename, func(t *testing.T) {
			_, err := SanitizeFilename(filename)
			assert.ErrorIs(t, err, ErrExtensionNotAllowed)
		})
	}
}

func TestSanitizeFilename_StripsPathTraversal(t *testing.T) {
	got, err := SanitizeFilename("../../etc/passwd.jpg")
	require.NoError(t, err)
	assert.Equal(t, "passwd.jpg", got)
}

func TestSanitizeFilename_SanitizesSpecialCharacters(t *testing.T) {
	got, err := SanitizeFilename("../../etc/my file!@#.jpg")
	require.NoError(t, err)
	assert.Equal(t, "my_file.jpg", got)
}

func TestSanitizeFilename_RejectsInvalidBasename(t *testing.T) {
	_, err := SanitizeFilename("../")
	assert.ErrorIs(t, err, ErrInvalidFilename)

	_, err = SanitizeFilename("")
	assert.ErrorIs(t, err, ErrInvalidFilename)
}
