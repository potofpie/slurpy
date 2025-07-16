package models

import (
	"encoding/json"
	"net/http"
	"time"
)

// LoggedRequest represents a complete HTTP request/response cycle
type LoggedRequest struct {
	ID        string            `json:"id"`
	Timestamp time.Time         `json:"timestamp"`
	Method    string            `json:"method"`
	URL       string            `json:"url"`
	Headers   map[string]string `json:"headers"`
	Body      string            `json:"body,omitempty"`
	Response  *LoggedResponse   `json:"response,omitempty"`
	Duration  time.Duration     `json:"duration"`
	Namespace string            `json:"namespace"`
	Error     string            `json:"error,omitempty"`
}

// LoggedResponse represents the HTTP response
type LoggedResponse struct {
	StatusCode int               `json:"status_code"`
	Headers    map[string]string `json:"headers"`
	Body       string            `json:"body,omitempty"`
	Size       int64             `json:"size"`
}

// RequestLog represents a collection of logged requests for a namespace
type RequestLog struct {
	Namespace string           `json:"namespace"`
	Requests  []*LoggedRequest `json:"requests"`
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt time.Time        `json:"updated_at"`
}

// ToJSON converts LoggedRequest to JSON bytes
func (lr *LoggedRequest) ToJSON() ([]byte, error) {
	return json.Marshal(lr)
}

// FromJSON creates LoggedRequest from JSON bytes
func FromJSON(data []byte) (*LoggedRequest, error) {
	var req LoggedRequest
	err := json.Unmarshal(data, &req)
	return &req, err
}

// HeadersFromHTTP converts http.Header to map[string]string
func HeadersFromHTTP(h http.Header) map[string]string {
	headers := make(map[string]string)
	for k, v := range h {
		if len(v) > 0 {
			headers[k] = v[0] // Take first value for simplicity
		}
	}
	return headers
}
