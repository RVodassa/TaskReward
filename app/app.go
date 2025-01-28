package app

import (
	"context"
	"fmt"
	"github.com/RVodassa/TaskReward/internal/handlers/http"
	"github.com/RVodassa/TaskReward/internal/infrastructure/postgres"
	"github.com/RVodassa/TaskReward/internal/infrastructure/postgres/repository"
	"github.com/RVodassa/TaskReward/internal/serve"
	"github.com/RVodassa/TaskReward/internal/services"
	"log"
	"os"
)

type App struct {
}

func NewApp() *App {
	return &App{}
}

func (app *App) Run() error {
	const op = "app.Run"

	database, err := postgres.ConnectDB()
	if err != nil {
		log.Fatal(err)
		return err
	}
	
	port := os.Getenv("SERVER_PORT")
	Repository := repository.NewRepo(database)
	Service := services.NewService(Repository)
	Controller := http_handlers.NewHandler(Service)
	router := http_handlers.NewRouter(Controller)
	newServe := serve.NewServe(port, router)

	// Генерация задач
	err = GenerateTask(10, Service)
	if err != nil {
		log.Println(err)
		return err
	}

	if err = newServe.RunServe(); err != nil {
		return fmt.Errorf("%s: %v", op, err)
	}

	return nil
}

func GenerateTask(count int, service *services.Service) error {
	for i := 0; i < count; i++ {
		err := service.AddTask(context.Background(), fmt.Sprintf("description%d", i), uint(i*10+10))
		if err != nil {
			log.Println(err)
			return err
		}
	}
	log.Printf("выполнено: генерация задач")
	return nil
}
