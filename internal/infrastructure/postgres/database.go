package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"os"
	"time"
)

// getDBConfig возвращает карту с переменными окружения для подключения к базе данных.
func getDBConfig() map[string]string {
	return map[string]string{
		"host":     getEnv("DB_HOST", "localhost"),
		"port":     getEnv("DB_PORT", "5432"),
		"user":     getEnv("DB_USER", ""),
		"password": getEnv("DB_PASSWORD", ""),
		"name":     getEnv("DB_NAME", ""),
		"ssl":      getEnv("DB_SSL", "disable"),
	}
}

// getEnv возвращает значение переменной окружения или значение по умолчанию.
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func ConnectDB() (*pgxpool.Pool, error) {
	const op = "postgres.ConnectDB"

	// Получаем конфигурацию базы данных
	config := getDBConfig()

	// Проверяем обязательные переменные
	if config["user"] == "" {
		return nil, fmt.Errorf("%s: переменная окружения DB_USER не задана", op)
	}
	if config["password"] == "" {
		return nil, fmt.Errorf("%s: переменная окружения DB_PASSWORD не задана", op)
	}
	if config["name"] == "" {
		return nil, fmt.Errorf("%s: переменная окружения DB_NAME не задана", op)
	}

	// Формируем строку подключения
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		config["user"], config["password"], config["host"], config["port"], config["name"], config["ssl"])

	// Подключаемся к базе данных
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db, err := pgxpool.New(ctx, connStr)
	if err != nil {
		return nil, fmt.Errorf("%s: не удалось подключиться к базе данных: %v", op, err)
	}

	// Проверяем соединение
	err = db.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: не удалось проверить подключение к базе данных: %v", op, err)
	}
	log.Println("выполнено: подключение к базе данных")

	log.Println("запуск миграций...")
	// Миграции
	err = runMigrations(connStr)
	if err != nil {
		log.Printf("%s: %s\n", op, err)
		return nil, err
	}
	log.Println("выполнено: миграции структур в базу данных")

	return db, nil
}

func runMigrations(connStr string) error {
	const op = "postgres.runMigrations"

	m, err := migrate.New("file://migrations", connStr)
	if err != nil {
		log.Printf("%s: %s\n", op, err)
		return fmt.Errorf("%s: не удалось создать объект миграции: %v", op, err)
	}
	defer func() {
		if m != nil {
			errSource, errDB := m.Close()
			if errSource != nil {
				log.Printf("%s: ошибка при закрытии миграций: %v\n", op, errSource)
			}
			if errDB != nil {
				log.Printf("%s: ошибка при закрытии миграций: %v\n", op, errDB)
			}
			return
		}
	}()

	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Printf("%s: %s\n", op, err)
		return fmt.Errorf("%s: не удалось применить миграции: %v", op, err)
	}

	return nil
}
