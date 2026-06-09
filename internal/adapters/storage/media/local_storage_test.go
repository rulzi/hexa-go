package media

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	mediaentity "github.com/rulzi/hexa-go/internal/domain/media/entity"
	mediaport "github.com/rulzi/hexa-go/internal/domain/media/port"
	"github.com/rulzi/hexa-go/internal/infrastructure/logger"
	"github.com/stretchr/testify/assert"
)

var testLogger logger.Logger = logger.NewSimpleLogger()

func TestNewLocalStorageAdapter_Success(t *testing.T) {
	tmpDir := t.TempDir()

	storage, err := NewLocalStorageAdapter(tmpDir, testLogger)

	assert.NoError(t, err)
	assert.NotNil(t, storage)
	assert.Equal(t, tmpDir, storage.basePath)

	// Verify directory was created
	info, err := os.Stat(tmpDir)
	assert.NoError(t, err)
	assert.True(t, info.IsDir())
}

func TestNewLocalStorageAdapter_CreatesDirectory(t *testing.T) {
	tmpDir := filepath.Join(t.TempDir(), "new", "storage", "path")

	storage, err := NewLocalStorageAdapter(tmpDir, testLogger)

	assert.NoError(t, err)
	assert.NotNil(t, storage)

	// Verify nested directory was created
	info, err := os.Stat(tmpDir)
	assert.NoError(t, err)
	assert.True(t, info.IsDir())
}

func TestLocalStorageAdapter_Save_Success(t *testing.T) {
	ctx := context.Background()
	tmpDir := t.TempDir()
	storage, err := NewLocalStorageAdapter(tmpDir, testLogger)
	assert.NoError(t, err)

	filename := "test.jpg"
	content := "test file content"
	file := strings.NewReader(content)

	path, err := storage.Save(ctx, filename, file)

	assert.NoError(t, err)
	assert.NotEmpty(t, path)
	// Path should contain the base filename (without extension) and the extension
	assert.Contains(t, path, "test")
	assert.Contains(t, path, ".jpg")
	assert.True(t, strings.HasSuffix(path, ".jpg"))

	// Verify file was created
	fullPath := filepath.Join(tmpDir, path)
	info, err := os.Stat(fullPath)
	assert.NoError(t, err)
	assert.False(t, info.IsDir())

	// Verify file content
	savedContent, err := os.ReadFile(fullPath)
	assert.NoError(t, err)
	assert.Equal(t, content, string(savedContent))
}

func TestLocalStorageAdapter_Save_CreatesDateDirectory(t *testing.T) {
	ctx := context.Background()
	tmpDir := t.TempDir()
	storage, err := NewLocalStorageAdapter(tmpDir, testLogger)
	assert.NoError(t, err)

	filename := "test.png"
	file := strings.NewReader("content")

	path, err := storage.Save(ctx, filename, file)
	assert.NoError(t, err)

	// Verify path contains date structure (YYYY/MM/DD)
	now := time.Now()
	expectedYear := now.Format("2006")
	expectedMonth := now.Format("01")
	expectedDay := now.Format("02")

	assert.Contains(t, path, expectedYear)
	assert.Contains(t, path, expectedMonth)
	assert.Contains(t, path, expectedDay)
}

func TestLocalStorageAdapter_Save_GeneratesUniqueFilename(t *testing.T) {
	ctx := context.Background()
	tmpDir := t.TempDir()
	storage, err := NewLocalStorageAdapter(tmpDir, testLogger)
	assert.NoError(t, err)

	filename := "test.jpg"
	file1 := strings.NewReader("content1")
	file2 := strings.NewReader("content2")

	// Save same filename twice in quick succession
	path1, err1 := storage.Save(ctx, filename, file1)
	path2, err2 := storage.Save(ctx, filename, file2)

	assert.NoError(t, err1)
	assert.NoError(t, err2)
	assert.NotEqual(t, path1, path2, "Paths should be different due to timestamp")

	// Both files should exist
	fullPath1 := filepath.Join(tmpDir, path1)
	fullPath2 := filepath.Join(tmpDir, path2)
	_, err1 = os.Stat(fullPath1)
	_, err2 = os.Stat(fullPath2)
	assert.NoError(t, err1)
	assert.NoError(t, err2)
}

