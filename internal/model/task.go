package model

type Task struct {
	ID             string            `json:"id"`
	Method         string            `json:"method"`
	URL            string            `json:"url"`
	Headers        map[string]string `json:"headers"`
	Status         string            `json:"status"`
	HTTPStatusCode int               `json:"httpStatusCode"`
	HeadersArray   []string          `json:"headersArray"`
	Length         int               `json:"length"`
}
