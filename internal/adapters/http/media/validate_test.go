package media

import (
	"bytes"
	"mime/multipart"
	"testing"

	mediastorage "github.com/rulzi/hexa-go/internal/adapters/storage/media"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateUpload_AcceptsValidJPEG(t *testing.T) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", "photo.jpg")
	require.NoError(t, err)
	_, err = part.Write(minimalJPEG)
	require.NoError(t, err)
	require.NoError(t, writer.Close())

	reader := multipart.NewReader(body, writer.Boundary())
	form, err := reader.ReadForm(MaxUploadSize)
	require.NoError(t, err)

	file := form.File["file"][0]
	src, err := file.Open()
	require.NoError(t, err)
	defer src.Close()

	filename, content, err := validateUpload(file, src)
	require.NoError(t, err)
	assert.Equal(t, "photo.jpg", filename)
	assert.Equal(t, minimalJPEG, mustReadAll(t, content))
}

func TestValidateUpload_RejectsOversizedFile(t *testing.T) {
	header := &multipart.FileHeader{
		Filename: "photo.jpg",
		Size:     MaxUploadSize + 1,
	}

	_, _, err := validateUpload(header, nopReader{})
	assert.ErrorContains(t, err, "exceeds maximum size")
}

func TestValidateUpload_RejectsDisallowedExtension(t *testing.T) {
	header := &multipart.FileHeader{
		Filename: "script.php",
		Size:     int64(len(minimalJPEG)),
	}

	_, _, err := validateUpload(header, bytes.NewReader(minimalJPEG))
	assert.ErrorIs(t, err, mediastorage.ErrExtensionNotAllowed)
}

func TestValidateUpload_RejectsMIMEMismatch(t *testing.T) {
	header := &multipart.FileHeader{
		Filename: "photo.jpg",
		Size:     12,
	}

	_, _, err := validateUpload(header, bytes.NewReader([]byte("plain text")))
	assert.ErrorContains(t, err, "not allowed")
}

type nopReader struct{}

func (nopReader) Read([]byte) (int, error) { return 0, nil }

func mustReadAll(t *testing.T, r interface{ Read([]byte) (int, error) }) []byte {
	t.Helper()
	var buf bytes.Buffer
	_, err := buf.ReadFrom(r)
	require.NoError(t, err)
	return buf.Bytes()
}