func TestLocalStorageAdapter_Save_RejectsFilesWithoutExtension(t *testing.T) {
	ctx := context.Background()
	tmpDir := t.TempDir()
	storage, err := NewLocalStorageAdapter(tmpDir, testLogger)
	assert.NoError(t, err)

	_, err = storage.Save(ctx, "testfile", strings.NewReader("test content"))
	assert.ErrorIs(t, err, ErrExtensionNotAllowed)
}

func TestLocalStorageAdapter_Save_HandlesEmptyFile(t *testing.T) {
	ctx := context.Background()
	tmpDir := t.TempDir()
	storage, err := NewLocalStorageAdapter(tmpDir, testLogger)
	assert.NoError(t, err)

	filename := "empty.jpg"
	file := strings.NewReader("")

	path, err := storage.Save(ctx, filename, file)

	assert.NoError(t, err)
	assert.NotEmpty(t, path)

	// Verify file was created (even if empty)
	fullPath := filepath.Join(tmpDir, path)
	info, err := os.Stat(fullPath)
	assert.NoError(t, err)
	assert.Equal(t, int64(0), info.Size())
}

func TestLocalStorageAdapter_Delete_Success(t *testing.T) {
	ctx := context.Background()
	tmpDir := t.TempDir()
	storage, err := NewLocalStorageAdapter(tmpDir, testLogger)
	assert.NoError(t, err)

	// First, save a file
	filename := "test.jpg"
	file := strings.NewReader("test content")
	path, err := storage.Save(ctx, filename, file)
	assert.NoError(t, err)

	// Verify file exists
	fullPath := filepath.Join(tmpDir, path)
	_, err = os.Stat(fullPath)
	assert.NoError(t, err)

	// Delete the file
	err = storage.Delete(ctx, path)
	assert.NoError(t, err)

	// Verify file was deleted
	_, err = os.Stat(fullPath)
	assert.True(t, os.IsNotExist(err))
}

func TestLocalStorageAdapter_Delete_FileNotExists(t *testing.T) {
	ctx := context.Background()
	tmpDir := t.TempDir()
	storage, err := NewLocalStorageAdapter(tmpDir, testLogger)
	assert.NoError(t, err)

	// Try to delete non-existent file (should be idempotent)
	err = storage.Delete(ctx, "nonexistent/file.jpg")
	assert.NoError(t, err) // Should not return error (idempotent)
}

func TestLocalStorageAdapter_Delete_WithNestedPath(t *testing.T) {
	ctx := context.Background()
	tmpDir := t.TempDir()
	storage, err := NewLocalStorageAdapter(tmpDir, testLogger)
	assert.NoError(t, err)

	// Save a file (which creates nested directory structure)
	filename := "nested.jpg"
	file := strings.NewReader("content")
	path, err := storage.Save(ctx, filename, file)
	assert.NoError(t, err)

	// Delete it
	err = storage.Delete(ctx, path)
	assert.NoError(t, err)

	// Verify file was deleted
	fullPath := filepath.Join(tmpDir, path)
	_, err = os.Stat(fullPath)
	assert.True(t, os.IsNotExist(err))
}

func TestLocalStorageAdapter_Get_Success(t *testing.T) {
	ctx := context.Background()
	tmpDir := t.TempDir()
	storage, err := NewLocalStorageAdapter(tmpDir, testLogger)
	assert.NoError(t, err)

	// First, save a file
	filename := "test.jpg"
	content := "test file content for reading"
	file := strings.NewReader(content)
	path, err := storage.Save(ctx, filename, file)
	assert.NoError(t, err)

	// Get the file
	reader, err := storage.Get(ctx, path)
	assert.NoError(t, err)
	assert.NotNil(t, reader)
	defer func() {
		if err := reader.Close(); err != nil {
			t.Fatalf("Failed to close reader: %v", err)
		}
	}()

	// Read content
	readContent, err := io.ReadAll(reader)
	assert.NoError(t, err)
	assert.Equal(t, content, string(readContent))
}

