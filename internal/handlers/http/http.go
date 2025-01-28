package http_handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	_ "github.com/RVodassa/TaskReward/docs"
	"github.com/RVodassa/TaskReward/internal/api"
	"github.com/RVodassa/TaskReward/internal/domain/models"
	"github.com/RVodassa/TaskReward/internal/services"
	"github.com/RVodassa/TaskReward/internal/services/auth"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"strconv"
)

var (
	ErrCredentialsRequired  = errors.New("ошибка: логин и пароль обязательны")
	ErrInvalidJSON          = errors.New("ошибка: неверный формат JSON")
	ErrUserAlreadyExist     = errors.New("ошибка: пользователь с таким логином уже существует")
	ErrUserNotFound         = errors.New("ошибка: пользователь найден")
	ErrInternalServer       = errors.New("ошибка: внутренняя ошибка сервера, обратитесь к администратору")
	ErrInvalidReferID       = errors.New("ошибка: некорректный refer_id")
	ErrInvalidID            = errors.New("ошибка: некорректный user_id")
	ErrInvalidTaskID        = errors.New("ошибка: некорректный task_id")
	ErrReferUserNotFound    = errors.New("ошибка: refer с указанным id не найден")
	ErrIncorrectPassword    = errors.New("ошибка: не правильный логин или пароль")
	ErrTaskNotFound         = errors.New("ошибка: задача не найдена")
	ErrTaskAlreadyCompleted = errors.New("ошибка: задача уже выполнена")
)

type UserServiceProvider interface {
	RegisterUser(ctx context.Context, login string, password string, referID uint) (*models.User, error)
	Login(ctx context.Context, login, password string) error
	StatusUser(ctx context.Context, userID uint) (*models.User, error)
	TaskComplete(ctx context.Context, taskID uint, userID uint) (*models.Task, error)
	GetListTopUsers(ctx context.Context) ([]*models.User, error)
	GetAllActiveTask(ctx context.Context) ([]*models.Task, error)
}

type Handler struct {
	userService UserServiceProvider
}

func NewHandler(userService UserServiceProvider) *Handler {
	return &Handler{userService: userService}
}

// GetAllActiveTask godoc
// @Summary Получить список активных задач
// @Description возвращает список активных задач
// @Tags Tasks
// @Produce json
// @Success 200 {object} api.GetAllTasksResponse "Успешно"
// @Failure 403 {object} api.ErrorResponse "Unauthorized"
// @Failure 400 {object} api.ErrorResponse "Ошибка клиента"
// @Failure 500 {object} api.ErrorResponse "Ошибка на сервере"
// @Router /users/tasks/activetasks [get]
// @security BearerAuth
func (h *Handler) GetAllActiveTask(w http.ResponseWriter, r *http.Request) {
	const op = "http_handlers.GetAllActiveTask"

	listTask, err := h.userService.GetAllActiveTask(r.Context())
	if err != nil {
		log.Printf("%s %v", op, err)
		Responder(w, http.StatusInternalServerError, api.ErrorResponse{Status: false, Message: ErrInternalServer.Error()})
		return
	}

	if len(listTask) == 0 {
		Responder(w, http.StatusOK, api.GetAllTasksResponse{
			Status:  true,
			Message: "Активные задачи не найдены",
			Tasks:   []*models.Task{},
		})
		return
	}

	resp := api.GetAllTasksResponse{
		Status:  true,
		Message: fmt.Sprintf("Список активных задач\nКол-во задач: %d", len(listTask)),
		Tasks:   listTask,
	}

	Responder(w, http.StatusOK, resp)
}

// LeaderBoard godoc
// @Summary Получить список лидеров
// @Description Возвращает топ 10 лидеров по балансу
// @Tags Users
// @Produce json
// @Success 200 {object} api.LeaderBoardResponse "Успешно"
// @Failure 403 {object} api.ErrorResponse "Unauthorized"
// @Failure 400 {object} api.ErrorResponse "Ошибка клиента"
// @Failure 500 {object} api.ErrorResponse "Ошибка на сервере"
// @Router /users/leaderboard [get]
// @security BearerAuth
func (h *Handler) LeaderBoard(w http.ResponseWriter, r *http.Request) {
	const op = "http_handlers.LeaderBoard"

	users, err := h.userService.GetListTopUsers(r.Context())
	if err != nil {
		log.Printf("%s %v", op, err)
		Responder(w, http.StatusInternalServerError, api.ErrorResponse{Status: false, Message: ErrInternalServer.Error()})
		return
	}

	if len(users) == 0 {
		Responder(w, http.StatusOK, api.LeaderBoardResponse{
			Status:     true,
			Message:    "Пользователи не найдены",
			ListLeader: []*models.User{},
		})
		return
	}

	resp := api.LeaderBoardResponse{
		Status:     true,
		Message:    fmt.Sprintf("Доска лидеров\nКол.во: %d", len(users)),
		ListLeader: users,
	}

	Responder(w, http.StatusOK, resp)
}

