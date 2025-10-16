package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
)
// go mod init app
// go get github.com/jackc/pgx/v5

func main() {
	// Параметры подключения
	host := "localhost"
	port := 15432
	user := "postgres_user"
	password := "postgres_ps"
	dbname := "postgres_db"

	// Формируем DSN в URL-формате
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		user, password, host, port, dbname)

	// Подключаемся
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, connStr)
	if err != nil {
		log.Fatalf("Не удалось подключиться к БД: %v", err)
	}
	defer conn.Close(ctx)

	// Выполняем запрос
	query := `INSERT INTO "auth_codes" ("code", "type", "uuid", "phone", "created_at", "updated_at")
	          VALUES ($1, $2, $3, $3, now(), now()) RETURNING "id"`

	var id int64
	err = conn.QueryRow(ctx, query, "123", "registration", nil).Scan(&id)
	if err != nil {
		log.Printf("Ошибка при вставке записи: %v", err)
		os.Exit(1)
	}

	fmt.Printf("Запись успешно добавлена с id=%d\n", id)
}