func TestLocalStorageAdapter_Get_FileNotFound(t *testing.T) {
	ctx := context.Background()
	tmpDir := t.TempDir()
	storage, err := NewLocalStorageAdapter(tmpDir, testLogger)
	assert.NoError(t, err)

	// Try to get non-existent file
	reader, err := storage.Get(ctx, "nonexistent/file.jpg")

	assert.Error(t, err)
	assert.True(t, mediaentity.IsMediaNotFound(err))
	assert.Nil(t, reader)
}

func TestLocalStorageAdapter_Get_WithNestedPath(t *testing.T) {
	ctx := context.Background()
	tmpDir := t.TempDir()
	storage, err := NewLocalStorageAdapter(tmpDir, testLogger)
	assert.NoError(t, err)

	// Save a file
	filename := "nested.jpg"
	content := "nested content"
	file := strings.NewReader(content)
	path, err := storage.Save(ctx, filename, file)
	assert.NoError(t, err)

	// Get the file
	reader, err := storage.Get(ctx, path)
	assert.NoError(t, err)
	assert.NotNil(t, reader)
	defer func() {
		if err := reader.Close(); err != nil {
			t.Fatalf("Failed to close reader: %v", err)
		}
	}()

	// Verify content
	readContent, err := io.ReadAll(reader)
	assert.NoError(t, err)
	assert.Equal(t, content, string(readContent))
}

func TestLocalStorageAdapter_RoundTrip(t *testing.T) {
	ctx := context.Background()
	tmpDir := t.TempDir()
	storage, err := NewLocalStorageAdapter(tmpDir, testLogger)
	assert.NoError(t, err)

	// Save
	filename := "roundtrip.jpg"
	originalContent := "original content for round trip test"
	file := strings.NewReader(originalContent)
	path, err := storage.Save(ctx, filename, file)
	assert.NoError(t, err)

	// Get
	reader, err := storage.Get(ctx, path)
	assert.NoError(t, err)
	defer func() {
		if err := reader.Close(); err != nil {
			t.Fatalf("Failed to close reader: %v", err)
		}
	}()

	// Verify content
	readContent, err := io.ReadAll(reader)
	assert.NoError(t, err)
	assert.Equal(t, originalContent, string(readContent))

	// Delete
	err = storage.Delete(ctx, path)
	assert.NoError(t, err)

	// Verify deleted
	_, err = storage.Get(ctx, path)
	assert.Error(t, err)
	assert.True(t, mediaentity.IsMediaNotFound(err))
}

func TestLocalStorageAdapter_ImplementsInterface(t *testing.T) {
	tmpDir := t.TempDir()
	storage, err := NewLocalStorageAdapter(tmpDir, testLogger)
	assert.NoError(t, err)

	// Verify that LocalStorageAdapter implements mediaport.StoragePort interface
	var _ mediaport.StoragePort = storage
}

func TestLocalStorageAdapter_Save_ReturnsRelativePath(t *testing.T) {
	ctx := context.Background()
	tmpDir := t.TempDir()
	storage, err := NewLocalStorageAdapter(tmpDir, testLogger)
	assert.NoError(t, err)

	filename := "relative.jpg"
	file := strings.NewReader("content")
	path, err := storage.Save(ctx, filename, file)
	assert.NoError(t, err)

	// Path should be relative (not absolute)
	assert.False(t, filepath.IsAbs(path))

	// But full path should exist
	fullPath := filepath.Join(tmpDir, path)
	_, err = os.Stat(fullPath)
	assert.NoError(t, err)
}