// TaskComplete godoc
// @Summary Выполнить задачу
// @Description Возвращает информацию о выполненной задаче
// @Tags Tasks
// Accept json
// @Produce json
// @Param userID path string true "ID пользователя"
// @Param taskID path string true "ID задачи"
// @Success 200 {object} api.TaskCompletedResponse "Успешно"
// @Failure 403 {object} api.ErrorResponse "Unauthorized"
// @Failure 400 {object} api.ErrorResponse "Ошибка клиента"
// @Failure 500 {object} api.ErrorResponse "Ошибка на сервере"
// @Router /users/{userID}/tasks/{taskID}/complete [post]
// @security BearerAuth
func (h *Handler) TaskComplete(w http.ResponseWriter, r *http.Request) {
	const op = "http_handlers.TaskComplete"

	userIDStr := chi.URLParam(r, "userID")
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil || userID <= 0 {
		log.Printf("%s %s %v", op, r.URL, err)
		Responder(w, http.StatusBadRequest, api.ErrorResponse{Status: false, Message: ErrInvalidID.Error()})
		return
	}

	taskIdStr := chi.URLParam(r, "taskID")
	taskID, err := strconv.ParseUint(taskIdStr, 10, 64)
	if err != nil || taskID <= 0 {
		log.Printf("%s %s %v", op, r.URL, err)
		Responder(w, http.StatusBadRequest, api.ErrorResponse{Status: false, Message: ErrInvalidTaskID.Error()})
		return
	}

	task, err := h.userService.TaskComplete(r.Context(), uint(taskID), uint(userID))
	if err != nil {
		switch {
		case errors.Is(err, services.ErrUserNotFound):
			Responder(w, http.StatusNotFound, api.ErrorResponse{Status: false, Message: ErrUserNotFound.Error()})
			return
		case errors.Is(err, services.ErrTaskNotFound):
			Responder(w, http.StatusNotFound, api.ErrorResponse{Status: false, Message: ErrTaskNotFound.Error()})
			return
		case errors.Is(err, services.ErrTaskAlreadyCompleted):
			Responder(w, http.StatusConflict, api.ErrorResponse{Status: false, Message: ErrTaskAlreadyCompleted.Error()})
			return
		default:
			log.Printf("%s %s %v", op, r.URL, err)
			Responder(w, http.StatusInternalServerError, api.ErrorResponse{Status: false, Message: ErrInternalServer.Error()})
			return
		}
	}

	resp := api.TaskCompletedResponse{
		Status:  true,
		Message: "Задача выполнена",
		Task:    task,
	}

	Responder(w, http.StatusOK, resp)
}

// StatusUser godoc
// @Summary Получить информацию о пользователе по ID
// @Description Возвращает информацию о пользователе в случае успешной операции
// @Tags Users
// Accept json
// @Produce json
// @Param userID path string true "ID пользователя"
// @Success 200 {object} api.StatusUserResponse "Успешно"
// @Failure 403 {object} api.ErrorResponse "Unauthorized"
// @Failure 400 {object} api.ErrorResponse "Ошибка клиента"
// @Failure 500 {object} api.ErrorResponse "Ошибка на сервере"
// @Router /users/{userID}/status [get]
// @security BearerAuth
func (h *Handler) StatusUser(w http.ResponseWriter, r *http.Request) {
	const op = "http_handlers.StatusUser"

	idStr := chi.URLParam(r, "userID")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil || id <= 0 {
		log.Printf("%s %s %v", op, r.URL, err)
		Responder(w, http.StatusBadRequest, api.ErrorResponse{Status: false, Message: ErrInvalidID.Error()})
		return
	}

	user, err := h.userService.StatusUser(r.Context(), uint(id))
	if err != nil {
		if errors.Is(err, services.ErrUserNotFound) {
			Responder(w, http.StatusNotFound, api.ErrorResponse{Status: false, Message: ErrUserNotFound.Error()})
			return
		}
		log.Printf("%s %s %v", op, r.URL, err)
		Responder(w, http.StatusInternalServerError, api.ErrorResponse{Status: false, Message: ErrInternalServer.Error()})
		return
	}

	resp := api.StatusUserResponse{
		Status:  true,
		Message: "OK",
		User:    user,
	}

	Responder(w, http.StatusOK, resp)
}

