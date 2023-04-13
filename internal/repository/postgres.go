package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"tasker/internal/handler"
	"tasker/internal/model"
	"time"

	"github.com/lib/pq"
)

type PostgresTaskCache struct {
	db *sql.DB
}

func NewPostgresTaskCache(c model.PostgresConfig) (*PostgresTaskCache, error) {
	db, err := sql.Open("postgres", c.DSN)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &PostgresTaskCache{
		db: db,
	}, nil
}

func (c *PostgresTaskCache) AddNew(task model.Task) string {
	task.Status = "new"
	task.ID = fmt.Sprintf("%d", time.Now().UnixNano())

	stmt, err := c.db.Prepare("INSERT INTO tasks(id, method, url, headers, status) VALUES($1, $2, $3, $4, $5)")
	if err != nil {
		log.Println("Could not prepare script to insert task :", task.ID)
		return ""
	}

	headers, err := json.Marshal(task.Headers)
	if err != nil {
		log.Println("Could not parse task Headers :", task.ID)
		return ""
	}

	_, err = stmt.Exec(task.ID, task.Method, task.URL, headers, task.Status)
	if err != nil {
		log.Println("Could not insert task :", task.ID)
		return ""
	}
	go c.execTask(&task)

	return task.ID
}

func (c *PostgresTaskCache) GetStatus(taskID string) (*model.Task, bool) {
	row := c.db.QueryRow("SELECT method, url, headers, status, http_status_code, headers_array, length FROM tasks WHERE id = $1", taskID)
	var task model.Task
	var headers []byte
	err := row.Scan(&task.Method, &task.URL, &headers, &task.Status, &task.HTTPStatusCode, pq.Array(&task.HeadersArray), &task.Length)
	if err != nil {
		return nil, false
	}

	err = json.Unmarshal(headers, &task.Headers)
	if err != nil {
		return nil, false
	}

	task.ID = taskID

	return &task, true
}

func (c *PostgresTaskCache) execTask(task *model.Task) {
	task.Status = "in_process"

	handler.ExecuteTask(task)

	stmt, err := c.db.Prepare("UPDATE tasks SET status = $1, http_status_code = $2, headers_array = $3, length = $4 WHERE id = $5")
	if err != nil {
		log.Println("Could not prepare script to update task :", task.ID, err)
		return
	}
	_, err = stmt.Exec(task.Status, task.HTTPStatusCode, pq.Array(task.HeadersArray), task.Length, task.ID)
	if err != nil {
		log.Println("Could not update task :", task.ID, err)
		return
	}
}
