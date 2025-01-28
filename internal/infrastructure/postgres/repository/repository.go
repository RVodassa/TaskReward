package repository

import (
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/RVodassa/TaskReward/internal/domain/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
	"time"
)

var (
	ErrUserAlreadyExist     = errors.New("ошибка: пользователь с таким именем уже существует")
	ErrReferUserNotFound    = errors.New("ошибка: refer с указанным id не найден")
	ErrUserNotFound         = errors.New("ошибка: пользователь найден")
	ErrTaskNotFound         = errors.New("ошибка: задача не найдена")
	ErrTaskAlreadyCompleted = errors.New("ошибка: задача уже выполнена")
)

const (
	StatusTaskClose = "завершено"
	StatusTaskOpen  = "не завершено"
)

type Repo struct {
	db      *pgxpool.Pool
	builder squirrel.StatementBuilderType
}

func NewRepo(db *pgxpool.Pool) *Repo {
	return &Repo{
		db:      db,
		builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (r *Repo) GetAllActiveTask(ctx context.Context) ([]*models.Task, error) {
	const op = "repository.GetAllActiveTask"

	query, args, err := r.builder.
		Select("id", "description", "bonus").
		From("tasks").
		Where(squirrel.Eq{"status": StatusTaskOpen}).
		ToSql()

	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}
	defer rows.Close()

	tasks := make([]*models.Task, 0)
	for rows.Next() {
		task := &models.Task{}
		if err = rows.Scan(&task.ID, &task.Description, &task.Bonus); err != nil {
			return nil, errors.Wrap(err, op)
		}
		tasks = append(tasks, task)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, op)
	}
	return tasks, nil
}

func (r *Repo) GetListTopUsers(ctx context.Context) ([]*models.User, error) {
	const op = "repository.GetListTopUsers"

	query, args, err := r.builder.
		Select("id", "balance").
		From("users").
		OrderBy("balance DESC").
		Limit(10).
		ToSql()

	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		var user models.User
		if err = rows.Scan(&user.ID, &user.Balance); err != nil {
			return nil, errors.Wrap(err, op)
		}
		users = append(users, &user)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, op)
	}

	return users, nil
}

func (r *Repo) TaskComplete(ctx context.Context, taskID uint, userID uint) (*models.Task, error) {
	const op = "repository.TaskComplete"

	if taskID == 0 {
		return nil, errors.New("invalid taskID: taskID cannot be zero")
	}
	if userID == 0 {
		return nil, errors.New("invalid userID: userID cannot be zero")
	}

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	var currentStatus string
	query, args, err := r.builder.Select("status").From("tasks").Where(squirrel.Eq{"id": taskID}).ToSql()
	if err != nil {
		return nil, errors.Wrap(err, op)
	}
	err = tx.QueryRow(ctx, query, args...).Scan(&currentStatus)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrTaskNotFound
		}
		return nil, errors.Wrap(err, op)
	}
	if currentStatus != StatusTaskOpen {
		return nil, ErrTaskAlreadyCompleted
	}

	// Обновляет статус задачи
	query, args, err = r.builder.Update("tasks").
		Set("user_id", userID).
		Set("completed_at", time.Now().UTC()).
		Set("status", StatusTaskClose).
		Where(squirrel.Eq{"id": taskID}).
		Suffix(`RETURNING "id", "user_id", "status", "description", "bonus", "completed_at", "created_at"`).
		ToSql()

	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	var task models.Task
	err = tx.QueryRow(ctx, query, args...).Scan(&task.ID, &task.UserID, &task.Status, &task.Description, &task.Bonus, &task.CompletedAt, &task.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrTaskNotFound
		}
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23503" {
			return nil, ErrUserNotFound
		}
		return nil, errors.Wrap(err, op)
	}

	err = r.IncreaseUserBalance(ctx, tx, userID, task.Bonus)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, errors.Wrap(err, op)
	}

	return &task, nil
}

