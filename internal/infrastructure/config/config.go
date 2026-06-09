package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

const (
	minJWTSecretLength = 32
	defaultJWTSecret   = "your-secret-key-change-in-production"
)

// Config holds all configuration for the application
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
	Storage  StorageConfig
	Logger   LoggerConfig
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port  string
	Host  string
	Debug bool
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

// RedisConfig holds Redis configuration
type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

// JWTConfig holds JWT configuration
type JWTConfig struct {
	Secret     string
	Expiration int // in hours
}

// StorageConfig holds storage configuration
type StorageConfig struct {
	Driver       string // local | s3
	BasePath     string
	BaseURL      string
	S3Bucket     string
	S3Region     string
	S3Endpoint   string
	S3PathStyle  bool
}

// LoggerConfig holds logger configuration
type LoggerConfig struct {
	Level         string // debug, info, warn, error
	Format        string // text, json
	FilePath      string // path to log file
	ReportCaller  bool
	EnableFile    bool
	EnableConsole bool
}

// Load loads configuration from environment variables and validates required values.
func Load() (*Config, error) {
	// Optional: load .env for local development (ignored if file is missing).
	_ = godotenv.Load()

	cfg := &Config{
		Server: ServerConfig{
			Port:  getEnv("SERVER_PORT", "8080"),
			Host:  getEnv("SERVER_HOST", "0.0.0.0"),
			Debug: getEnvBool("DEBUG", false),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "3306"),
			User:     getEnv("DB_USER", "root"),
			Password: os.Getenv("DB_PASSWORD"),
			DBName:   getEnv("DB_NAME", "hexa_go"),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: os.Getenv("REDIS_PASSWORD"),
		},
		JWT: JWTConfig{
			Secret:     os.Getenv("JWT_SECRET"),
			Expiration: getEnvInt("JWT_EXPIRATION", 24),
		},
		Storage: StorageConfig{
			Driver:      getEnv("STORAGE_DRIVER", "local"),
			BasePath:    getEnv("STORAGE_BASE_PATH", "./storage"),
			BaseURL:     getEnv("STORAGE_BASE_URL", "http://localhost:8080"),
			S3Bucket:    os.Getenv("S3_BUCKET"),
			S3Region:    getEnv("S3_REGION", "us-east-1"),
			S3Endpoint:  os.Getenv("S3_ENDPOINT"),
			S3PathStyle: getEnvBool("S3_USE_PATH_STYLE", false),
		},
		Logger: LoggerConfig{
			Level:         getEnv("LOG_LEVEL", "info"),
			Format:        getEnv("LOG_FORMAT", "text"),
			FilePath:      getEnv("LOG_FILE_PATH", "./logs/app.log"),
			EnableFile:    getEnvBool("LOG_ENABLE_FILE", true),
			EnableConsole: getEnvBool("LOG_ENABLE_CONSOLE", true),
			ReportCaller:  getEnvBool("LOG_REPORT_CALLER", false),
		},
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// Validate checks that required environment variables are set with safe values.
func (c *Config) Validate() error {
	var errs []string

	if c.Server.Debug {
		if c.JWT.Secret == "" {
			c.JWT.Secret = defaultJWTSecret
		}
		return nil
	}

	if c.Database.Password == "" {
		errs = append(errs, "DB_PASSWORD is required when DEBUG=false")
	}

	if err := validateJWTSecret(c.JWT.Secret); err != nil {
		errs = append(errs, err.Error())
	}

	if err := validateStorageConfig(c.Storage); err != nil {
		errs = append(errs, err.Error())
	}

	if len(errs) == 0 {
		return nil
	}

	return fmt.Errorf("configuration validation failed:\n  - %s", strings.Join(errs, "\n  - "))
}

func validateStorageConfig(storage StorageConfig) error {
	switch storage.Driver {
	case "local", "":
		return nil
	case "s3":
		if storage.S3Bucket == "" {
			return fmt.Errorf("S3_BUCKET is required when STORAGE_DRIVER=s3")
		}
		if storage.S3Region == "" {
			return fmt.Errorf("S3_REGION is required when STORAGE_DRIVER=s3")
		}
		return nil
	default:
		return fmt.Errorf("unsupported STORAGE_DRIVER %q (supported: local, s3)", storage.Driver)
	}
}

func validateJWTSecret(secret string) error {
	if secret == "" {
		return fmt.Errorf("JWT_SECRET is required when DEBUG=false")
	}
	if secret == defaultJWTSecret {
		return fmt.Errorf("JWT_SECRET must be changed from the default placeholder value")
	}
	if len(secret) < minJWTSecretLength {
		return fmt.Errorf("JWT_SECRET must be at least %d characters", minJWTSecretLength)
	}
	return nil
}

// getEnvInt gets an environment variable as integer or returns a default value
func getEnvInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	var result int
	if _, err := fmt.Sscanf(value, "%d", &result); err != nil {
		return defaultValue
	}
	return result
}

// GetDSN returns the MySQL DSN string
func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.User,
		c.Password,
		c.Host,
		c.Port,
		c.DBName,
	)
}

// GetAddr returns the Redis address string
func (c *RedisConfig) GetAddr() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvBool gets an environment variable as boolean or returns a default value
func getEnvBool(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value == "true" || value == "1" || value == "TRUE" || value == "True"
}
