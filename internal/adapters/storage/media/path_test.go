package media

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateStoragePath_SanitizesTraversal(t *testing.T) {
	path, err := generateStoragePath("../../etc/passwd.jpg")
	require.NoError(t, err)
	assert.False(t, strings.Contains(path, ".."))
	assert.Contains(t, path, "passwd_")
	assert.True(t, strings.HasSuffix(path, ".jpg"))
}

func TestGenerateStoragePath_RejectsDangerousExtension(t *testing.T) {
	_, err := generateStoragePath("shell.php")
	assert.ErrorIs(t, err, ErrExtensionNotAllowed)
}

func TestValidateResolvedPath_RejectsTraversal(t *testing.T) {
	base := t.TempDir()
	err := validateResolvedPath(base, base+"/../outside.txt")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "path traversal")
}
