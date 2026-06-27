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

func TestLoad_DebugMode(t *testing.T) {
	t.Setenv("DEBUG", "true")
	t.Setenv("DB_PASSWORD", "")
	t.Setenv("JWT_SECRET", "")
	t.Setenv("SERVER_PORT", "9090")
	t.Setenv("CORS_ALLOWED_ORIGINS", " http://a.com , http://b.com ")
	t.Setenv("JWT_EXPIRATION", "48")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.Server.Port != "9090" {
		t.Errorf("Server.Port = %q, want 9090", cfg.Server.Port)
	}
	if cfg.JWT.Expiration != 48 {
		t.Errorf("JWT.Expiration = %d, want 48", cfg.JWT.Expiration)
	}
	if len(cfg.Server.CORSAllowedOrigins) != 2 {
		t.Errorf("CORSAllowedOrigins = %v, want 2 entries", cfg.Server.CORSAllowedOrigins)
	}
}

func TestDatabaseConfig_GetDSN(t *testing.T) {
	cfg := DatabaseConfig{
		User:     "root",
		Password: "secret",
		Host:     "localhost",
		Port:     "3306",
		DBName:   "hexa_go",
	}

	dsn := cfg.GetDSN()
	expected := "root:secret@tcp(localhost:3306)/hexa_go?charset=utf8mb4&parseTime=True&loc=Local"
	if dsn != expected {
		t.Errorf("GetDSN() = %q, want %q", dsn, expected)
	}
}

func TestRedisConfig_GetAddr(t *testing.T) {
	cfg := RedisConfig{Host: "redis.local", Port: "6379"}
	if got := cfg.GetAddr(); got != "redis.local:6379" {
		t.Errorf("GetAddr() = %q, want redis.local:6379", got)
	}
}

func TestValidate_ProductionRequiresS3Region(t *testing.T) {
	cfg := &Config{
		Server: ServerConfig{Debug: false},
		Database: DatabaseConfig{
			Password: "secure-db-password",
		},
		JWT: JWTConfig{
			Secret: "this-is-a-secure-jwt-secret-with-enough-length",
		},
		Storage: StorageConfig{
			Driver:   "s3",
			S3Bucket: "my-bucket",
			S3Region: "",
		},
	}

	err := cfg.Validate()
	if err == nil {
		t.Fatal("expected validation error for missing S3 region")
	}
	if !strings.Contains(err.Error(), "S3_REGION") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidate_ProductionAcceptsS3Config(t *testing.T) {
	cfg := &Config{
		Server: ServerConfig{Debug: false},
		Database: DatabaseConfig{
			Password: "secure-db-password",
		},
		JWT: JWTConfig{
			Secret: "this-is-a-secure-jwt-secret-with-enough-length",
		},
		Storage: StorageConfig{
			Driver:   "s3",
			S3Bucket: "my-bucket",
			S3Region: "us-east-1",
		},
	}

	if err := cfg.Validate(); err != nil {
		t.Fatalf("expected validation to pass, got: %v", err)
	}
}

func TestGetEnvHelpers(t *testing.T) {
	t.Setenv("TEST_STRING_ENV", "value")
	t.Setenv("TEST_BOOL_ENV", "true")
	t.Setenv("TEST_INT_ENV", "not-int")
	t.Setenv("TEST_INT_VALID", "42")

	if got := getEnv("TEST_STRING_ENV", "default"); got != "value" {
		t.Errorf("getEnv() = %q, want value", got)
	}
	if got := getEnv("TEST_MISSING_ENV", "default"); got != "default" {
		t.Errorf("getEnv() = %q, want default", got)
	}
	if got := getEnvBool("TEST_BOOL_ENV", false); !got {
		t.Error("getEnvBool() = false, want true")
	}
	if got := getEnvBool("TEST_MISSING_BOOL", true); !got {
		t.Error("getEnvBool() = false, want default true")
	}
	if got := getEnvInt("TEST_INT_ENV", 10); got != 10 {
		t.Errorf("getEnvInt() = %d, want default 10", got)
	}
	if got := getEnvInt("TEST_INT_VALID", 10); got != 42 {
		t.Errorf("getEnvInt() = %d, want 42", got)
	}
	if got := parseCSVEnv("TEST_MISSING_CSV"); got != nil {
		t.Errorf("parseCSVEnv() = %v, want nil", got)
	}
}
