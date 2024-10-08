package migrations

import (
	"context"
	"database/sql"
	"testing"

	_ "github.com/lib/pq" // Импортируйте драйвер базы данных, например PostgreSQL
)

func TestMigrations(t *testing.T) {
	// Создаем соединение с базой данных
	db, err := sql.Open("postgres", "user=username password=userpassword dbname=dbname sslmode=disable")
	if err != nil {
		t.Fatalf("Failed to open DB: %v", err)
	}
	defer db.Close()

	// Создаем контекст
	ctx := context.Background()

	// Запускаем транзакцию
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		t.Fatalf("Failed to begin transaction: %v", err)
	}

	// Запускаем миграцию
	if err := upInitTable(ctx, tx); err != nil {
		tx.Rollback() // Откатываем транзакцию в случае ошибки
		t.Fatalf("Migration up failed: %v", err)
	}

	// Проверяем, что таблица была создана
	var exists bool
	err = tx.QueryRowContext(ctx, "SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'metrics')").Scan(&exists)
	if err != nil {
		tx.Rollback()
		t.Fatalf("Failed to check if table exists: %v", err)
	}
	if !exists {
		tx.Rollback()
		t.Fatal("Expected table metrics to exist, but it does not")
	}

	// Откатываем миграцию (rollback)
	if err := downInitTable(ctx, tx); err != nil {
		t.Fatalf("Migration down failed: %v", err)
	}

	// Завершение транзакции
	if err := tx.Commit(); err != nil {
		t.Fatalf("Failed to commit transaction: %v", err)
	}
}
