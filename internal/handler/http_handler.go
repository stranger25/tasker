package handler

import (
	"fmt"
	"io"
	"net/http"
	"tasker/internal/model"
)

// ----------------------------------------------------------------------------------------------------------------------
// ExecuteTask takes a pointer to a model.Task and updates the task's status, HTTPStatusCode, HeadersArray, and Length.
// It creates an HTTP request with the task's Method and URL, and adds the task's Headers to the request.
// It then makes the request and reads the response body, updating the task's fields accordingly.
// If an error occurs, the task's status is set to "error".
func ExecuteTask(task *model.Task) {
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
