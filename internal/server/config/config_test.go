package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseFlags(t *testing.T) {
	// Сохраняем оригинальные аргументы командной строки
	originalArgs := os.Args
	defer func() { os.Args = originalArgs }() // Восстанавливаем оригинальные аргументы после теста

	// Задаем тестовые аргументы
	os.Args = []string{
		"cmd",
		"--address", "localhost:9090",
		"--LogLevel", "info",
		"--StoreInterval", "10",
		"--FileStoragePath", "/tmp/test-db.json",
		"--Restore", "false",
		"--DBDSN", "postgres://testuser:testpassword@localhost:5432/testdb",
		"--Key", "testkey",
		"--crypto-key", "/path/to/test/crypto.pem",
	}

	flags, err := ParseFlags()
	if err != nil {
		t.Fatalf("ParseFlags returned an error: %v", err)
	}

	// Проверяем значения
	if flags.FlagRunAddr != "localhost:9090" {
		t.Errorf("Expected FlagRunAddr to be 'localhost:9090', got '%s'", flags.FlagRunAddr)
	}
	if flags.LogLevel != "info" {
		t.Errorf("Expected LogLevel to be 'info', got '%s'", flags.LogLevel)
	}
	if flags.StoreInterval != 10 {
		t.Errorf("Expected StoreInterval to be 10, got %d", flags.StoreInterval)
	}
	if flags.FileStoragePath != "/tmp/test-db.json" {
		t.Errorf("Expected FileStoragePath to be '/tmp/test-db.json', got '%s'", flags.FileStoragePath)
	}
	if flags.DBDSN != "postgres://testuser:testpassword@localhost:5432/testdb" {
		t.Errorf("Expected DBDSN to be 'postgres://testuser:testpassword@localhost:5432/testdb', got '%s'", flags.DBDSN)
	}
	if flags.Key != "testkey" {
		t.Errorf("Expected Key to be 'testkey', got '%s'", flags.Key)
	}
	if flags.CryptoKey != "/path/to/test/crypto.pem" {
		t.Errorf("Expected CryptoKey to be '/path/to/test/crypto.pem', got '%s'", flags.CryptoKey)
	}
}

func TestReadConfig(t *testing.T) {
	// Create a temporary JSON config file for testing
	configData := `{"config": "/path/to/config.json"}`
	tempFile, err := os.CreateTemp("", "config.json")
	assert.NoError(t, err)
	defer os.Remove(tempFile.Name())

	_, err = tempFile.Write([]byte(configData))
	assert.NoError(t, err)
	tempFile.Close()

	// Set the config flag to the temporary file path
	flags := &ClientFlags{}
	flags.Config = tempFile.Name() // Set the config path directly for testing

	result, err := ReadConfig(flags)

	assert.NoError(t, err)
	assert.NotNil(t, result)
}
