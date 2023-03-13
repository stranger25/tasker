package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"tasker/internal/model"
	"time"
)

var taskMap = make(map[string]*model.Task)
var mutex = &sync.Mutex{}

// ----------------------------------------------------------------------------------------------------------------------
// @Summary Create task
// @Description Add and execute new task
// @Accept  json
// @Produce  json
// @Param task  body      model.Task  true  "Add new task"
// @Success 200
// @Failure 400
// @Failure 500
// @Router /task [post]
func CreateTask(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var task model.Task
	err = json.Unmarshal(body, &task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	task.Status = "new"
	task.ID = fmt.Sprintf("%d", time.Now().UnixNano())

	mutex.Lock()
	taskMap[task.ID] = &task
	mutex.Unlock()

	go ExecuteTask(&task)

	response, err := json.Marshal(map[string]string{"id": task.ID})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

// ----------------------------------------------------------------------------------------------------------------------
// @Summary Get task status
// @Description Return task status and details
// @Param        taskid   path      int  true  "Task ID"  Format(int)
// @Success 200
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /task/ [get]
func GetTaskStatus(w http.ResponseWriter, r *http.Request) {
	taskID := r.URL.Path[len("/task/"):]

	mutex.Lock()
	task, ok := taskMap[taskID]
	mutex.Unlock()

	if !ok {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	response, err := json.Marshal(task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

// ----------------------------------------------------------------------------------------------------------------------
func ExecuteTask(task *model.Task) {
	task.Status = "in_process"

	client := &http.Client{}
	request, err := http.NewRequest(task.Method, task.URL, nil)
	if err != nil {
		task.Status = "error"
		return
	}

	for key, value := range task.Headers {
		request.Header.Add(key, value)
	}

	response, err := client.Do(request)
	if err != nil {
		task.Status = "error"
		return
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		task.Status = "error"
		return
	}

	task.HTTPStatusCode = response.StatusCode
	task.HeadersArray = []string{}
	for key, values := range response.Header {
		for _, value := range values {
			task.HeadersArray = append(task.HeadersArray, fmt.Sprintf("%s: %s", key, value))
		}
	}
	task.Length = len(body)
	task.Status = "done"
}
