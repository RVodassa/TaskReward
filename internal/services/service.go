package services

import (
	"context"
	"fmt"
	"github.com/RVodassa/TaskReward/internal/domain/interfaces"
	"github.com/RVodassa/TaskReward/internal/domain/models"
	repo "github.com/RVodassa/TaskReward/internal/infrastructure/postgres/repository"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"log"
)

var (
	ErrCredentialsRequired  = errors.New("ошибка: логин и пароль обязательны")
	ErrUserAlreadyExist     = errors.New("ошибка: пользователь с таким именем уже существует")
	ErrReferUserNotFound    = errors.New("пользователь с указанным refer_id не найден")
	ErrIncorrectPassword    = errors.New("ошибка: не правильный пароль")
	ErrUserNotFound         = errors.New("ошибка: пользователь найден")
	ErrTaskNotFound         = errors.New("ошибка: задача не найдена")
	ErrTaskAlreadyCompleted = errors.New("ошибка: задача уже выполнена")
)

type Service struct {
	repo interfaces.RepositoryProvider
}

func NewService(repo interfaces.RepositoryProvider) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) GetAllActiveTask(ctx context.Context) ([]*models.Task, error) {
	const op = "service.GetAllActiveTask"

	tasks, err := s.repo.GetAllActiveTask(ctx)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}
	return tasks, nil
}

func (s *Service) GetListTopUsers(ctx context.Context) ([]*models.User, error) {
	const op = "services.GetListTopUsers"

	users, err := s.repo.GetListTopUsers(ctx)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	return users, nil
}

func (s *Service) TaskComplete(ctx context.Context, taskID uint, userID uint) (*models.Task, error) {
	const op = "services.TaskComplete"

	task, err := s.repo.TaskComplete(ctx, taskID, userID)
	if err != nil {
		switch {
		case errors.Is(err, repo.ErrTaskNotFound):
			return nil, ErrTaskNotFound
		case errors.Is(err, repo.ErrUserNotFound):
			return nil, ErrUserNotFound
		case errors.Is(err, repo.ErrTaskAlreadyCompleted):
			return nil, ErrTaskAlreadyCompleted
		default:
			return nil, errors.Wrap(err, op)
		}
	}

	return task, nil
}

func (s *Service) StatusUser(ctx context.Context, userID uint) (*models.User, error) {
	const op = "services.StatusUser"

	getUser, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, repo.ErrUserNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, errors.Wrap(err, op)
	}

	return getUser, nil
}

func (s *Service) Login(ctx context.Context, login, password string) error {
	const op = "services.Service.Login"

	// Получение пользователя по логину
	getUser, err := s.repo.GetUserByLogin(ctx, login)
	if err != nil {
		if errors.Is(err, repo.ErrUserNotFound) {
			return ErrUserNotFound
		}
		return errors.Wrap(err, op)
	}

	// Проверка пароля
	err = checkPassword(getUser.PasswordHash, password)
	if err != nil {
		if errors.Is(err, ErrIncorrectPassword) {
			return ErrIncorrectPassword
		}
		return errors.Wrap(err, op)
	}

	return nil
}

func (s *Service) RegisterUser(ctx context.Context, login, password string, referID uint) (*models.User, error) {
	const op = "services.Service.RegisterUser"

	if login == "" || password == "" {
		return nil, ErrCredentialsRequired
	}

	// Хэшируем пароль
	hashedPassword, err := hashPassword(password)
	if err != nil {
		log.Printf("%s: ошибка при хэшировании пароля: %v", op, err)
		return nil, errors.Wrap(err, op)
	}

	// Новый инстанс пользователя
	user := models.NewUser(login, hashedPassword, referID)

	// Регистрируем пользователя в репозитории
	if err = s.repo.RegisterUser(ctx, user); err != nil {
		switch {
		case errors.Is(err, repo.ErrUserAlreadyExist): // если уже существует
			return nil, ErrUserAlreadyExist
		case errors.Is(err, repo.ErrReferUserNotFound): // если refer_id не найден
			return nil, ErrReferUserNotFound
		default: // другая ошибка
			log.Printf("%s: ошибка при регистрации пользователя: %v", op, err)
			return nil, errors.Wrap(err, op)
		}
	}

	return user, nil
}

func (s *Service) AddTask(ctx context.Context, description string, bonus uint) error {
	const op = "services.Service.AddTask"

	// Новый инстанс пользователя
	task := models.NewTask(description, bonus)

	err := s.repo.AddTask(ctx, task)
	if err != nil {
		return errors.Wrap(err, op)
	}

	return nil
}

func checkPassword(hashedPassword, password string) error {
	const op = "services.checkPassword"

	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return ErrIncorrectPassword
		}
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

// HashPassword хэширует пароль с использованием bcrypt.
// Возвращает хэшированный пароль или ошибку, если хэширование не удалось.
func hashPassword(password string) (string, error) {
	const op = "services.HashPassword"

	// Генерация хэша с использованием bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return string(hashedPassword), nil
}
