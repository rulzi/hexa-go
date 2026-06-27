package logger

import (
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func discardStdout(t *testing.T) {
	t.Helper()

	orig := os.Stdout
	r, w, err := os.Pipe()
	require.NoError(t, err)
	os.Stdout = w

	t.Cleanup(func() {
		require.NoError(t, w.Close())
		_, _ = io.Copy(io.Discard, r)
		require.NoError(t, r.Close())
		os.Stdout = orig
	})
}

func TestNewLogger_ConsoleOnly(t *testing.T) {
	discardStdout(t)

	log, err := NewLogger(Config{
		Level:         "debug",
		Format:        "text",
		EnableConsole: true,
		EnableFile:    false,
	})
	require.NoError(t, err)
	require.NotNil(t, log)
	defer log.Close()

	log.Debug("debug msg")
	log.DebugWithFields("debug fields", map[string]interface{}{"k": "v"})
	log.Info("info msg")
	log.InfoWithFields("info fields", map[string]interface{}{"k": "v"})
	log.Warn("warn msg")
	log.WarnWithFields("warn fields", map[string]interface{}{"k": "v"})
	log.Error("error msg")
	log.ErrorWithFields("error fields", map[string]interface{}{"k": "v"})
}

func TestNewLogger_JSONFormat(t *testing.T) {
	discardStdout(t)

	log, err := NewLogger(Config{
		Level:         "info",
		Format:        "json",
		EnableConsole: true,
		EnableFile:    false,
		ReportCaller:  true,
	})
	require.NoError(t, err)
	defer log.Close()

	log.Info("json log")
}

func TestNewLogger_InvalidLevelDefaultsToInfo(t *testing.T) {
	discardStdout(t)

	log, err := NewLogger(Config{
		Level:         "not-a-level",
		Format:        "text",
		EnableConsole: true,
	})
	require.NoError(t, err)
	defer log.Close()

	log.Info("still works")
}

func TestNewLogger_FileOutput(t *testing.T) {
	tmpDir := t.TempDir()
	logPath := filepath.Join(tmpDir, "logs", "app.log")

	log, err := NewLogger(Config{
		Level:         "info",
		Format:        "text",
		FilePath:      logPath,
		EnableFile:    true,
		EnableConsole: false,
	})
	require.NoError(t, err)

	log.Info("written to file")

	err = log.Close()
	require.NoError(t, err)

	info, err := os.Stat(logPath)
	require.NoError(t, err)
	assert.True(t, info.Size() > 0)
}

func TestNewLogger_NoWritersDefaultsToStdout(t *testing.T) {
	discardStdout(t)

	log, err := NewLogger(Config{
		Level:         "info",
		Format:        "text",
		EnableConsole: false,
		EnableFile:    false,
	})
	require.NoError(t, err)
	defer log.Close()

	log.Info("stdout fallback")
}

func TestNewSimpleLogger(t *testing.T) {
	discardStdout(t)

	log := NewSimpleLogger()
	require.NotNil(t, log)
	log.Info("simple logger")
	assert.NoError(t, log.Close())
}

func TestLogrusLogger_WithFields(t *testing.T) {
	discardStdout(t)

	log := NewSimpleLogger()
	child := log.WithFields(map[string]interface{}{"component": "test"})
	child.Info("with preset fields")
}

func TestLogrusLogger_CloseWithoutFile(t *testing.T) {
	log := NewSimpleLogger()
	assert.NoError(t, log.Close())
}