func TestLocalStorageAdapter_Save_DifferentFileTypes(t *testing.T) {
	ctx := context.Background()
	tmpDir := t.TempDir()
	storage, err := NewLocalStorageAdapter(tmpDir, testLogger)
	assert.NoError(t, err)

	testCases := []struct {
		name    string
		content string
	}{
		{"image.jpg", "jpeg content"},
		{"image.png", "png content"},
		{"document.pdf", "pdf content"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			file := strings.NewReader(tc.content)
			path, err := storage.Save(ctx, tc.name, file)
			assert.NoError(t, err)
			assert.NotEmpty(t, path)
			// Path should contain the base filename (without extension) and the extension
			ext := filepath.Ext(tc.name)
			nameWithoutExt := tc.name[:len(tc.name)-len(ext)]
			assert.Contains(t, path, nameWithoutExt)
			if ext != "" {
				assert.True(t, strings.HasSuffix(path, ext))
			}

			// Verify file exists and has correct content
			fullPath := filepath.Join(tmpDir, path)
			savedContent, err := os.ReadFile(fullPath)
			assert.NoError(t, err)
			assert.Equal(t, tc.content, string(savedContent))
		})
	}
}

func TestLocalStorageAdapter_Save_LargeFile(t *testing.T) {
	ctx := context.Background()
	tmpDir := t.TempDir()
	storage, err := NewLocalStorageAdapter(tmpDir, testLogger)
	assert.NoError(t, err)

	// Create a larger file (1MB)
	largeContent := strings.Repeat("A", 1024*1024)
	filename := "large.pdf"
	file := strings.NewReader(largeContent)

	path, err := storage.Save(ctx, filename, file)
	assert.NoError(t, err)
	assert.NotEmpty(t, path)

	// Verify file size
	fullPath := filepath.Join(tmpDir, path)
	info, err := os.Stat(fullPath)
	assert.NoError(t, err)
	assert.Equal(t, int64(1024*1024), info.Size())
}

func TestLocalStorageAdapter_Get_ClosesFile(t *testing.T) {
	ctx := context.Background()
	tmpDir := t.TempDir()
	storage, err := NewLocalStorageAdapter(tmpDir, testLogger)
	assert.NoError(t, err)

	// Save a file
	filename := "closable.jpg"
	file := strings.NewReader("content")
	path, err := storage.Save(ctx, filename, file)
	assert.NoError(t, err)

	// Get the file
	reader, err := storage.Get(ctx, path)
	assert.NoError(t, err)
	assert.NotNil(t, reader)

	// Close the reader
	err = reader.Close()
	assert.NoError(t, err)

	// Try to read after close should fail
	_, err = io.ReadAll(reader)
	assert.Error(t, err)
}

func TestLocalStorageAdapter_Save_HandlesSpecialCharactersInFilename(t *testing.T) {
	ctx := context.Background()
	tmpDir := t.TempDir()
	storage, err := NewLocalStorageAdapter(tmpDir, testLogger)
	assert.NoError(t, err)

	testCases := []struct {
		filename       string
		expectedInPath string
	}{
		{"file with spaces.jpg", "file_with_spaces"},
		{"file-with-dashes.jpg", "file-with-dashes"},
		{"file_with_underscores.jpg", "file_with_underscores"},
		{"file123.jpg", "file123"},
	}

	for _, tc := range testCases {
		t.Run(tc.filename, func(t *testing.T) {
			content := "test content"
			file := strings.NewReader(content)
			path, err := storage.Save(ctx, tc.filename, file)
			assert.NoError(t, err)
			assert.NotEmpty(t, path)
			assert.Contains(t, path, tc.expectedInPath)

			// Verify file exists
			fullPath := filepath.Join(tmpDir, path)
			_, err = os.Stat(fullPath)
			assert.NoError(t, err)
		})
	}
}

func TestLocalStorageAdapter_Save_RejectsDisallowedExtension(t *testing.T) {
	ctx := context.Background()
	tmpDir := t.TempDir()
	storage, err := NewLocalStorageAdapter(tmpDir, testLogger)
	assert.NoError(t, err)

	_, err = storage.Save(ctx, "malware.php", strings.NewReader("<?php"))
	assert.Error(t, err)
}
