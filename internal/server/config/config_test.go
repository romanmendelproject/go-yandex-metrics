// Модуль для объявления конфигурации сервера
package config

import (
	"os"
	"testing"
)

func TestActivateEnvFlags(t *testing.T) {
	defer os.Clearenv() // Clear environment variables after test

	os.Setenv("ADDRESS", ":7070")
	os.Setenv("STORE_INTERVAL", "15")
	os.Setenv("FILE_STORAGE_PATH", "/tmp/env-db.json")
	os.Setenv("RESTORE", "true")
	os.Setenv("DATABASE_DSN", "user:password@/dbname")
	os.Setenv("KEY", "secretkey")

	activateEnvFlags()

	if FlagRunAddr != ":7070" {
		t.Errorf("expected FlagRunAddr %q, got %q", ":7070", FlagRunAddr)
	}
	if StoreInterval != 15 {
		t.Errorf("expected StoreInterval %d, got %d", 15, StoreInterval)
	}
	if FileStoragePath != "/tmp/env-db.json" {
		t.Errorf("expected FileStoragePath %q, got %q", "/tmp/env-db.json", FileStoragePath)
	}
	if !Restore {
		t.Error("expected Restore to be true")
	}
	if DBDSN != "user:password@/dbname" {
		t.Errorf("expected DBDSN %q, got %q", "user:password@/dbname", DBDSN)
	}
	if Key != "secretkey" {
		t.Errorf("expected Key %q, got %q", "secretkey", Key)
	}
}

func TestParseFlags(t *testing.T) {
	// Устанавливаем флаги командной строки, как если бы мы запустили программу
	os.Args = []string{
		"cmd",
		"-a", ":9090",
		"-l", "info",
		"-i", "10",
		"-f", "/tmp/test-metrics-db.json",
		"-r", "true",
	}

	// Вызов функции ParseFlags
	ParseFlags()

	// Проверяем значения переменных после парсинга
	if FlagRunAddr != ":9090" {
		t.Errorf("Expected FlagRunAddr to be :9090, got %s", FlagRunAddr)
	}
	if LogLevel != "info" {
		t.Errorf("Expected LogLevel to be info, got %s", LogLevel)
	}
	if StoreInterval != 10 {
		t.Errorf("Expected StoreInterval to be 10, got %d", StoreInterval)
	}
	if FileStoragePath != "/tmp/test-metrics-db.json" {
		t.Errorf("Expected FileStoragePath to be /tmp/test-metrics-db.json, got %s", FileStoragePath)
	}
	if Restore != true {
		t.Errorf("Expected Restore to be false, got %t", Restore)
	}
}
