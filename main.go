package main

import (
	"github.com/RVodassa/TaskReward/app"
	"github.com/joho/godotenv"
	"log"
)

// @title TaskReward API
// @version 1.0
// @description API для работы с пользователями и выполнения задач

// @contact.name API Support
// @contact.email assadov.spb@bk.ru

// @license.name Free open source

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Укажите свой токен 'Bearer JWT_TOKEN'.
// @host localhost:8080

func main() {
	const op = "main.main"

	// Загружаем переменные окружения из .env
	log.Printf("%s: загрузка переменных окружения", op)
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("%s ошибка загрузки файла .env: %v", op, err)
	}

	log.Printf("%s: инициализация и запуск приложения", op)
	newApp := app.NewApp()
	err = newApp.Run()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%s: приложение запущено", op)
}
