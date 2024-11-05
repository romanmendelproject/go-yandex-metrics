package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestParseFlags(t *testing.T) {
	// Сохраняем оригинальные аргументы командной строки
	originalArgs := os.Args

	// Задаем тестовые аргументы
	defer func() { os.Args = originalArgs }() // Восстанавливаем оригинальные аргументы после теста
	os.Args = []string{"cmd", "--address", "localhost:9090", "--report-single-interval", "10", "--report-batch-interval", "60", "--poll-interval", "3", "--key", "test-key", "--rateLimit", "5", "--crypto-key", "test/path/to/key.pem"}

	flags, err := ParseFlags()
	if err != nil {
		t.Fatalf("ParseFlags returned an error: %v", err)
	}

	// Проверяем значения
	if flags.FlagReqAddr != "localhost:9090" {
		t.Errorf("Expected FlagReqAddr to be 'localhost:9090', got '%s'", flags.FlagReqAddr)
	}
	if flags.ReportSingleInterval != 10 {
		t.Errorf("Expected ReportSingleInterval to be 10, got %d", flags.ReportSingleInterval)
	}
	if flags.ReportBatchInterval != 60 {
		t.Errorf("Expected ReportBatchInterval to be 60, got %d", flags.ReportBatchInterval)
	}
	if flags.PollInterval != 3 {
		t.Errorf("Expected PollInterval to be 3, got %d", flags.PollInterval)
	}
	if flags.Key != "test-key" {
		t.Errorf("Expected Key to be 'test-key', got '%s'", flags.Key)
	}
	if flags.RateLimit != 5 {
		t.Errorf("Expected RateLimit to be 5, got %d", flags.RateLimit)
	}
	if flags.CryptoKey != "test/path/to/key.pem" {
		t.Errorf("Expected CryptoKey to be 'test/path/to/key.pem', got '%s'", flags.CryptoKey)
	}
}

// Mocking os.ReadFile function
type MockOS struct {
	mock.Mock
}

var mockOS = &MockOS{}

// Function to replace os.ReadFile for testing
func (m *MockOS) ReadFile(filename string) ([]byte, error) {
	args := m.Called(filename)
	return args.Get(0).([]byte), args.Error(1)
}

func TestReadConfig_Success(t *testing.T) {
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
