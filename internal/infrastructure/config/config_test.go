package config

import (
	"strings"
	"testing"
)

func TestValidate_ProductionRequiresSecrets(t *testing.T) {
	cfg := &Config{
		Server: ServerConfig{Debug: false},
		Database: DatabaseConfig{
			Password: "",
		},
		JWT: JWTConfig{
			Secret: "",
		},
	}

	err := cfg.Validate()
	if err == nil {
		t.Fatal("expected validation error for missing production secrets")
	}

	msg := err.Error()
	if !strings.Contains(msg, "DB_PASSWORD") {
		t.Errorf("expected DB_PASSWORD error, got: %s", msg)
	}
	if !strings.Contains(msg, "JWT_SECRET") {
		t.Errorf("expected JWT_SECRET error, got: %s", msg)
	}
}

func TestValidate_ProductionRejectsDefaultJWTSecret(t *testing.T) {
	cfg := &Config{
		Server: ServerConfig{Debug: false},
		Database: DatabaseConfig{
			Password: "secure-db-password",
		},
		JWT: JWTConfig{
			Secret: defaultJWTSecret,
		},
	}

	err := cfg.Validate()
	if err == nil {
		t.Fatal("expected validation error for default JWT secret")
	}
	if !strings.Contains(err.Error(), "default placeholder") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidate_ProductionRejectsShortJWTSecret(t *testing.T) {
	cfg := &Config{
		Server: ServerConfig{Debug: false},
		Database: DatabaseConfig{
			Password: "secure-db-password",
		},
		JWT: JWTConfig{
			Secret: "too-short",
		},
	}

	err := cfg.Validate()
	if err == nil {
		t.Fatal("expected validation error for short JWT secret")
	}
	if !strings.Contains(err.Error(), "at least 32 characters") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidate_ProductionAcceptsValidSecrets(t *testing.T) {
	cfg := &Config{
		Server: ServerConfig{Debug: false},
		Database: DatabaseConfig{
			Password: "secure-db-password",
		},
		JWT: JWTConfig{
			Secret: "this-is-a-secure-jwt-secret-with-enough-length",
		},
	}

	if err := cfg.Validate(); err != nil {
		t.Fatalf("expected validation to pass, got: %v", err)
	}
}

func TestValidate_ProductionRequiresS3BucketWhenDriverIsS3(t *testing.T) {
	cfg := &Config{
		Server: ServerConfig{Debug: false},
		Database: DatabaseConfig{
			Password: "secure-db-password",
		},
		JWT: JWTConfig{
			Secret: "this-is-a-secure-jwt-secret-with-enough-length",
		},
		Storage: StorageConfig{
			Driver: "s3",
		},
	}

	err := cfg.Validate()
	if err == nil {
		t.Fatal("expected validation error for missing S3 bucket")
	}
	if !strings.Contains(err.Error(), "S3_BUCKET") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidate_ProductionRejectsUnsupportedStorageDriver(t *testing.T) {
	cfg := &Config{
		Server: ServerConfig{Debug: false},
		Database: DatabaseConfig{
			Password: "secure-db-password",
		},
		JWT: JWTConfig{
			Secret: "this-is-a-secure-jwt-secret-with-enough-length",
		},
		Storage: StorageConfig{
			Driver: "gcs",
		},
	}

	err := cfg.Validate()
	if err == nil {
		t.Fatal("expected validation error for unsupported storage driver")
	}
	if !strings.Contains(err.Error(), "STORAGE_DRIVER") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidate_DebugAllowsMissingSecrets(t *testing.T) {
	cfg := &Config{
		Server: ServerConfig{Debug: true},
		Database: DatabaseConfig{
			Password: "",
		},
		JWT: JWTConfig{
			Secret: "",
		},
	}

	if err := cfg.Validate(); err != nil {
		t.Fatalf("expected debug mode to allow missing secrets, got: %v", err)
	}
	if cfg.JWT.Secret != defaultJWTSecret {
		t.Fatalf("expected debug mode to use default JWT secret, got: %q", cfg.JWT.Secret)
	}
}
