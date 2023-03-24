package tasker_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"tasker/internal/handler"
	"tasker/internal/model"
	"tasker/internal/repository"
	"tasker/internal/service"
	"testing"
)

func TestExecuteTask(t *testing.T) {
	// Create a test task
	task := &model.Task{
		ID:     "test-task",
		Method: "GET",
		URL:    "http://example.com",
		Headers: map[string]string{
			"User-Agent": "Test",
		},
	}

	// Execute the task
	handler.ExecuteTask(task)

	// Check the task status
	if task.Status != "done" {
		t.Errorf("Task status is not 'done', got %s", task.Status)
	}

	// Check the HTTP status code
	if task.HTTPStatusCode != 200 {
		t.Errorf("HTTP status code is not 200, got %d", task.HTTPStatusCode)
	}

	// Check the length of the response body
	if task.Length == 0 {
		t.Errorf("Response body is empty")
	}

	// Check the headers array
	if len(task.HeadersArray) == 0 {
		t.Errorf("Headers array is empty")
	}
}

func TestGetTaskStatus(t *testing.T) {

	var taskHandler *repository.TaskHandler
	c, err := service.InitConfig()

	if err != nil {
		t.Fatalf("Could not read config.yaml %v: ", err)
	}

	s := service.NewService(c, taskHandler)

	err = s.InitTaskHandler()
	if err != nil {
		t.Fatalf("Could not init TaskHandler %v: ", err)
	}

	reqBody := map[string]interface{}{
		"method": "GET",
		"url":    "https://example.com",
	}
	jsonReq, err := json.Marshal(reqBody)
	if err != nil {
		t.Fatalf("failed to marshal JSON request: %v", err)
	}

	createReq, err := http.NewRequest("POST", "/task", bytes.NewReader(jsonReq))
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	createRR := httptest.NewRecorder()
	createHandler := http.HandlerFunc(s.CreateTask)
	createHandler.ServeHTTP(createRR, createReq)

	var createRespBody map[string]interface{}
	if err := json.Unmarshal(createRR.Body.Bytes(), &createRespBody); err != nil {
		t.Fatalf("failed to unmarshal JSON response: %v", err)
	}

	taskID := createRespBody["id"].(string)

	getReq, err := http.NewRequest("GET", "/task/"+taskID, nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	getRR := httptest.NewRecorder()
	getHandler := http.HandlerFunc(s.GetTaskStatus)
	getHandler.ServeHTTP(getRR, getReq)

	if status := getRR.Code; status != http.StatusOK {
		t.Errorf("unexpected status code: got %v, want %v", status, http.StatusOK)
	}

	var respBody map[string]interface{}
	if err := json.Unmarshal(getRR.Body.Bytes(), &respBody); err != nil {
		t.Fatalf("failed to unmarshal JSON response: %v", err)
	}

	if respBody["id"] != taskID {
		t.Errorf("unexpected 'id' field value in response: got %v, want %v", respBody["id"], taskID)
	}

	if _, ok := respBody["status"]; !ok {
		t.Error("expected 'status' field in response body")
	}
}
