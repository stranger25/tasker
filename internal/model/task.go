package model

// ----------------------------------------------------------------------------------------------------------------------
// Task is a struct that contains the data for an individual task.
type Task struct {
	ID             string            `json:"id"`             // ID is a string that contains the task's ID.
	Method         string            `json:"method"`         // Method is a string that contains the HTTP method of the task.
	URL            string            `json:"url"`            // URL is a string that contains the URL of the task.
	Headers        map[string]string `json:"headers"`        // Headers is a map of strings that contains the headers of the task.
	Status         string            `json:"status"`         // Status is a string that contains the status of the task.
	HTTPStatusCode int               `json:"httpStatusCode"` // HTTPStatusCode is an int that contains the HTTP status code of the task.
	HeadersArray   []string          `json:"headersArray"`   // HeadersArray is an array of strings that contains the headers of the task.
	Length         int               `json:"length"`         // Length is an int that contains the length of the task.
}