func (r *Repo) IncreaseUserBalance(ctx context.Context, tx pgx.Tx, userID uint, amount uint) error {
	const op = "repository.IncreaseUserBalance"

	if amount <= 0 {
		return errors.New("invalid amount: amount must be positive")
	}

	query, args, err := r.builder.
		Update("users").
		Set("balance", squirrel.Expr("balance + ?", amount)).
		Where("id = ?", userID).
		ToSql()

	if err != nil {
		return errors.Wrap(err, op)
	}

	_, err = tx.Exec(ctx, query, args...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrUserNotFound
		}
		return errors.Wrap(err, op)
	}

	return nil
}

func (r *Repo) AddTask(ctx context.Context, task *models.Task) error {
	const op = "repository.AddTask"

	query, args, err := r.builder.
		Insert("tasks").
		Columns("description", "bonus", "created_at", "status").
		Values(task.Description, task.Bonus, task.CreatedAt, StatusTaskOpen).
		Suffix(`RETURNING "id"`).
		ToSql()

	if err != nil {
		return errors.Wrap(err, op)
	}

	err = r.db.QueryRow(ctx, query, args...).Scan(&task.ID)
	if err != nil {
		return errors.Wrap(err, op)
	}

	return nil
}

func (r *Repo) GetUserByID(ctx context.Context, id uint) (*models.User, error) {
	const op = "repository.GetUserByID"

	query, args, err := r.builder.
		Select("login", "password_hash", "id", "refer_id", "balance", "created_at").
		From("users").
		Where("id = ?", id).ToSql()

	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	var user models.User
	err = r.db.QueryRow(ctx, query, args...).Scan(
		&user.Login,
		&user.PasswordHash,
		&user.ID,
		&user.ReferID,
		&user.Balance,
		&user.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, errors.Wrap(err, op)
	}
	return &user, nil
}

func (r *Repo) GetUserByLogin(ctx context.Context, login string) (*models.User, error) {
	const op = "repository.GetUserByLogin"

	query, args, err := r.builder.
		Select("login", "password_hash", "id", "refer_id", "balance", "created_at").
		From("users").
		Where("login = ?", login).ToSql()

	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	var user models.User
	err = r.db.QueryRow(ctx, query, args...).Scan(
		&user.Login,
		&user.PasswordHash,
		&user.ID,
		&user.ReferID,
		&user.Balance,
		&user.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, errors.Wrap(err, op)
	}

	return &user, nil
}

func (r *Repo) RegisterUser(ctx context.Context, user *models.User) error {
	const op = "repository.RegisterUser"

	// Проверка входных данных
	if user.Login == "" || user.PasswordHash == "" {
		return errors.New("login and password_hash are required")
	}

	// Начинаем транзакцию
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return errors.Wrap(err, op)
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	// Проверяет существование пользователя с id = refer_id
	if user.ReferID != 0 {
		_, err = r.GetUserByID(ctx, user.ReferID)
		if err != nil {
			if errors.Is(err, ErrUserNotFound) {
				return ErrReferUserNotFound
			}
			return errors.Wrap(err, op)
		}
	}

	// Вставляем нового пользователя
	query, args, err := r.builder.
		Insert("users").
		Columns("login", "password_hash", "refer_id", "created_at").
		Values(user.Login, user.PasswordHash, user.ReferID, user.CreatedAt).
		Suffix(`RETURNING "id"`).
		ToSql()

	if err != nil {
		return errors.Wrap(err, op)
	}

	err = tx.QueryRow(ctx, query, args...).Scan(&user.ID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return fmt.Errorf("%w: login %s already exists", ErrUserAlreadyExist, user.Login)
		}
		return errors.Wrap(err, op)
	}

	// Фиксируем транзакцию
	if err = tx.Commit(ctx); err != nil {
		return errors.Wrap(err, op)
	}

	return nil
}
