package http_handlers

import (
	"github.com/RVodassa/TaskReward/internal/services/auth"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	httpSwagger "github.com/swaggo/http-swagger"
	"log"
	"os"
)

func NewRouter(controller *Handler) *chi.Mux {
	const op = "http_handlers.NewRouter"

	jwtAuth, err := auth.InitJWTAuth([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		log.Printf("%s %s", op, err)
	}

	r := chi.NewRouter()

	// Публичные маршруты (без авторизации)
	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", controller.Register)
		r.Post("/login", controller.Login)
	})
	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(jwtAuth))       // Извлекает токен из запроса
		r.Use(jwtauth.Authenticator(jwtAuth))  // Проверяет токен
		r.Route("/users", func(r chi.Router) { //
			r.Get("/{userID}/status", controller.StatusUser)
			r.Post("/{userID}/tasks/{taskID}/complete", controller.TaskComplete)
			r.Get("/leaderboard", controller.LeaderBoard)
			r.Get("/tasks/activetasks", controller.GetAllActiveTask)
		})
	})

	// Маршрут для Swagger UI (публичный)
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"), // URL для swagger.json
	))

	return r
}
