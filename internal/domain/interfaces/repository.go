package interfaces

import (
	"context"
	"github.com/RVodassa/TaskReward/internal/domain/models"
)

type RepositoryProvider interface {
	RegisterUser(ctx context.Context, user *models.User) error
	GetUserByLogin(ctx context.Context, login string) (*models.User, error)
	GetUserByID(ctx context.Context, id uint) (*models.User, error)
	AddTask(ctx context.Context, task *models.Task) error
	TaskComplete(ctx context.Context, taskID uint, userID uint) (*models.Task, error)
	GetListTopUsers(ctx context.Context) ([]*models.User, error)
	GetAllActiveTask(ctx context.Context) ([]*models.Task, error)
}
