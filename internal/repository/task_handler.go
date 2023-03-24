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

// NewTaskHandler creates a new TaskHandler based on the given repository type and configuration.
// It will create a new TaskCache if the repoType is "memory", a new RedisTaskCache if the repoType is "redis",
// or a new PostgresTaskCache if the repoType is "postgres".
// If an unsupported repoType is given, an error is returned.
// It returns a pointer to the newly created TaskHandler, or an error if one occurred.
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
		handler, err = NewPostgresTaskCache(c.Postgres)
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