// Register godoc
// @Summary Регистрация пользователя
// @Description Создает нового пользователя, возвращает информацию о новом пользователе.
// @Tags auth
// @Accept json
// @Produce json
// @Param referID query string true "ID реферала, если нет укажите 0"
// @Param request body api.AuthRequest true "Логин и пароль"
// @Success 200 {object} api.StatusUserResponse "Успешная регистрация"
// @Failure 403 {object} api.ErrorResponse "Unauthorized"
// @Failure 400 {object} api.ErrorResponse "Ошибка клиента"
// @Failure 500 {object} api.ErrorResponse "Ошибка на сервере"
// @Router /auth/register [post]
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	const op = "http_handlers.Register"

	var request api.AuthRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Printf("%s: ошибка при декодировании запроса: %v", op, err)
		Responder(w, http.StatusBadRequest, api.ErrorResponse{Status: false, Message: ErrInvalidJSON.Error()})
		return
	}

	// Первичная валидация
	if request.Login == "" || request.Password == "" {
		Responder(w, http.StatusBadRequest, api.ErrorResponse{Status: false, Message: ErrCredentialsRequired.Error()})
		return
	}

	var referID uint64
	var err error

	referIdStr := r.URL.Query().Get("referID")
	if referIdStr != "" {
		referID, err = strconv.ParseUint(referIdStr, 10, 64)
		if err != nil {
			Responder(w, http.StatusBadRequest, api.ErrorResponse{Status: false, Message: ErrInvalidReferID.Error()})
			return
		}
	}

	// Передаем данные для регистрации в сервис
	regUser, err := h.userService.RegisterUser(r.Context(), request.Login, request.Password, uint(referID))
	if err != nil {
		switch {
		case errors.Is(err, services.ErrReferUserNotFound):
			Responder(w, http.StatusBadRequest, api.ErrorResponse{Status: false, Message: ErrReferUserNotFound.Error()})
			return
		case errors.Is(err, services.ErrUserAlreadyExist):
			Responder(w, http.StatusConflict, api.ErrorResponse{Status: false, Message: ErrUserAlreadyExist.Error()})
			return
		default:
			log.Printf("%s: ошибка при регистрации пользователя: %v", op, err)
			Responder(w, http.StatusInternalServerError, api.ErrorResponse{Status: false, Message: ErrInternalServer.Error()})
			return
		}
	}

	// Успешный ответ
	resp := api.StatusUserResponse{
		Status:  true,
		Message: "Пользователь успешно зарегистрирован",
		User:    regUser,
	}

	Responder(w, http.StatusCreated, resp)
}

// Login godoc
// @Summary Аутентификация пользователя
// @Description Возвращает JWT токен для доступа к защищенным маршрутам.
// @Tags auth
// @Produce json
// @Param request body api.AuthRequest true "Логин и пароль"
// @Success 200 {object} api.LoginResponse "Успешная аутентификация"
// @Failure 403 {object} api.ErrorResponse "Unauthorized"
// @Failure 400 {object} api.ErrorResponse "Ошибка клиента"
// @Failure 500 {object} api.ErrorResponse "Ошибка на сервере"
// @Router /auth/login [post]
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	const op = "http_handlers.Login"

	var request api.AuthRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Printf("%s: ошибка при декодировании запроса: %v", op, err)
		Responder(w, http.StatusBadRequest, api.ErrorResponse{Status: false, Message: ErrInvalidJSON.Error()})
		return
	}

	// Первичная валидация
	if request.Login == "" || request.Password == "" {
		Responder(w, http.StatusBadRequest, api.ErrorResponse{Status: false, Message: ErrCredentialsRequired.Error()})
		return
	}

	err := h.userService.Login(r.Context(), request.Login, request.Password)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrUserNotFound):
			Responder(w, http.StatusNotFound, api.ErrorResponse{Status: false, Message: ErrUserNotFound.Error()})
			return
		case errors.Is(err, services.ErrIncorrectPassword):
			Responder(w, http.StatusUnauthorized, api.ErrorResponse{Status: false, Message: ErrIncorrectPassword.Error()})
			return
		default:
			log.Printf("%s: ошибка при авторизации: %v", op, err)
			Responder(w, http.StatusInternalServerError, api.ErrorResponse{Status: false, Message: ErrInternalServer.Error()})
			return
		}
	}

	tokenStr, err := auth.GenerateToken(request.Login)
	if err != nil {
		log.Printf("%s: ошибка при создании токена jwt: %v", op, err)
		Responder(w, http.StatusInternalServerError, api.ErrorResponse{Status: false, Message: ErrInternalServer.Error()})
		return
	}

	// Успешный ответ
	resp := api.LoginResponse{
		Status:  true,
		Message: "Пользователь успешно авторизирован",
		JWToken: tokenStr,
	}

	Responder(w, http.StatusCreated, resp)
}
