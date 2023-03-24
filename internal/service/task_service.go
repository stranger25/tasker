package service

import (
	"encoding/json"
	"io"
	"net/http"
	"tasker/internal/model"
	"tasker/internal/repository"
)

type Service struct {
	cfg         *model.Config
	taskHandler *repository.TaskHandler
}

// ----------------------------------------------------------------------------------------------------------------------
func NewService(config *model.Config, taskHandler *repository.TaskHandler) *Service {
	return &Service{
		cfg:         config,
		taskHandler: taskHandler,
	}
}

func (s *Service) InitTaskHandler() error {
	var err error
	s.taskHandler, err = repository.NewTaskHandler(s.cfg.Server.Storage, *s.cfg)
	return err
}

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
func (s *Service) CreateTask(w http.ResponseWriter, r *http.Request) {
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

	taskID := s.taskHandler.AddNew(task)

	response, err := json.Marshal(map[string]string{"id": taskID})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//go s.taskHandler.ExecTask(taskID)

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
func (s *Service) GetTaskStatus(w http.ResponseWriter, r *http.Request) {
	taskID := r.URL.Path[len("/task/"):]

	task, ok := s.taskHandler.GetStatus(taskID)

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
