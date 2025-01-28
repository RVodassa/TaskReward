package api

import "github.com/RVodassa/TaskReward/internal/domain/models"

type TaskCompletedResponse struct {
	Status  bool
	Message string
	Task    *models.Task
}

type StatusUserResponse struct {
	Status  bool
	Message string
	User    *models.User
}

type LoginResponse struct {
	Status  bool
	Message string
	JWToken string
}

type LeaderBoardResponse struct {
	Status     bool
	Message    string
	ListLeader []*models.User
}

type GetAllTasksResponse struct {
	Status  bool
	Message string
	Tasks   []*models.Task
}

type ErrorResponse struct {
	Status  bool
	Message string
}
