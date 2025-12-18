package config_test

import (
	"os"
	"testing"
	"time"

	"github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestDSN(t *testing.T) {
	cfg := &config.DatabaseConfig{
		Host:     "localhost",
		Port:     "5432",
		User:     "user",
		Password: "password",
		DBName:   "dbname",
		SSLMode:  "disable",
	}

	expectedDSN := "postgres://user:password@localhost:5432/dbname?sslmode=disable"
	assert.Equal(t, expectedDSN, cfg.DSN())
}

func TestGetDatabaseURL(t *testing.T) {
	cfg := &config.DatabaseConfig{
		Host:     "localhost",
		Port:     "5432",
		User:     "user",
		Password: "password",
		DBName:   "dbname",
		SSLMode:  "disable",
	}

	expectedURL := "postgres://user:password@localhost:5432/dbname?sslmode=disable"
	assert.Equal(t, expectedURL, cfg.GetDatabaseURL())
}

func TestLoadConfig(t *testing.T) {
	// Set environment variables for testing
	os.Setenv("SERVER_PORT", "8080")
	os.Setenv("SERVER_HOST", "localhost")
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_USER", "user")
	os.Setenv("DB_PASSWORD", "password")
	os.Setenv("DB_NAME", "dbname")
	os.Setenv("DB_SSLMODE", "disable")
	os.Setenv("JWT_SECRET", "secret")
	os.Setenv("JWT_EXPIRATION_HOURS", "1")

	cfg, err := config.LoadConfig()
	assert.NoError(t, err)

	assert.Equal(t, "8080", cfg.Server.Port)
	assert.Equal(t, "localhost", cfg.Server.Host)
	assert.Equal(t, "localhost", cfg.Database.Host)
	assert.Equal(t, "5432", cfg.Database.Port)
	assert.Equal(t, "user", cfg.Database.User)
	assert.Equal(t, "password", cfg.Database.Password)
	assert.Equal(t, "dbname", cfg.Database.DBName)
	assert.Equal(t, "disable", cfg.Database.SSLMode)
	assert.Equal(t, "secret", cfg.JWT.Secret)
	assert.Equal(t, time.Hour, cfg.JWT.Expiration)

	// Unset environment variables after testing
	os.Unsetenv("SERVER_PORT")
	os.Unsetenv("SERVER_HOST")
	os.Unsetenv("DB_HOST")
	os.Unsetenv("DB_PORT")
	os.Unsetenv("DB_USER")
	os.Unsetenv("DB_PASSWORD")
	os.Unsetenv("DB_NAME")
	os.Unsetenv("DB_SSLMODE")
	os.Unsetenv("JWT_SECRET")
	os.Unsetenv("JWT_EXPIRATION_HOURS")
}

func TestLoadConfigError(t *testing.T) {
	// Set an invalid value for JWT_EXPIRATION_HOURS
	os.Setenv("JWT_EXPIRATION_HOURS", "invalid")

	cfg, err := config.LoadConfig()
	assert.NoError(t, err)
	assert.Equal(t, time.Duration(24)*time.Hour, cfg.JWT.Expiration)

	os.Unsetenv("JWT_EXPIRATION_HOURS")
}

func TestGetEnv(t *testing.T) {
	key := "TEST_ENV_VAR"
	expectedValue := "test_value"

	// Test when environment variable is set
	os.Setenv(key, expectedValue)
	value := config.GetEnvForTest(key, "default")
	assert.Equal(t, expectedValue, value)
	os.Unsetenv(key)

	// Test when environment variable is not set
	value = config.GetEnvForTest(key, "default")
	assert.Equal(t, "default", value)
}
