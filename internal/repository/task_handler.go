package repository

import (
	"fmt"
	"tasker/internal/model"
)

type TaskHandlerInterface interface {
	AddNew(task model.Task) string
	GetStatus(taskID string) (*model.Task, bool)
}

type TaskHandler struct {
	handler TaskHandlerInterface
}

func NewTaskHandler(repoType string, c model.Config) (*TaskHandler, error) {
	var handler TaskHandlerInterface
	var err error

	switch repoType {
	case "memory":
		handler = NewTaskCache()
	case "redis":
		handler, err = NewRedisTaskCache(c.Redis)
		if err != nil {
			return nil, err
		}
	case "postgres":
		handler, err = NewPostgresTaskCache(c.Database)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unsupported cache type: %s", repoType)
	}

	return &TaskHandler{
		handler: handler,
	}, nil
}

func (c *TaskHandler) AddNew(task model.Task) string {
	return c.handler.AddNew(task)
}

func (c *TaskHandler) GetStatus(taskID string) (*model.Task, bool) {
	return c.handler.GetStatus(taskID)
}
