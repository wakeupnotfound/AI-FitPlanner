package config

import (
	"os"
	"testing"
)

func TestSetDefaults(t *testing.T) {
	// Test that defaults are set correctly
	setDefaults()

	// This test just ensures setDefaults doesn't panic
	// Actual values are tested through InitConfig
}

func TestInitConfigWithDefaults(t *testing.T) {
	// Create a temporary config file
	configContent := `
app:
  name: "Test App"
  port: 9090
  mode: "test"
  secret_key: "test-secret"

database:
  mysql:
    host: "localhost"
    port: 3306
    user: "test_user"
    password: "test_pass"
    dbname: "test_db"
  redis:
    host: "localhost"
    port: 6379

jwt:
  secret: "test-jwt-secret"

log:
  level: "debug"
`

	// Create temp directory and file
	tmpDir := t.TempDir()
	configPath := tmpDir + "/config.yaml"

	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}

	// Set config path
	os.Setenv("FITNESS_CONFIG_PATH", tmpDir)
	defer os.Unsetenv("FITNESS_CONFIG_PATH")

	// Note: We can't fully test InitConfig without mocking viper
	// This is a placeholder for when we implement proper config testing
}

func TestGetDSN(t *testing.T) {
	GlobalConfig = &Config{
		Database: DatabaseConfig{
			MySQL: MySQLConfig{
				Host:     "testhost",
				Port:     3306,
				User:     "testuser",
				Password: "testpass",
				DBName:   "testdb",
			},
		},
	}

	expected := "testuser:testpass@tcp(testhost:3306)/testdb?charset=utf8mb4&parseTime=True&loc=Local"
	actual := GetDSN()

	if actual != expected {
		t.Errorf("GetDSN() = %v, want %v", actual, expected)
	}
}

func TestGetRedisAddr(t *testing.T) {
	GlobalConfig = &Config{
		Database: DatabaseConfig{
			Redis: RedisConfig{
				Host: "redishost",
				Port: 6380,
			},
		},
	}

	expected := "redishost:6380"
	actual := GetRedisAddr()

	if actual != expected {
		t.Errorf("GetRedisAddr() = %v, want %v", actual, expected)
	}
}
