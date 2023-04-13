package repository

import (
	"fmt"
	"sync"
	"tasker/internal/handler"
	"tasker/internal/model"
	"time"
)

type TaskCache struct {
	tasks map[string]*model.Task
	mutex *sync.Mutex
}

func NewTaskCache() *TaskCache {
	return &TaskCache{
		tasks: make(map[string]*model.Task),
		mutex: &sync.Mutex{},
	}
}

func (c *TaskCache) AddNew(task model.Task) string {
	task.Status = "new"
	task.ID = fmt.Sprintf("%d", time.Now().UnixNano())

	c.mutex.Lock()
	c.tasks[task.ID] = &task
	c.mutex.Unlock()

	go c.execTask(&task)

	return task.ID
}

func (c *TaskCache) GetStatus(taskID string) (*model.Task, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	task, ok := c.tasks[taskID]

	return task, ok
}

func (c *TaskCache) execTask(task *model.Task) {
	c.mutex.Lock()
	task.Status = "in_process"
	c.mutex.Unlock()

	handler.ExecuteTask(task)
	c.tasks[task.ID] = task
}
