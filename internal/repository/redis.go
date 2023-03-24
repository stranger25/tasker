package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"tasker/internal/handler"
	"tasker/internal/model"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisTaskCache struct {
	client *redis.Client
}

func NewRedisTaskCache(c model.RedisConfig) (*RedisTaskCache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     c.Addr,
		Password: c.Password,
		DB:       c.DB,
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	return &RedisTaskCache{
		client: client,
	}, nil
}

func (c *RedisTaskCache) AddNew(task model.Task) string {
	task.Status = "new"
	task.ID = fmt.Sprintf("%d", time.Now().UnixNano())

	jsonTask, err := json.Marshal(task)
	if err != nil {
		return ""
	}

	err = c.client.Set(context.Background(), task.ID, jsonTask, 0).Err()
	if err != nil {
		return ""
	}

	go c.execTask(&task)
	return task.ID
}

func (c RedisTaskCache) GetStatus(taskID string) (*model.Task, bool) {
	jsonTask, err := c.client.Get(context.Background(), taskID).Result()
	if err != nil {
		return nil, false
	}

	var task model.Task
	err = json.Unmarshal([]byte(jsonTask), &task)
	if err != nil {
		return nil, false
	}

	return &task, true
}

func (c *RedisTaskCache) execTask(task *model.Task) {
	task.Status = "in_process"

	handler.ExecuteTask(task)

	jsonBytes, err := json.Marshal(task)
	if err != nil {
		log.Println("Could not parse task :", task.ID)
		return
	}

	err = c.client.Set(context.Background(), task.ID, string(jsonBytes), 0).Err()
	if err != nil {
		log.Println("Could not update task:", task.ID)
		return
	}
}
